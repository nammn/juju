from datetime import (
    datetime,
    timedelta,
)
from mock import (
    Mock,
    patch,
)
from unittest import TestCase

from winazure import (
    DATETIME_PATTERN,
    delete_unused_disks,
    is_old_deployment,
    list_services,
    main,
    ServiceManagementService,
    SUCCEEDED,
    wait_for_success,
)


class WinAzureTestCase(TestCase):

    def test_main_delete_unused_disks(self):
        with patch('winazure.delete_unused_disks', autospec=True) as d_mock:
            with patch('winazure.ServiceManagementService',
                       autospec=True, return_value='sms') as sms_mock:
                main(['-d', '-v', '-c' 'cert.pem', '-s', 'secret',
                      'delete-unused-disks'])
        sms_mock.assert_called_once_with('secret', 'cert.pem')
        d_mock.assert_called_once_with('sms', dry_run=True, verbose=True)

    def test_main_list_services(self):
        with patch('winazure.list_services', autospec=True) as l_mock:
            with patch('winazure.ServiceManagementService',
                       autospec=True, return_value='sms') as sms_mock:
                main(['-d', '-v', '-c' 'cert.pem', '-s', 'secret',
                      'list-services', 'juju-deploy*'])
        sms_mock.assert_called_once_with('secret', 'cert.pem')
        l_mock.assert_called_once_with(
            'sms', glob='juju-deploy*', verbose=True)

    def test_main_delete_services(self):
        with patch('winazure.delete_services', autospec=True) as l_mock:
            with patch('winazure.ServiceManagementService',
                       autospec=True, return_value='sms') as sms_mock:
                main(['-d', '-v', '-c' 'cert.pem', '-s', 'secret',
                      'delete-services', '-o', '2', 'juju-deploy*'])
        sms_mock.assert_called_once_with('secret', 'cert.pem')
        l_mock.assert_called_once_with(
            'sms', glob='juju-deploy*', old_age=2, verbose=True, dry_run=True)

    def test_delete_unused_disks(self):
        sms = ServiceManagementService('secret', 'cert.pem')
        disk1 = Mock()
        disk1.name = 'disk1'
        disk1.attached_to = None
        disk2 = Mock()
        disk2.name = 'disk2'
        disk2.attached_to.hosted_service_name = ''
        disk3 = Mock()
        disk3.name = 'disk3'
        disk3.attached_to.hosted_service_name = 'hs3'
        disks = [disk1, disk2, disk3]
        with patch.object(sms, 'list_disks', autospec=True,
                          return_value=disks) as ld_mock:
            with patch.object(sms, 'delete_disk', autospec=True) as dd_mock:
                delete_unused_disks(sms, dry_run=False, verbose=False)
        ld_mock.assert_called_once_with()
        dd_mock.assert_any_call('disk1', delete_vhd=True)
        dd_mock.assert_any_call('disk2', delete_vhd=True)
        self.assertEqual(2, dd_mock.call_count)

    def test_list_services(self):
        sms = ServiceManagementService('secret', 'cert.pem')
        hs1 = Mock()
        hs1.service_name = 'juju-upgrade-foo'
        hs2 = Mock()
        hs2.service_name = 'juju-deploy-foo'
        services = [hs1, hs2]
        with patch.object(sms, 'list_hosted_services', autospec=True,
                          return_value=services) as ls_mock:
            services = list_services(sms, 'juju-deploy-*', verbose=False)
        ls_mock.assert_called_once_with()
        self.assertEqual([hs2], services)

    def test_is_old_deployment(self):
        now = datetime.utcnow()
        ago = timedelta(hours=2)
        d1 = Mock()
        d1.created_time = (now - timedelta(hours=3)).strftime(DATETIME_PATTERN)
        self.assertTrue(
            is_old_deployment([d1], now, ago, verbose=False))
        d2 = Mock()
        d2.created_time = (now - timedelta(hours=1)).strftime(DATETIME_PATTERN)
        self.assertFalse(
            is_old_deployment([d2], now, ago, verbose=False))
        self.assertTrue(
            is_old_deployment([d1, d2], now, ago, verbose=False))

    def test_wait_for_success(self):
        sms = ServiceManagementService('secret', 'cert.pem')
        request = Mock()
        request.request_id = 'foo'
        op1 = Mock()
        op1.status = 'Pending'
        op2 = Mock()
        op2.status = SUCCEEDED
        op3 = Mock()
        op3.status = 'Not Reachable'
        with patch.object(sms, 'get_operation_status', autospec=True,
                          side_effect=[op1, op2, op3]) as gs_mock:
            wait_for_success(sms, request, pause=0, verbose=False)
        self.assertEqual(2, gs_mock.call_count)
