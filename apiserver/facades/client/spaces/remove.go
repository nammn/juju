/*
 * Copyright 2020 Canonical Ltd.
 * Licensed under the AGPLv3, see LICENCE file for details.
 */

package spaces

import (
	"github.com/juju/errors"
	"gopkg.in/juju/names.v3"
	"gopkg.in/mgo.v2/txn"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/controller"
	"github.com/juju/juju/core/permission"
	"github.com/juju/juju/state"
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
}

type removeSpaceStateShim struct {
	*state.State
}

type spaceRemoveModelOp struct {
	st    RemoveSpaceState
	space RemoveSpace
}

func (o *spaceRemoveModelOp) Done(err error) error {
	return err
}

func NewRemoveSpaceModelOp(st RemoveSpaceState, space RemoveSpace) *spaceRemoveModelOp {
	return &spaceRemoveModelOp{
		st:    st,
		space: space,
	}
}

func (sp *spaceRemoveModelOp) Build(attempt int) ([]txn.Op, error) {
	// get subnets
	// get move subnets ops
	// get remove space ops
	totalops := sp.space.RemoveSpaceOps()
	return totalops, nil
}

// RefreshSpaces refreshes spaces from substrate
func (api *API) RemoveSpace(entities params.Entities) (params.RemoveSpaceResults, error) {
	isAdmin, err := api.auth.HasPermission(permission.AdminAccess, api.backing.ModelTag())
	if err != nil && !errors.IsNotFound(err) {
		return params.RemoveSpaceResults{}, errors.Trace(err)
	}
	if !isAdmin {
		return params.RemoveSpaceResults{}, common.ServerError(common.ErrPerm)
	}
	if err := api.check.ChangeAllowed(); err != nil {
		return params.RemoveSpaceResults{}, errors.Trace(err)
	}
	if err = api.checkSupportsProviderSpaces(); err != nil {
		return params.RemoveSpaceResults{}, common.ServerError(errors.Trace(err))
	}

	results := params.RemoveSpaceResults{
		Results: make([]params.RemoveSpaceResult, len(entities.Entities)),
	}
	for i, entity := range entities.Entities {
		skip := false
		spacesTag, err := names.ParseSpaceTag(entity.Tag)
		if err != nil {
			results.Results[i].Error = common.ServerError(errors.Trace(err))
			continue
		}
		space, err := api.backing.SpaceByName(spacesTag.Id())
		if err != nil {
			results.Results[i].Error = common.ServerError(errors.Trace(err))
			continue
		}
		bindingTags, err := api.getApplicationTagsPerSpace(space.Id())
		// bindings check
		if err != nil {
			results.Results[i].Error = common.ServerError(errors.Trace(err))
			continue
		}
		if len(bindingTags) != 0 {
			results.Results[i].Bindings = convertTagsToEntities(bindingTags)
			skip = true
		}
		// constraints check
		constraintTags, err := api.filterConstraints(space.Name())
		if err != nil {
			results.Results[i].Error = common.ServerError(errors.Trace(err))
			continue
		}
		if len(constraintTags) != 0 {
			results.Results[i].Constraints = convertTagsToEntities(constraintTags)
			skip = true
		}

		// controller check
		matches, err := api.checkControllerForSpace(space.Name())
		if err != nil {
			results.Results[i].Error = common.ServerError(errors.Trace(err))
			continue
		}
		if len(matches) != 0 {
			results.Results[i].ControllerSettings = matches
			skip = true
		}

		if skip {
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

func (api *API) getApplicationTagsPerSpace(spaceID string) ([]names.Tag, error) {
	applications, err := api.getApplicationsBindSpace(spaceID)
	if err != nil {
		return nil, errors.Trace(nil)
	}
	tags := make([]names.Tag, len(applications))
	for i, app := range applications {
		tags[i] = names.NewApplicationTag(app)
	}
	return tags, nil
}

func convertTagsToEntities(tags []names.Tag) []params.Entity {
	entities := make([]params.Entity, len(tags))
	for i, tag := range tags {
		entities[i].Tag = tag.String()
	}

	return entities
}

func (api *API) filterConstraints(spaceName string) ([]names.Tag, error) {
	tags, err := api.backing.ConstraintsTagForSpaceName(spaceName)
	var notSkipping []names.Tag
	if err != nil {
		return nil, errors.Trace(err)
	}
	for _, tag := range tags {
		if tag.Kind() == names.UnitTagKind {
			continue
		}
		if tag.Kind() == names.MachineTagKind {
			continue
		}
		notSkipping = append(notSkipping, tag)
	}
	return notSkipping, nil
}

func (api *API) checkControllerForSpace(spaceName string) ([]string, error) {
	var matches []string
	is, err := api.backing.IsControllerModel()
	if err != nil {
		return matches, errors.Trace(err)
	}
	if !is {
		return matches, nil
	}

	currentControllerConfig, err := api.backing.ControllerConfig()
	if err != nil {
		return matches, errors.Trace(err)
	}

	if mgmtSpace := currentControllerConfig.JujuManagementSpace(); mgmtSpace == spaceName {
		matches = append(matches, controller.JujuManagementSpace)
	}
	if haSpace := currentControllerConfig.JujuHASpace(); haSpace == spaceName {
		matches = append(matches, controller.JujuHASpace)
	}
	return matches, nil
}
