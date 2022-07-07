##  little-vm-helper

little-vm-helper is a VM management toolset, targeting testing and development, especially with
features that are dependent the kernel, such as BPF. It should not be used for running production
VMs. That is, it is only intended for VMs that have a "fire-once" lifetime. Booting fast, building
images fast, and being storage efficient are the main goals.

It uses [qemu](https://www.qemu.org/) and [libguestfs tools](https://libguestfs.org/).

### TODO
 - [ ] development workflow for MacOS X
 - [ ] images: configuration option for using different deb distros (hardcoded to sid now)
 - [ ] images: build tetragon images
   - [ ] unit tests
   - [ ] e2e tests (kind)
 - [ ] images: docker image with required binaries (libguestfs, mmdebstrap, etc.) to run the tool
        - [ ]  is that possible? libguestfs needs to boot a mini-VM
 - [x] kernels: add suport for buidling kernels
 - [ ] runner: qemu runner wrapper
 - [x] images bootable VMs: running qemu with --kernel is convinient for development. If we want to store images externally (e.g., AWS), it might make sense to support bootable VMs.
 - [ ] improve boot time: minimal init, use qemu microvm (https://qemu.readthedocs.io/en/latest/system/i386/microvm.html, https://mergeboard.com/blog/2-qemu-microvm-docker/)
 - [ ] images: on a failed run, save everything in a image-failed-$(date) directory

## FAQ

### Why not use packer to build images?

Existing packer builders
(e.g,.https://github.com/cilium/packer-ci-build/blob/710ad61e7d5b0b6872770729a30bcdade2ee1acb/cilium-ubuntu.json#L19,
https://www.packer.io/plugins/builders/qemu) are meant to manage VMs with
longer lifetimes than a single use, and use facilities that introduce unnecessary overhead for our use-case.

Also, packer does not seem to have a way to provision images without booting a
machine. There is an outdated chroot package
https://github.com/summerwind/packer-builder-qemu-chroot, and cloud chroot builders
(e.g., https://www.packer.io/plugins/builders/amazon/chroot).

That being said, if we need packer functionality we can create a packer plugin
(https://www.packer.io/docs/plugins/creation#developing-plugins).

### Why not use vagrant (or libvirt-based tools)?

These tools also target production VMs with lifetime streching beyond a single
use. As a result, they introduce overhead in booting time, provisioning time,
and storage.

## Example:

For an example script, see [scripts/example.sh](scripts/example.sh)

### Images

Build example images:
```
$ mkdir -p _data/images
$ go run cmd/lvh images example-config > _data/images/conf.json
% go run cmd/lvh images build --dir _data/images
```

The first command will create a configuration file:
```
jq . < _data/images/conf.json 
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
    "packages": [
      "docker.io"
    ]
  }
]
```

The configuration file includes:
 * a set of packages for the image
 * an optioanl parent image
 * a set of actions to be performed after the installation of the packets. For now, only a
   "run-command" action  is supported, where a command is executed inside the image (using
   `virt-customize`) but other ones can be easily added (see
   [pkg/images/actions.go](pkg/images/actions.go))

After the `build-images` command completes, there will be two images in the images directory. Note
that the images are stored as sparse files so they take less space:

```
$ ls -sh1 _data/images/*.img
341M _data/images/base.img
1.2G _data/images/k8s.img
```


### Kernels

```
$ mkdir -p _data/kernels
$ go run cmd/lvh kernels --dir _data/kernels init
$ go run cmd/lvh kernels --dir _data/kernels add bpf-next git://git.kernel.org/pub/scm/linux/kernel/git/bpf/bpf-next.git --fetch
$ go run cmd/lvh kernels --dir _data/kernels build bpf-next
```

The configuration file keeps the url for a kernel, togther with its configuration options:
```
$ jq . < _data/kernels/conf.json  
{
  "kernels": [
    {
      "name": "bpf-next",
      "url": "git://git.kernel.org/pub/scm/linux/kernel/git/bpf/bpf-next.git",
      "opts": [
        [
          "--enable",
          "CONFIG_LOCALVERSION_AUTO"
        ],
	... more options ...
        [
          "--disable",
          "CONFIG_SOUND"
        ]
      ]
    },
  ]
}
```

The kernels are kept in [worktrees](https://git-scm.com/docs/git-worktree). Specifically, There is a
git bare directory (`git`) that holds all the objects, and one worktree per kernel. This allows
efficient fetching and, also, having each kernel on its own seperate directory.

For example:
```
$ ls -1 _data/kernels 
5.18/
bpf-next/
conf.json
git/
```

Currently, kernels are built using the `bzImage` and `dir-pkg` targets (see [pkg/kernels/conf.go](pkg/kernels/conf.go)).

### Booting images

The goal is to have some wrappers for running qemu, but until then, here is an example:

```
qemu-system-x86_64 -enable-kvm -m 4G -hda _data/images/base.img -nographic -kernel _data/kernels/bpf-next/arch/x86_64/boot/bzImage  -append "root=/dev/sda console=ttyS0"
```

Or, even:

```
qemu-system-x86_64 -enable-kvm -m 4G -hda _data/images/base.img -nographic
```



## Notes

Previous attempt of doing this was: https://github.com/kkourt/kvm-dev-scripts
