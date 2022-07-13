package images

import (
	"context"
	"os/exec"
	"path"

	"github.com/cilium/little-vm-helper/pkg/logcmd"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

// VirtCustomizeStep is a step implemented a set of arguments in virt-customize
//
// NB: we can maybe merge multiple VirtCustomizeStep in a single virt-customize invocation.
// The idea here would be that virt-customize performs the actions in its
// arguments sequentially.
//
// NB: we can probably do the same with guestfish as well
type VirtCustomizeStep struct {
	*StepConf
	Args []string
}

func (s *VirtCustomizeStep) Run(ctx context.Context, b multistep.StateBag) multistep.StepAction {
	imgFname := path.Join(s.imagesDir, s.imgCnf.Name)
	args := []string{"-a", imgFname}
	args = append(args, s.Args...)
	cmd := exec.CommandContext(ctx, "virt-customize", args...)
	err := logcmd.RunAndLogCommand(cmd, s.log)
	if err != nil {
		s.log.WithField("image", s.imgCnf.Name).WithError(err).Error("error executing command")
		b.Put("err", err)
		return multistep.ActionHalt
	}
	return multistep.ActionContinue
}

func (s *VirtCustomizeStep) Cleanup(b multistep.StateBag) {
}
