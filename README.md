# All Go Microservice framework using gRPC

This is based upon various google tutorials / examples and is for my own learning benefit, it's provided
as-is and you're free to download and do your own thing.

Basically I wanted a framework that could run locally both inside and outside of docker, be tested to run
on kubernetes locally before deploying to kubernetes in the cloud.  This is so I don't have to rely upon
 an internet connection to do stuff.

Deploys to Google Cloud Kubernetes Engine (GKE) or local kubernetes engine Minikube by deploying
docker containers to Kubernetes pods, one pod per service.  Can just run docker containers locally
or run as local executables, no docker or kubernetes.

One goal was to make it so that each microservice truly is independent of the other microservices even though
they are all in the same GO project and re-use a common library of functions without manual copy/paste.  So actually
it's an automatic copy/paste (haha) with a rudimentary version control.  I also wanted to stay away from the
pre GO11 project structure of pkg/cmd.

## The four microservices

| Service                                               | Port | Description                                    |
| ----------------------------------------------| ----- | ----------------------------------------------|
| [frontend](./services/frontend)               | 8080  | Exposes an HTTP server to serve the website.  |
| [systemservice](./services/systemservice)     | 8082  | Just provides system info to the front end.   |
| [routeguide](./services/routeguide)           | 10000 | [grpc tutorial](https://www.grpc.io/docs/languages/go/basics/) |
| [book](./services/book)                       | 4000  | [microserver version of bookshelf](https://github.com/GoogleCloudPlatform/golang-samples/tree/master/getting-started/bookshelf) |

## Features

- The services run in 4 deployment modes:
   - GO RUN, locally, one process per microservice, check them with `ps -ef | grep simplems`
   - DOCKER, use local docker to build images and deploy containers, all connected in docker network with
   port exposure for access if `EXPOSE` exists in `Dockerfile`
   - Local Kubernetes cluster with one node per microservice using minikube
   - [Kubernetes on the Google cloud with GKE]((https://kubernetes.io)/[GKE](https://cloud.google.com/kubernetes-engine/)),
   although any Kubernetes should work but I've only used gcloud.
- **Use of [Skaffold](https://skaffold.dev)** Application to deploy to any Kubernetes cluster.  Same scripts are used
with minikube and GKE.

## Installation on Ubuntu 20 with GoLang 14
* Install GoLang
* Install Docker - make sure you can run docker without root privileges by adding your user name to docker groups.
* git clone the repository somewhere on your PC, it's GO 11+ (uses modules) so `go/src` isn't mandatory

### Intall kubernetes controller (kubectl)
  - Can be installed via `gcloud components install kubectl` or `sudo apt-get install kubectl`

### Install minikube to run a Kubernetes cluster locally:
  - First update the system, no need for a virtual box as we're on linux
  - just run minikube with the `docker` driver
  - Download and install minikube
```bash
sudo apt-get update
sudo apt-get install apt-transport-https # not sure if this was needed
sudo apt-get upgrade
cd ~/Downloads
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
#chmod +x minikube
sudo install minikube /usr/local/bin/
minikube start --driver=docker
minikube status
```
Should provide a working status e.g.:
```
minikube
type: Control Plane
host: Running
kubelet: Running
apiserver: Running
kubeconfig: Configured
```
Now set the default driver to docker so we don't need `--driver=docker` anymore
```bash
minikube config set driver docker
```

### Install Skaffold
Skaffold is the tool that automatically builds docker files and then deploys
them to kubernetes.  The kubernetes instance it deploys to is defined by the configuration
of the kubernetes controller (kubectl)
```bash
cd ~/Downloads
curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
sudo install skaffold /usr/local/bin/
skaffold version
```
If all is well
```
v1.11.0
```

## Deploying
The `deploy.sh` script in the main project directory supports all 4 deployment options.

The order of the command line instructions shouldn't matter, `./deploy.sh --help` for latest instructions

* Run without docker/minikube: `./deploy.sh local`
* Stop running the local version: `./deploy.sh local stop`
* Run on local docker installation, no minikube/KBE: `./deploy.sh docker`
* Stop running the docker installation: `./deploy.sh docker stop`
* Clean up the docker installation: `./deploy.sh docker clean`
* Run on minikube Cluster: `./deploy.sh`
* Stop minikube Cluster: `./deploy.sh stop`
* Run on GKE Cluster: `./deploy.sh gke`
* Stop GKE Cluster: `./deploy.sh gke stop`

## Other notes
### To start **Minikube** with more resources:
 - 4 CPU's
 - 4.0 Gb memory

Need to rebuild the minikube container for cpu/memory config change
```bash
minikube stop
minikube delete
minikube start --cpus=4 --memory 4096
```
Minikube automatically configures kubectl commands to talk to the minikube kubernetes cluster. Check
that `kubectl` is connected to minikube
```bash
kubectl get nodes
```
Should see something like this
```
NAME       STATUS   ROLES    AGE    VERSION
minikube   Ready    master   109s   v1.18.3
```

### Using skaffold to deploy to kubernetes
Once a cluster is created, Skaffold, facilitates continuous deployment to Kubernetes applications. You can
develop & test your application source code locally then deploy to local or remote Kubernetes clusters. Skaffold
handles the workflow for building, pushing and deploying the application based upon kubernetes manifest files. If
you only change code for one microservice then only that is re-deployed.

#### Running skaffold
From the project directory where the `skaffold.yaml` file is run the following command to deploy each
microservice to a separate node within the cluster.  As long as kubectl works and a cluster exists this
will build and deploy the different microservices, 

The `-vdebug` option will show you everything going on and the `--tail` options logs the microservices output
on the console.
```bash
skaffold run -vdebug --tail
```
The first time will be the slowest as it builds the containers for the first time.  

To rebuild the images automatically as you refactor the code, run `skaffold dev` command instead.  This is
somewhat more useful when working in combination with minikube but works equally well for GKE.

### Verify the Deployments, Pods & services are up and running
```bash
ubectl get deployment -o wide
kubectl get pods -o wide
kubectl get services -o wide
```
#### Access the web frontend through your browser
**Minikube** requires you to run a command to access the frontend service:

```bash
minikube service frontend-external
```

The browser should open automatically.  The main page shows the system details, '/version' shows the
version of the frontend microservice.

### Quick clean up of the services/deployments and pods
Go to the ./src directory so that you have a list of all the services.  If you clean up the deployments then you have
to rebuild everything.  The pods are deleted with the deployment but the services stay.
```bash
kubectl delete deployment *
kubectl delete service *
```

## Installing firebase emulators (e.g. local firestore) 
Assume gcloud installed, google account etc.
### Install JAVA 8+
Check if java loaded
```bash
java --version
## install?
sudo apt install openjdk-14-jdk
## installed
java --version
```

### Install/Upgrade Firebase CLI
Run the following cURL command:
```bash
curl -sL https://firebase.tools | bash
```
Create a directory to use as a firebase directory
```bash
mkdir ~/firebase && cd ~/firebase
```

### Install Firebase Emulator(s)
From the firebase directory, initialise the emulator:
```bash
firebase init emulators
```

Setup the emulator suite, 4 options presented:
1. Use an existing project
2. Create a new project
3. Add Firebase to an existing Google Cloud Platform project
4. Don't set up a default project 

Select **4** "Don't setup a default project" (we can do this later).

#### 5 Emulators are available
1. Functions
2. **Firestore** is Firebase's newest database for mobile app development. It builds on the 
successes of the Realtime Database with a new, more intuitive data model. Cloud Firestore 
also features richer, faster queries and scales further than the Realtime Database.
3. **Database** is Firebase's original database. It's an efficient, low-latency solution for
mobile apps that require synced states across clients in realtime.
4. Hosting
5. Pubsub

We'll go with firestore, option **2**.

Use port 4000 for the firestore emulator and enable the emulaator UI because it's 
[new](https://firebase.googleblog.com/2020/05/local-firebase-emulator-ui.html) and will
be good to try.  We'll use port 4080 for the UI.

This should install the emulator and create a `firebase.json` file like this
```json
{
  "emulators": {
    "firestore": {
      "port": 4000
    },
    "ui": {
      "enabled": true,
      "port": 4080
    }
  }
}
```
## Start the emulator
From your firebase project directory

```bash
cd ~/firebase
firebase emulators:start
```

It's also possible to specify what emulator to start
```bash
cd ~/firebase
sfirebase emulators:start --only firestore
```
