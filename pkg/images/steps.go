package images

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/cilium/little-vm-helper/pkg/logcmd"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/sirupsen/logrus"
)

var (
	// DelImageIfExists: if set to true, image will be deleted at Cleanup() by the CreateImage step
	DelImageIfExists = "DelImageIfExist"
)

// Approach for creating images:
// - Base (root) images are build using mmdebstrap and copying files using guestfish.
// - Non-root images are build using virt-customize, by copying the parent.
// - All images use the raw format (not qcow2)
// - Images are read-only. Users can use them to create other images (by copying or via qcow2)
//
// Alternative options I considred and may be useful for future reference:
//  - using qemu-nbd+chroot, would probably be a faster way to do this, but it requires root.
//  - using either debootstrap, or multistrap (with fakeroot and fakechroot) instead of mmdebstrap.
//    The latter seems faster, so I thought I'd use it. If something breaks, we can always go another
//    route.
//  - using the go bindings for libguestfs (https://libguestfs.org/guestfs-golang.3.html). Using the
//    CLI seemed simpler.
//  - having bootable images. I don't think we need this since we can specify --kernel and friends
//    in qemu.
//  - having the images in qcow2 so that we save some space. I think the sparsity of the files is
//    enough, so decided to keep things simple. Note that we can use virt-sparsify if we want to (e.g.,
//    when downloading images).

// CreateImage is a step for creating an image. It's cleanup will delete it if DelImageIfExists is set.
type CreateImage struct {
	imageDir string
	imgCnf   *ImageConf
	log      *logrus.Logger
}

func (s *CreateImage) makeRootImage(ctx context.Context) error {
	tarFname := path.Join(s.imageDir, fmt.Sprintf("%s.tar", s.imgCnf.Name))
	cmd := exec.CommandContext(ctx, Mmdebstrap,
		"sid",
		"--include", strings.Join(s.imgCnf.Packages, ","),
		tarFname,
	)
	err := logcmd.RunAndLogCommand(cmd, s.log)
	if err != nil {
		return err
	}
	defer func() {
		err := os.Remove(tarFname)
		if err != nil {
			s.log.WithError(err).Info("failed to remove tarfile")
		}
	}()

	imgFname := path.Join(s.imageDir, fmt.Sprintf("%s.img", s.imgCnf.Name))
	// example: guestfish -N foo.img=disk:8G -- mkfs ext4 /dev/sda : mount /dev/sda / : tar-in /tmp/foo.tar /
	cmd = exec.CommandContext(ctx, GuestFish,
		"-N", fmt.Sprintf("%s=disk:%s", imgFname, DefaultImageSize),
		"--",
		"mkfs", "ext4", "/dev/sda",
		":",
		"mount", "/dev/sda", "/",
		":",
		"tar-in", tarFname, "/",
	)
	err = logcmd.RunAndLogCommand(cmd, s.log)
	if err != nil {
		os.Remove(imgFname)
	}
	return err
}

func (s *CreateImage) makeDerivedImage(ctx context.Context) error {
	parFname := path.Join(s.imageDir, fmt.Sprintf("%s.img", s.imgCnf.Parent))
	imgFname := path.Join(s.imageDir, fmt.Sprintf("%s.img", s.imgCnf.Name))

	// NB: cp has detection for sparse files, so just use that for now
	// -n: don't override.
	cmd := exec.CommandContext(ctx, "cp", "--sparse", "always", "-n", parFname, imgFname)
	err := logcmd.RunAndLogCommand(cmd, s.log)
	if err != nil {
		return err
	}

	cmd = exec.CommandContext(ctx, "virt-customize",
		"-a", imgFname,
		"--install", strings.Join(s.imgCnf.Packages, ","),
	)
	err = logcmd.RunAndLogCommand(cmd, s.log)
	if err != nil {
		os.Remove(imgFname)
		return err
	}

	return nil
}

func (s *CreateImage) Run(ctx context.Context, b multistep.StateBag) multistep.StepAction {

	var err error
	if s.imgCnf.Parent == "" {
		err = s.makeRootImage(ctx)
	} else {
		err = s.makeDerivedImage(ctx)
	}

	if err != nil {
		s.log.WithField("image", s.imgCnf.Name).WithError(err).Error("error buiding image")
		b.Put("err", err)
		return multistep.ActionHalt
	}
	return multistep.ActionContinue

}

func (s *CreateImage) Cleanup(b multistep.StateBag) {
}
