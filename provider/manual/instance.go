// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package manual

import (
	"github.com/juju/juju/environs/manual"
	"github.com/juju/juju/instance"
)

type manualBootstrapInstance struct {
	host string
}

func (manualBootstrapInstance) Id() instance.Id {
	// The only way to bootrap is via manual bootstrap.
	return manual.BootstrapInstanceId
}

func (manualBootstrapInstance) Status() string {
	return ""
}

func (manualBootstrapInstance) Refresh() error {
	return nil
}

func (inst manualBootstrapInstance) Addresses() (addresses []instance.Address, err error) {
	addr, err := manual.HostAddress(inst.host)
	if err != nil {
		return nil, err
	}
	return []instance.Address{addr}, nil
}

func (manualBootstrapInstance) OpenPorts(machineId string, ports []instance.Port) error {
	return nil
}

func (manualBootstrapInstance) ClosePorts(machineId string, ports []instance.Port) error {
	return nil
}

func (manualBootstrapInstance) Ports(machineId string) ([]instance.Port, error) {
	return []instance.Port{}, nil
}
