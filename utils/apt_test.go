// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package utils_test

import (
	"fmt"

	gc "launchpad.net/gocheck"

	"launchpad.net/juju-core/juju/osenv"
	jc "launchpad.net/juju-core/testing/checkers"
	"launchpad.net/juju-core/testing/testbase"
	"launchpad.net/juju-core/utils"
)

type AptSuite struct {
	testbase.LoggingSuite
}

var _ = gc.Suite(&AptSuite{})

func (s *AptSuite) TestOnePackage(c *gc.C) {
	cmdChan := s.HookCommandOutput(&utils.AptCommandOutput, []byte{}, nil)
	err := utils.AptGetInstall("test-package")
	c.Assert(err, gc.IsNil)
	cmd := <-cmdChan
	c.Assert(cmd.Args, gc.DeepEquals, []string{
		"apt-get", "--option=Dpkg::Options::=--force-confold",
		"--option=Dpkg::options::=--force-unsafe-io", "--assume-yes", "--quiet",
		"install", "test-package",
	})
	c.Assert(cmd.Env[len(cmd.Env)-1], gc.Equals, "DEBIAN_FRONTEND=noninteractive")
}

func (s *AptSuite) TestAptGetError(c *gc.C) {
	const expected = `E: frobnicator failure detected`
	cmdError := fmt.Errorf("error")
	cmdExpectedError := fmt.Errorf("apt-get failed: error")
	cmdChan := s.HookCommandOutput(&utils.AptCommandOutput, []byte(expected), cmdError)
	err := utils.AptGetInstall("foo")
	c.Assert(err, gc.DeepEquals, cmdExpectedError)
	cmd := <-cmdChan
	c.Assert(cmd.Args, gc.DeepEquals, []string{
		"apt-get", "--option=Dpkg::Options::=--force-confold",
		"--option=Dpkg::options::=--force-unsafe-io", "--assume-yes", "--quiet",
		"install", "foo",
	})
}

func (s *AptSuite) TestConfigProxyEmpty(c *gc.C) {
	cmdChan := s.HookCommandOutput(&utils.AptCommandOutput, []byte{}, nil)
	out, err := utils.AptConfigProxy()
	c.Assert(err, gc.IsNil)
	cmd := <-cmdChan
	c.Assert(cmd.Args, gc.DeepEquals, []string{
		"apt-config", "dump", "Acquire::http::Proxy",
		"Acquire::https::Proxy", "Acquire::ftp::Proxy",
	})
	c.Assert(out, gc.Equals, "")
}

func (s *AptSuite) TestConfigProxyConfigured(c *gc.C) {
	const expected = `Acquire::http::Proxy "10.0.3.1:3142";
Acquire::https::Proxy "false";`
	cmdChan := s.HookCommandOutput(&utils.AptCommandOutput, []byte(expected), nil)
	out, err := utils.AptConfigProxy()
	c.Assert(err, gc.IsNil)
	cmd := <-cmdChan
	c.Assert(cmd.Args, gc.DeepEquals, []string{
		"apt-config", "dump", "Acquire::http::Proxy",
		"Acquire::https::Proxy", "Acquire::ftp::Proxy",
	})
	c.Assert(out, gc.Equals, expected)
}

func (s *AptSuite) TestDetectAptProxy(c *gc.C) {
	const output = `CommandLine::AsString "apt-config dump";
Acquire::http::Proxy  "10.0.3.1:3142";
Acquire::https::Proxy "false";
Acquire::ftp::Proxy "none";
Acquire::magic::Proxy "none";
`
	_ = s.HookCommandOutput(&utils.AptCommandOutput, []byte(output), nil)

	proxy, err := utils.DetectAptProxies()
	c.Assert(err, gc.IsNil)
	c.Assert(proxy, gc.DeepEquals, osenv.ProxySettings{
		Http:  "10.0.3.1:3142",
		Https: "false",
		Ftp:   "none",
	})
}

func (s *AptSuite) TestDetectAptProxyNone(c *gc.C) {
	_ = s.HookCommandOutput(&utils.AptCommandOutput, []byte{}, nil)
	proxy, err := utils.DetectAptProxies()
	c.Assert(err, gc.IsNil)
	c.Assert(proxy, gc.DeepEquals, osenv.ProxySettings{})
}

func (s *AptSuite) TestConfigProxyConfiguredFilterOutput(c *gc.C) {
	const (
		output = `CommandLine::AsString "apt-config dump";
Acquire::http::Proxy  "10.0.3.1:3142";
Acquire::https::Proxy "false";`
		expected = `Acquire::http::Proxy  "10.0.3.1:3142";
Acquire::https::Proxy "false";`
	)
	cmdChan := s.HookCommandOutput(&utils.AptCommandOutput, []byte(output), nil)
	out, err := utils.AptConfigProxy()
	c.Assert(err, gc.IsNil)
	cmd := <-cmdChan
	c.Assert(cmd.Args, gc.DeepEquals, []string{
		"apt-config", "dump", "Acquire::http::Proxy",
		"Acquire::https::Proxy", "Acquire::ftp::Proxy",
	})
	c.Assert(out, gc.Equals, expected)
}

func (s *AptSuite) TestConfigProxyError(c *gc.C) {
	const expected = `E: frobnicator failure detected`
	cmdError := fmt.Errorf("error")
	cmdExpectedError := fmt.Errorf("apt-config failed: error")
	cmdChan := s.HookCommandOutput(&utils.AptCommandOutput, []byte(expected), cmdError)
	out, err := utils.AptConfigProxy()
	c.Assert(err, gc.DeepEquals, cmdExpectedError)
	cmd := <-cmdChan
	c.Assert(cmd.Args, gc.DeepEquals, []string{
		"apt-config", "dump", "Acquire::http::Proxy",
		"Acquire::https::Proxy", "Acquire::ftp::Proxy",
	})
	c.Assert(out, gc.Equals, "")
}

func (s *AptSuite) patchLsbRelease(c *gc.C, name string) {
	content := fmt.Sprintf("#!/bin/bash --norc\necho %s", name)
	patchExecutable(s, c.MkDir(), "lsb_release", content)
}

func (s *AptSuite) TestIsUbuntu(c *gc.C) {
	s.patchLsbRelease(c, "Ubuntu")
	c.Assert(utils.IsUbuntu(), jc.IsTrue)
}

func (s *AptSuite) TestIsNotUbuntu(c *gc.C) {
	s.patchLsbRelease(c, "Windows NT")
	c.Assert(utils.IsUbuntu(), jc.IsFalse)
}

func (s *AptSuite) patchDpkgQuery(c *gc.C, installed bool) {
	rc := 0
	if !installed {
		rc = 1
	}
	content := fmt.Sprintf("#!/bin/bash --norc\nexit %v", rc)
	patchExecutable(s, c.MkDir(), "dpkg-query", content)
}

func (s *AptSuite) TestIsPackageInstalled(c *gc.C) {
	s.patchDpkgQuery(c, true)
	c.Assert(utils.IsPackageInstalled("foo-bar"), jc.IsTrue)
}

func (s *AptSuite) TestIsPackageNotInstalled(c *gc.C) {
	s.patchDpkgQuery(c, false)
	c.Assert(utils.IsPackageInstalled("foo-bar"), jc.IsFalse)
}
