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

func (kd *KernelsDir) ConfigureKernel(ctx context.Context, log *logrus.Logger, dir, kernName string) error {
	kc := kd.KernelConfig(kernName)
	if kc == nil {
		return fmt.Errorf("kernel '%s' not found", kernName)
	}
	return kd.configureKernel(ctx, log, dir, kc)
}

func (kd *KernelsDir) configureKernel(ctx context.Context, log *logrus.Logger, dir string, kc *KernelConf) error {
	srcDir := filepath.Join(dir, kc.Name)

	oldPath, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(srcDir)
	if err != nil {
		return err
	}
	defer os.Chdir(oldPath)

	if err := logcmd.RunAndLogCmdContext(ctx, log, MakeBinary, "defconfig", "prepare"); err != nil {
		return err
	}

	configCmd := filepath.Join(".", "scripts", "config")

	// common options first
	for _, opts := range kd.Conf.CommonOpts {
		if err := logcmd.RunAndLogCmdContext(ctx, log, configCmd, opts...); err != nil {
			return err
		}
	}

	// then kernel-specific options
	for _, opts := range kc.Opts {
		if err := logcmd.RunAndLogCmdContext(ctx, log, configCmd, opts...); err != nil {
			return err
		}
	}

	// run make olddefconfig to clean up the config file
	if err := logcmd.RunAndLogCmdContext(ctx, log, MakeBinary, "olddefconfig"); err != nil {
		return err
	}

	return nil
}

func (kd *KernelsDir) buildKernel(ctx context.Context, log *logrus.Logger, dir string, kc *KernelConf) error {
	if err := CheckEnvironment(); err != nil {
		return err
	}

	srcDir := filepath.Join(dir, kc.Name)
	configFname := filepath.Join(srcDir, ".config")

	if exists, err := regularFileExists(configFname); err != nil {
		return err
	} else if !exists {
		log.Info("Configuring kernel")
		err = kd.configureKernel(ctx, log, dir, kc)
		if err != nil {
			return fmt.Errorf("failed to configure kernel: %w", err)
		}
	}

	ncpus := fmt.Sprintf("%d", runtime.NumCPU())
	if err := logcmd.RunAndLogCmdContext(ctx, log, MakeBinary, "-C", srcDir, "-j", ncpus, "bzImage"); err != nil {
		return fmt.Errorf("buiding bzImage failed: %w", err)
	}

	if err := logcmd.RunAndLogCmdContext(ctx, log, MakeBinary, "-C", srcDir, "dir-pkg"); err != nil {
		return fmt.Errorf("build dir failed: %w", err)
	}

	return nil
}
