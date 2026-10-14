// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package runner

import (
	"testing"

	"github.com/cilium/little-vm-helper/pkg/slogger"
	"github.com/stretchr/testify/assert"
)

func TestBuildQemuArgsKernelModules(t *testing.T) {
	rcnf := RunConf{
		Image:            "image.qcow2",
		RootDev:          "vda",
		CPU:              1,
		Mem:              "1G",
		Daemonize:        true,
		KernelModulesDir: "/host/modules",
	}

	args, err := BuildQemuArgs(slogger.NewDiscard(), &rcnf)
	assert.Nil(t, err)
	assert.Contains(t, args, "-fsdev")
	assert.Contains(t, args, "local,id=kernel_modules_id,path=/host/modules,security_model=none,readonly=on")
	assert.Contains(t, args, "-device")
	assert.Contains(t, args, "virtio-9p-pci,fsdev=kernel_modules_id,mount_tag=kernel_modules")
}

func TestBuildQemuArgsNoKernelModules(t *testing.T) {
	rcnf := RunConf{
		Image:     "image.qcow2",
		RootDev:   "vda",
		CPU:       1,
		Mem:       "1G",
		Daemonize: true,
	}

	args, err := BuildQemuArgs(slogger.NewDiscard(), &rcnf)
	assert.Nil(t, err)
	assert.NotContains(t, args, "virtio-9p-pci,fsdev=kernel_modules_id,mount_tag=kernel_modules")
}
