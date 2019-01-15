package vagrant

import (
	"log"
	"strings"

	"github.com/hashicorp/packer/helper/multistep"
)

type StepHalt struct {
	TeardownMethod string
}

func (s *StepHalt) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)

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

func (s *StepHalt) Cleanup(state multistep.StateBag) {
}
