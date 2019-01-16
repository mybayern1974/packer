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
}
