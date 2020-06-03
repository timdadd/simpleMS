#!/usr/bin/env bash

set -euo pipefail
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo $SCRIPTDIR

log() { echo "$1" >&2; }

TAG="${TAG:?TAG env variable must be specified}"
REPO_PREFIX="${REPO_PREFIX:?REPO_PREFIX env variable must be specified}"

while IFS= read -d $'\0' -r dir; do
    # build image
    svcname="$(basename "${dir}")"
    image="${REPO_PREFIX}/$svcname:$TAG"
    (
        cd "${dir}"
        log "Building: ${image}"
        docker build -t "${image}" .

        log "Pushing: ${image}"
        docker push "${image}"
    )
done < <(find "${SCRIPTDIR}/../src" -mindepth 1 -maxdepth 1 -type d -print0)

log "Successfully built and pushed all images."
