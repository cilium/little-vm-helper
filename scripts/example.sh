#!/bin/bash

set -e
set -o pipefail
set -u

LVH="go run ./cmd/lvh"
# NB: uding _data so that this directory is ignored by the go build system
DATADIR="_data"
IMAGESDIR="$DATADIR/images"
KERNELSDIR="$DATADIR/kernels"

mkdir -p $IMAGESDIR
# first create an example configuration
$LVH images example-config > $IMAGESDIR/conf.json
# then, build the images using it
$LVH images build --dir $IMAGESDIR

mkdir -p $KERNELSDIR
$LVH kernels --dir $KERNELSDIR init
# add a kernel by specifying a git URL, and the build it
# (a default configuration is used, but this can be changed)
$LVH kernels --dir $KERNELSDIR add bpf-next git://git.kernel.org/pub/scm/linux/kernel/git/bpf/bpf-next.git --fetch
$LVH kernels --dir $KERNELSDIR build bpf-next
$LVH kernels --dir $KERNELSDIR add 5.18 'git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git#linux-5.18.y' --fetch
$LVH kernels --dir $KERNELSDIR build 5.18
