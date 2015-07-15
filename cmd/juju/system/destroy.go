// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package system

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/juju/cmd"
	"github.com/juju/errors"
	"launchpad.net/gnuflag"

	"github.com/juju/juju/api/systemmanager"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/cmd/envcmd"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/environs/configstore"
	"github.com/juju/juju/juju"
)

// DestroyCommand destroys the specified system.
type DestroyCommand struct {
	envcmd.SysCommandBase
	systemName  string
	assumeYes   bool
	destroyEnvs bool

	// The following fields are for mocking out
	// api behavior for testing.
	api       destroySystemAPI
	apierr    error
	clientapi destroyClientAPI
}

var destroyDoc = `Destroys the specified system`
var destroySysMsg = `
WARNING! This command will destroy the %q system.
This includes all machines, services, data and other resources.

Continue [y/N]? `[1:]

// destroySystemAPI defines the methods on the system manager API endpoint
// that the destroy command calls.
type destroySystemAPI interface {
	Close() error
	EnvironmentConfig() (map[string]interface{}, error)
	DestroySystem(destroyEnvs bool, ignoreBlocks bool) error
}

// destroyClientAPI defines the methods on the client API endpoint that the
// destroy command might call.
type destroyClientAPI interface {
	Close() error
	EnvironmentGet() (map[string]interface{}, error)
	DestroyEnvironment() error
}

// Info implements Command.Info.
func (c *DestroyCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "destroy",
		Args:    "<system name>",
		Purpose: "terminate all machines and other associated resources for a system environment",
		Doc:     destroyDoc,
	}
}

// SetFlags implements Command.SetFlags.
func (c *DestroyCommand) SetFlags(f *gnuflag.FlagSet) {
	f.BoolVar(&c.assumeYes, "y", false, "Do not ask for confirmation")
	f.BoolVar(&c.assumeYes, "yes", false, "")
	f.BoolVar(&c.destroyEnvs, "destroy-all-environments", false, "destroy all hosted environments on the system")
}

// Init implements Command.Init.
func (c *DestroyCommand) Init(args []string) error {
	switch len(args) {
	case 0:
		return errors.New("no system specified")
	case 1:
		c.systemName = args[0]
		return nil
	default:
		return cmd.CheckEmpty(args[1:])
	}
}

func (c *DestroyCommand) getSystemAPI() (destroySystemAPI, error) {
	if c.api != nil {
		return c.api, c.apierr
	}
	root, err := juju.NewAPIFromName(c.systemName)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return systemmanager.NewClient(root), nil
}

