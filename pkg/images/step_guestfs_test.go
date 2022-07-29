package images

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestACtionMerge(t *testing.T) {
	cnf := StepConf{}
	s1 := (&CopyInCommand{LocalPath: "a", RemoteDir: "/"}).ToStep(&cnf)
	s2 := (&CopyInCommand{LocalPath: "b", RemoteDir: "/"}).ToStep(&cnf)
	err := mergeSteps(s1, s2)
	assert.Nil(t, err)
}
