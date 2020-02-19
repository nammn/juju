// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package space

import (
	"fmt"
	"github.com/juju/gnuflag"
	"github.com/juju/juju/core/network"
	"strings"

	"github.com/juju/cmd"
	"github.com/juju/collections/set"
	"github.com/juju/errors"

	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/modelcmd"
)

// NewMoveCommand returns a command used to update subnets in a space.
func NewMoveCommand() modelcmd.ModelCommand {
	return modelcmd.Wrap(&MoveCommand{})
}

// MoveCommand calls the API to update an existing network space.
type MoveCommand struct {
	SpaceCommandBase
	Name  string
	CIDRs set.Strings

	Force bool
}

const updateCommandDoc = `
Replaces the list of associated subnets of the space. Since subnets
can only be part of a single space, all specified subnets (using their
CIDRs) "leave" their current space and "enter" the one we're updating.

Examples:

To move a list of CIDRs from their space to a new space:

  juju move-to-space db-space 172.31.1.0/28 172.31.16.0/20

See also:
	add-space
	list-spaces
	reload-spaces
	rename-space
	show-space
	remove-space
`

// TODO: update DOC.
// Info is defined on the cmd.Command interface.
func (c *MoveCommand) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name:    "move-to-space",
		Args:    "<name> <CIDR1> [ <CIDR2> ...]",
		Purpose: "Update a network space's CIDRs.",
		Doc:     strings.TrimSpace(updateCommandDoc),
	})
}

func (c *MoveCommand) SetFlags(f *gnuflag.FlagSet) {
	f.BoolVar(&c.Force, "force", false, "Allow to force a move of subnets to a space even if they are in use on another machine.")
}

// Init is defined on the cmd.Command interface. It checks the
// arguments for sanity and sets up the command to run.
func (c *MoveCommand) Init(args []string) error {
	var err error
	c.Name, c.CIDRs, err = ParseNameAndCIDRs(args, false)
	return errors.Trace(err)
}

// Run implements Command.Run.
func (c *MoveCommand) Run(ctx *cmd.Context) error {
	return c.RunWithAPI(ctx, func(api SpaceAPI, ctx *cmd.Context) error {
		moved, err := api.MoveToSpace(c.Name, c.CIDRs.SortedValues(), c.Force)
		if err != nil {
			return errors.Annotatef(err, "cannot update space %q", c.Name)
		}
		for _, change := range createMovementsChangelog(moved) {
			ctx.Infof(change)
		}
		return nil
	})
}

func createMovementsChangelog(moved []network.MovedSpace) []string {
	var changelog []string
	for _, movedSubnet := range moved {
		changelog = append(changelog, fmt.Sprintf("Subnet %q moved from %q to %q", movedSubnet.CIDR, movedSubnet.SpaceFrom, movedSubnet.SpaceTo))
	}
	return changelog

}
