#!/bin/bash -eu
#Tim Dadd: 06-06-2020: First version

# check https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html#The-Set-Builtin for options
# Exit immediately if error, unset varible usage is an error, if a pipe command fails then use output of previous
set -o nounset
set -o errexit
set -o pipefail

RED=$(tput setaf 1)
YELLOW=$(tput setaf 3)
CYAN=$(tput setaf 6)

projd=$PWD # make a note of the project directory
if [ ! -d $projd/services ];then
  echo "${RED}Why is there no ${YELLOW}services${RED} directory?"
  echo "${CYAN}This script should be run from project root${WHITE}"
  exit 1
fi

source lib/scripts/env.sh

show_help () {
    echo "Usage: $0 {local|docker|minikube|gke} [option...] {start|stop|dev|test|build|clean} [service names ...]" >&2
    echo "   local ..... without docker, kubernetes etc."
    echo "   docker .... local docker only (using Dockerfile)"
    echo "   minikube .. minikube using skaffold - this is the default"
    echo "   gke ....... Google Kubernetes Engine"
    echo
    echo "   -h, --help"
    echo "   -v, --verbose"
    echo "   -s, --silent if you want less output to terminal"
    echo
    echo "   Guess what start, stop and restart means"
    echo
    echo " Roll your own command order:"
    echo "   $0 local start -o tim.txt"
    echo "   $0 start local -o tim.txt"
    echo "   $0 -o tim.txt gke restart"
    echo
    echo "You can only specify service names for local & docker deployment"
}

SERVICES=($(ls $projd/services))

# Parse the command line
# Initialize our own variables
verbose=0
silent=0
terminate=0
deploy_to="minikube"
deploy_mode="start"
services=()
IFS=",";svc_csv=${SERVICES[*]};IFS=""
while :
do
    if [ $# = 0 ]; then break ; fi
    case "$1" in
#      -o | --output)
#          if [ $# -ne 0 ]; then
#            output_file="$2"   # Should we check if $2 is valid - maybe later
#          fi
#          shift 2
#          ;;
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
      clean | stop)
        deploy_mode="$1"
        terminate=1
        shift 1
        ;;
      start | dev | build | test)
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
        svc_exists=$(echo ",$svc_csv," | grep -ic ",$1," || true)
        if [ "$svc_exists" == "1" ] ; then
          services+=("$1")
          shift 1
        else
          echo "${RED}Error: Unknown value: $YELLOW$1$CYAN"
          echo "Service names: ${POWDER_BLUE}$svc_csv${WHITE}"
          show_help
          echo $WHITE
          exit 1
        fi
        ;;
    esac
done
#echo "${services[@]}"
if [  ${#services[@]} -ne 0 ]; then
  case "$deploy_to" in
    local | docker)
      SERVICES=()
      SERVICES=("${services[@]}")
      ;;
    *)
      echo "${RED}You cannot specify a service name when deploying to $YELLOW$deploy_to$CYAN"
      show_help
      echo $WHITE
      exit 1
      ;;
  esac
fi
echo "Deploy to $deploy_to, mode=$deploy_mode, verbose=$verbose, silent=$silent"  #, output_file='$output_file'"
echo "Services:" "${SERVICES[*]}"

