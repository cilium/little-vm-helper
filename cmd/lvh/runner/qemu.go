package runner

import (
	"fmt"
	"os"
	"strings"

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

	qemuArgs = append(qemuArgs,
		"-hda", rcnf.testImageFname(),
	)

	if rcnf.KernelFname != "" {
		appendArgs := []string{
			"root=/dev/sda",
			"console=ttyS0",
			"earlyprintk=ttyS0",
			"panic=-1",
		}
		qemuArgs = append(qemuArgs,
			"-kernel", rcnf.KernelFname,
			"-append", fmt.Sprintf("%s", strings.Join(appendArgs, " ")),
		)
	}

	if !rcnf.DisableNetwork {
		netdev := "user,id=user.0"
		for _, fwd := range rcnf.ForwardedPorts {
			netdev = fmt.Sprintf("%s,hostfwd=%s::%d-:%d", netdev, fwd.Protocol, fwd.HostPort, fwd.VMPort)
		}

		qemuArgs = append(qemuArgs,
			"-netdev", netdev,
			"-device", "virtio-net-pci,netdev=user.0",
		)
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

	qemuArgs = append(qemuArgs,
		"-fsdev", fmt.Sprintf("local,id=host_id,path=%s,security_model=none", rcnf.HostMount),
		"-device", "virtio-9p-pci,fsdev=host_id,mount_tag=host_mount",
	)

	return qemuArgs, nil
}
