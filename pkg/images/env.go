package images

import (
	"fmt"
	"os/exec"
)

func (ib *Builder) CheckEnvironment() error {
	for _, cmd := range Binaries {
		_, err := exec.LookPath(cmd)
		if err != nil {
			return fmt.Errorf("required cmd '%s' not found", cmd)
		}
	}
	return nil
}
