##  little-vm-helper

little-vm-helper (lvh) is a VM management tool, aimed for testing and development of features that
depend on the kernel, such as BPF. It is used in [cilium](https://github.com/cilium/cilium),
[tetragon](https://github.com/cilium/tetragon), and [pwru](https://github.com/cilium/pwru). It can also be used for kernel development. It is not
meant, and should not be used for running production VMs. Fast booting and image building, as well
as being storage efficient are the main goals.

It uses [qemu](https://www.qemu.org/) and [libguestfs tools](https://libguestfs.org/). See [dependencies](#what-are-the-dependencies-of-lvh).

Configurations for specific images used in the Cilium project can be found in:
https://github.com/cilium/little-vm-helper-images.

## Usage

For an example script, see [scripts/example.sh](scripts/example.sh).

LVH can be used to:
 * build root images for VMs
 * build kernels
 * boot VMs using above

### Root images

Build example images:
```bash
$ mkdir _data
$ go run ./cmd/lvh images example-config > _data/images.json
$ go run ./cmd/lvh images build --dir _data # this may require sudo as relies on /dev/kvm
```

The first command will create a configuration file:
```jsonc
jq . < _data/images.json
[
  {
    "name": "base",
    "packages": [
      "less",
      "vim",
      "sudo",
      "openssh-server",
      "curl"
    ],
    "actions": [
      {
        "comment": "disable password for root",
        "op": {
          "Cmd": "passwd -d root"
        },
        "type": "run-command"
      }
    ]
  },
  {
    "name": "k8s",
    "parent": "base",
    "image_size": "20G",
    "packages": [
      "docker.io"
    ]
  }
]
```

The configuration file includes:
 * a set of packages for the image
 * an optional parent image
 * a set of actions to be performed after the installation of the packets. There are multiple
   actions supported, see [pkg/images/actions.go](pkg/images/actions.go).

Once the `build-images` command completes, the two images described in the configuration file will
be present in the images directory. ote that the images are stored as sparse files so they take less
space:

```bash
$ ls -sh1 _data/images/*.img
856M _data/images/base.img
1.7G _data/images/k8s.img
```

### Kernels

```bash
$ mkdir -p _data/kernels
$ go run ./cmd/lvh kernels --dir _data init
$ go run ./cmd/lvh kernels --dir _data add bpf-next git://git.kernel.org/pub/scm/linux/kernel/git/bpf/bpf-next.git --fetch
$ go run ./cmd/lvh kernels --dir _data build bpf-next
```

Please note, to cross-build for a different architecture, you can use the
`--arch=arm64` or `--arch=amd64` flag.

The configuration file keeps the url for a kernel, together with its configuration options:
```jsonc
$ jq . < _data/kernel.json
{
  "kernels": [
    {
      "name": "bpf-next",
      "url": "git://git.kernel.org/pub/scm/linux/kernel/git/bpf/bpf-next.git"
    }
  ],
  "common_opts": [
    [
      "--enable",
      "CONFIG_LOCALVERSION_AUTO"
    ],
     ... more options ...
  ]
}
```

There are options that are applied to all kernels (`common_opts`) as well as
kernel-specific options.

The kernels are kept in [worktrees](https://git-scm.com/docs/git-worktree). Specifically, there is a
git bare directory (`git`) that holds all the objects, and one worktree per kernel. This allows
efficient fetching and, also, having each kernel on its own separate directory.

For example:
```bash
$ ls -1 _data/kernels
5.18/
bpf-next/
git/
```

Currently, kernels are built using the `bzImage` for x86\_64 or `Image.gz` for
arm64, and `tar-pkg` targets (see [pkg/kernels/conf.go](pkg/kernels/conf.go)).

### Booting images

You can use the `run` subcommand to start images.

For example:
```bash
go run ./cmd/lvh run --image _data/images/base.qcow2 --kernel _data/kernels/bpf-next/arch/x86_64/boot/bzImage
```

Or, to with the kernel installed in the image:
```bash
go run ./cmd/lvh run --image _data/images/base.qcow2
```

OCI images are also supported:
```bash
go run ./cmd/lvh run --image quay.io/lvh-images/root-images:main
```

**Note**: Building images and kernels is only supported on Linux. On the other hand, images and kernels already build on Linux can be booted in MacOS (both x86 and Arm). The only requirement is ```qemu-system-x86_64```. As MacOS does not support KVM, the commands to boot images are:
```bash
go run ./cmd/lvh run --image _data/images/base.qcow2 --qemu-disable-kvm
```

## FAQ

### Why not use packer to build images?

Existing packer builders
(e.g,.https://github.com/cilium/packer-ci-build/blob/710ad61e7d5b0b6872770729a30bcdade2ee1acb/cilium-ubuntu.json#L19,
https://www.packer.io/plugins/builders/qemu) are meant to manage VMs with
longer lifetimes than a single use, and use facilities that introduce unnecessary overhead for our use-case.

Also, packer does not seem to have a way to provision images without booting a
machine. There is an outdated chroot package
https://github.com/summerwind/packer-builder-qemu-chroot, and cloud chroot builders
(e.g., https://www.packer.io/plugins/builders/amazon/chroot that uses https://github.com/hashicorp/packer-plugin-sdk/tree/main/chroot).

That being said, if we need packer functionality we can create a packer plugin
(https://www.packer.io/docs/plugins/creation#developing-plugins).

### Why not use vagrant (or libvirt-based tools)?

These tools also target production VMs with lifetime streching beyond a single
use. As a result, they introduce overhead in booting time, provisioning time,
and storage.

### What are the dependencies of LVH?

On debian distribution, here is a list of packages needed for LVH to work.

| Action                         | Debian packages                                                                                                        |
| --------                       | -------                                                                                                                |
| Building images                | `qemu-kvm mmdebstrap debian-archive-keyring libguestfs-tools`                                                          |
| Building the Linux kernel      | `libncurses-dev gawk flex bison openssl libssl-dev dkms libelf-dev libudev-dev libpci-dev libiberty-dev autoconf llvm` |
| Cross-compile arm64 on x86\_64 | `gcc-aarch64-linux-gnu`                                                                                                |
| Cross-compile x86\_64 on arm64 | `gcc-x86-64-linux-gnu`                                                                                                 |


### TODO

 - [ ] development workflow for MacOS X
 - [ ] images: configuration option for using different deb distros (hardcoded to sid now)
 - [x] images: build tetragon images
   - [x] unit tests
   - [x] e2e tests (kind)
 - [x] images: docker image with required binaries (libguestfs, mmdebstrap, etc.) to run the tool
        - [x]  is that possible? libguestfs needs to boot a mini-VM
 - [x] kernels: add support for building kernels
 - [x] runner: qemu runner wrapper
 - [x] images bootable VMs: running qemu with --kernel is convinient for development. If we want to store images externally (e.g., AWS), it might make sense to support bootable VMs.
 - [ ] improve boot time: minimal init, use qemu microvm (https://qemu.readthedocs.io/en/latest/system/i386/microvm.html, https://mergeboard.com/blog/2-qemu-microvm-docker/)
 - [ ] images: on a failed run, save everything in a image-failed-$(date) directory
 - [ ] use `guestfish --listen` (see
   https://github.com/libbpf/ci/blob/cbb3b92facbad705bbb619b496d0debb4b3d806f/prepare-rootfs/run.sh#L345)

## Notes

- earlier attempt: https://github.com/kkourt/kvm-dev-scripts
