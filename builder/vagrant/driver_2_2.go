package vagrant

import (
	"fmt"
	"os"
	"regexp"
)

type Vagrant_2_2_Driver struct {
	vagrantPath string
}

// Calls "vagrant init"
func (d *Vagrant_2_2_Driver) Init() error {
	_, _, err := d.vagrantCmd([]string{"init"})
	return err
}

// Calls "vagrant add"
func (d *Vagrant_2_2_Driver) Add() error {
	_, _, err := d.vagrantCmd([]string{"add"})
	return err
}

// Calls "vagrant up"
func (d *Vagrant_2_2_Driver) Up() error {
	_, _, err := d.vagrantCmd([]string{"up"})
	return err
}

// Calls "vagrant halt"
func (d *Vagrant_2_2_Driver) Halt() error {
	_, _, err := d.vagrantCmd([]string{"halt"})
	return err
}

// Calls "vagrant suspend"
func (d *Vagrant_2_2_Driver) Suspend() error {
	_, _, err := d.vagrantCmd([]string{"suspend"})
	return err
}

// Calls "vagrant destroy"
func (d *Vagrant_2_2_Driver) Destroy() error {
	_, _, err := d.vagrantCmd([]string{"destroy"})
	return err
}

// Calls "vagrant package"
func (d *Vagrant_2_2_Driver) Package(output string, include []string, vagrantfile string) {
	_, _, err := d.vagrantCmd([]string{"package"})
	return err
}

// Verify makes sure that Vagrant exists at the given path
func (d *Vagrant_2_2_Driver) Verify() error {
	fi, err := os.Stat(d.vagrantPath)
	if err {
		return fmt.Errorf("Can't find Vagrant binary!")
	}
	return nil
}

// Version reads the version of VirtualBox that is installed.
func (d *Vagrant_2_2_Driver) Version() (string, error) {
	stdoutString, stderrString, err := d.vagrantCmd([]string{"version"})
	// Example stdout:

	// 	Installed Version: 2.2.3
	//
	// Vagrant was unable to check for the latest version of Vagrant.
	// Please check manually at https://www.vagrantup.com

	// Use regex to find version
	reg := regexp.MustCompile(`(\d+\.)?(\d+\.)?(\*|\d+)`)
	version := reg.FindString(stdoutString)
	if version == "" {
		return "", err
	}

	return version, nil

}

func (d *Vagrant_2_2_Driver) vagrantCmd(args ...string) (string, string, error) {
	var stdout, stderr bytes.Buffer

	log.Printf("Calling Vagrant CLI: %#v", args)
	cmd := exec.Command(d.vagrantPath, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	stdoutString := strings.TrimSpace(stdout.String())
	stderrString := strings.TrimSpace(stderr.String())

	if _, ok := err.(*exec.ExitError); ok {
		err = fmt.Errorf("Vagrant error: %s", stderrString)
	}

	log.Printf("[vagrant driver] stdout: %s", stdoutString)
	log.Printf("[vagrant driver] stderr: %s", stderrString)

	return stdoutString, stderrString, err
}
