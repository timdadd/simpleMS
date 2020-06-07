#!/bin/bash -eu
#Tim Dadd: 06-06-2020: First version

# check https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html#The-Set-Builtin for options
# Exit immediately if error, unset varible usage is an error, if a pipe command fails then use output of previous
set -o errexit
#set -o unset
set -o pipefail

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

show_help () {
    echo "Usage: $0 {local|docker|minikube|gke} [option...] {start|stop|dev|test|clean} [service names ...]" >&2
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

projd=$PWD
if [ ! -d $projd/services ];then
  echo "${RED}Why is there no ${YELLOW}services${RED} directory?"
  echo "${CYAN}This script should be run from project root${WHITE}"
  exit 1
fi
SERVICES=($(ls $projd/services))

# Parse the command line
# Initialize our own variables and make a note of the project directory:
output_file=""
verbose=0
silent=0
deploy_to="minikube"
deploy_mode="start"
services=()
while :
do
    if [ $# = 0 ]; then break ; fi
#    if [[ " ${SERVICES[@]} " =~ " $1 " ]]; then
#      services+=($1)
#    fi

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
          echo "${RED}Error: Unknown option: $YELLOW$1$CYAN" >&2
          show_help
          echo $WHITE
          exit 1
          ;;
      *)  # Could this be a service name?
          if 'echo "${SERVICES[@]}" | grep -q "$1"' ; then
              services+=("$1")
              shift 1
          else
            echo "${RED}Error: Unknown value: $YELLOW$1$CYAN" >&2
            show_help
            echo $WHITE
            exit 1
          fi
#          break
          ;;
    esac
done
echo "${services[@]}"
if [  ${#services[@]} -ne 0 ]; then SERVICES=(); SERVICES=("${services[@]}"); fi

echo "Deploy to $deploy_to, mode=$deploy_mode, verbose=$verbose, silent=$silent, output_file='$output_file'"
echo -n "Services: "
echo "${SERVICES[@]}"

# Get a list of services
# Collect test targets
# 2>&1 combines stderr and stdout into the stdout stream
SERVICES_COUNT=0
GO_SERVICES=""
GO_SERVICE_COUNT=0
export GO111MODULE=on
for servicename in "${SERVICES[@]}"
do
  # Enable C code, as it is needed for SQLite3 database binary
  # Enable go vendor because we need everything in one place to build the docker image
  export CGO_ENABLED=1
  export GOFLAGS="-mod=vendor"

  let SERVICES_COUNT+=1
  svc_dir="$projd/services/$servicename"
      ## If we find main.go then OK
  if [ -f "$svc_dir/main.go" ]; then
    GO_SERVICES[$GO_SERVICE_COUNT]=$servicename
    let GO_SERVICE_COUNT+=1
    if [ "$deploy_mode" != "stop" ]; then
      echo "${ORANGE}Running tests on $LIME_YELLOW$servicename $POWDER_BLUE"
      cd $svc_dir
      echo "${YELLOW}Cleaning up go.mod & running go test"
      go mod edit -module $servicename
      go mod edit -replace lib/common@v0.0.0=./lib/common
      go mod edit -require lib/common@v0.0.0
      go mod tidy
      go mod vendor
      go mod tidy
      go test ./... 2>&1
      echo $WHITE

      # Collect all `.go` files and `gofmt` against them. If some need formatting - print them.
      echo -n "${CYAN}Checking go fmt: "
      ERRS=$(find "$@" -type f -name \*.go | xargs gofmt -l 2>&1 || true)
      if [ -n "${ERRS}" ]; then
  #        echo "${YELLOW}Formatting the following files:"
          for e in ${ERRS}; do
            case "$e" in
            *vendor/gopkg.in/yaml.v2/*)
  #            echo "Ignore formatting of $e"
              ;;
            *)
              echo " fmt $e:"
              go fmt $e
            esac
          done
      fi
      echo "${GREEN}ALL FORMATTED$WHITE"

      # Run `go vet` against all targets. If problems are found - print them.
      echo -n "${CYAN}Checking go vet: "
      ERRS=$(go vet ./... 2>&1 || true)
      if [ -n "${ERRS}" ]; then
          echo "${RED}FAIL"
          echo "${ERRS}"
          echo
          exit 1
      fi
      echo "${GREEN}PASS$WHITE"
      echo

      if [ "$deploy_mode" != "test" ]; then
        echo -n "${CYAN}Building $servicename: "
        # Disable C code, enable Go modules
        export CGO_ENABLED=0
        echo $PWD
        go build
        echo "${GREEN}DONE$WHITE"
      fi
    fi
  fi
done
echo "Found $SERVICES_COUNT services & $GO_SERVICE_COUNT GO services"

if [ "$deploy_mode" == "test" ]; then exit 0; fi

## Keep these as separate tests for future thinking
if [ "$deploy_mode" == "build" ]; then exit 0; fi

## OK let's build the command line
case "$deploy_to" in
  minikube)
    echo "Running minkube locally"
    cd $projd
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
#        bash -c "$cmd"
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
        echo "Checking all the services working as expected"
          ## 6 Columns: NAME:TYPE:CLUSTER-IP:EXXTERNAL-IP:PORT(S):AGE
        OK_Count=0
        i=0
        for sk in "${ks[@]}"
        do
#          echo "$sk, $i"
          # Up the service count?
          for s in "${GO_SERVICES[@]}"
          do
            if [ "$s" == "$sk" ]; then
              let OK_Count+=1
              for j in {0..5};
              do printf "${ks[0+j]}=${ks[i+j]},"; done
              echo""
            fi
          done
          i=$((i+1))
        done
        echo "$OK_Count Services Running"
        echo "Looking for $GO_SERVICE_COUNT services"
        if [[ "$OK_Count" != "$GO_SERVICE_COUNT" ]]; then
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
    for s in "${GO_SERVICES[@]}"
    do
      svc_dir="$projd/services/$s"
          ## If we find main.go then OK
      echo "$deploy_mode service $s"
      killall -qw $s || true
      case "$deploy_mode" in
      start | restart | dev)
        bash -c "cd $svc_dir && ./$s &"
        sleep 2 # Let any early logs appear in the right order
        ;;
      clean)
        bash -c "rm $svc_dir/&s"
        ;;
      esac
    done
    echo "${GREEN}$deploy_mode complete$WHITE"
    ;;

  docker)
    echo "Running locally with Docker"
    case "$deploy_mode" in
      start | restart | dev | stop)
          ## Network name is top level directoryNname_net
        dknet="${projd##*/}_net"
        dkname=() ## List of docker containers to start
        dkcmd=() ## list of docker commands to start the containers

        all_env="" ## List of environment variables to pass into containers
#        docker network rm $dknet || true >/dev/null
        if ! docker network ls | grep -q "$dknet" ; then
          echo "${POWDER_BLUE}Creating docker network $dknet"
          docker network create $dknet
        fi
        # Start any service independent support services here (e.g. kafka, MQ, logging)
#        echo "rabbitmq"
#        dkname+=("rabbitmq")
#        dkcmd+=("container run -d --name rabbitmq --network $dknet $all_env rabbitmq:3-management")
#        all_env+="--env AMQP_BROKER_URL=amqp://guest:guest@rabbitmq:5672/"

        ## Now loop through each service and get them fired up
        for service in "${GO_SERVICES[@]}"
        do
          echo "${CYAN}Docker image(s) for $service"
          # Disable C code, enable Go modules
          export CGO_ENABLED=0
          svc_dir="$projd/services/$service"
          cd $svc_dir
          if [ $deploy_mode != "stop" ]; then
            echo "Building container $service at $svc_dir"
            docker rmi -f $service || true
            docker image build -t $service .
            echo "${GREEN}DONE$WHITE"
          fi

          ms_env=""
          if grep -q "mongo" go.mod ; then
            cname="$service-mongo"
            echo $cname
            dkname+=($cname)
            dkcmd+=("container run -d --name $cname --network $dknet mongo")
            ms_env+="--env MONGO_URL=mongodb://$cname/$service"
          fi
          dkopts=""
          portMap=$(grep "EXPOSE" Dockerfile | cut -d ' ' -f 2)
          if [ $portMap -a $deploy_mode != "stop" ] ;then
            echo "Expose Port $portMap"
            dkopts+="-p $portMap:$portMap"
          fi
          # -p 8181:8181 - need to fix port mappings
          dkname+=("$service")
          dkcmd+=("container run -d --name $service --network $dknet $dkopts $all_env $ms_env $service")
        done
        echo "${GREEN}All Images Determined${WHITE}"

        echo "${YELLOW}Stopping existing containers"
          ## Now we're going to turn all existing containers off
        cnames=$(docker container ls --format="{{.Names}}")
        for c in "${dkname[@]}"
        do
          if [[ $cnames == *"$c"* ]]; then
            docker stop $c
            docker rm $c
            echo "${GREEN}Containers $c Stopped and removed${WHITE}"
          fi
        done

        if [ $deploy_mode != "stop" ]; then
            ## Now deploy and run containers with the freshly built images
          i=0
          for cmd in "${dkcmd[@]}"
          do
            echo "${ORANGE}Creating and starting container for ${dkname[i]}"
            echo $cmd
            bash -c "docker $cmd"
            sleep 3 ## Wait for everything to stabilise
            docker container logs ${dkname[i]}
            let i+=1
            echo
          done

          sleep 3 ## Wait for everything to stabilise
          echo "${GREEN}All Containers Started${WHITE}"
        fi
        if [ $silent == 0 ]; then
          echo $CYAN
          docker ps
          echo $WHITE
          docker network inspect simplems_net
        fi
      ;;

      clean)
        echo "${CYAN}Cleaning up docker containers & images"
        echo "${GREEN}DONE$WHITE"
        docker image prune -f
        docker container prune -f
        docker network prune -f
        echo "${GREEN}Cleaned up docker containers & images"
      ;;
    esac
  ;;
esac