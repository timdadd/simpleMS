#!/bin/bash -eu

# check https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html#The-Set-Builtin for options
# Exit immediately if error, unset varible usage is an error, if a pipe command fails then use output of previous
#set -o errexit
#set -o unset
#set -o pipefail

BLACK=$(tput setaf 0)
RED=$(tput setaf 1)
GREEN=$(tput setaf 2)
LIME_YELLOW=$(tput setaf 190)
YELLOW=$(tput setaf 3)
POWDER_BLUE=$(tput setaf 153)
BLUE=$(tput setaf 4)
MAGENTA=$(tput setaf 5)
CYAN=$(tput setaf 6)
WHITE=$(tput setaf 7)
BRIGHT=$(tput bold)
NORMAL=$(tput sgr0)
BLINK=$(tput blink)
REVERSE=$(tput smso)
UNDERLINE=$(tput smul)

show_help () {
    echo "Usage: $0 {local|docker|minikube|gke} [option...] {start|stop|dev|test|clean}" >&2
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
SERVICES=($(ls services))
SERVICES_COUNT=0
for s in "${SERVICES[@]}"
do
  echo "Test service $s"
  SERVICES_COUNT=$((SERVICES_COUNT+1))
done
echo "Found $SERVICES_COUNT services"
## Before any deployment I need to test everything
##  Need to write some tests first - let's do that bit later
# Enable C code, as it is needed for SQLite3 database binary
# Enable go modules
export CGO_ENABLED=1
export GO111MODULE=on
#export GOFLAGS="-mod=vendor"

if [ "$deploy_mode" == "test" ]; then exit 0; fi


if [ "$deploy_mode" == "build" ]; then exit 0; fi


## OK let's build the command line
case "$deploy_to" in
  minikube)
    echo "Running minkube locally"
    case "$deploy_mode" in
      start | restart | dev)
        mk_host_state=$(minikube status -f={{.Host}}) || true
        echo "Minikube host is $mk_host_state"
        if [ "$mk_host_state" != "Running" ]; then
          minikube start
        fi
        run_mode=${deploy_mode/restart/run}
        run_mode=${run_mode/start/run}
        cmd="skaffold $run_mode"
        if [ $verbose == 1 ]; then
          cmd="$cmd -vdebug"
        fi
        if [ $silent == 0 ]; then
          cmd="$cmd --tail"
        fi
        bash -c $cmd
        ks=($(kubectl get services))
        ## Is the output in the right format?
        expected_fmt="NAME:TYPE:EXTERNAL-IP"
        actual_fmt="${ks[0]}:${ks[1]}:${ks[3]}"
        if [ "$expected_fmt" != "$actual_fmt" ] ; then
          echo 'Unexpected format for "kubectl get services"'
          echo "Expected: $expected_fmt"
          echo "Received: $actual_fmt"
          kubectl get services
          exit 1
        fi
          ## Check all the services working as expected
          ## 6 Columns: NAME:TYPE:CLUSTER-IP:EXXTERNAL-IP:PORT(S):AGE
        OK_Count=0
        i=0
        for sk in "${ks[@]}"
        do
#          echo "$sk, $i"
          # Up the service count?
          for s in "${SERVICES[@]}"
          do
            if [ "$s" == "$sk" ]; then
              OK_Count=$((OK_Count+1))
              for j in {0..5};
              do printf "${ks[0+j]}=${ks[i+j]},"; done
              echo""
            fi
          done
          i=$((i+1))
        done
        echo "$OK_Count Services Running"
        if [[ "$OK_Count" != "$SERVICES_COUNT" ]]; then
          echo -e "${RED}Not all the services are deployed!!! $NORMAL"
        else
            ## Now we need to make the frontend-external service public in Minikube
          minikube service frontend-external
        fi
        ;;

      clean | stop)
        bash -c "skaffold delete"
        read -p "Do you want to stop minikube ${YELLOW}(y to stop)${NORMAL} ? " yn
        if [ "$yn" == "y" ]; then minikube stop; fi
        ;;
    esac
    ;;

  local)
    echo "Running locally"
    for s in "${SERVICES[@]}"
    do
      ms_dir="$PWD/services/$s"
          ## If we find main.go then OK
      if test -f "$ms_dir/main.go"; then
        echo "$deploy_mode service $s"
        killall -qw $s || true
        case "$deploy_mode" in
        start | restart | dev)
          bash -c "cd $ms_dir && go build"
          bash -c "cd $ms_dir && ./$s &"
          ;;
        clean)
          bash -c "rm $ms_dir/&s"
          ;;
        esac
      fi
    done

esac
