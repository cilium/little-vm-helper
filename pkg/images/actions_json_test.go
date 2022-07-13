package images

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActionsJSON(t *testing.T) {
	acts := []Action{
		{Op: &RunCommand{Cmd: "echo hello!"}, Comment: "hello"},
		{Op: &CopyInCommand{LocalPath: "/foo", RemoteDir: "/"}, Comment: "copy"},
	}

	for i := range acts {
		act := acts[i]
		actb, err := json.Marshal(&act)
		assert.Nil(t, err)
		var xact Action
		err = json.Unmarshal(actb, &xact)
		assert.Nil(t, err)
		assert.Equal(t, act, xact)
	}
}
