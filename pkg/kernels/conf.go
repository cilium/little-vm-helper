package kernels

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/cilium/little-vm-helper/pkg/logcmd"
	"github.com/sirupsen/logrus"
)

// ConfigOption are switches passed to scripts/config in a kernel dir
type ConfigOption []string

// KernelConf is the configuration of a kernel (to build from source)
type KernelConf struct {
	Name string `json:"name"`
	// URL of the kernel source
	URL string `json:"url"`
	// config options
	Conf []ConfigOption `json:"opts"`
}

type Conf struct {
	Kernels []KernelConf `json:"kernels"`
}

var ConfigOptGroups = map[string][]ConfigOption{
	"basic": []ConfigOption{
		{"--enable", "CONFIG_LOCALVERSION_AUTO"},
		{"--enable", "CONFIG_DEBUG_INFO"},
		{"--disable", "CONFIG_WERROR"},
	},
	"minimize": []ConfigOption{
		{"--disable", "CONFIG_DRM"},
		{"--disable", "CONFIG_GPU"},
		{"--disable", "CONFIG_CDROM"},
		{"--disable", "CONFIG_ISO9669_FS"},
		// wireless
		{"--disable", "CONFIG_CFG80211"},
		{"--disable", "CONFIG_RFKILL"},
		{"--disable", "CONFIG_MACINTOSH_DRIVERS"},
		{"--disable", "CONFIG_SOUND"},
	},
	"bpf": []ConfigOption{
		{"--enable", "CONFIG_BPF"},
		{"--enable", "CONFIG_HAVE_BPF_JIT"},
		{"--enable", "CONFIG_HAVE_EBPF_JIT"},
		{"--enable", "CONFIG_BPF_EVENTS"},
		{"--enable", "CONFIG_TEST_BPF"},
	},
	"virtio": []ConfigOption{
		{"--enable", "CONFIG_VIRTIO"},
		{"--enable", "CONFIG_VIRTIO_MENU"},
		{"--enable", "CONFIG_VIRTIO_PCI_LIB"},
		{"--enable", "CONFIG_VIRTIO_PCI"},
		{"--enable", "CONFIG_VIRTIO_NET"},
		{"--enable", "CONFIG_NET_9P_VIRTIO"},
		{"--enable", "CONFIG_VIRTIO_BLK"},
	},
}

func GetConfigGroupNames() []string {
	ret := make([]string, 0, len(ConfigOptGroups))
	for k := range ConfigOptGroups {
		ret = append(ret, k)
	}
	return ret
}

func (cnf *Conf) SaveTo(dir string) error {
	fname := path.Join(dir, ConfigFname)
	confb, err := json.Marshal(cnf)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	err = os.WriteFile(fname, confb, 0666)
	if err != nil {
		return fmt.Errorf("error writing configuration: %w", err)
	}

	return nil
}

func (kc *KernelConf) Validate() error {
	_, err := ParseURL(kc.URL)
	return err
}

func (kc *KernelConf) AddGroups(gs ...string) error {
	newOpts := make([]ConfigOption, 0)

	for _, g := range gs {
		opts, ok := ConfigOptGroups[g]
		if !ok {
			return fmt.Errorf("unknown group %s", g)
		}
		for _, opt := range opts {
			newOpts = append(newOpts, opt)
		}
	}

	for _, opt := range newOpts {
		kc.Conf = append(kc.Conf, opt)
	}

	return nil
}

func (kc *KernelConf) Configure(ctx context.Context, log *logrus.Logger, dir string) error {
	srcDir := filepath.Join(dir, kc.Name)
	if err := logcmd.RunAndLogCmdContext(ctx, log, "make", "-C", srcDir, "defconfig", "prepare"); err != nil {
		return err
	}

	configCmd := filepath.Join(dir, kc.Name, "scripts", "config")
	for _, opts := range kc.Conf {
		if err := logcmd.RunAndLogCmdContext(ctx, log, configCmd, opts...); err != nil {
			return err
		}
	}

	return nil
}

func (kc *KernelConf) Build(ctx context.Context, log *logrus.Logger, dir string) error {
	srcDir := filepath.Join(dir, kc.Name)
	err := logcmd.RunAndLogCmdContext(ctx, log, "make", "-C", srcDir, "-j", fmt.Sprintf("%d", runtime.NumCPU()))
	return err
}
