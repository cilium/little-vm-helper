package images

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/sirupsen/logrus"
)

// doBuildImageDryRun just creates an empty file for the image.
func (f *ImageForest) doBuildImageDryRun(image string) error {
	_, ok := f.confs[image]
	if !ok {
		return fmt.Errorf("building image '%s' failed, configuration not found", image)
	}

	fname := f.imageFilename(image)
	file, err := os.Create(fname)
	defer file.Close()

	return err
}

func (f *ImageForest) doBuildImage(ctx context.Context, log *logrus.Logger, image string) error {
	cnf, ok := f.confs[image]
	if !ok {
		return fmt.Errorf("building image '%s' failed, configuration not found", image)
	}

	stepConf := &StepConf{
		imagesDir: f.imagesDir,
		imgCnf:    cnf,
		log:       log,
	}

	state := new(multistep.BasicStateBag)
	steps := make([]multistep.Step, 1+len(cnf.Actions))
	steps[0] = NewCreateImage(stepConf)
	fmt.Printf("%d\n", len(cnf.Actions))
	fmt.Printf("%v\n", cnf.Actions)
	for i := 0; i < len(cnf.Actions); i++ {
		steps[1+i] = cnf.Actions[i].Op.ToStep(stepConf)
	}

	runner := &multistep.BasicRunner{Steps: steps}
	runner.Run(ctx, state)
	err := state.Get("err")
	if err != nil {
		imgFname := f.imageFilename(image)
		log.Warnf("image file '%s' left for inspection", imgFname)
		return err.(error)
	}
	return nil
}
