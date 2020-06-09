#!/bin/bash -eu

## functions that can be used by other scripts, use source lib/scripts/common.sh to use
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

function install_protogen() {
    ## Is protoc installed somewhere?  This gets the first directory where it exists
  #protocLocation=$(whereis -b protoc | sed -ne 's/^\(protoc: \)\([/|a-z]*\)protoc\(.*\)$/\2/p')
  protocLocation=$(whatdir protoc)
  if [ ! "$protocLocation" ]; then
    echo "${RED}YOU MUST INSTALL PROTOC, CHECK https://github.com/protocolbuffers/protobuf/releases$WHITE"
    exit 1
  fi

    ## Is the file already installed somewhere
  if [ "$(whatdir $2)" != "" ]; then return; fi
  echo "${GREEN}Found protoc at $CYAN$protocLocation"
  echo "${YELLOW}Installing $2$WHITE"
    ## Load the executable into GOBIN
  go install google.golang.org/$1/cmd/$2
    ## Can we see the program now?
  if [ "$(whatdir $2)" != "" ]; then return; fi

  ## Dig it out of GOBIN and install in the same location as protoc
  gobin=$(go env GOBIN)
  if [ "$gobin" = "" ]; then gobin="$HOME/go/bin"; fi
  echo "GOBIN=$gobin"
  pgc="$gobin/$2"
  if [ ! -f $pgc ]; then
    echo "${RED}Cannot find $2 in $gobin!!!, Sorry you have to install yourself$WHITE}"
    exit 1
  fi
  cp $pgc $protocLocation
    ## Can we see it now?
  if [ "$(whatdir $2)" = "" ]; then
    echo "${RED}You must install $2 in a location where it can be found with $YELLOWwhereis$WHITE}"
    exit 1
  fi
}
