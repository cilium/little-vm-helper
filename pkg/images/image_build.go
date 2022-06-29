package images

import (
	"context"
	"fmt"
	"os"

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

	state := new(multistep.BasicStateBag)
	steps := []multistep.Step{
		&CreateImage{
			imageDir: ib.imageDir,
			imgCnf:   cnf,
			log:      log,
		},
	}

	runner := &multistep.BasicRunner{Steps: steps}
	runner.Run(ctx, state)
	err := state.Get("err")
	if err != nil {
		return err.(error)
	}
	return nil
}
