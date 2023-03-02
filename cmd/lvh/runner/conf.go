// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package runner

import (
	"github.com/cilium/little-vm-helper/pkg/runner"
	"github.com/sirupsen/logrus"
)

type RunConf struct {
	// Image filename
	Image string
	// kernel filename to boot with. (if empty no -kernel option will be passed to qemu)
	KernelFname string
	// Do not run the qemu command, just print it
	QemuPrint bool
	// Do not use KVM acceleration, even if /dev/kvm exists
	DisableKVM bool
	// Daemonize QEMU after initializing
	Daemonize bool

	// Disable the network connection to the VM
	DisableNetwork bool
	ForwardedPorts runner.PortForwards

	Logger *logrus.Logger

	HostMount string

	SerialPort int

	CPU int
	Mem string
	// Kind of CPU to use (e.g. host or kvm64)
	CPUKind string

	QemuMonitorPort int
}

func (rc *RunConf) testImageFname() string {
	return rc.Image
}
