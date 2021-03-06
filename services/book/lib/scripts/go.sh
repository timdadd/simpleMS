#!/bin/bash -eu

## GO related functions that can be used by other scripts, use source lib/scripts/go.sh to use
## Tim Dadd : Genesis June 2020

function install_protogen() {
    ## Is protoc installed somewhere?  This gets the first directory where it exists
  #protocLocation=$(whereis -b protoc | sed -ne 's/^\(protoc: \)\([/|a-z]*\)protoc\(.*\)$/\2/p')
  protocLocation=$(whatdir protoc)
  if [ ! "$protocLocation" ]; then
    echo "${RED}YOU MUST INSTALL PROTOC, CHECK https://github.com/protocolbuffers/protobuf/releases$WHITE"
    exit 1
  fi

  PlugIn=$2
  PlugInRepo=$1

    ## Is the protoc PlugIn already installed somewhere?
  if [ "$(whatdir $2)" != "" ]; then return; fi
  echo "${GREEN}Found protoc at $CYAN$protocLocation"
  echo "${YELLOW}Installing $PlugIn$WHITE"
    ## Load the executable into GOBIN
  (
    set -e
    go install $PlugInRepo/$PlugIn
    go get -u $PlugInRepo/$PlugIn
  )

    ## Can we see the program now?
  if [ "$(whatdir $PlugIn)" != "" ]; then return; fi

  ## Dig it out of GOBIN and install in the same location as protoc
  gobin=$(go env GOBIN)
  if [ "$gobin" = "" ]; then gobin="$HOME/go/bin"; fi
  echo "GOBIN=$gobin"
  pgc="$gobin/$PlugIn"
  if [ ! -f $pgc ]; then
    echo "${RED}Cannot find $PlugIn in $gobin!!!, Sorry you have to install yourself$WHITE}"
    exit 1
  fi
  cp $pgc $protocLocation
    ## Can we see it now?
  if [ "$(whatdir $PlugIn)" = "" ]; then
    echo "${RED}You must install $PlugIn in a location where it can be found with $YELLOWwhereis$WHITE}"
    exit 1
  fi
}