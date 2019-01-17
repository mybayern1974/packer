package vagrant

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type StepInitializeVagrant struct {
	BoxVersion string
	Minimal    bool
	Template   string
	BoxName    string
}

func (s *StepInitializeVagrant) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(VagrantDriver)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Initializing Vagrant in build directory...")

	// Prepare arguments
	initArgs := []string{}

	if s.BoxVersion != "" {
		initArgs = append(initArgs, "--box-version", s.BoxVersion)
	}

	if s.Minimal {
		initArgs = append(initArgs, "-m")
	}

	if s.Template != "" {
		initArgs = append(initArgs, "--template", s.Template)
	}

	initArgs = append(initArgs, s.BoxName)

	// Move Packer execution into the output directory.
	if !SkipPackage {
		os.Chdir(config.OutputDir)
	}
	// Call vagrant using prepared arguments
	err := driver.Init(initArgs)
	if err != nil {
		if strings.Contains(err.Error(), "already exists in this directory") {
			log.Println("Vagrantfile already exists; using present Vagrantfile rather than initializing.")
			return multistep.ActionContinue
		}
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepInitializeVagrant) Cleanup(state multistep.StateBag) {
}
