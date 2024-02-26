package arch

import (
	"fmt"
	"runtime"
)

var ErrUnsupportedArch = fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)

// Target returns the Linux Makefile target to build the kernel, for historical
// reasons, those are different between architectures.
func Target() (string, error) {
	switch runtime.GOARCH {
	case "amd64":
		return "bzImage", nil
	case "arm64":
		return "Image.gz", nil
	default:
		return "", ErrUnsupportedArch
	}
}

func QemuBinary() (string, error) {
	switch runtime.GOARCH {
	case "amd64":
		return "qemu-system-x86_64", nil
	case "arm64":
		return "qemu-system-aarch64", nil
	default:
		return "", ErrUnsupportedArch
	}
}

// Console returns the name of the device for the first serial port.
func Console() (string, error) {
	switch runtime.GOARCH {
	case "amd64":
		return "ttyS0", nil
	case "arm64":
		return "ttyAMA0", nil
	default:
		return "", ErrUnsupportedArch
	}
}

// AppendArchSpecificQemuArgs appends Qemu arguments to the input that are
// specific to the architecture lvh is running on. For example on ARM64, Qemu
// needs some precision on the -machine option to start.
func AppendArchSpecificQemuArgs(qemuArgs []string) []string {
	switch runtime.GOARCH {
	case "arm64":
		return append(qemuArgs, "-machine", "virt", "-cpu", "max")
	default:
		return qemuArgs
	}
}
