package arch

import (
	"fmt"
	"runtime"
)

// Target returns the Linux Makefile target to build the kernel, for historical
// reasons, those are different between architectures.
func Target() (string, error) {
	switch runtime.GOARCH {
	case "amd64":
		return "bzImage", nil
	case "arm64":
		return "Image.gz", nil
	default:
		return "", fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}
}
