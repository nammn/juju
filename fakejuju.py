from argparse import ArgumentParser
from base64 import b64encode
from contextlib import contextmanager
import copy
from hashlib import sha512
from itertools import count
import json
import logging
import subprocess

import yaml

from jujupy import (
    EnvJujuClient,
    JESNotSupported,
    JujuData,
)

__metaclass__ = type


def check_juju_output(func):
    def wrapper(*args, **kwargs):
        result = func(*args, **kwargs)
        if 'service' in result:
            raise AssertionError('Result contained service')
        return result
    return wrapper


class ControllerOperation(Exception):

    def __init__(self, operation):
        super(ControllerOperation, self).__init__(
            'Operation "{}" is only valid on controller models.'.format(
                operation))


def assert_juju_call(test_case, mock_method, client, expected_args,
                     call_index=None):
    if call_index is None:
        test_case.assertEqual(len(mock_method.mock_calls), 1)
        call_index = 0
    empty, args, kwargs = mock_method.mock_calls[call_index]
    test_case.assertEqual(args, (expected_args,))


class FakeControllerState:

    def __init__(self):
        self.state = 'not-bootstrapped'
        self.models = {}

    def add_model(self, name):
        state = FakeEnvironmentState()
        state.name = name
        self.models[name] = state
        state.controller = self
        state.controller.state = 'created'
        return state

    def require_controller(self, operation, name):
        if name != self.controller_model.name:
            raise ControllerOperation(operation)

    def bootstrap(self, model_name, config, separate_controller):
        default_model = self.add_model(model_name)
        default_model.name = model_name
        if separate_controller:
            controller_model = default_model.controller.add_model('controller')
        else:
            controller_model = default_model
        self.controller_model = controller_model
        controller_model.state_servers.append(controller_model.add_machine())
        self.state = 'bootstrapped'
        default_model.model_config = copy.deepcopy(config)
        self.models[default_model.name] = default_model
        return default_model