func (c *DestroyCommand) getClientAPI() (destroyClientAPI, error) {
	if c.clientapi != nil {
		return c.clientapi, nil
	}
	root, err := juju.NewAPIFromName(c.systemName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return root.Client(), nil
}

// Run implements Command.Run
func (c *DestroyCommand) Run(ctx *cmd.Context) error {
	store, err := configstore.Default()
	if err != nil {
		return errors.Annotate(err, "cannot open system info storage")
	}

	cfgInfo, err := store.ReadInfo(c.systemName)
	if err != nil {
		return errors.Annotate(err, "cannot read system info")
	}

	// Verify that we're destroying a system
	apiEndpoint := cfgInfo.APIEndpoint()
	if apiEndpoint.ServerUUID != "" && apiEndpoint.EnvironUUID != apiEndpoint.ServerUUID {
		return errors.Errorf("%q is not a system; use juju environment destroy to destroy it", c.systemName)
	}

	if !c.assumeYes {
		if err = confirmDestruction(ctx, c.systemName); err != nil {
			return err
		}
	}

	// Attempt to connect to the API.  If we can't, fail the destroy.  Users will
	// need to use the system kill command if we can't connect.
	api, err := c.getSystemAPI()
	if err != nil {
		return c.ensureUserFriendlyErrorLog(errors.Annotate(err, "cannot connect to API"))
	}
	defer api.Close()

	// Obtain bootstrap / system environ information
	systemEnviron, err := c.getSystemEnviron(cfgInfo, api)
	if err != nil {
		return errors.Annotate(err, "cannot obtain bootstrap information")
	}

	// Attempt to destroy the system.
	err = api.DestroySystem(c.destroyEnvs, false)
	if params.IsCodeNotImplemented(err) {
		// Fall back to using the client endpoint to destroy the system,
		// sending the info we were already able to collect.
		return c.destroySystemViaClient(ctx, cfgInfo, systemEnviron, store)
	}
	if err != nil {
		return c.ensureUserFriendlyErrorLog(errors.Annotate(err, "cannot destroy system"))
	}

	return environs.Destroy(systemEnviron, store)
}

// getSystemEnviron gets the bootstrap information required to destroy the environment
// by first checking the config store, then querying the API if the information is not
// in the store.
func (c *DestroyCommand) getSystemEnviron(info configstore.EnvironInfo, sysAPI destroySystemAPI) (_ environs.Environ, err error) {
	bootstrapCfg := info.BootstrapConfig()
	if bootstrapCfg == nil {
		if sysAPI == nil {
			return nil, errors.New("unable to get bootstrap information from API")
		}
		bootstrapCfg, err = sysAPI.EnvironmentConfig()
		if params.IsCodeNotImplemented(err) {
			// Fallback to the client API. Better to encapsulate the logic for
			// old servers than worry about connecting twice.
			client, err := c.getClientAPI()
			if err != nil {
				return nil, errors.Trace(err)
			}
			defer client.Close()
			bootstrapCfg, err = client.EnvironmentGet()
			if err != nil {
				return nil, errors.Trace(err)
			}
		} else if err != nil {
			return nil, errors.Trace(err)
		}
	}

	cfg, err := config.New(config.NoDefaults, bootstrapCfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return environs.New(cfg)
}

// destroySystemViaClient attempts to destroy the system using the client
// endpoint for older juju systems which do not implement systemmanager.DestroySystem
func (c *DestroyCommand) destroySystemViaClient(ctx *cmd.Context, info configstore.EnvironInfo, systemEnviron environs.Environ, store configstore.Storage) error {
	api, err := c.getClientAPI()
	if err != nil {
		return c.ensureUserFriendlyErrorLog(errors.Annotate(err, "cannot connect to API"))
	}
	defer api.Close()

	err = api.DestroyEnvironment()
	if err != nil {
		return c.ensureUserFriendlyErrorLog(errors.Annotate(err, "cannot destroy system"))
	}

	return environs.Destroy(systemEnviron, store)
}

// ensureUserFriendlyErrorLog ensures that error will be logged and displayed
// in a user-friendly manner with readable and digestable error message.
func (c *DestroyCommand) ensureUserFriendlyErrorLog(err error) error {
	if err == nil {
		return nil
	}
	if params.IsCodeOperationBlocked(err) {
		logger.Errorf(`there are blocks preventing system destruction
To remove all blocks in the system, please run:
    juju system remove-blocks
`)
		return cmd.ErrSilent
	}
	logger.Errorf(stdFailureMsg, c.systemName)
	return err
}

func confirmDestruction(ctx *cmd.Context, systemName string) error {
	// Get confirmation from the user that they want to continue
	fmt.Fprintf(ctx.Stdout, destroySysMsg, systemName)

	scanner := bufio.NewScanner(ctx.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil && err != io.EOF {
		return errors.Annotate(err, "system destruction aborted")
	}
	answer := strings.ToLower(scanner.Text())
	if answer != "y" && answer != "yes" {
		return errors.New("system destruction aborted")
	}

	return nil
}

var stdFailureMsg = `failed to destroy system %q

If the system is unusable, then you may run

    juju system kill

to forcefully destroy the system. Upon doing so, review
your environment provider console for any resources that need
to be cleaned up.
`
