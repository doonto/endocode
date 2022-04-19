## The goal

To create a simplest http web-server in Golang, package it in the Docker image and deploy to k8s as a container.

## Usage cases

Build and deploy to GKE or locally. Access the created host with :8080/helloworld. All the other URLs give the 404
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

Builds and runs the code in the docker container locally. You shouldn't run binaries locally, if you suppose to run
container in the cloud. Brings some other non-executing by default commands as well.

#### Docker run for local build and IDE

```docker build -t latest . && docker run -p localhost:8080:8080 --name endocode latest -rm```

### Deploy

Deployment pipelines are defined in ```.github/workflows```. There are kustomize - based pipeline, which utilizes the
overlay way of making updates and helm-based pipeline, templating all the values. kust-pipeline.yml is set up to be run
on push to master. helm-pipeline.yaml must be triggered manually.

When creating the pipelines priority was given to a predefined tasks over the inline scripts.

### Debug

To debug the Go-code in the container, you want to run it with the following command,
see ```k8s/base/deployment.yaml```:

```- tail - "-f" - /dev/null```

then, open the bash in it:

```kubectl exec --stdin --tty <POD-NAME> -- /bin/bash```