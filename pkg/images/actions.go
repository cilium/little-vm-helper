package images

import (
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

// ActionOp is the interface that actions operations need to implement.
//
// Note:
// If you create an instance of ActionOp, you need to add it to
// actionOpInstances so that JSON marshaling/unmarshaling works. Please also
// consider adding a test case in actions_json_test.go to ensure that all
// works.
type ActionOp interface {
	ActionOpName() string
	ToSteps(s *StepConf) ([]multistep.Step, error)
}

type Action struct {
	Comment string
	Op      ActionOp
}

var actionOpInstances = []ActionOp{
	&RunCommand{},
	&CopyInCommand{},
	&SetHostnameCommand{},
	&MkdirCommand{},
	&UploadCommand{},
	&ChmodCommand{},
	&AppendLineCommand{},
}

type VirtCustomizeAction struct {
	OpName  string
	getArgs func() []string
}

// RunCommand runs a script in a path specified by a string
type RunCommand struct {
	Cmd string
}

func (rc *RunCommand) ActionOpName() string {
	return "run-command"
}

func (rc *RunCommand) ToSteps(s *StepConf) ([]multistep.Step, error) {
	return []multistep.Step{&VirtCustomizeStep{
		StepConf: s,
		Args:     []string{"--run-command", rc.Cmd},
	}}, nil
}

// CopyInCommand copies local files in the image (recursively)
type CopyInCommand struct {
	LocalPath string
	RemoteDir string
}

func (c *CopyInCommand) ActionOpName() string {
	return "copy-in"
}

func (c *CopyInCommand) ToSteps(s *StepConf) ([]multistep.Step, error) {
	return []multistep.Step{&VirtCustomizeStep{
		StepConf: s,
		Args:     []string{"--copy-in", fmt.Sprintf("%s:%s", c.LocalPath, c.RemoteDir)},
	}}, nil
}

// SetHostnameCommand sets the hostname
type SetHostnameCommand struct {
	Hostname string
}

func (c *SetHostnameCommand) ActionOpName() string {
	return "set-hostname"
}

func (c *SetHostnameCommand) ToSteps(s *StepConf) ([]multistep.Step, error) {
	return []multistep.Step{&VirtCustomizeStep{
		StepConf: s,
		Args:     []string{"--hostname", c.Hostname},
	}}, nil
}

// MkdirCommand creates a directory
type MkdirCommand struct {
	Dir string
}

func (c *MkdirCommand) ActionOpName() string {
	return "mkdir"
}

func (c *MkdirCommand) ToSteps(s *StepConf) ([]multistep.Step, error) {
	return []multistep.Step{&VirtCustomizeStep{
		StepConf: s,
		Args:     []string{"--mkdir", c.Dir},
	}}, nil
}

// UploadCommand copies a file to the vim
type UploadCommand struct {
	File string
	Dest string
}

func (c *UploadCommand) ActionOpName() string {
	return "upload"
}

func (c *UploadCommand) ToSteps(s *StepConf) ([]multistep.Step, error) {
	return []multistep.Step{&VirtCustomizeStep{
		StepConf: s,
		Args:     []string{"--upload", fmt.Sprintf("%s:%s", c.File, c.Dest)},
	}}, nil
}

// ChmodCommand
type ChmodCommand struct {
	Permissions string
	File        string
}

func (c *ChmodCommand) ActionOpName() string {
	return "chmod"
}

func (c *ChmodCommand) ToSteps(s *StepConf) ([]multistep.Step, error) {
	return []multistep.Step{&VirtCustomizeStep{
		StepConf: s,
		Args:     []string{"--chmod", fmt.Sprintf("%s:%s", c.Permissions, c.File)},
	}}, nil
}

// AppendLineCommand
type AppendLineCommand struct {
	File string
	Line string
}

func (c *AppendLineCommand) ActionOpName() string {
	return "append-line"
}

func (c *AppendLineCommand) ToSteps(s *StepConf) ([]multistep.Step, error) {
	return []multistep.Step{&VirtCustomizeStep{
		StepConf: s,
		Args:     []string{"--append-line", fmt.Sprintf("%s:%s", c.File, c.Line)},
	}}, nil
}
