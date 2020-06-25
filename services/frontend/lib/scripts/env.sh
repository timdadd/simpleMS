#!/bin/bash -eu

## environment related functions that can be used by other scripts, use source lib/scripts/env.sh to use
## V0.0.2 : Tim Dadd : First version

BLACK=$(tput setaf 0)
RED=$(tput setaf 1)
GREEN=$(tput setaf 2)
LIME_YELLOW=$(tput setaf 190)
YELLOW=$(tput setaf 3)
POWDER_BLUE=$(tput setaf 153)
BLUE=$(tput setaf 4)
MAGENTA=$(tput setaf 5)
CYAN=$(tput setaf 6)
ORANGE=$(tput setaf 10)
WHITE=$(tput setaf 7)
BRIGHT=$(tput bold)
NORMAL=$(tput sgr0)
BLINK=$(tput blink)
REVERSE=$(tput smso)
UNDERLINE=$(tput smul)

function whatdir() {
  # $(whereis -b $1 | sed -n -e 's/^$1: //p') - this doesn't get first location
  w="whereis -b $1 | sed -ne 's/^\($1: \)\([/|a-z]*\)$1\(.*\)$/\2/p'"
  bash -c "$w"
}

function whatAsc() {
  str=$1
  while   [ 0 -ne "${#str}" ]
  do      printf '%c(%x) ' "$str" "'$str"    #identical results for %.1s
          str=${str#?}
  done
  echo
  }