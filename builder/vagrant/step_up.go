package vagrant

import (
	"context"

	"github.com/hashicorp/packer/helper/multistep"
)

type StepUp struct{}

func (s *StepUp) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(VagrantDriver)

	_, _, err := driver.Up()

	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepUp) Cleanup(state multistep.StateBag) {
	driver := state.Get("driver").(VagrantDriver)

	var err error
	if s.TeardownMethod == "halt" {
		err = driver.Halt()
	} else if s.TeardownMethod == "suspend" {
		err = driver.Suspend()
	} else if s.TeardownMethod == "destroy" {
		err = driver.Destroy()
	} else {
		state.Put("error", fmt.Errorf("Invalid teardown method selected; must be either halt, suspend, or destory."))
		return multistep.ActionHalt
	}
	if err != nil {
		state.Put("error", fmt.Errorf("Error halting Vagrant machine; please try to do this manually"))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}
}
