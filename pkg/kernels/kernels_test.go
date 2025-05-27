// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package kernels

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type testLogger struct {
	*testing.T
}

func (tl testLogger) Write(p []byte) (n int, err error) {
	tl.Log((string)(p))
	return len(p), nil
}

var (
	testKconfs = []KernelConf{
		{
			Name: "bpf-next",
			URL:  "git://git.kernel.org/pub/scm/linux/kernel/git/bpf/bpf-next.git",
			Opts: []ConfigOption{
				{"--enable", "CONFIG_DEBUG_INFO"},
				{"--disable", "CONFIG_DEBUG_KERNEL"},
			},
		}, {
			Name: "5.18",
			URL:  "git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git#linux-5.18.y",
			Opts: []ConfigOption{
				{"--enable CONFIG_BPF"},
				{"--enable CONFIG_BPF_SYSCALL"},
			},
		},
	}
)

func TestDir(t *testing.T) {
	xlog := logrus.New()
	xlog.SetOutput(testLogger{t})
	configs := []*Conf{
		nil,
		{
			Kernels:    testKconfs,
			CommonOpts: []ConfigOption{{"--disable", "CONFIG_WERROR"}},
		},
	}

	for _, conf := range configs {
		// NB: anonymous function so that os.RemoveAll() is called in all iterations
		func() {
			dir, err := ioutil.TempDir("", "test_kernel")
			assert.Nil(t, err)
			defer os.RemoveAll(dir)
			err = InitDir(xlog, dir, conf, InitDirFlags{Force: false, BackupConf: false})
			assert.Nil(t, err)

			if conf == nil {
				conf = &Conf{
					Kernels: make([]KernelConf, 0),
				}
			}

			kd, err := LoadDir(dir)
			assert.Nil(t, err)
			assert.Equal(t, &kd.Conf, conf)
		}()
	}
}
