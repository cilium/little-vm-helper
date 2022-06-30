package images

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// BuildImageResult describes the result of building a single image
type BuildResult struct {
	// Error is not nil, if building the image failed.
	Error error

	// CachedImageUsed is set to true if a cached image was found and no
	// actual build happened.
	CachedImageUsed bool

	// CachedImageDeleted is set to an non empty string if the image file
	// was deleted. The string describes the reason.
	CachedImageDeleted string
}

// BuilderResult  encodes the result of image builds
type BuilderResult struct {
	// Error is not nil if an error happened outside of buiding an image
	Error error
	// ImageResults results of building images
	ImageResults map[string]BuildResult
}

// Err() returns a summary error or nil if no errors were encountered
func (r *BuilderResult) Err() error {
	var imgErr strings.Builder
	imgErr.WriteString("images errors:")
	errCount := 0
	for image, res := range r.ImageResults {
		if res.Error != nil {
			if errCount > 0 {
				imgErr.WriteString("; ")
			}
			imgErr.WriteString(fmt.Sprintf("%s: %v", image, res.Error))
			errCount++
		}
	}

	if errCount == 0 {
		return r.Error
	}

	if r.Error == nil {
		return errors.New(imgErr.String())
	} else {
		return fmt.Errorf("builder error:%w %s", r.Error, imgErr.String())
	}
}

// BuildConf configures image builds
type BuildConf struct {
	Log *logrus.Logger

	// if DryRun set, no actual images will be build. Instead, empty files will be created
	DryRun bool
	// if ForceRebuild is set, images will be build even if they exist already
	ForceRebuild bool
}

type buildState struct {
	ib        *Builder
	bldConf   *BuildConf
	bldResult BuilderResult
}

func newBuildState(ib *Builder, cnf *BuildConf) *buildState {
	return &buildState{
		ib:      ib,
		bldConf: cnf,
		bldResult: BuilderResult{
			ImageResults: make(map[string]BuildResult),
		},
	}
}

func (ib *Builder) BuildAllImages(bldConf *BuildConf) *BuilderResult {

	log := bldConf.Log
	st := newBuildState(ib, bldConf)

	if err := ib.CheckEnvironment(); err != nil {
		st.bldResult.Error = fmt.Errorf("environment check failed: %w", err)
		return &st.bldResult
	}

	queue := ib.RootImages()
	log.WithFields(logrus.Fields{
		"queue": strings.Join(queue, ","),
	}).Info("starting to build all images")
	for {
		var image string
		if len(queue) == 0 {
			break
		}
		image, queue = queue[0], queue[1:]
		imgRes := st.buildImage(image)
		if imgRes.Error == nil {
			children := ib.children[image]
			queue = append(queue, children...)
		}

		xlog := log.WithFields(logrus.Fields{
			"image":  image,
			"queue":  strings.Join(queue, ","),
			"result": fmt.Sprintf("%+v", imgRes),
		})
		if imgRes.Error == nil {
			xlog.Info("image built succesfully")
		} else {
			xlog.Warn("image build failed")
		}
	}

	return &st.bldResult
}

func (b *buildState) buildImage(image string) BuildResult {
	res := b.doBuildImage(image)
	b.bldResult.ImageResults[image] = res
	return res
}

// skipRebuild checks if an image is not required to be build, because it
// already exists.
func (b *buildState) skipRebuild(image string) BuildResult {
	imageFnamePrefix, err := b.ib.ImageFilenamePrefix(image)
	if err != nil {
		return BuildResult{
			Error: err,
		}
	}
	imageFname := fmt.Sprintf("%s.%s", imageFnamePrefix, DefaultImageExt)
	if fi, err := os.Stat(imageFname); err == nil {
		mode := fi.Mode()
		if !mode.IsRegular() {
			// NB: we could do something like os.RemoveAll() here
			// but this is a weird case, so we just bail out
			err := fmt.Errorf("'%s' is not a regular file. Bailing out.", imageFname)
			return BuildResult{
				Error: err,
			}
		}

		if b.bldConf.ForceRebuild {
			os.Remove(imageFname)
			return BuildResult{
				CachedImageDeleted: fmt.Sprintf("image '%s' was deleted because a rebuild was forced", imageFname),
			}
		}

		if !b.bldConf.DryRun && fi.Size() == 0 {
			os.Remove(imageFname)
			return BuildResult{
				CachedImageDeleted: fmt.Sprintf("image '%s' was an empty file, and this was not a dry run", imageFname),
			}
		}

		if parent := b.ib.confs[image].Parent; parent != "" && !b.bldResult.ImageResults[parent].CachedImageUsed {
			os.Remove(imageFname)
			return BuildResult{
				CachedImageDeleted: fmt.Sprintf("image '%s' existed, but parent '%s' did not use the cache", imageFname, parent),
			}
		}

		return BuildResult{
			CachedImageUsed: true,
		}
	}
	// NB: we might want to check the error that this is an actual ENOENT error
	return BuildResult{}
}

func (b *buildState) doBuildImage(image string) BuildResult {

	imgRes := b.skipRebuild(image)
	if imgRes.Error != nil || imgRes.CachedImageUsed {
		return imgRes
	}

	buildImage := func(image string) error {
		return b.ib.doBuildImage(context.Background(), b.bldConf.Log, image)
	}
	if b.bldConf.DryRun {
		buildImage = b.ib.doBuildImageDryRun
	}
	imgRes.Error = buildImage(image)
	return imgRes
}