class FakeEnvironmentState:
    """A Fake environment state that can be used by multiple FakeClients."""

    def __init__(self, controller=None):
        self._clear()
        if controller is not None:
            self.controller = controller
        else:
            self.controller = FakeControllerState()

    def _clear(self):
        self.name = None
        self.machine_id_iter = count()
        self.state_servers = []
        self.services = {}
        self.machines = set()
        self.containers = {}
        self.relations = {}
        self.token = None
        self.exposed = set()
        self.machine_host_names = {}
        self.current_bundle = None
        self.model_config = None

    @property
    def state(self):
        return self.controller.state

    def add_machine(self, host_name=None, machine_id=None):
        if machine_id is None:
            machine_id = str(self.machine_id_iter.next())
        self.machines.add(machine_id)
        if host_name is None:
            host_name = '{}.example.com'.format(machine_id)
        self.machine_host_names[machine_id] = host_name
        return machine_id

    def add_ssh_machines(self, machines):
        for machine in machines:
            self.add_machine()

    def add_container(self, container_type, host=None, container_num=None):
        if host is None:
            host = self.add_machine()
        host_containers = self.containers.setdefault(host, set())
        if container_num is None:
            same_type_containers = [x for x in host_containers if
                                    container_type in x]
            container_num = len(same_type_containers)
        container_name = '{}/{}/{}'.format(host, container_type, container_num)
        host_containers.add(container_name)
        host_name = '{}.example.com'.format(container_name)
        self.machine_host_names[container_name] = host_name

    def remove_container(self, container_id):
        for containers in self.containers.values():
            containers.discard(container_id)

    def remove_machine(self, machine_id):
        self.machines.remove(machine_id)
        self.containers.pop(machine_id, None)

    def remove_state_server(self, machine_id):
        self.remove_machine(machine_id)
        self.state_servers.remove(machine_id)

    def destroy_environment(self):
        self._clear()
        self.controller.state = 'destroyed'
        return 0

    def kill_controller(self):
        self._clear()
        self.controller.state = 'controller-killed'

    def destroy_model(self):
        del self.controller.models[self.name]
        self._clear()
        self.controller.state = 'model-destroyed'

    def restore_backup(self):
        self.controller.require_controller('restore', self.name)
        if len(self.state_servers) > 0:
            exc = subprocess.CalledProcessError('Operation not permitted', 1,
                                                2)
            exc.stderr = 'Operation not permitted'
            raise exc

    def enable_ha(self):
        self.controller.require_controller('enable-ha', self.name)
        for n in range(2):
            self.state_servers.append(self.add_machine())

    def deploy(self, charm_name, service_name):
        self.add_unit(service_name)

    def deploy_bundle(self, bundle_path):
        self.current_bundle = bundle_path

    def add_unit(self, service_name):
        machines = self.services.setdefault(service_name, set())
        machines.add(
            ('{}/{}'.format(service_name, str(len(machines))),
             self.add_machine()))

    def remove_unit(self, to_remove):
        for units in self.services.values():
            for unit_id, machine_id in units:
                if unit_id == to_remove:
                    self.remove_machine(machine_id)
                    units.remove((unit_id, machine_id))
                    break

    def destroy_service(self, service_name):
        for unit, machine_id in self.services.pop(service_name):
            self.remove_machine(machine_id)

    def get_status_dict(self):
        machines = {}
        for machine_id in self.machines:
            machine_dict = {'juju-status': {'current': 'idle'}}
            hostname = self.machine_host_names.get(machine_id)
            machine_dict['instance-id'] = machine_id
            if hostname is not None:
                machine_dict['dns-name'] = hostname
            machines[machine_id] = machine_dict
            if machine_id in self.state_servers:
                machine_dict['controller-member-status'] = 'has-vote'
        for host, containers in self.containers.items():
            container_dict = dict((c, {}) for c in containers)
            for container in containers:
                dns_name = self.machine_host_names.get(container)
                if dns_name is not None:
                    container_dict[container]['dns-name'] = dns_name

            machines[host]['containers'] = container_dict
        services = {}
        for service, units in self.services.items():
            unit_map = {}
            for unit_id, machine_id in units:
                unit_map[unit_id] = {
                    'machine': machine_id,
                    'juju-status': {'current': 'idle'}}
            services[service] = {
                'units': unit_map,
                'relations': self.relations.get(service, {}),
                'exposed': service in self.exposed,
                }
        return {'machines': machines, 'applications': services}


class AutoloadCredentials:

    def __init__(self, backend, juju_home, extra_env):
        self.backend = backend
        self.juju_home = juju_home
        self.extra_env = extra_env
        self.last_expect = None
        self.cloud = None

    def expect(self, line):
        self.last_expect = line

    def sendline(self, line):
        if self.last_expect == (
                'Enter cloud to which the credential belongs, or Q to quit.*'):
            self.cloud = line

    def isalive(self):
        juju_data = JujuData('foo', juju_home=self.juju_home)
        juju_data.load_yaml()
        creds = juju_data.credentials.setdefault('credentials', {})
        creds.update({self.cloud: {self.extra_env['OS_USERNAME']: {
            'domain-name': '',
            'auth-type': 'userpass',
            'username': self.extra_env['OS_USERNAME'],
            'password': self.extra_env['OS_PASSWORD'],
            'tenant-name': self.extra_env['OS_TENANT_NAME'],
            }}})
        juju_data.dump_yaml(self.juju_home, {})
        return False


