package vagrant

type StepSyncedFolder struct {
	SyncedFolder string
}

func (s *StepSyncedFolder) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Currently you must set your Synced Folder via your vagrantfile.")
	return multistep.ActionContinue
	// vagrantfile := state.Get("vagrantfile").(string)
	// fi, err := os.Open(vagrantfile)
	// if err != nil {
	// 	err := fmt.Errorf("Error getting vagrantfile for new box")
	// 	state.Put("error", err)
	// 	ui.Error(err.Error())
	// 	return multistep.ActionHalt
	// }

	// TODO Modify vagrantfile to sync folder.
	if s.SyncedFolder == "" {
		return multistep.ActionContinue
	}

	return multistep.ActionContinue
}

func (s *StepSyncedFolder) Cleanup(state multistep.StateBag) {
}
