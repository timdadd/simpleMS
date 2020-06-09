#!/bin/bash -eu

pb_file="routeguide"  ## Name of your .proto fie here

RED=$(tput setaf 1)
CYAN=$(tput setaf 6)
if [ ! -d ../../services ]; then
  echo "${RED}This needs to be run from a microservice directory$WHITE"
  exit 1
fi

## Does this directory have the lib files?
if [ ! -d ./lib ]; then
  echo "${CYAN}Copying the lib directory$WHITE"
  cp -r ../../lib .
fi

/bin/bash lib/scripts/genpb.sh $pb_file