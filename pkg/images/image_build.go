package images

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/sirupsen/logrus"
)

// doBuildImageDryRun just creates an empty file for the image.
func (ib *Builder) doBuildImageDryRun(image string) error {
	_, ok := ib.confs[image]
	if !ok {
		return fmt.Errorf("building image '%s' failed, configuration not found", image)
	}

	fname := fmt.Sprintf("%s.%s", ib.imageFilenamePrefix(image), DefaultImageExt)
	f, err := os.Create(fname)
	defer f.Close()

	return err
}

func (ib *Builder) doBuildImage(ctx context.Context, log *logrus.Logger, image string) error {
	cnf, ok := ib.confs[image]
	if !ok {
		return fmt.Errorf("building image '%s' failed, configuration not found", image)
	}

	stepConf := &StepConf{
		imageDir: ib.imageDir,
		imgCnf:   cnf,
		log:      log,
	}

	state := new(multistep.BasicStateBag)
	steps := make([]multistep.Step, 1+len(cnf.Actions))
	steps[0] = &CreateImage{stepConf}
	for i := 0; i < len(cnf.Actions); i++ {
		steps[1+i] = cnf.Actions[i].Op.ToStep(stepConf)
	}

	runner := &multistep.BasicRunner{Steps: steps}
	runner.Run(ctx, state)
	err := state.Get("err")
	if err != nil {
		imgFname := path.Join(ib.imageDir, fmt.Sprintf("%s.%s", cnf.Name, DefaultImageExt))
		log.Warnf("image file '%s' left for inspection", imgFname)
		return err.(error)
	}
	return nil
}