# Get a list of services
# Collect test targets
# 2>&1 combines stderr and stdout into the stdout stream
SERVICES_COUNT=0
GO_SERVICES=""
GO_SERVICE_COUNT=0
GO111MODULE=on # Everything about this uses go modules
# Get the master library version, IFS = Internal Field Separator, versioning as per https://semver.org/
libVer=$(grep -Po 'var VERSION = "\K.*(?=")' $projd/lib/version.go)
IFS=".";libV=($libVer);IFS=""
#echo "Master library Version: $libVer =>" "${libV[@]}"
for servicename in "${SERVICES[@]}"
do
  # Enable C code, as it is needed for SQLite3 database binary testing
  export CGO_ENABLED=1
  let SERVICES_COUNT+=1
  svc_dir="$projd/services/$servicename"
      ## If we find main.go then OK
  if [ -f "$svc_dir/main.go" ]; then
    GO_SERVICES[$GO_SERVICE_COUNT]=$servicename
    let GO_SERVICE_COUNT+=1
    if [ $terminate -eq 0 ]; then
      cd $svc_dir
      echo "${ORANGE}Verifying $LIME_YELLOW$servicename$WHITE"
      # We need lib for version control
      if [ ! -d lib ]; then
        cp -r ../../lib lib
      fi
      # We need go.mod to use lib
      if [ ! -f go.mod ]; then
        go mod init
        go mod edit -module $servicename
        go mod edit -require lib@v0.0.0
        go mod edit -replace lib@v0.0.0=./lib
      fi
      svcLibVer=$(grep -Po 'var VERSION = "\K.*(?=")' lib/version.go)
      IFS=".";svcLibV=($svcLibVer);IFS=""
      echo "Service library version: $svcLibVer =>" "${svcLibV[@]}"
      ## Are the library versions different?
      libEquality="equal"
      for ((i=0; i<${#libV[@]}; i++))
      do
          if [ ${libV[i]} -gt ${svcLibV[i]} ]; then
            libEquality="gt"
            break
          elif [ ${libV[i]} -lt ${svcLibV[i]} ]; then
            libEquality="lt"
            break
          fi
      done
      export GOFLAGS="-mod=mod"  # Use the library within the microservice directory (default prior to 14)
      rm -Rf $svc_dir/vendor
      case "$libEquality" in
        equal)
          echo "${GREEN}Library versions equal"
          export GOFLAGS="-mod=vendor"  ## Use the library within the vendor directory (default since 14)
          go mod vendor
          ;;
        gt)
          echo "${ORANGE}The master library version ($libVer) is NEWER than the version ($svcLibVer) being used by $LIME_YELLOW$servicename$WHITE"
          read -p "Do you want to update this library from master? ${YELLOW}(y to update)${NORMAL} ? " yn
          if [ "$yn" == "y" ]; then
            rm -Rf lib
            rm -Rf vendor/lib
            cp -r ../../lib lib
          fi
          ;;
        lt)
          echo "${RED}The master library version ($libVer) is OLDER than the version ($svcLibVer) being used by $LIME_YELLOW$servicename$WHITE"
          read -p "Do you want to update the master library with this library? ${YELLOW}(y to update)${NORMAL} ? " yn
          if [ "$yn" == "y" ]; then
            mv ../../lib ../../lib_V$libVer
#            rm -Rf ../../lib
            cp -r lib ../../lib
            libVer=$(grep -Po 'var VERSION = "\K.*(?=")' $projd/lib/version.go)
            IFS=".";libV=($libVer);IFS=""
            echo "Master library Updated to Version: $libVer =>" "${libV[@]}"
          fi
          ;;
      esac
      echo "${YELLOW}Cleaning up go.mod & running go test $WHITE"
      go mod tidy

      echo "${ORANGE}Running tests on $LIME_YELLOW$servicename$WHITE"
      go test ./... 2>&1
      echo $WHITE

      # Collect all `.go` files and `gofmt` against them. If some need formatting - print them.
      echo -n "${CYAN}Format check: "
