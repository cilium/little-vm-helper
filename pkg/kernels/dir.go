package kernels

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cilium/little-vm-helper/pkg/logcmd"
	"github.com/sirupsen/logrus"
)

type KernelsDir struct {
	Dir  string
	Conf Conf
}

func (kd *KernelsDir) KernelConfig(name string) *KernelConf {
	for i := range kd.Conf.Kernels {
		if kd.Conf.Kernels[i].Name == name {
			return &kd.Conf.Kernels[i]
		}
	}

	return nil
}

// RemoveKernelConfig returns the removed kernel config if it was found
func (kd *KernelsDir) RemoveKernelConfig(name string) *KernelConf {
	for i := range kd.Conf.Kernels {
		if kd.Conf.Kernels[i].Name == name {
			ret := &kd.Conf.Kernels[i]
			kd.Conf.Kernels = append(kd.Conf.Kernels[:i], kd.Conf.Kernels[i+1:]...)
			return ret
		}
	}

	return nil
}

func (kd *KernelsDir) ConfigureKernel(ctx context.Context, log *logrus.Logger, kernName string) error {
	kc := kd.KernelConfig(kernName)
	if kc == nil {
		return fmt.Errorf("kernel '%s' not found", kernName)
	}
	return kd.configureKernel(ctx, log, kc)
}

func (kd *KernelsDir) configureKernel(ctx context.Context, log *logrus.Logger, kc *KernelConf) error {
	srcDir := filepath.Join(kd.Dir, kc.Name)

	oldPath, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(srcDir)
	if err != nil {
		return err
	}
	defer os.Chdir(oldPath)

	if err := logcmd.RunAndLogCommandContext(ctx, log, MakeBinary, "defconfig", "prepare"); err != nil {
		return err
	}

	configCmd := filepath.Join(".", "scripts", "config")
	for _, opts := range kd.Conf.getOptions(kc) {
		// NB: we could do this in a single command, but doing it one-by-one makes it easier to debug things
		if err := logcmd.RunAndLogCommandContext(ctx, log, configCmd, opts...); err != nil {
			return err
		}
	}

	// run make olddefconfig to clean up the config file
	if err := logcmd.RunAndLogCommandContext(ctx, log, MakeBinary, "olddefconfig"); err != nil {
		return err
	}

	return nil
}

func (kd *KernelsDir) buildKernel(ctx context.Context, log *logrus.Logger, kc *KernelConf) error {
	if err := CheckEnvironment(); err != nil {
		return err
	}

	srcDir := filepath.Join(kd.Dir, kc.Name)
	configFname := filepath.Join(srcDir, ".config")

	if exists, err := regularFileExists(configFname); err != nil {
		return err
	} else if !exists {
		log.Info("Configuring kernel")
		err = kd.configureKernel(ctx, log, kc)
		if err != nil {
			return fmt.Errorf("failed to configure kernel: %w", err)
		}
	}

	ncpus := fmt.Sprintf("%d", runtime.NumCPU())
	if err := logcmd.RunAndLogCommandContext(ctx, log, MakeBinary, "-C", srcDir, "-j", ncpus, "bzImage"); err != nil {
		return fmt.Errorf("buiding bzImage failed: %w", err)
	}

	if err := logcmd.RunAndLogCommandContext(ctx, log, MakeBinary, "-C", srcDir, "dir-pkg"); err != nil {
		return fmt.Errorf("build dir failed: %w", err)
	}

	return nil
}
