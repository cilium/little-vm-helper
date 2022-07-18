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
		confs []ImgConf
		test  func(*ImageForest, error)
	}{
		{
			confs: []ImgConf{
				{Name: "base"},
				{Name: "base"},
			},
			test: func(f *ImageForest, err error) {
				assert.Nil(t, f)
				assert.Error(t, err)
			},
		},
		{
			confs: []ImgConf{
				{Name: "base"},
			},
			test: func(f *ImageForest, err error) {
				assert.NotNil(t, f)
				assert.Nil(t, err)
			},
		},
		{
			confs: []ImgConf{
				{Name: "base"},
				{Name: "image1", Parent: "base"},
				{Name: "image2", Parent: "image1"},
			},
			test: func(f *ImageForest, err error) {
				assert.NotNil(t, f)
				assert.Nil(t, err)
				{
					deps, err := f.getDependencies("image1")
					assert.Nil(t, err, "unexpected error: %v", err)
					assert.Equal(t, deps, []string{"base"})
				}
				{
					deps, err := f.getDependencies("image2")
					assert.Nil(t, err, "unexpected error: %v", err)
					assert.Equal(t, deps, []string{"base", "image1"})
				}
				assert.Equal(t, f.LeafImages(), []string{"image2"})
				assert.Equal(t, f.RootImages(), []string{"base"})
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
			conf := &ImagesConf{
				Dir:    dir,
				Images: test.confs,
			}
			f, err := NewImageForest(conf, true)
			test.test(f, err)
			// if no errors, verify that conf.json matches the configuration
			if err == nil {
				data, err := os.ReadFile(path.Join(dir, DefaultConfFile))
				assert.Nil(t, err)
				var fcnf ImagesConf
				err = json.Unmarshal(data, &fcnf)
				assert.Nil(t, err)
				assert.Equal(t, &conf.Images, &fcnf.Images)
			}
		}()
	}
}
