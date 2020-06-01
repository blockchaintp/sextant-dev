# sextant-dev

This is a supporting repository to help developers work on sextant related codes, which includes the:

* [sextant](https://github.com/catenasys/sextant) frontend and;
* [sextant-api](https://github.com/catenasys/sextant-api) api stack.

## Pre-requisite

Install:

 * [docker](https://docs.docker.com/install/)
 * [docker-compose](https://docs.docker.com/compose/install/)
 * [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

You should also activate a local install of Kubernetes. Use the version of `docker` that you have installed and follow this [instruction](https://rominirani.com/tutorial-getting-started-with-kubernetes-with-docker-on-mac-7f58467203fd) to verify that it works.

Clone the following repos under the same root folder:

 * [sextant](https://github.com/catenasys/sextant)
 * [sextant-api](https://github.com/catenasys/sextant-api)
 * [sextant-dev](https://github.com/catenasys/sextant-dev)

PLEASE NOTE: If you are running this under WSL with Docker-Win the folder needs to be addressable from windows under the same basic name, and the drive needs to be shared to Docker-Win

For example:

```bash
mkdir sextant-project
cd sextant-project
git clone git@github.com:catenasys/sextant-dev.git
git clone git@github.com:catenasys/sextant.git
git clone git@github.com:catenasys/sextant-api.git
```

## Running the project for development

STEP 1: Set up your AWS credentials <br>
STEP 2: Build artifacts and prep docker containers for execution <br>
STEP 3: Start up api and sextant frontend <br>

#### STEP 1 - AWS credentials

You will need a `credentials.env` file inside the `sextant-dev` folder that has the following variables:

```
AWS_ACCESS_KEY_ID=XXX
AWS_SECRET_ACCESS_KEY=XXX
```

You will need to set it in a profile (e.g. in `.bashrc` or `.bash_profile` or using `direnv`) that is independent of terminal session.

These credentials will be used by the api container when connecting to AWS in instances where you need to connect to a BTP AWS kubernetes instance.

#### STEP 2 - Build executable artifacts

From within the `sextant-dev` folder:

STEP 2.1: Open a shell terminal <br>
STEP 2.2: Set the environment variable `export MANUALRUN=1` <br>
STEP 2.3: Run the command `make dev` <br>

For example, open one bash terminal and run the following sequences of commands:

```bash
export MANUALRUN=1
make dev
```

#### STEP 3 - Start up api and sextant frontend

All Sextant artifacts run from within docker containers. To start and stop sextant, you will need to manipulate the running containers accordingly.

**Docker operations**

You should consult docker documentation to ensure that you find this execution sequence appropriate for your needs.

The following is an example where you wish to ensure a **completely** clean state with no containers, images and postgress db in your system. **NOTE:** This is a highly destructive action.

```bash
docker rm -f $(docker ps -aq)
docker rmi -f $(docker images -q)
docker volume rm sextant-dev_postgres-data
```
**Starting api:**

STEP 3.1.1: Open a shell terminal. <br>
STEP 3.1.2: Set the environment variable `export MANUALRUN=1` <br>
STEP 3.1.3: Access the internals of the sextant-api container by running the `make api.cli` script. <br>
STEP 3.1.4: Activate the api code.<br>

```bash
export MANUALRUN=1
make api.cli
npm run serve
```

**Starting frontend:**

STEP 3.2.1: Open a shell terminal (one that is separate from the one you use to start the api).<br>
STEP 3.2.2: Set the environment variable `export MANUALRUN=1`<br>
STEP 3.2.3: Open a shell terminal (different from the one for the API)<br>
STEP 3.2.4: Run the following command sequence in the terminal<br>

```bash
export MANUALRUN=1
make frontend.cli
npm run develop
```
## Choose an Edition Module
By default, sextant builds in 'Dev-Mode'. If you'd like to build other editions, you'll need to edit the volume being copied in `docker-compose.yml` Instead of dev.js, copy the edition you want to build.
```yaml
volumes:
  - ../sextant-api/src:/app/api/src
  - ../sextant-api/test:/app/api/test
  - ../sextant-api/migrations:/app/api/migrations
  - ../sextant-api/editions/dev.js:/app/api/src/edition.js
  ```

## Support for NodeJS GRPC smoke testing
Please refer to instruction [./daml-node-grpc-smoke/README.md](./daml-node-grpc-smoke/README.md)

## Boot a Kubernetes cluster locally

It can be useful to have a Kubernetes cluster running locally on your laptop to test against.

To do this, we will use [kind](https://github.com/kubernetes-sigs/kind).

[install kind](https://github.com/kubernetes-sigs/kind#installation-and-usage)

Boot the stack as normal (so you have sextant and sextant-api up and running).

Create a new kind cluster:

```bash
kind create cluster
```

This will have created a container called `kind-control-plane`

It will also have adjusted your kubeconfig to point at the local cluster.

```bash
kubectl get ns
```

You might have done some work in the meantime that involved connecting to a different cluster.  sIf you want to re-point your kubeconfig at the local cluster:

```bash
kind export kubeconfig
```

Before we can use this cluster with sextant - we need to connect the api container to the kind-control-plane.

```bash
docker network connect sextant-dev_default kind-control-plane
```

We can then test this:

```bash
make api.cli
apt-get install telnet
telnet kind-control-plane 6443
```

You can launch sextant, and create a cluster.

Run the script to create the service account and download the credentials.

You must change the API server URL to `https://kind-control-plane:6443`.

## Running the taekion-tp locally

To test the taekion tp locally with kind - first spin up a kind cluster by following the guide above.

Then get access to the taekion docker repo and:

```bash
docker pull taekion/taekion-fs-tp:latest
```

Then import this image to the kind cluster:

```bash
kind load docker-image taekion/taekion-fs-tp:latest
```

The taekion deployment yaml has got `imagePullPolicy: Never` for the taekion tp.

**IMPORTANT** make sure we remove `imagePullPolicy: Never` before trying to deploy this to production.

The current setup will only work locally in kind by doing the trick above.

TODO: work out how to get imagePullSecrets to work with kind.