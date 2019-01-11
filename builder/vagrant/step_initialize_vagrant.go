package vagrant

type StepInitializeVagrant struct {
	BoxVersion string
	Minimal    bool
	Template   string
	BoxName    string
}

func (s *StepInitializeVagrant) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)

	// Prepare arguments
	initArgs := []string{}

	if s.BoxVersion != "" {
		initArgs = append(initArgs, "--box-version", s.BoxVersion)
	}

	if s.Minimal {
		initArgs = append(initArgs, "-m")
	}

	if s.Template {
		initArgs = append(initArgs, "--template", s.Template)
	}

	initArgs = append(initArgs, s.BoxName)

	// Call vagrant using prepared arguments
	err := driver.Init(initArgs)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepInitializeVagrant) Cleanup(state multistep.StateBag) {
}
