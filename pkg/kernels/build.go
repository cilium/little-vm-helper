package kernels

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

func BuildKernel(ctx context.Context, log *logrus.Logger, dir string, name string) error {
	kd, err := LoadDir(dir)
	if err != nil {
		return err
	}

	kc := kd.KernelConfig(name)
	if kc == nil {
		return fmt.Errorf("kernel `%s` does not exist", name)
	}

	kurl, err := ParseURL(kc.URL)
	if err != nil {
		return err
	}

	if err = kurl.Fetch(ctx, log, dir, kc.Name); err != nil {
		return err
	}

	return nil
}
