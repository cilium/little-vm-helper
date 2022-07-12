#!/bin/bash

set -e
set -o pipefail
set -u
set -x

LVH="go run ./cmd/lvh"
# NB: using _data so that this directory is ignored by the go build system
DATADIR="_data"

mkdir -p $DATADIR
# first create an example configuration
$LVH images example-config > $DATADIR/images.json
# then, build the images using it
$LVH images build --dir $DATADIR

$LVH kernels --dir $DATADIR init
# add a kernel by specifying a git URL, and the build it
# (a default configuration is used, but this can be changed)
$LVH kernels --dir $DATADIR add bpf-next git://git.kernel.org/pub/scm/linux/kernel/git/bpf/bpf-next.git --fetch
$LVH kernels --dir $DATADIR build bpf-next
$LVH kernels --dir $DATADIR add 5.18 'git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git#linux-5.18.y' --fetch
$LVH kernels --dir $DATADIR build 5.18