#      go fmt ./...
      fmtErrors=$(find "$@" -type f -name \*.go | xargs gofmt -l 2>&1 || true)
      if [ "${fmtErrors}" ]; then
          for f in ${fmtErrors}; do
            case "$f" in
            *vendor/gopkg.in/yaml.v2/* | *lib/*)
  #            echo "Ignore formatting of $e"
              ;;
            *)
              echo -n "($f): "
              go fmt $f
            esac
          done
      fi
      echo "${GREEN}ALL FORMATTED$WHITE"

      # Run `go vet` against all targets. If problems are found - print them.
      echo -n "${CYAN}Vetting..."
      VetErrors=$(go vet ./... 2>&1 || true)
      if [ "${VetErrors}" ]; then
          echo "${RED}FAIL"
          echo "${ERRS}${WHITE}"
          exit 1
      fi
      echo "${GREEN}VETTED OK$WHITE"
      if [ "$deploy_mode" != "test" ]; then
        echo -n "${CYAN}Building $servicename: "
        # Disable C code, enable Go modules
        export CGO_ENABLED=0
        echo $PWD
        go build # since 14 we need this to ignore vendor directory
        echo "${GREEN}DONE$WHITE"
      fi
    fi
  fi
done
echo "Found $SERVICES_COUNT services & $GO_SERVICE_COUNT GO services"
IFS=" \t\n"

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
        bash -c "$cmd"
        ks=($(kubectl get services))
        echo $ks
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
        if [ $verbose == 1 ]; then
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

  gke)
    echo "${BRIGHT}Deploying to Google Cloud Google Kubernetes Engine (GKE)${NORMAL}"
    gcloudProject="$(gcloud config get-value project)" # Get the project name
    if [ ! $gcloudProject ]; then
      echo "${RED}Please configure the gcloud project ID. ${YELLOW}gcloud config set project ${BRIGHT}project-id$NORMAL$WHITE"
      exit
    fi
    echo "${CYAN}Project: $gcloudProject$WHITE"

    ke_clusterName=${projd##*/} # make the cluster name same as project name
    ke_clusterZone=$(gcloud config get-value compute/zone) # Use compute zone for cluster zone
    cd "$projd"
    case "$deploy_mode" in
      start | restart | dev)
        bash -c "gcloud services enable container.googleapis.com"
        ## Is the cluster already created
        #WARNING: Currently VPC-native is not the default mode during cluster creation. In the future, this will become the default mode and can be disabled using `--no-enable-ip-alias` flag. Use `--[no-]enable-ip-alias` flag to suppress this warning.
        #WARNING: Starting with version 1.18, clusters will have shielded GKE nodes by default.
        #WARNING: Your Pod address range (`--cluster-ipv4-cidr`) can accommodate at most 1008 node(s).
        IFS=$' \n'  # 0x0A=/n, default: IFS=$' \t\n'
        kc=($(gcloud container clusters list --filter="$ke_clusterName"))
        IFS=""
        ## 8 Columns[0-7]: NAME:LOCATION:MASTER_VERSION:MASTER_IP:MACHINE_TYPE:NODE_VERSION:NUM_NODES:STATUS
        if [ $kc ]; then
          ## Is the output in the right format?
          expected_fmt="NAME:LOCATION:MASTER_IP:STATUS:$ke_clusterName"
          actual_fmt="${kc[0]}:${kc[1]}:${kc[3]}:${kc[7]}:${kc[8]}"
          if [ "$expected_fmt" != "$actual_fmt" ] ; then
            echo 'Unexpected format for "kubectl get container clusters list"'
            echo "Expected: $expected_fmt"
            echo "Received: $actual_fmt"
            gcloud container clusters list --filter=$ke_clusterName
            exit 1
          fi
          if [ "${kc[15]}" != "RUNNING" ]; then
            echo "${RED}Cluster $ke_clusterName is not running!$WHITE"
            exit 1
          fi
          echo "${GREEN}Cluster $ke_clusterName is running${WHITE}"
          cmd="gcloud container clusters get-credentials $ke_clusterName"
        else
          cmd="gcloud container clusters create $ke_clusterName --enable-autoupgrade \
                --enable-autoscaling --min-nodes=3 --max-nodes=10 --num-nodes=5 --zone=$ke_clusterZone"
        fi
        echo "$cmd"
        bash -c "$cmd"
        kn=($(kubectl get nodes))  ## Check all is well
        # Do we need to add a check here?
        echo $kn
        # Enable Google Container Registry (GCR) on the GC project
        gcloud services enable containerregistry.googleapis.com
        # Configure the docker CLI to authenticate to GCR
        gcloud auth configure-docker -q
        # Use skaffold to do all the heavy lifting, build and push to GCR as per the kubernetes manifest
        run_mode=${deploy_mode/restart/run}
        run_mode=${run_mode/start/run}
        cmd="skaffold $run_mode --default-repo=gcr.io/$gcloudProject"
        if [ $verbose == 1 ]; then
          cmd="$cmd -vdebug"
        fi
        if [ $silent == 0 ]; then
          cmd="$cmd --tail"
        fi
        bash -c "$cmd"
        IFS=" \n"
        ks=($(kubectl get services))
        IFS=""
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
        echo "Checking all the services are working as expected"
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
          exit 1
        fi
        kubectl get service frontend-external
        ;;

      clean | stop)
        gcloud auth configure-docker -q
        bash -c "skaffold delete"
        read -p "Do you want to delete GKE Cluster $ke_clusterName ${YELLOW}(y to stop)${NORMAL} ? " yn
        if [ "$yn" == "y" ]; then
          cmd="gcloud container clusters delete $ke_clusterName"
          bash -c "$cmd"
        fi
        ;;
    esac
  ;;
esac