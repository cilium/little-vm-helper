package images

import (
	"fmt"
	"os/exec"
)

func (ib *ImageBuilder) CheckEnvironment() error {
	for _, cmd := range []string{Debootstrap} {
		_, err := exec.LookPath(cmd)
		if err != nil {
			return fmt.Errorf("required cmd '%s' not found", cmd)
		}
	}
	return nil
}