class FakeBackend:
    """A fake juju backend for tests.

    This is a partial implementation, but should be suitable for many uses,
    and can be extended.

    The state is provided by controller_state, so that multiple clients and
    backends can manipulate the same state.
    """

    def __init__(self, controller_state, feature_flags=None, version=None,
                 full_path=None, debug=False):
        assert isinstance(controller_state, FakeControllerState)
        self.controller_state = controller_state
        if feature_flags is None:
            feature_flags = set()
        self.feature_flags = feature_flags
        self.version = version
        self.full_path = full_path
        self.debug = debug
        self.juju_timings = {}
        self.log = logging.getLogger('jujupy')

    def clone(self, full_path=None, version=None, debug=None,
              feature_flags=None):
        if version is None:
            version = self.version
        if full_path is None:
            full_path = self.full_path
        if debug is None:
            debug = self.debug
        if feature_flags is None:
            feature_flags = set(self.feature_flags)
        controller_state = self.controller_state
        return self.__class__(controller_state, feature_flags, version,
                              full_path, debug)

    def set_feature(self, feature, enabled):
        if enabled:
            self.feature_flags.add(feature)
        else:
            self.feature_flags.discard(feature)

    def is_feature_enabled(self, feature):
        if feature == 'jes':
            return True
        return bool(feature in self.feature_flags)

    def deploy(self, model_state, charm_name, service_name=None, series=None):
        if service_name is None:
            service_name = charm_name.split(':')[-1].split('/')[-1]
        model_state.deploy(charm_name, service_name)

    def bootstrap(self, args):
        parser = ArgumentParser()
        parser.add_argument('controller_name')
        parser.add_argument('cloud_name_region')
        parser.add_argument('--constraints')
        parser.add_argument('--config')
        parser.add_argument('--default-model')
        parser.add_argument('--agent-version')
        parser.add_argument('--bootstrap-series')
        parser.add_argument('--upload-tools', action='store_true')
        parsed = parser.parse_args(args)
        with open(parsed.config) as config_file:
            config = yaml.safe_load(config_file)
        cloud_region = parsed.cloud_name_region.split('/', 1)
        cloud = cloud_region[0]
        # Although they are specified with specific arguments instead of as
        # config, these values are listed by get-model-config:
        # name, region, type (from cloud).
        config['type'] = cloud
        if len(cloud_region) > 1:
            config['region'] = cloud_region[1]
        config['name'] = parsed.default_model
        if parsed.bootstrap_series is not None:
            config['default-series'] = parsed.bootstrap_series
        self.controller_state.bootstrap(parsed.default_model, config,
                                        self.is_feature_enabled('jes'))

    def quickstart(self, model_name, config, bundle):
        default_model = self.controller_state.bootstrap(
            model_name, config, self.is_feature_enabled('jes'))
        default_model.deploy_bundle(bundle)

    def destroy_environment(self, model_name):
        try:
            state = self.controller_state.models[model_name]
        except KeyError:
            return 0
        state.destroy_environment()
        return 0

    def add_machines(self, model_state, args):
        if len(args) == 0:
            return model_state.add_machine()
        ssh_machines = [a[4:] for a in args if a.startswith('ssh:')]
        if len(ssh_machines) == len(args):
            return model_state.add_ssh_machines(ssh_machines)
        parser = ArgumentParser()
        parser.add_argument('host_placement', nargs='*')
        parser.add_argument('-n', type=int, dest='count', default='1')
        parsed = parser.parse_args(args)
        if len(parsed.host_placement) == 1:
            split = parsed.host_placement[0].split(':')
            if len(split) == 1:
                container_type = split[0]
                host = None
            else:
                container_type, host = split
            for x in range(parsed.count):
                model_state.add_container(container_type, host=host)
        else:
            for x in range(parsed.count):
                model_state.add_machine()

    def get_controller_model_name(self):
        return self.controller_state.controller_model.name

    def make_controller_dict(self, controller_name):
        controller_model = self.controller_state.controller_model
        server_id = list(controller_model.state_servers)[0]
        server_hostname = controller_model.machine_host_names[server_id]
        api_endpoint = '{}:23'.format(server_hostname)
        return {controller_name: {'details': {'api-endpoints': [
            api_endpoint]}}}

    def list_models(self):
        model_names = [state.name for state in
                       self.controller_state.models.values()]
        return {'models': [{'name': n} for n in model_names]}

    def _log_command(self, command, args, model, level=logging.INFO):
        full_args = ['juju', command]
        if model is not None:
            full_args.extend(['-m', model])
        full_args.extend(args)
        self.log.log(level, ' '.join(full_args))

    def juju(self, command, args, used_feature_flags,
             juju_home, model=None, check=True, timeout=None, extra_env=None):
        if 'service' in command:
            raise Exception('Command names must not contain "service".')
        if isinstance(args, basestring):
            args = (args,)
        self._log_command(command, args, model)
        if model is not None:
            if ':' in model:
                model = model.split(':')[1]
            model_state = self.controller_state.models[model]
            if command == 'enable-ha':
                model_state.enable_ha()
            if (command, args[:1]) == ('set-config', ('dummy-source',)):
                name, value = args[1].split('=')
                if name == 'token':
                    model_state.token = value
            if command == 'deploy':
                parser = ArgumentParser()
                parser.add_argument('charm_name')
                parser.add_argument('service_name', nargs='?')
                parser.add_argument('--to')
                parser.add_argument('--series')
                parsed = parser.parse_args(args)
                self.deploy(model_state, parsed.charm_name,
                            parsed.service_name, parsed.series)
            if command == 'remove-application':
                model_state.destroy_service(*args)
            if command == 'add-relation':
                if args[0] == 'dummy-source':
                    model_state.relations[args[1]] = {'source': [args[0]]}
            if command == 'expose':
                (service,) = args
                model_state.exposed.add(service)
            if command == 'unexpose':
                (service,) = args
                model_state.exposed.remove(service)
            if command == 'add-unit':
                (service,) = args
                model_state.add_unit(service)
            if command == 'remove-unit':
                (unit_id,) = args
                model_state.remove_unit(unit_id)
            if command == 'add-machine':
                return self.add_machines(model_state, args)
            if command == 'remove-machine':
                parser = ArgumentParser()
                parser.add_argument('machine_id')
                parser.add_argument('--force', action='store_true')
                parsed = parser.parse_args(args)
                machine_id = parsed.machine_id
                if '/' in machine_id:
                    model_state.remove_container(machine_id)
                else:
                    model_state.remove_machine(machine_id)
            if command == 'quickstart':
                parser = ArgumentParser()
                parser.add_argument('--constraints')
                parser.add_argument('--no-browser', action='store_true')
                parser.add_argument('bundle')
                parsed = parser.parse_args(args)
                # Released quickstart doesn't seem to provide the config via
                # the commandline.
                self.quickstart(model, {}, parsed.bundle)
        else:
            if command == 'bootstrap':
                self.bootstrap(args)
            if command == 'kill-controller':
                if self.controller_state.state == 'not-bootstrapped':
                    return
                model = args[0]
                model_state = self.controller_state.models[model]
                model_state.kill_controller()
            if command == 'destroy-model':
                if not self.is_feature_enabled('jes'):
                    raise JESNotSupported()
                model = args[0]
                model_state = self.controller_state.models[model]
                model_state.destroy_model()
            if command == 'add-model':
                if not self.is_feature_enabled('jes'):
                    raise JESNotSupported()
                parser = ArgumentParser()
                parser.add_argument('-c', '--controller')
                parser.add_argument('--config')
                parser.add_argument('model_name')
                parsed = parser.parse_args(args)
                self.controller_state.add_model(parsed.model_name)

    @contextmanager
    def juju_async(self, command, args, used_feature_flags,
                   juju_home, model=None, timeout=None):
        yield
        self.juju(command, args, used_feature_flags,
                  juju_home, model, timeout=timeout)

    @check_juju_output
    def get_juju_output(self, command, args, used_feature_flags,
                        juju_home, model=None, timeout=None):
        if 'service' in command:
            raise Exception('No service')
        self._log_command(command, args, model, logging.DEBUG)
        if model is not None:
            if ':' in model:
                model = model.split(':')[1]
            model_state = self.controller_state.models[model]
        from deploy_stack import GET_TOKEN_SCRIPT
        if (command, args) == ('ssh', ('dummy-sink/0', GET_TOKEN_SCRIPT)):
            return model_state.token
        if (command, args) == ('ssh', ('0', 'lsb_release', '-c')):
            return 'Codename:\t{}\n'.format(
                model_state.model_config['default-series'])
        if command == 'get-model-config':
            return yaml.safe_dump(model_state.model_config)
        if command == 'restore-backup':
            model_state.restore_backup()
        if command == 'show-controller':
            return yaml.safe_dump(self.make_controller_dict(args[0]))
        if command == 'list-models':
            return yaml.safe_dump(self.list_models())
        if command == 'add-user':
            permissions = 'read'
            if set(["--acl", "write"]).issubset(args):
                permissions = 'write'
            username = args[0]
            model = args[2]
            info_string = \
                'User "{}" added\nUser "{}"granted {} access to model "{}\n"' \
                .format(username, username, permissions, model)
            register_string = get_user_register_command_info(username)
            return info_string + register_string
        if command == 'show-status':
            status_dict = model_state.get_status_dict()
            # Parsing JSON is much faster than parsing YAML, and JSON is a
            # subset of YAML, so emit JSON.
            return json.dumps(status_dict)
        if command == 'create-backup':
            self.controller_state.require_controller('backup', model)
            return 'juju-backup-0.tar.gz'
        return ''

    def expect(self, command, args, used_feature_flags, juju_home, model=None,
               timeout=None, extra_env=None):
        if command == 'autoload-credentials':
            return AutoloadCredentials(self, juju_home, extra_env)

    def pause(self, seconds):
        pass


