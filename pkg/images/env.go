package images

import (
	"fmt"
	"os"
	"os/exec"
)

var ignoreEnviromentCheckForTesting = false

func CheckEnvironment() error {
	if ignoreEnviromentCheckForTesting {
		return nil
	}

	for _, cmd := range Binaries {
		_, err := exec.LookPath(cmd)
		if err != nil {
			return fmt.Errorf("required cmd '%s' not found", cmd)
		}
	}

	// libguestfs requires access to KVM
	f, err := os.OpenFile("/dev/kvm", os.O_RDWR, 0755)
	if err != nil {
		return fmt.Errorf("Unable to open /dev/kvm")
	}
	f.Close()

	return nil
}
