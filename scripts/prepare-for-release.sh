#!/bin/bash

set -exo pipefail

usage() {
	echo "$0 [-n] [release]"
	echo "-n: dry run (do not push)"
}


while getopts "n" opt; do
	case ${opt} in
		n)
			DRY_RUN=1
			;;
	esac
done
shift $((OPTIND-1))

if [ -z "${1+x}" ]; then
	git fetch --tags
	lastver=$(git tag | grep '^v' | sort -V | tail -1)
	v=$(echo $lastver | perl -n -e 'if (/v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$/) {print "v","$1.$2.",$3+1,"\n" }')
	if [ -z "$v" ]; then
		set +x
		echo "cannot suggest a new version from $lastver"
		usage
		exit 0
	fi
	set +x
	echo "latest version: $lastver, how about releasing $v?"
	echo "press <enter> to accept or <Ctrl-C> and use an explicit version as the first argument."
	read
	set -x
else
	v="${1}"
fi

set -u
echo "Releasing $v"
RELEASE=$v

alias yq='docker run --rm -v "${PWD}":/workdir --user "$(id -u):$(id -g)" mikefarah/yq:4.40.5'
yq ".inputs.lvh-version.default = \"$RELEASE\"" -i action.yaml
git checkout -b pr/prepare-$RELEASE
git add action.yaml
git commit -s -m "Prepare for $RELEASE release"

if [ -z "${DRY_RUN+x}" ]; then
	git push origin HEAD
else
	echo "dry run, skipping push"
fi
