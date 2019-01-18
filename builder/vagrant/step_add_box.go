package vagrant

import (
	"context"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type StepAddBox struct {
	BoxVersion   string
	CACert       string
	CAPath       string
	DownloadCert string
	Clean        bool
	Force        bool
	Insecure     bool
	Provider     string
	Address      string
}

func (s *StepAddBox) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	// driver := state.Get("driver").(VagrantDriver)
	ui := state.Get("ui").(packer.Ui)

	// Prepare arguments
	addArgs := []string{}

	if s.BoxVersion != "" {
		addArgs = append(addArgs, "--box-version", s.BoxVersion)
	}

	if s.CACert != "" {
		addArgs = append(addArgs, "--cacert", s.CACert)
	}

	if s.CAPath != "" {
		addArgs = append(addArgs, "--capath", s.CAPath)
	}

	if s.DownloadCert != "" {
		addArgs = append(addArgs, "--cert", s.DownloadCert)
	}

	if s.Clean {
		addArgs = append(addArgs, "--clean")
	}

	if s.Force {
		addArgs = append(addArgs, "--force")
	}

	if s.Insecure {
		addArgs = append(addArgs, "--insecure")
	}

	if s.Provider != "" {
		addArgs = append(addArgs, "--provider", s.Provider)
	}

	addArgs = append(addArgs, s.Address)

	log.Printf("[vagrant] Calling box add with following args %s", strings.Join(addArgs, " "))
	// Call vagrant using prepared arguments
	err := driver.Add(addArgs)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepAddBox) Cleanup(state multistep.StateBag) {
}
