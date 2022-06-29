package images

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

func TestImageBuilds(t *testing.T) {
	xlog := logrus.New()
	xlog.SetOutput(testLogger{t})

	tests := []struct {
		confs   []ImageConf
		prepare func(imagesDir string)
		test    func(imagesDir string, res *BuilderResult)
	}{
		{
			confs: []ImageConf{
				ImageConf{Name: "base"},
			},
			test: func(dir string, r *BuilderResult) {
				assert.Nil(t, r.Err())
				assert.Equal(t, 1, len(r.ImageResults))
				assert.False(t, r.ImageResults["base"].CachedImageUsed)
				assert.Equal(t, "", r.ImageResults["base"].CachedImageDeleted)
				assert.FileExists(t, path.Join(dir, fmt.Sprintf("%s.%s", "base", DefaultImageExt)))
			},
		}, {
			confs: []ImageConf{
				ImageConf{Name: "base"},
				ImageConf{Name: "image1", Parent: "base"},
				ImageConf{Name: "image2", Parent: "image1"},
			},
			test: func(dir string, r *BuilderResult) {
				assert.Nil(t, r.Err())
				assert.Equal(t, 3, len(r.ImageResults))
				for _, fname := range []string{"base", "image1", "image2"} {
					assert.False(t, r.ImageResults[fname].CachedImageUsed)
					assert.Equal(t, "", r.ImageResults[fname].CachedImageDeleted)
					assert.FileExists(t, path.Join(dir, fmt.Sprintf("%s.%s", fname, DefaultImageExt)))
				}
			},
		}, {
			confs: []ImageConf{
				ImageConf{Name: "base"},
				ImageConf{Name: "image1", Parent: "base"},
				ImageConf{Name: "image2", Parent: "image1"},
			},
			prepare: func(dir string) {
				fname := path.Join(dir, fmt.Sprintf("%s.%s", "base", DefaultImageExt))
				f, err := os.Create(fname)
				assert.Nil(t, err)
				defer f.Close()
			},
			test: func(dir string, r *BuilderResult) {
				assert.Nil(t, r.Err())
				assert.Equal(t, 3, len(r.ImageResults))

				assert.True(t, r.ImageResults["base"].CachedImageUsed)
				assert.Equal(t, "", r.ImageResults["base"].CachedImageDeleted)
				assert.FileExists(t, path.Join(dir, fmt.Sprintf("%s.%s", "base", DefaultImageExt)))

				assert.False(t, r.ImageResults["image1"].CachedImageUsed)
				assert.Equal(t, "", r.ImageResults["image1"].CachedImageDeleted)
				assert.FileExists(t, path.Join(dir, fmt.Sprintf("%s.%s", "image1", DefaultImageExt)))

				assert.False(t, r.ImageResults["image2"].CachedImageUsed)
				assert.Equal(t, "", r.ImageResults["image2"].CachedImageDeleted)
				assert.FileExists(t, path.Join(dir, fmt.Sprintf("%s.%s", "image2", DefaultImageExt)))
			},
		}, {
			confs: []ImageConf{
				ImageConf{Name: "base"},
				ImageConf{Name: "image1", Parent: "base"},
				ImageConf{Name: "image2", Parent: "image1"},
			},
			prepare: func(dir string) {
				fname := path.Join(dir, fmt.Sprintf("%s.%s", "image1", DefaultImageExt))
				f, err := os.Create(fname)
				assert.Nil(t, err)
				defer f.Close()
			},
			test: func(dir string, r *BuilderResult) {
				assert.Nil(t, r.Err())
				assert.Equal(t, 3, len(r.ImageResults))

				assert.False(t, r.ImageResults["base"].CachedImageUsed)
				assert.Equal(t, "", r.ImageResults["base"].CachedImageDeleted)
				assert.FileExists(t, path.Join(dir, fmt.Sprintf("%s.%s", "base", DefaultImageExt)))

				assert.False(t, r.ImageResults["image1"].CachedImageUsed)
				assert.NotEqual(t, "", r.ImageResults["image1"].CachedImageDeleted)
				assert.FileExists(t, path.Join(dir, fmt.Sprintf("%s.%s", "image1", DefaultImageExt)))

				assert.False(t, r.ImageResults["image2"].CachedImageUsed)
				assert.Equal(t, "", r.ImageResults["image2"].CachedImageDeleted)
				assert.FileExists(t, path.Join(dir, fmt.Sprintf("%s.%s", "image2", DefaultImageExt)))
			},
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
			assert.Nil(t, err)
			defer os.RemoveAll(dir)
			test := &tests[i]
			if test.prepare != nil {
				test.prepare(dir)
			}
			conf := &BuilderConf{
				ImageDir: dir,
				Images:   test.confs,
			}
			ib, err := NewImageBuilder(conf)
			assert.Nil(t, err)
			res := ib.BuildAllImages(&bldConf)
			test.test(dir, res)
		}()
	}
}
