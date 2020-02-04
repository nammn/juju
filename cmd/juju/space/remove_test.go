// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package space_test

import (
	"github.com/juju/cmd"
	"github.com/juju/cmd/cmdtesting"
	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/cmd/juju/space"
	"github.com/juju/juju/feature"
)

type RemoveSuite struct {
	BaseSpaceSuite
}

var _ = gc.Suite(&RemoveSuite{})

func (s *RemoveSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetFeatureFlags(feature.PostNetCLIMVP)
	s.BaseSpaceSuite.SetUpTest(c)
}

func (s *RemoveSuite) TestInit(c *gc.C) {
	ctrl, api := setUpMocks(c)
	defer ctrl.Finish()
	validSpaceName := "myspace"
	api.EXPECT().RemoveSpace(validSpaceName).Return(nil)

	for i, test := range []struct {
		about     string
		args      []string
		expectErr string
	}{{
		about:     "no arguments",
		expectErr: "space name is required",
	}, {
		about:     "invalid space name",
		args:      s.Strings("%inv$alid", "new-name"),
		expectErr: `"%inv\$alid" is not a valid space name`,
	}, {
		about:     "multiple space names aren't allowed",
		args:      s.Strings("a-space", "another-space"),
		expectErr: `unrecognized args: \["another-space"\]`,
	}, {
		about: "delete a valid space name",
		args:  s.Strings(validSpaceName),
	}} {
		c.Logf("test #%d: %s", i, test.about)
		_, err := s.runCommand(c, api, test.args...)
		if test.expectErr != "" {
			prefixedErr := "invalid arguments specified: " + test.expectErr
			c.Check(err, gc.ErrorMatches, prefixedErr)
		}
	}
}

func (s *RemoveSuite) TestRunWithValidSpaceSucceeds(c *gc.C) {
	ctrl, api := setUpMocks(c)
	defer ctrl.Finish()
	spaceName := "myspace"
	expectedStdout := `removed space "myspace"`
	api.EXPECT().RemoveSpace(spaceName).Return(nil)

	ctx, err := s.runCommand(c, api, spaceName)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(cmdtesting.Stdout(ctx), gc.Equals, expectedStdout)
}

func (s *RemoveSuite) TestRunWhenSpacesAPIFails(c *gc.C) {
	ctrl, api := setUpMocks(c)
	defer ctrl.Finish()
	spaceName := "myspace"
	bam := errors.New("bam")
	api.EXPECT().RemoveSpace(spaceName).Return(bam)

	ctx, err := s.runCommand(c, api, spaceName)
	c.Assert(err, gc.ErrorMatches, bam.Error())
	c.Assert(cmdtesting.Stdout(ctx), gc.Equals, "")
	c.Assert(cmdtesting.Stderr(ctx), gc.Equals, "")
}

func (s *RemoveSuite) runCommand(c *gc.C, api space.SpaceAPI, name ...string) (*cmd.Context, error) {
	base := space.NewSpaceCommandBase(api)
	command := space.RemoveCommand{
		SpaceCommandBase: base,
	}
	return cmdtesting.RunCommand(c, &command, name...)
}
