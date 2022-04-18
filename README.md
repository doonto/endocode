## The goal

To create a simplest http web-server in Golang, package it in the Docker image and deploy to k8s as a container.

## Usage cases

Build and deploy to GKE or locally. Access the created host with :8080/helloworld. All the other URLs must give the 404
error.

### Dependencies

- chocolatey latest for windows. Install from https://chocolatey.org/install#install-step2
- Git >=2.35.3. Install for win by ```choco install git```
- go >=1.18. Install using https://go.dev/doc/install or for win by ```choco install go```
- Docker >=20.10.xx. Install from https://www.docker.com/products/docker-desktop/
- make for local builds. Install for win by ```choco install make```
- kubectl >=1.23. Install for win by ```choco install kubernetes-cli```
- kustomize >=4.4.1. Install for win by ```choco install kustomize```
- helm >=3.8.2. Install for win by ```choco install kubernetes-helm```

### Build

The code supposed to be run only in container. Ensure, you have all the dependencies in place.

#### Makefile for local builds

```make -f makefile```

Builds and run the code in the docker container. Brings some other non-executing by default commands as well.

####