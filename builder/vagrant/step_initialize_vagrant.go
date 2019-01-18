package vagrant

import (
	"context"
	"interpolate"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type StepInitializeVagrant struct {
	BoxName    string
	BoxVersion string
	Minimal    bool
	Template   string
	SourceBox  string
	OutputDir  string
}

var DEFAULT_TEMPLATE = `Vagrant.configure("2") do |config|
  config.vm.box = "<%= @box_name %>"
  {{ if ne .SyncedFolder "" -}}
  		config.vm.synced_folder "{{.SyncedFolder}}", "/vagrant"
  {{ - else - }}
  		config.vm.synced_folder ".", "/vagrant", disabled: true
  {{ - end}}
end`

type VagrantfileOptions struct {
	SyncedFolder string
}

func (s *StepInitializeVagrant) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(VagrantDriver)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Initializing Vagrant in build directory...")

	// Prepare arguments
	initArgs := []string{}

	if s.BoxName {
		initArgs = append(initArgs, s.BoxName)
	}

	initArgs = append(initArgs, s.SourceBox)

	if s.BoxVersion != "" {
		initArgs = append(initArgs, "--box-version", s.BoxVersion)
	}

	if s.Minimal {
		initArgs = append(initArgs, "-m")
	}

	if s.Template == "" {
		// Generate vagrantfile template based on our default
		tpl := template.Must(template.New("VagrantTpl").Parse(DEFAULT_TEMPLATE))
		tplPath := filepath.Join(s.OutputDir, "packer-vagrantfile-template.erb")
	} else {
		// Read in the template from provided file.
		tpl, err := template.ParseFiles(s.Template)
		if err != nil {
			log.Print(err)
			return
		}
	}

	templateFile, err := os.Create(tplPath)
	if err != nil {
		log.Println("Error creating vagrantfile ", err)
		return
	}

	opts := &VagrantfileOptions{
		SyncedFolder: s.SyncedFolder,
	}

	err = tpl.Execute(templateFile, opts)
	if err != nil {
		return "", err
	}
	initArgs = append(initArgs, "--template", s.Template)

	os.Chdir(s.OutputDir)
	// Call vagrant using prepared arguments
	err = driver.Init(initArgs)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepInitializeVagrant) Cleanup(state multistep.StateBag) {
}
