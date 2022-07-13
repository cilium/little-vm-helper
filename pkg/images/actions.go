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
// add a test case in actions_json_test.go to ensure that all works.
type ActionOp interface {
	ActionOpName() string
	ToStep(s *StepConf) multistep.Step
}

type Action struct {
	Comment string
	Op      ActionOp
}

var actionOpInstances = []ActionOp{
	&RunCommand{},
	&CopyInCommand{},
}

// RunCommand runs a script in a path specified by a string
type RunCommand struct {
	Cmd string
}

func (rc *RunCommand) ActionOpName() string {
	return "run-command"
}

func (rc *RunCommand) ToStep(s *StepConf) multistep.Step {
	return &VirtCustomizeStep{
		StepConf: s,
		Args:     []string{"--run-command", rc.Cmd},
	}
}

// CopyInCommand copies local files in the image (recursively)
type CopyInCommand struct {
	LocalPath string
	RemoteDir string
}

func (c *CopyInCommand) ActionOpName() string {
	return "copy-in"
}

func (c *CopyInCommand) ToStep(s *StepConf) multistep.Step {
	return &VirtCustomizeStep{
		StepConf: s,
		Args:     []string{"--copy-in", fmt.Sprintf("%s:%s", c.LocalPath, c.RemoteDir)},
	}
}
