##  little-vm-helper

little-vm-helper is a VM management toolset, targeting testing and development, especially with
features that are dependent the kernel, such as BPF. It should not be used for running production
VMs. That is, it is only intended for VMs that have a "fire-once" lifetime. Booting fast, building
images fast, and being storage efficient are the main goals.

It uses [qemu](https://www.qemu.org/) and [libguestfs tools](https://libguestfs.org/).

### TODO
 * images: configuration option for using different deb distros (hardcoded to sid now)
 * images: build tetragon images
     * unit tests
     * e2e tests (kind)
 * images: docker image with required binaries (libguestfs, mmdebstrap, etc.) to run the tool
    * is that possible? libguestfs needs to boot a mini-VM
 * kernels: add suport for buidling kernels
 * runner: qemu runner wrapper
 * images bootable VMs: running qemu with --kernel is convinient for
   development. If we want to store images externally (e.g., AWS), it might
   make sense to support bootable VMs.
 * improve boot time: minimize init, use qemu microvm
   (https://qemu.readthedocs.io/en/latest/system/i386/microvm.html,
   https://mergeboard.com/blog/2-qemu-microvm-docker/)
 * images: on a failed run, save everything in a image-failed-$(date) directory
 * development workflow for MacOS X

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

### Images


```
$ mkdir images
$ go run cmd/little-vm-helper example-config > images/conf.json
$ go run cmd/little-vm-helper build-images --dir ./images
INFO[0000] starting to build all images                  queue=base
INFO[0000] starting command                              args="[mmdebstrap sid --include less,vim,sudo,openssh-server,curl images/base.tar]" path=/usr/bin/mmdebstrap
WARN[0000] stderr> I: automatically chosen mode: unshare
WARN[0000] stderr> I: chroot architecture amd64 is equal to the host's architecture
WARN[0000] stderr> I: Reading sources.list from standard input...
WARN[0000] stderr> I: using /tmp/mmdebstrap.Iixrf4Anko as tempdir
WARN[0000] stderr> I: running apt-get update...
WARN[0004] stderr> I: downloading packages with apt...
WARN[0006] stderr> I: extracting archives...
WARN[0008] stderr> I: installing packages...
WARN[0016] stderr> I: downloading apt...
WARN[0017] stderr> I: installing apt...
WARN[0020] stderr> I: installing remaining packages inside the chroot...
WARN[0038] stderr> I: cleaning package lists and apt cache...
WARN[0038] stderr> I: creating tarball...
WARN[0039] stderr> I: done
WARN[0040] stderr> I: removing tempdir /tmp/mmdebstrap.Iixrf4Anko...
INFO[0040] starting command                              args="[guestfish -N images/base.img=disk:8G -- mkfs ext4 /dev/sda : mount /dev/sda / : tar-in images/base.tar /]" path=/usr/bin/guestfish
INFO[0043] image built succesfully                       image=base queue=k8s result="{Error:<nil> CachedImageUsed:false CachedImageDeleted:}"
INFO[0043] starting command                              args="[cp --sparse always -n images/base.img images/k8s.img]" path=/usr/bin/cp
INFO[0043] starting command                              args="[virt-customize -a images/k8s.img --install docker.io]" path=/usr/bin/virt-customize
INFO[0043] stdout> [   0.0] Examining the guest ...
INFO[0046] stdout> [   2.6] Setting a random seed
INFO[0046] stdout> virt-customize: warning: random seed could not be set for this type of
INFO[0046] stdout> guest
INFO[0046] stdout> [   2.6] Installing packages: docker.io
INFO[0079] stdout> [  35.6] Finishing off
INFO[0079] image built succesfully                       image=k8s queue= result="{Error:<nil> CachedImageUsed:false CachedImageDeleted:}"
INFO[0079] images built succesfully                      time-elapsed=1m19.680405705s
image:base       cachedImageUsed:false cachedImageDeleted:
image:k8s        cachedImageUsed:false cachedImageDeleted:
```

## Notes

Based on: https://github.com/kkourt/kvm-dev-scripts
