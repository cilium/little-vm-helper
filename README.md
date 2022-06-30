##  little-vm-helper

little-vm-helper is a VM management toolset, targeting testing and development. It should
not be used for running production VMs. That is, it is only intended for VMs
that have a "fire-once" lifetime. Booting fast, building images fast, and being
storage efficient are the main goals.

It uses [qemu](https://www.qemu.org/) and [libguestfs tools](https://libguestfs.org/).

### TODO
 * images: configuration option for using different deb distros (hardcoded to sid now)
 * images: build tetragon images
     * unit tests
     * e2e tests (kind)
 * images: docker image with required binaries (libguestfs, mmdebstrap, etc.) to run the tool
 * kernels: add suport for buidling kernels
 * runner: qemu runner wrapper
 * images bootable VMs: running qemu with --kernel is convinient for
   development. If we want to store images externally (e.g., AWS), it might
   make sense to support bootable VMs.
 * improve boot time: minimize init, use qemu microvm
   (https://qemu.readthedocs.io/en/latest/system/i386/microvm.html,
   https://mergeboard.com/blog/2-qemu-microvm-docker/)

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


```
$ mkdir images                                    
$ go run main.go example-config > images/conf.json                                                    
$ go run main.go build-images --dir ./images                                                                                                                                             
INFO[0000] starting to build all images                  queue="base,k8s"
INFO[0000] starting command                              args="[mmdebstrap sid --include less,vim,sudo,openssh-server,curl images/base.tar]" path=/usr/bin/mmdebstrap
WARN[0000] stderr> I: automatically chosen mode: unshare 
WARN[0000] stderr> I: chroot architecture amd64 is equal to the host's architecture 
WARN[0000] stderr> I: Reading sources.list from standard input... 
WARN[0000] stderr> I: using /tmp/mmdebstrap.lEkVaNpEDh as tempdir 
WARN[0000] stderr> I: running apt-get update...         
WARN[0004] stderr> I: downloading packages with apt...  
WARN[0005] stderr> I: extracting archives...            
WARN[0007] stderr> I: installing packages...            
WARN[0015] stderr> I: downloading apt...                
WARN[0016] stderr> I: installing apt...                 
WARN[0019] stderr> I: installing remaining packages inside the chroot... 
WARN[0037] stderr> I: cleaning package lists and apt cache... 
WARN[0037] stderr> I: creating tarball...               
WARN[0038] stderr> I: done                              
WARN[0038] stderr> I: removing tempdir /tmp/mmdebstrap.lEkVaNpEDh... 
INFO[0038] starting command                              args="[guestfish -N images/base.img=disk:8G -- mkfs ext4 /dev/sda : mount /dev/sda / : tar-in images/base.tar /]" path=/usr/bin/guestfish
INFO[0041] image built succesfully                       image=base queue=k8s result="{Error:<nil> CachedImageUsed:false CachedImageDeleted:}"
INFO[0041] starting command                              args="[mmdebstrap sid --include docker.io images/k8s.tar]" path=/usr/bin/mmdebstrap
WARN[0042] stderr> I: automatically chosen mode: unshare 
WARN[0042] stderr> I: chroot architecture amd64 is equal to the host's architecture 
WARN[0042] stderr> I: Reading sources.list from standard input... 
WARN[0042] stderr> I: using /tmp/mmdebstrap.JmbUcgB42f as tempdir 
WARN[0042] stderr> I: running apt-get update...         
WARN[0045] stderr> I: downloading packages with apt...  
WARN[0047] stderr> I: extracting archives...            
WARN[0049] stderr> I: installing packages...            
WARN[0057] stderr> I: downloading apt...                
WARN[0058] stderr> I: installing apt...                 
WARN[0060] stderr> I: installing remaining packages inside the chroot... 
WARN[0085] stderr> I: cleaning package lists and apt cache... 
WARN[0086] stderr> I: creating tarball...               
WARN[0087] stderr> I: done                              
WARN[0087] stderr> I: removing tempdir /tmp/mmdebstrap.JmbUcgB42f... 
INFO[0088] starting command                              args="[guestfish -N images/k8s.img=disk:8G -- mkfs ext4 /dev/sda : mount /dev/sda / : tar-in images/k8s.tar /]" path=/usr/bin/guestfish
INFO[0092] image built succesfully                       image=k8s queue= result="{Error:<nil> CachedImageUsed:false CachedImageDeleted:}"
INFO[0092] images built succesfully                      time-elapsed=1m32.100372026s
image:base       cachedImageUsed:false cachedImageDeleted:
image:k8s        cachedImageUsed:false cachedImageDeleted:
```

## Notes

Previous attempt: https://github.com/kkourt/kvm-dev-scripts
