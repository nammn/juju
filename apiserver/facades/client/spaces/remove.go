/*
 * Copyright 2020 Canonical Ltd.
 * Licensed under the AGPLv3, see LICENCE file for details.
 */

package spaces

import (
	"github.com/juju/errors"
	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/params"
	jujucontroller "github.com/juju/juju/controller"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/permission"
	"github.com/juju/juju/state"
	"gopkg.in/juju/names.v3"
	"gopkg.in/mgo.v2/txn"
)

// RemoveSpace describes a space that can be removed.
type RemoveSpace interface {
	Refresh() error
	Id() string
	Name() string
	RemoveSpaceOps() []txn.Op
}

// RemoveSpaceState describes state operations required
// to execute the renameSpace operation.
// * This allows us to indirect state at the operation level instead of the
// * whole API level as currently done in interface.go
type RemoveSpaceState interface {
	// ControllerConfig returns current ControllerConfig.
	ControllerConfig() (jujucontroller.Config, error)

	// ConstraintsOpsForSpaceNameChange returns all the database transaction operation required
	// to transform a constraints spaces from `a` to `b`
	ConstraintsForSpaceName(name string) ([]constraints.Value, error)
}

type removeSpaceStateShim struct {
	*state.State
}

type spaceRemoveModelOp struct {
	st           RemoveSpaceState
	isController bool
	space        RemoveSpace
}

func (o *spaceRemoveModelOp) Done(err error) error {
	return err
}

func NewRemoveSpaceModelOp(isController bool, st RemoveSpaceState, space RemoveSpace) *spaceRemoveModelOp {
	return &spaceRemoveModelOp{
		st:           st,
		space:        space,
		isController: isController,
	}
}

func (sp *spaceRemoveModelOp) Build(attempt int) ([]txn.Op, error) {
	if sp.isController {
		//	check if space is used
		// controller setting if model == controller model
	}

	// check constraint -> ConstraintsOpsForSpaceNameChange
	// TODO: problem, we need id/name to show what constraints are missing
	constraints, err := sp.st.ConstraintsForSpaceName(sp.space.Name())
	if err != nil {
		return nil, errors.Trace(err)
	}
	if len(constraints) != 0 {
		return nil, errors.New("removing space is not possible because of existing constraints: %q")
	}
	totalops := sp.space.RemoveSpaceOps()
	return totalops, nil
}

// RefreshSpaces refreshes spaces from substrate
func (api *API) RemoveSpace(entities params.Entities) (params.ErrorResults, error) {
	// TODO: don' allow provider provided spaces
	isAdmin, err := api.auth.HasPermission(permission.AdminAccess, api.backing.ModelTag())
	if err != nil && !errors.IsNotFound(err) {
		return params.ErrorResults{}, errors.Trace(err)
	}
	if !isAdmin {
		return params.ErrorResults{}, common.ServerError(common.ErrPerm)
	}
	if err := api.check.ChangeAllowed(); err != nil {
		return params.ErrorResults{}, errors.Trace(err)
	}
	if err := api.checkSupportsSpaces(); err != nil {
		return params.ErrorResults{}, common.ServerError(errors.Trace(err))
	}
	results := params.ErrorResults{
		Results: make([]params.ErrorResult, len(entities.Entities)),
	}
	for i, entity := range entities.Entities {
		spacesTag, err := names.ParseSpaceTag(entity.Tag)
		if err != nil {
			results.Results[i].Error = common.ServerError(errors.Trace(err))
		}
		space, err := api.backing.SpaceByName(spacesTag.Id())
		if err != nil {
			results.Results[i].Error = common.ServerError(errors.Trace(err))
			continue
		}
		applications, err := api.getApplicationsBindSpace(space.Id())
		if err != nil {
			results.Results[i].Error = common.ServerError(errors.Trace(err))
			continue
		}
		// endpoint bindings
		// TODO: maybe move?
		if len(applications) != 0 {
			newErr := errors.Errorf("removing space is not possible, applications %q are bind to it.", applications)
			results.Results[i].Error = common.ServerError(newErr)
			continue
		}

		// create operation factory
		operation, err := api.opFactory.NewRemoveSpaceModelOp(spacesTag.Id())
		if err != nil {
			results.Results[i].Error = common.ServerError(errors.Trace(err))
			continue
		}

		// apply operations
		if err = api.backing.ApplyOperation(operation); err != nil {
			results.Results[i].Error = common.ServerError(errors.Trace(err))
			continue
		}
	}
	return results, nil
}
