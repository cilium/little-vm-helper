// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package kernels

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigGroups(t *testing.T) {
	for _, g := range DefaultConfigGroups {
		if _, ok := ConfigOptGroups[g]; !ok {
			t.Fatalf("default config group '%s' does not exist", g)
		}
	}
}

func TestGetOptions(t *testing.T) {
	conf := Conf{
		CommonOpts: []ConfigOption{{"--disable", "CONFIG_WERROR"}},
		Kernels: []KernelConf{{
			Name: "",
			URL:  "",
			Opts: []ConfigOption{
				{"--enable CONFIG_BPF"},
				{"--enable CONFIG_BPF_SYSCALL"},
			},
		}},
	}

	opts := conf.getOptions(&conf.Kernels[0])
	assert.Equal(t, opts, []ConfigOption{
		{"--disable", "CONFIG_WERROR"},
		{"--enable CONFIG_BPF"},
		{"--enable CONFIG_BPF_SYSCALL"},
	})
}
