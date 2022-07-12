package kernels

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var ConfigFname = "conf.json"

type InitDirFlags struct {
	Force      bool
	BackupConf bool
}

// Initalizes a new directory (it will create it if it does not exist).
// the provided conf will be saved in the directory.
// if conf is nil, an empty configuration will be used.
func InitDir(log *logrus.Logger, dir string, conf *Conf, flags InitDirFlags) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create directory '%s': %w", dir, err)
	}

	confFname := path.Join(dir, ConfigFname)
	if !flags.Force {
		if _, err := os.Stat(confFname); err == nil {
			return fmt.Errorf("config file `%s` already exists", dir)
		}
	}

	if conf == nil {
		conf = &Conf{
			Kernels:    make([]KernelConf, 0),
			CommonOpts: make([]ConfigOption, 0),
		}
	}

	return conf.SaveTo(log, dir, flags.BackupConf)
}

// Load configuration from a directory
func LoadDir(dir string) (*KernelsDir, error) {
	data, err := os.ReadFile(path.Join(dir, ConfigFname))
	if err != nil {
		return nil, err
	}

	ks := KernelsDir{Dir: dir}
	err = json.Unmarshal(data, &ks.Conf)
	if err != nil {
		return nil, err
	}
	return &ks, nil
}

func AddKernel(log *logrus.Logger, dir string, cnf *KernelConf, backupConf bool) error {
	kd, err := LoadDir(dir)
	if err != nil {
		return err
	}

	if kd.KernelConfig(cnf.Name) != nil {
		return fmt.Errorf("kernel `%s` already exists", cnf.Name)
	}

	kd.Conf.Kernels = append(kd.Conf.Kernels, *cnf)
	return kd.Conf.SaveTo(log, dir, backupConf)
}

func RemoveKernel(ctx context.Context, log *logrus.Logger, dir string, name string, backupConf bool) error {
	kd, err := LoadDir(dir)
	if err != nil {
		return err
	}

	if kd.RemoveKernelConfig(name) == nil {
		log.Warnf("kernel `%s` does not exist in configuration", name)
	} else {
		defer kd.Conf.SaveTo(log, dir, backupConf)
	}

	gitRemoveWorkdir(ctx, log, &gitRemoveWorkdirArg{
		workDir:     name,
		bareDir:     filepath.Join(dir, MainGitDir),
		remoteName:  name,
		localBranch: fmt.Sprintf("lvh-%s", name),
	})

	return nil
}

func BuildKernel(ctx context.Context, log *logrus.Logger, dir, kname string) error {
	kd, err := LoadDir(dir)
	if err != nil {
		return err
	}

	kconf := kd.KernelConfig(kname)
	if kconf == nil {
		return fmt.Errorf("kernel `%s` not found", kname)
	}

	kURL, err := ParseURL(kconf.URL)
	if err != nil {
		return err
	}

	err = kURL.Fetch(context.Background(), log, dir, kconf.Name)
	if err != nil {
		return err
		log.Fatal(err)
	}

	return kd.buildKernel(context.Background(), log, dir, kconf)
}
