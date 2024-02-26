// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package runner

import (
	"fmt"
	"os"
	"strings"

	"github.com/cilium/little-vm-helper/pkg/arch"
	"github.com/sirupsen/logrus"
)

func BuildQemuArgs(log *logrus.Logger, rcnf *RunConf) ([]string, error) {
	qemuArgs := []string{
		// no need for all the default devices
		"-nodefaults",
		// no need display (-nographics seems a bit slower)
		"-display", "none",
		// don't reboot, just exit
		"-no-reboot",
		// cpus, memory
		"-smp", fmt.Sprintf("%d", rcnf.CPU), "-m", rcnf.Mem,
	}

	// quick-and-dirty kvm detection
	if !rcnf.DisableKVM {
		if f, err := os.OpenFile("/dev/kvm", os.O_RDWR, 0755); err == nil {
			qemuArgs = append(qemuArgs, "-enable-kvm", "-cpu", rcnf.CPUKind)
			f.Close()
		} else {
			log.Info("KVM disabled")
		}
	}

	if rcnf.SerialPort != 0 {
		qemuArgs = append(qemuArgs,
			"-serial",
			fmt.Sprintf("telnet:localhost:%d,server,nowait", rcnf.SerialPort))
	}

	if rcnf.ConsoleLogFile != "" {
		qemuArgs = append(qemuArgs,
			"-serial",
			fmt.Sprintf("file:%s", rcnf.ConsoleLogFile))
	}

	var kernelRoot string
	switch rcnf.RootDev {
	case "hda":
		qemuArgs = append(qemuArgs, "-hda", rcnf.testImageFname())
		kernelRoot = "/dev/sda"
	case "vda":
		qemuArgs = append(qemuArgs, "-drive", fmt.Sprintf("file=%s,if=virtio,index=0,media=disk", rcnf.testImageFname()))
		kernelRoot = "/dev/vda"
	default:
		return nil, fmt.Errorf("invalid root device: %s", rcnf.RootDev)
	}

	if rcnf.KernelFname != "" {
		console, err := arch.Console()
		if err != nil {
			return nil, fmt.Errorf("failed retrieving console name: %w", err)
		}

		appendArgs := []string{
			fmt.Sprintf("root=%s", kernelRoot),
			fmt.Sprintf("console=%s", console),
			"earlyprintk=ttyS0",
			"panic=-1",
		}
		qemuArgs = append(qemuArgs,
			"-kernel", rcnf.KernelFname,
			"-append", fmt.Sprintf("%s", strings.Join(appendArgs, " ")),
		)
	}

	if !rcnf.DisableNetwork {
		qemuArgs = append(qemuArgs, rcnf.ForwardedPorts.QemuArgs()...)
	}

	if !rcnf.Daemonize {
		qemuArgs = append(qemuArgs,
			"-serial", "mon:stdio",
			"-device", "virtio-serial-pci",
		)
	} else {
		qemuArgs = append(qemuArgs, "-daemonize")
	}

	if rcnf.QemuMonitorPort != 0 {
		arg := fmt.Sprintf("tcp:localhost:%d,server,nowait", rcnf.QemuMonitorPort)
		qemuArgs = append(qemuArgs, "-monitor", arg)
	}

	if len(rcnf.HostMount) > 0 {
		qemuArgs = append(qemuArgs,
			"-fsdev", fmt.Sprintf("local,id=host_id,path=%s,security_model=none", rcnf.HostMount),
			"-device", "virtio-9p-pci,fsdev=host_id,mount_tag=host_mount",
		)
	}

	return qemuArgs, nil
}
