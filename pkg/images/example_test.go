package images

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExampleJSON(t *testing.T) {
	orig := &ExampleImagesConf
	origb, err := json.Marshal(&orig)
	assert.Nil(t, err)
	var parsed []ImgConf
	err = json.Unmarshal(origb, &parsed)
	assert.Nil(t, err)
	assert.Equal(t, orig, &parsed)
}
