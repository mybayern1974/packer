package vagrant

import (
	"github.com/hashicorp/packer/helper/multistep"
)

func CommHost(config *SSHConfig) func(multistep.StateBag) (string, error) {
	return func(state multistep.StateBag) (string, error) {
		return config.Comm.SSHHost, nil
	}
}

func SSHPort(config *SSHConfig) func(multistep.StateBag) (int, error) {
	return func(state multistep.StateBag) (int, error) {
		return config.Comm.SSHPort, nil
	}
}
