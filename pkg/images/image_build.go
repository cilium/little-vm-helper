package images

import (
	"fmt"
	"os"
)

func (ib *ImageBuilder) doBuildImage(image string) error {
	_, ok := ib.confs[image]
	if !ok {
		return fmt.Errorf("building image '%s' failed, configuration not found", image)
	}

	return fmt.Errorf("NYI: %d", 1)
}

func (ib *ImageBuilder) doBuildImageDryRun(image string) error {
	_, ok := ib.confs[image]
	if !ok {
		return fmt.Errorf("building image '%s' failed, configuration not found", image)
	}

	fname := ib.imageFilename(image)
	f, err := os.Create(fname)
	defer f.Close()

	return err
}
