package images

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type testLogger struct {
	*testing.T
}

func (tl testLogger) Write(p []byte) (n int, err error) {
	tl.Logf((string)(p))
	return len(p), nil
}

func init() {
	ignoreEnviromentCheckForTesting = true
}

func TestImageBuilds(t *testing.T) {
	xlog := logrus.New()
	xlog.SetOutput(testLogger{t})

	tests := []struct {
		confs        []ImgConf
		prepare      func(imagesDir string)
		test         func(imagesDir string, res *BuilderResult)
		forceRebuild bool
	}{
		{
			confs: []ImgConf{
				{Name: "base"},
			},
			test: func(dir string, r *BuilderResult) {
				assert.Nil(t, r.Err())
				assert.Equal(t, 1, len(r.ImageResults))
				assert.False(t, r.ImageResults["base"].CachedImageUsed)
				assert.Equal(t, "", r.ImageResults["base"].CachedImageDeleted)
				assert.FileExists(t, filepath.Join(dir, "base"))
			},
		}, {
			confs: []ImgConf{
				{Name: "base"},
				{Name: "image1", Parent: "base"},
				{Name: "image2", Parent: "image1"},
			},
			test: func(dir string, r *BuilderResult) {
				assert.Nil(t, r.Err())
				assert.Equal(t, 3, len(r.ImageResults))
				for _, fname := range []string{"base", "image1", "image2"} {
					assert.False(t, r.ImageResults[fname].CachedImageUsed)
					assert.Equal(t, "", r.ImageResults[fname].CachedImageDeleted)
					assert.FileExists(t, filepath.Join(dir, fname))
				}
			},
		}, {
			confs: []ImgConf{
				{Name: "base"},
				{Name: "image1", Parent: "base"},
				{Name: "image2", Parent: "image1"},
			},
			prepare: func(dir string) {
				fname := filepath.Join(dir, "base")
				f, err := os.Create(fname)
				assert.Nil(t, err)
				defer f.Close()
			},
			test: func(dir string, r *BuilderResult) {
				assert.Nil(t, r.Err())
				assert.Equal(t, 3, len(r.ImageResults))

				assert.True(t, r.ImageResults["base"].CachedImageUsed)
				assert.Equal(t, "", r.ImageResults["base"].CachedImageDeleted)
				assert.FileExists(t, filepath.Join(dir, "base"))

				assert.False(t, r.ImageResults["image1"].CachedImageUsed)
				assert.Equal(t, "", r.ImageResults["image1"].CachedImageDeleted)
				assert.FileExists(t, filepath.Join(dir, "image1"))

				assert.False(t, r.ImageResults["image2"].CachedImageUsed)
				assert.Equal(t, "", r.ImageResults["image2"].CachedImageDeleted)
				assert.FileExists(t, filepath.Join(dir, "image2"))
			},
		}, {
			confs: []ImgConf{
				{Name: "base"},
				{Name: "image1", Parent: "base"},
				{Name: "image2", Parent: "image1"},
			},
			prepare: func(dir string) {
				fname := filepath.Join(dir, "image1")
				f, err := os.Create(fname)
				assert.Nil(t, err)
				defer f.Close()
			},
			test: func(dir string, r *BuilderResult) {
				assert.Nil(t, r.Err())
				assert.Equal(t, 3, len(r.ImageResults))

				assert.False(t, r.ImageResults["base"].CachedImageUsed)
				assert.Equal(t, "", r.ImageResults["base"].CachedImageDeleted)
				assert.FileExists(t, filepath.Join(dir, "base"))

				assert.False(t, r.ImageResults["image1"].CachedImageUsed)
				assert.NotEqual(t, "", r.ImageResults["image1"].CachedImageDeleted)
				assert.FileExists(t, filepath.Join(dir, "image1"))

				assert.False(t, r.ImageResults["image2"].CachedImageUsed)
				assert.Equal(t, "", r.ImageResults["image2"].CachedImageDeleted)
				assert.FileExists(t, filepath.Join(dir, "image2"))
			},
		}, {
			confs: []ImgConf{
				{Name: "base"},
			},
			prepare: func(dir string) {
				fname := filepath.Join(dir, "base")
				f, err := os.Create(fname)
				assert.Nil(t, err)
				defer f.Close()
			},
			test: func(dir string, r *BuilderResult) {
				assert.Nil(t, r.Err())
				assert.Equal(t, 1, len(r.ImageResults))
				assert.False(t, r.ImageResults["base"].CachedImageUsed)
				assert.NotEqual(t, "", r.ImageResults["base"].CachedImageDeleted)
			},
			forceRebuild: true,
		},
	}

	bldConf := BuildConf{
		Log:    xlog,
		DryRun: true,
	}

	for i := range tests {
		// NB: anonymous function so that os.RemoveAll() is called in all iterations
		func() {
			dir, err := ioutil.TempDir("", "test_build_images")
			imagesDir := filepath.Join(dir, DefaultImagesDir)
			os.Mkdir(imagesDir, 0755)
			assert.Nil(t, err)
			defer os.RemoveAll(dir)
			test := &tests[i]
			if test.prepare != nil {
				test.prepare(imagesDir)
			}
			conf := &ImagesConf{
				Dir:    dir,
				Images: test.confs,
			}
			ib, err := NewImageForest(conf, false)
			assert.Nil(t, err)
			bldConf.ForceRebuild = test.forceRebuild
			res := ib.BuildAllImages(&bldConf)
			test.test(imagesDir, res)
		}()
	}
}
