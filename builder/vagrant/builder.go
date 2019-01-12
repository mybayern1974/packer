package vagrant

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

// Builder implements packer.Builder and builds the actual VirtualBox
// images.
type Builder struct {
	config *Config
	runner multistep.Runner
}

type Config struct {
	common.PackerConfig    `mapstructure:",squash"`
	common.HTTPConfig      `mapstructure:",squash"`
	common.ISOConfig       `mapstructure:",squash"`
	common.FloppyConfig    `mapstructure:",squash"`
	bootcommand.BootConfig `mapstructure:",squash"`

	// This is the name of the new virtual machine.
	// By default this is "packer-BUILDNAME", where "BUILDNAME" is the name of the build.
	OutputDir string `mapstructure:"output_dir"`
	VMName    string `mapstructure:"vm_name"`

	Communicator string `mapstructure:"communicator"`

	// What vagrantfile to use
	Vagrantfile string `mapstructure:"vagrantfile"`

	// Whether to Halt, Suspend, or Destroy the box
	TeardownMethod string `mapstructure:"teardown_method"`

	// Override the default provider
	Provider string `mapstructure:"provider"`

	// Don't package the Vagrant box after build.
	SkipPackage bool `mapstructure:"skip_package"`

	ctx interpolate.Context
}

// Prepare processes the build configuration parameters.
func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
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

	if b.config.OutputDir == "" {
		b.config.OutputDir = fmt.Sprintf("output-%s", c.PackerBuildName)
	}

	if b.config.TeardownMethod == "" {
		b.config.TeardownMethod = "halt"
	} else {
		matches := false
		for _, name := range []string{"halt", "suspend", "destroy"} {
			if strings.ToLower(b.config.TeradownMethod) == name {
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
			BoxName:    b.config.BoxName,
		},
		&StepAddBox{
			BoxVersion:   b.config.BoxVersion,
			CACert:       b.config.AddCACert,
			CAPath:       b.config.AddCAPath,
			DownloadCert: b.config.AddDownloadCert,
			Clean:        b.config.AddClean,
			Force:        b.config.AddForce,
			Insecure:     b.config.AddInsecure,
			Provider:     b.config.Provider,
			Address:      b.config.BoxName,
		},
		&StepUp{},
		// step load box

		// step provision
		new(common.StepProvision),

		// step shutdown

		// step package box
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
