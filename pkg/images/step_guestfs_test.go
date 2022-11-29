// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package images

import (
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/stretchr/testify/assert"
)

func TestACtionMerge(t *testing.T) {
	cnf := StepConf{}

	getSingleStep := func(op ActionOp) multistep.Step {
		steps, err := op.ToSteps(&cnf)
		assert.Nil(t, err)
		assert.Len(t, steps, 1)
		return steps[0]
	}

	s1 := getSingleStep(&CopyInCommand{LocalPath: "a", RemoteDir: "/"})
	s2 := getSingleStep(&CopyInCommand{LocalPath: "b", RemoteDir: "/"})
	err := mergeSteps(s1, s2)
	assert.Nil(t, err)
}
