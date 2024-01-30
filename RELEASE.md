# little-vm-helper Release Process

## Set `RELEASE` environment variable

Set `RELEASE` environment variable to the new version. For example:

    export RELEASE=v0.0.16

## Open a pull request

Bump up the default version in action.yaml and open a pull request against `main` branch:

    alias yq='docker run --rm -v "${PWD}":/workdir --user "$(id -u):$(id -g)" mikefarah/yq:4.40.5'
    yq ".inputs.lvh-version.default = \"$RELEASE\"" -i action.yaml
    git checkout -b pr/prepare-$RELEASE
    git add action.yaml
    git commit -s -m "Prepare for $RELEASE release"
    git push origin HEAD

Wait for the PR to be reviewed and merged.

## Tag a release

Checkout `main` branch:

    git checkout main
    git pull origin main

Set the commit you want to tag. Usually this is the most recent commit on `main`, i.e.

    export COMMIT_SHA=$(git rev-parse origin/main)

Then tag and push the release:

    git tag -a $RELEASE -m "$RELEASE release" $COMMIT_SHA && git push origin $RELEASE
