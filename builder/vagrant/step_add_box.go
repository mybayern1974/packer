package vagrant

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
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)

	// Prepare arguments
	addArgs := []string{}

	if s.BoxVersion != "" {
		addArgs = append(initArgs, "--box-version", s.BoxVersion)
	}

	if s.CACert != "" {
		addArgs = append(initArgs, "--cacert", s.CACert)
	}

	if s.CAPath != "" {
		addArgs = append(initArgs, "--capath", s.CAPath)
	}

	if s.DownloadCert != "" {
		addArgs = append(initArgs, "--cert", s.DownloadCert)
	}

	if s.Clean {
		addArgs = append(initArgs, "--clean")
	}

	if s.Force {
		addArgs = append(initArgs, "--force")
	}

	if s.Insecure {
		addArgs = append(initArgs, "--insecure")
	}

	if s.Provider != "" {
		addArgs = append(initArgs, "--provider", s.Provider)
	}

	addArgs = append(initArgs, Address)

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