def get_user_register_command_info(username):
    code = get_user_register_token(username)
    return 'Please send this command to {}\n    juju register {}'.format(
        username, code)


def get_user_register_token(username):
    return b64encode(sha512(username).digest())


class FakeBackend2B7(FakeBackend):

    def juju(self, command, args, used_feature_flags,
             juju_home, model=None, check=True, timeout=None, extra_env=None):
        if model is not None:
            model_state = self.controller_state.models[model]
        if command == 'destroy-service':
            model_state.destroy_service(*args)
        if command == 'remove-service':
            model_state.destroy_service(*args)
        return super(FakeBackend2B7).juju(command, args, used_feature_flags,
                                          juju_home, model, check, timeout,
                                          extra_env)


class FakeBackendOptionalJES(FakeBackend):

    def is_feature_enabled(self, feature):
        return bool(feature in self.feature_flags)


def fake_juju_client(env=None, full_path=None, debug=False, version='2.0.0',
                     _backend=None, cls=EnvJujuClient):
    if env is None:
        env = JujuData('name', {
            'type': 'foo',
            'default-series': 'angsty',
            'region': 'bar',
            }, juju_home='foo')
    juju_home = env.juju_home
    if juju_home is None:
        juju_home = 'foo'
    if _backend is None:
        backend_state = FakeControllerState()
        _backend = FakeBackend(
            backend_state, version=version, full_path=full_path,
            debug=debug)
        _backend.set_feature('jes', True)
    client = cls(
        env, version, full_path, juju_home, debug, _backend=_backend)
    client.bootstrap_replaces = {}
    return client


def fake_juju_client_optional_jes(env=None, full_path=None, debug=False,
                                  jes_enabled=True, version='2.0.0',
                                  _backend=None):
    if _backend is None:
        backend_state = FakeControllerState()
        _backend = FakeBackendOptionalJES(
            backend_state, version=version, full_path=full_path,
            debug=debug)
        _backend.set_feature('jes', jes_enabled)
    client = fake_juju_client(env, full_path, debug, version, _backend,
                              cls=FakeJujuClientOptionalJES)
    client.used_feature_flags = frozenset(['address-allocation', 'jes'])
    return client


class FakeJujuClientOptionalJES(EnvJujuClient):

    def get_controller_model_name(self):
        return self._backend.controller_state.controller_model.name
