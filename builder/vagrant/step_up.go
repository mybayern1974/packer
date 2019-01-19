package vagrant

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
)

type StepUp struct {
	TeardownMethod string
}

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
		// Should never get here because of template validation
		state.Put("error", fmt.Errorf("Invalid teardown method selected; must be either halt, suspend, or destory."))
	}
	if err != nil {
		state.Put("error", fmt.Errorf("Error halting Vagrant machine; please try to do this manually"))
	}
}
