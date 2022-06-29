package images

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageBuilderConfs(t *testing.T) {
	tests := []struct {
		confs []ImageConf
		test  func(*Builder, error)
	}{
		{
			confs: []ImageConf{
				ImageConf{Name: "base"},
				ImageConf{Name: "base"},
			},
			test: func(ib *Builder, err error) {
				assert.Nil(t, ib)
				assert.Error(t, err)
			},
		},
		{
			confs: []ImageConf{
				ImageConf{Name: "base"},
			},
			test: func(ib *Builder, err error) {
				assert.NotNil(t, ib)
				assert.Nil(t, err)
			},
		},
		{
			// error: parent is not defined anywhere
			confs: []ImageConf{
				ImageConf{Name: "image1", Parent: "base"},
			},
			test: func(ib *Builder, err error) {
				assert.Nil(t, ib)
				assert.NotNil(t, err)
			},
		},
		{
			confs: []ImageConf{
				ImageConf{Name: "base"},
				ImageConf{Name: "image1", Parent: "base"},
				ImageConf{Name: "image2", Parent: "image1"},
			},
			test: func(ib *Builder, err error) {
				assert.NotNil(t, ib)
				assert.Nil(t, err)
				{
					deps, err := ib.getDependencies("image1")
					assert.Nil(t, err, "unexpected error: %v", err)
					assert.Equal(t, deps, []string{"base"})
				}
				{
					deps, err := ib.getDependencies("image2")
					assert.Nil(t, err, "unexpected error: %v", err)
					assert.Equal(t, deps, []string{"base", "image1"})
				}
				assert.Equal(t, ib.LeafImages(), []string{"image2"})
				assert.Equal(t, ib.RootImages(), []string{"base"})
			},
		},
	}

	for i := range tests {
		// NB: anonymous function so that os.RemoveAll() is called in all iterations
		func() {
			dir, err := ioutil.TempDir("", "test_images")
			assert.Nil(t, err)
			defer os.RemoveAll(dir)
			test := &tests[i]
			conf := &BuilderConf{
				ImageDir:     dir,
				Images:       test.confs,
				saveConfFile: true,
			}
			ib, err := NewImageBuilder(conf)
			test.test(ib, err)
			// if no errors, verify that conf.json matches the configuration
			if err == nil {
				data, err := os.ReadFile(path.Join(dir, DefaultConfFile))
				assert.Nil(t, err)
				var fcnf BuilderConf
				err = json.Unmarshal(data, &fcnf)
				assert.Nil(t, err)
				assert.Equal(t, &conf.Images, &fcnf.Images)
			}
		}()
	}
}
