// Copyright 2015 Cloudbase Solutions
// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// +build !linux windows

package windows_test

import (
	"code.google.com/p/winsvc/svc"
	"github.com/juju/errors"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/service/common"
	"github.com/juju/juju/service/windows"
	coretesting "github.com/juju/juju/testing"
)

type serviceManagerSuite struct {
	coretesting.BaseSuite

	stub       *testing.Stub
	passwdStub *testing.Stub
	conn       *windows.StubMgr
	getPasswd  *windows.StubGetPassword

	name string
	conf common.Conf

	mgr windows.ServiceManagerInterface

	execPath string
}

var _ = gc.Suite(&serviceManagerSuite{})

func (s *serviceManagerSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)
	var err error
	s.execPath = `C:\juju\bin\jujud.exe`
	s.stub = &testing.Stub{}
	s.passwdStub = &testing.Stub{}
	s.conn = windows.PatchMgrConnect(s, s.stub)
	s.getPasswd = windows.PatchGetPassword(s, s.passwdStub)

	// Set up the service.
	s.name = "machine-1"
	s.conf = common.Conf{
		Desc:      "service for " + s.name,
		ExecStart: s.execPath + " " + s.name,
	}

	s.mgr, err = windows.NewServiceManager()
	c.Assert(err, gc.IsNil)

	// Clear services
	s.conn.Clear()
}

func (s *serviceManagerSuite) TestCreate(c *gc.C) {
	s.getPasswd.SetPasswd("fake")
	err := s.mgr.Create(s.name, s.conf)
	c.Assert(err, gc.IsNil)

	c.Assert(s.getPasswd.Calls, gc.HasLen, 1)

	exists := s.conn.Exists(s.name)
	c.Assert(exists, jc.IsTrue)

	svcs := s.conn.List()
	c.Assert(svcs, gc.HasLen, 1)

	m, ok := s.mgr.(*windows.SvcManager)
	c.Assert(ok, jc.IsTrue)

	cfg, err := m.Config(s.name)
	c.Assert(err, gc.IsNil)
	c.Assert(cfg.ServiceStartName, gc.Equals, windows.JujudUser)

	running, err := s.mgr.Running(s.name)
	c.Assert(err, gc.IsNil)
	c.Assert(running, jc.IsFalse)
}

func (s *serviceManagerSuite) TestCreateInvalidPassword(c *gc.C) {
	passwdError := errors.New("Failed to get password")
	s.passwdStub.SetErrors(passwdError)

	err := s.mgr.Create(s.name, s.conf)
	c.Assert(err.Error(), gc.Equals, passwdError.Error())

	c.Assert(s.getPasswd.Calls, gc.HasLen, 1)

	exists := s.conn.Exists(s.name)
	c.Assert(exists, jc.IsFalse)
}

func (s *serviceManagerSuite) TestCreateInvalidUser(c *gc.C) {
	s.getPasswd.SetPasswd("fake")

	err := s.mgr.Create(s.name, s.conf)
	c.Assert(err, gc.IsNil)

	c.Assert(s.getPasswd.Calls, gc.HasLen, 1)

	m, ok := s.mgr.(*windows.SvcManager)
	c.Assert(ok, jc.IsTrue)

	cfg, err := m.Config(s.name)

	c.Assert(err, gc.IsNil)
	c.Assert(cfg.ServiceStartName, gc.Equals, windows.JujudUser)
}

func (s *serviceManagerSuite) TestCreateAlreadyExists(c *gc.C) {
	err := s.mgr.Create(s.name, s.conf)
	c.Assert(err, gc.IsNil)
	exists := s.conn.Exists(s.name)
	c.Assert(exists, jc.IsTrue)
	err = s.mgr.Create(s.name, s.conf)
	c.Assert(err.Error(), gc.Equals, windows.ERROR_SERVICE_EXISTS.Error())
}

func (s *serviceManagerSuite) TestCreateMultipleServices(c *gc.C) {
	err := s.mgr.Create("test-service", common.Conf{})
	c.Assert(err, gc.IsNil)
	exists := s.conn.Exists("test-service")
	c.Assert(exists, jc.IsTrue)

	err = s.mgr.Create("another-test-service", common.Conf{})
	c.Assert(err, gc.IsNil)
	exists = s.conn.Exists("another-test-service")
	c.Assert(exists, jc.IsTrue)

	svcs := s.conn.List()
	c.Assert(svcs, gc.HasLen, 2)
}

func (s *serviceManagerSuite) TestStart(c *gc.C) {
	windows.AddService(s.name, s.execPath, s.stub, svc.Status{State: svc.Stopped})

	err := s.mgr.Start(s.name)
	c.Assert(err, gc.IsNil)

	running, err := s.mgr.Running(s.name)
	c.Assert(err, gc.IsNil)
	c.Assert(running, jc.IsTrue)
}

func (s *serviceManagerSuite) TestStartTwice(c *gc.C) {
	windows.AddService(s.name, s.execPath, s.stub, svc.Status{State: svc.Stopped})

	err := s.mgr.Start(s.name)
	c.Assert(err, gc.IsNil)

	err = s.mgr.Start(s.name)
	c.Assert(err, gc.IsNil)

	running, err := s.mgr.Running(s.name)
	c.Assert(err, gc.IsNil)
	c.Assert(running, jc.IsTrue)
}

func (s *serviceManagerSuite) TestStartStop(c *gc.C) {
	windows.AddService(s.name, s.execPath, s.stub, svc.Status{State: svc.Stopped})

	err := s.mgr.Start(s.name)
	c.Assert(err, gc.IsNil)

	running, err := s.mgr.Running(s.name)
	c.Assert(err, gc.IsNil)
	c.Assert(running, jc.IsTrue)

	err = s.mgr.Stop(s.name)
	c.Assert(err, gc.IsNil)

	running, err = s.mgr.Running(s.name)
	c.Assert(err, gc.IsNil)
	c.Assert(running, jc.IsFalse)
}

func (s *serviceManagerSuite) TestStop(c *gc.C) {
	windows.AddService(s.name, s.execPath, s.stub, svc.Status{State: svc.Running})

	running, err := s.mgr.Running(s.name)
	c.Assert(err, gc.IsNil)
	c.Assert(running, jc.IsTrue)

	err = s.mgr.Stop(s.name)
	c.Assert(err, gc.IsNil)

	running, err = s.mgr.Running(s.name)
	c.Assert(err, gc.IsNil)
	c.Assert(running, jc.IsFalse)
}

func (s *serviceManagerSuite) TestDelete(c *gc.C) {
	windows.AddService(s.name, s.execPath, s.stub, svc.Status{State: svc.Running})

	err := s.mgr.Delete(s.name)
	c.Assert(err, gc.IsNil)
	exists := s.conn.Exists(s.name)
	c.Assert(exists, jc.IsFalse)
}

func (s *serviceManagerSuite) TestDeleteInexistent(c *gc.C) {
	exists := s.conn.Exists(s.name)
	c.Assert(exists, jc.IsFalse)

	err := s.mgr.Delete(s.name)
	c.Assert(err.Error(), gc.Equals, windows.ERROR_SERVICE_DOES_NOT_EXIST.Error())
}
