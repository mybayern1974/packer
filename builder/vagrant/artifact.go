package vagrant

import (
	"fmt"

	"github.com/hashicorp/packer/packer"
)

// This is the common builder ID to all of these artifacts.
const BuilderId = "vagrant"

// Artifact is the result of running the vagrant builder, namely a set
// of files associated with the resulting machine.
type artifact struct {
	boxName  string
	globalID string
}

// NewArtifact returns a hyperv artifact containing the files
// in the given directory.
func NewArtifact(dir string) (packer.Artifact, error) {

	return &artifact{}, nil
}

func (*artifact) BuilderId() string {
	return BuilderId
}

func (a *artifact) Files() []string {
	return []string{""}
}

func (*artifact) Id() string {
	return "Box"
}

func (a *artifact) String() string {
	return fmt.Sprintf("Vagrant box global ID is: %s", a.globalID)
}

func (a *artifact) State(name string) interface{} {
	return nil
}

func (a *artifact) Destroy() error {
	return nil
}
