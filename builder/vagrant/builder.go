package vagrant

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/common/bootcommand"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

// Builder implements packer.Builder and builds the actual VirtualBox
// images.
type Builder struct {
	config *Config
	runner multistep.Runner
}

type SSHConfig struct {
	Comm communicator.Config `mapstructure:",squash"`
}

type Config struct {
	common.PackerConfig    `mapstructure:",squash"`
	common.HTTPConfig      `mapstructure:",squash"`
	common.ISOConfig       `mapstructure:",squash"`
	common.FloppyConfig    `mapstructure:",squash"`
	bootcommand.BootConfig `mapstructure:",squash"`
	SSHConfig              `mapstructure:",squash"`

	// This is the name of the new virtual machine.
	// By default this is "packer-BUILDNAME", where "BUILDNAME" is the name of the build.
	OutputDir string `mapstructure:"output_dir"`
	VMName    string `mapstructure:"vm_name"`

	Communicator string `mapstructure:"communicator"`

	// What vagrantfile to use
	Vagrantfile string `mapstructure:"vagrantfile"`

	// Whether to Halt, Suspend, or Destroy the box
	TeardownMethod string `mapstructure:"teardown_method"`

	// Options for the "vagrant init" command
	BoxVersion        string `mapstructure:"box_version"`
	Minimal           bool   `mapstructure:"init_minimal"`
	OutputVagrantfile string `mapstructure:"output_vagrantfile"`
	Template          string `mapstructure:"template"`

	// Options for the "vagrant box add" command
	AddCACert       string `mapstructure:"add_cacert"`
	AddCAPath       string `mapstructure:"add_capath"`
	AddDownloadCert string `mapstructure:"add_cert"`
	AddClean        bool   `mapstructure:"add_clean"`
	AddForce        bool   `mapstructure:"add_force"`
	AddInsecure     bool   `mapstructure:"add_insecure"`

	// what folder to sync. Defaults to current build dir.
	SyncedFolder string `mapstructure:"synced_folder"`

	// Don't package the Vagrant box after build.
	SkipPackage bool `mapstructure:"skip_package"`

	ctx interpolate.Context
}

// Prepare processes the build configuration parameters.
func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	b.config = new(Config)
	err := config.Decode(&b.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &b.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"boot_command",
			},
		},
	}, raws...)
	if err != nil {
		return nil, err
	}

	// Accumulate any errors and warnings
	var errs *packer.MultiError
	warnings := make([]string, 0)

	if b.config.OutputDir == "" {
		b.config.OutputDir = fmt.Sprintf("output-%s", b.config.PackerBuildName)
	}

	if b.config.Comm.SSHTimeout == 0 {
		b.config.Comm.SSHTimeout = 10 * time.Minute
	}

	if b.config.Comm.Type != "ssh" {
		errs = packer.MultiErrorAppend(errs,
			fmt.Errorf(`The Vagrant builder currently only supports the ssh communicator"`))
	}

	if b.config.TeardownMethod == "" {
		b.config.TeardownMethod = "halt"
	} else {
		matches := false
		for _, name := range []string{"halt", "suspend", "destroy"} {
			if strings.ToLower(b.config.TeardownMethod) == name {
				matches = true
			}
		}
		if !matches {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf(`TeardownMethod must be "halt", "suspend", or "destroy"`))
		}
	}

	if errs != nil && len(errs.Errors) > 0 {
		return warnings, errs
	}

	return warnings, nil
}

// Run executes a Packer build and returns a packer.Artifact representing
// a VirtualBox appliance.
func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	// Create the driver that we'll use to communicate with VirtualBox
	driver, err := NewDriver()
	if err != nil {
		return nil, fmt.Errorf("Failed creating VirtualBox driver: %s", err)
	}

	// Set up the state.
	state := new(multistep.BasicStateBag)
	state.Put("config", b.config)
	state.Put("debug", b.config.PackerDebug)
	state.Put("driver", driver)
	state.Put("cache", cache)
	state.Put("hook", hook)
	state.Put("ui", ui)

	// Build the steps.
	steps := []multistep.Step{
		&common.StepOutputDir{
			Force: b.config.PackerForce,
			Path:  b.config.OutputDir,
		},
		&StepInitializeVagrant{
			BoxVersion: b.config.BoxVersion,
			Minimal:    b.config.Minimal,
			Template:   b.config.Template,
			BoxName:    b.config.VMName,
		},
		// &StepAddBox{
		// 	BoxVersion:   b.config.BoxVersion,
		// 	CACert:       b.config.AddCACert,
		// 	CAPath:       b.config.AddCAPath,
		// 	DownloadCert: b.config.AddDownloadCert,
		// 	Clean:        b.config.AddClean,
		// 	Force:        b.config.AddForce,
		// 	Insecure:     b.config.AddInsecure,
		// 	Provider:     b.config.Provider,
		// 	Address:      b.config.VMName,
		// },
		// Don't need an http server when vagrant does sharing for us.
		&StepSyncedFolder{
			SyncedFolder: b.config.SyncedFolder,
		},
		&StepUp{},
		// In StepUp, we get ssh information from the vagrant up command stdout.
		// and save it to state. This function wraps communicator.StepConnect
		// so that we can pass in the information we need.
		&StepSSHConfig{},
		&communicator.StepConnect{
			Config:    &b.config.SSHConfig.Comm,
			Host:      CommHost(),
			SSHConfig: b.config.SSHConfig.Comm.SSHConfigFunc(),
		},
		new(common.StepProvision),
		&StepHalt{
			b.config.TeardownMethod,
		},

		// Step package box
	}

	// Run the steps.
	b.runner = common.NewRunnerWithPauseFn(steps, b.config.PackerConfig, ui, state)
	b.runner.Run(state)

	// Report any errors.
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	// If we were interrupted or cancelled, then just exit.
	if _, ok := state.GetOk(multistep.StateCancelled); ok {
		return nil, errors.New("Build was cancelled.")
	}

	if _, ok := state.GetOk(multistep.StateHalted); ok {
		return nil, errors.New("Build was halted.")
	}

	return NewArtifact(b.config.OutputDir)
}

// Cancel.
func (b *Builder) Cancel() {
	if b.runner != nil {
		log.Println("Cancelling the step runner...")
		b.runner.Cancel()
	}
}
