#!/bin/bash -eu

## GO related functions that can be used by other scripts, use source lib/scripts/go.sh to use
## V0.0.2 : Tim Dadd : First version

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