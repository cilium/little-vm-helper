#!/usr/bin/env bash

CTR=""
IMAGES_DIR="${IMAGES_DIR:-.}"

cleanup() {
    docker rm $CTR
}

# $1 - image
main() {
    if [ $# -ne 1 ]; then
        >&2 echo "usage: $0 <quay.io/lvh-images/kind:TAG>"
        exit 1
    fi

    mkdir -p $IMAGES_DIR
    CTR=$(docker create $1)
    trap cleanup EXIT
    docker cp $CTR:/data/images $IMAGES_DIR >/dev/null

    files=($(find $IMAGES_DIR/images/*.zst -type f))
    for f in ${files[@]}; do
        zstd --decompress $f
    done
}

main "$@"
