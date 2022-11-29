// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package kernels

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testDirectoryExist(t *testing.T) {
	dir, err := os.MkdirTemp("", "dir-exists-test-")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	dirname := filepath.Join(dir, "dir")
	exists, err := directoryExists(dirname)
	assert.Nil(t, err)
	assert.False(t, exists)

	err = os.Mkdir(dirname, 0755)
	assert.Nil(t, err)
	exists, err = directoryExists(dirname)
	assert.Nil(t, err)
	assert.True(t, exists)

	filename := filepath.Join(dir, "filename")
	err = ioutil.WriteFile(filename, []byte("hello"), 0644)
	assert.Nil(t, err)

	exists, err = directoryExists(filename)
	assert.NotNil(t, err)
}
