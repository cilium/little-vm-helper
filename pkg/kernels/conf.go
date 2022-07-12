package kernels

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

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
	Opts []ConfigOption `json:"opts,omitempty"`
}

type Conf struct {
	Kernels    []KernelConf   `json:"kernels"`
	CommonOpts []ConfigOption `json:"common_opts,omitempty"`
}

func confAddGroups(opts []ConfigOption, gs ...string) ([]ConfigOption, error) {
	newOpts := make([]ConfigOption, 0)
	for _, g := range gs {
		nopts, ok := ConfigOptGroups[g]
		if !ok {
			return nil, fmt.Errorf("unknown group %s", g)
		}
		for _, opt := range nopts {
			newOpts = append(newOpts, opt)
		}
	}

	for _, opt := range newOpts {
		opts = append(opts, opt)
	}

	return opts, nil
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
		{"--disable", "CONFIG_AGP"},
		{"--disable", "CONFIG_USB_SUPPORT"},
		{"--disable", "CONFIG_USB"},
		{"--disable", "CONFIG_WLAN"},
		{"--disable", "CONFIG_HID"},
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

var DefaultConfigGroups = []string{"basic", "bpf", "virtio", "minimize"}

func GetConfigGroupNames() []string {
	ret := make([]string, 0, len(ConfigOptGroups))
	for k := range ConfigOptGroups {
		ret = append(ret, k)
	}
	return ret
}

func (cnf *Conf) SaveTo(log *logrus.Logger, dir string, backup bool) error {
	fname := path.Join(dir, ConfigFname)
	confb, err := json.MarshalIndent(cnf, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if backup {
		// rename configuration file if it exists
		if ok, _ := regularFileExists(fname); ok {
			dateStr := time.Now().Format("20060102.150405000000")
			fnameOld := fmt.Sprintf("%s.%s", fname, dateStr)
			err := os.Rename(fname, fnameOld)
			if err != nil {
				log.Infof("failed to rename %s to %s", fname, fnameOld)
			} else {
				log.Infof("renamed %s to %s", fname, fnameOld)
			}
		}
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

func (kc *KernelConf) AddGroupsOpts(gs ...string) error {
	opts, err := confAddGroups(kc.Opts, gs...)
	if err != nil {
		return err
	}
	kc.Opts = opts
	return nil
}

func (c *Conf) AddGroupsCommonOpts(gs ...string) error {
	opts, err := confAddGroups(c.CommonOpts, gs...)
	if err != nil {
		return err
	}
	c.CommonOpts = opts
	return nil
}
