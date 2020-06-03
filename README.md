#Simple (docker/kubernetes) (local/GKE) deployable microservice demo

This is based upon various google tutorials / examples

## Service Architecture

Composed of two microservices, frontend and systemservice that talk to each other over http.

[![Architecture of the two microservices]

| Service                                              | Description                                     |                                                                                  |
| ---------------------------------------------------- | ------------- | --------------------------------|
| [frontend](./services/frontend)                      | Exposes an HTTP server to serve the website.    |
| [systemservice](./services/systemservice)            | Just provides system info to the front end.     |

## Features

- **[Kubernetes](https://kubernetes.io)/[GKE](https://cloud.google.com/kubernetes-engine/):**
  The services run in the following modes:
   - GO RUN, one terminal window per microservice
   - DOCKER, build docker containers and then connect together in a network
   - Kubernets using minikube
   - Kuberneters on the cloud with GKE
- **[Skaffold](https://skaffold.dev):** Application
  is deployed to Kubernetes with a single command using Skaffold.

## Installation on Ubuntu 20 with GoLang 14

* Install GoLang
* Install Docker - make sure you can run docker without root privileges by adding your user name to docker groups.
* git clone the repository somewhere on your PC

### Running locally

> 💡 Recommended if you're planning to develop the application or giving it a try on your local cluster.

### Install tools to run a Kubernetes cluster locally:
* kubectl
  - Can be installed via `gcloud components install kubectl`) or `sudo apt-get install kubectl`
* minikube
First update the system, no need for a virtual box as we're on linux, just run minikube with
the `docker` driver. Download and install minikube
```bash
sudo apt-get update
sudo apt-get install apt-transport-https
sudo apt-get upgrade
cd ~/Downloads
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
#chmod +x minikube
sudo install minikube /usr/local/bin/
minikube start --driver=docker
minikube status
```
Should provide a working status
```
minikube
type: Control Plane
host: Running
kubelet: Running
apiserver: Running
kubeconfig: Configured
```
Now set the default driver to docker
```bash
minikube config set driver docker
```
### Launch the minikube Kubernetes cluster
Start **Minikube** with a Kubernetes cluster, here we configure with at least:
 - 4 CPU's
 - 4.0 Gb memory

Need to rebuild the container for cpu/memeory config change
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
### Rebuild the container images and deploy to minikube
To do this we will use Skaffold which is a command line tool that facilitates continuous
development for Kubernetes applications. You can iterate on your application source code
locally then deploy to local or remote Kubernetes clusters. Skaffold handles the workflow
for building, pushing and deploying the application based upon kubernetes manifest files.

#### Install Skaffold
```bash
cd ~/Downloads
curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
sudo install skaffold /usr/local/bin/
skaffold version
```
If all is well
```
v1.10.1
```

#### Run skaffold
From the project directory where the `skaffold.yaml` file is.  The first time will be very slow as it builds the containers for
the first time.  Hopefully this will build and deploy the application.

This will show you everything going on and also the logs from the microservices on the console.
```bash
skaffold run -vdebug --tail
```

To rebuild the images automatically as you refactor the code, run `skaffold dev` command.

Verify the Deployments, Pods & services are up and running
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

The browser should open automatically.  The main page shows the system details, '/version' shows the version of the frontend
microservice.

### Quick clean up of the services/deployments and pods
Go to the ./services directory so that you have a list of all the services.  If you clean up the deployments then you have
to rebuild everything.  The pods are deleted with the deployment but the services stay.
```bash
kubectl delete deployment *
kubectl delete service *
```