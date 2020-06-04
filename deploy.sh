#!/bin/bash -eu

# check https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html#The-Set-Builtin for options
# Exit immediately if error, unset varible usage is an error, if a pipe command fails then use output of previous
#set -o errexit
#set -o unset
#set -o pipefail

show_help () {
    echo "Usage: $0 {local|docker|minikube|gke} [option...] {start|stop|dev|test}" >&2
    echo "   local ..... without docker, kubernetes etc."
    echo "   docker .... local docker only (using docker.yaml?)"
    echo "   minikube .. minikube using skaffold - this is the default"
    echo "   gke ....... Google Kubernetes Engine"
    echo
    echo "   -h, --help"
    echo "   -v, --verbose"
    echo "   -o, --output filename"
    echo "   -s, --silent if you don't want any output to terminal"
    echo
    echo "   Guess what start, stop and restart means"
    echo
    echo " Roll your own ommand examples:"
    echo "   $0 local start -o tim.txt"
    echo "   $0 start local -o tim.txt"
    echo "   $0 -o tim.txt gke restart"
}

# Parse the command line
# Initialize our own variables:
output_file=""
verbose=0
silent=0
deploy_to="minikube"
deploy_mode="start"
while :
do
    if [ $# = 0 ]; then break ; fi
    case "$1" in
      -o | --output)
          if [ $# -ne 0 ]; then
            output_file="$2"   # Should we check if $2 is valid - maybe later
          fi
          shift 2
          ;;
      -h | --help)
          show_help
          exit 0
          ;;
      -v | --verbose)
          verbose=1
          shift 1
           ;;
      -s | --silent)
          silent=1
          shift 1
           ;;

      local | docker | minikube | gke)
          deploy_to="$1"
          shift 1
           ;;

      start | stop | dev | test)
          deploy_mode="$1"
          shift 1
           ;;

      --) # End of all options
          shift
          break
          ;;
      -*)
          echo "Error: Unknown option: $1" >&2
          show_help
          exit 1
          ;;
      *)  # No more options
          break
          ;;
    esac
done

echo "Deploy to $deploy_to, mode=$deploy_mode, verbose=$verbose, silent=$silent, output_file='$output_file'"

# Get a list of services
# Collect test targets
SERVICES=$(for d in "services"; do echo ./$d/...;done)
echo $SERVICES

## Before any deployment I need to test everything
##  Need to write some tests first - let's do that bit later
# Enable C code, as it is needed for SQLite3 database binary
# Enable go modules
export CGO_ENABLED=1
export GO111MODULE=on
#export GOFLAGS="-mod=vendor"


## OK let's build the command line
case "$deploy_to" in
  minikube)  ## Deploy to minkube
    mk_host_state=$(minikube status -f={{.Host}})
    echo "Minikube host is $mk_host_state"
    if [ "$mk_host_state" != "Running" ]; then
      minikube start
    fi
    run_mode=${deploy_mode/start/run}
    cmd="skaffold $run_mode"
    if [ $verbose == 1 ]; then
      cmd="$cmd -vdebug"
    fi
    if [ $silent == 0 ]; then
      cmd="$cmd --tail"
    fi
esac
bash -c "$cmd"