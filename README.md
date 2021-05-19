# **sextant-dev**

This repository supports sextant development which includes:

* [sextant](https://github.com/catenasys/sextant) frontend and;
* [sextant-api](https://github.com/catenasys/sextant-api) API stack.

## **Prerequisites**

---

### **Install:**

* [docker](https://docs.docker.com/install/)
* [docker-compose](https://docs.docker.com/compose/install/)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

You should also activate a local install of Kubernetes.

Using the version of `docker` you have installed, [verify it is activated](https://docs.docker.com/desktop/kubernetes/).

### **Clone the following repos under the same root folder:**

* [sextant](https://github.com/catenasys/sextant)
* [sextant-api](https://github.com/catenasys/sextant-api)
* [sextant-dev](https://github.com/catenasys/sextant-dev)

**PLEASE NOTE:** If you are running this under WSL with Docker-Win
the folder needs to be addressable from windows under the same
basic name, and the drive needs to be shared to Docker-Win

For example:

```bash
mkdir sextant-project
cd sextant-project
git clone git@github.com:catenasys/sextant-dev.git
git clone git@github.com:catenasys/sextant.git
git clone git@github.com:catenasys/sextant-api.git
```

&nbsp;

## **Running the project for development**

---

1. Set up your AWS credentials
1. Build artifacts and prep docker containers for execution
1. Start up API and sextant frontend

### **1. Set up your AWS credentials**

You will need to create a `credentials.env` file inside the `sextant-dev`.

* Replace `XXX` with your `AWS_ACCESS_KEY_ID`
* Replace `YYY` with your `AWS_SECRET_ACCESS_KEY`

```bash
echo -e "AWS_ACCESS_KEY_ID=XXX\nAWS_SECRET_ACCESS_KEY=YYY" > credentials.env
```

These credentials will be used by the API container when connecting
to AWS in instances where you need to connect to a BTP AWS
kubernetes instance.

### **2. Build executable artifacts**

From within the `sextant-dev` folder:

* Open a shell terminal
* Set the environment variable `export MANUALRUN=1`
* Run the command `make dev`

```bash
export MANUALRUN=1; make dev
```

### **3. Start up API and sextant frontend**

All Sextant artifacts run from within [docker containers](https://docs.docker.com/get-started/#what-is-a-container).

To start and stop sextant, you will need to manipulate the running
containers accordingly.

#### **Docker operations**

To ensure that you find these execution sequences appropriate
for your needs, consult [docker's documentation.](https://docs.docker.com/engine/reference/commandline/rm/)

The following command will simply stop all docker containers:

```bash
docker stop $(docker ps -aq)
```

The following is a sequence of commands that will remove all
containers and images and delete any sextant related postgres data:

**NOTE:** This is a **highly destructive** action.

```bash
docker rm -f $(docker ps -aq)
docker rmi -f $(docker images -q)
docker volume rm sextant-dev_postgres-data
```

#### **Starting API:**

* Open a shell terminal.
* Run the command `env` to double check that `MANUALRUN=1`
* If not, run this command `export MANUALRUN=1`
* Access the sextant-api container by running `make api.cli`
* If you're starting a totally clean slate, run `npm run preserve`
* If you are just re-starting stopped containers skip this step.
* Activate the API code.

Assuming a completely clean state, execute the following sequence of commands:

```bash
make api.cli
npm run preserve
node src/index.js
```

Running the command `npm run preserve` populates the postgres db
with the appropriate schema.

Alternatively, if you have already executed the above sequence
previously, all you need to do is to run this sequence:

```bash
make api.cli
node src/index.js
```

#### **Starting frontend:**

* Open a new shell terminal (a separate from the one you use to start the API).
* Run the command `env` to double check that `MANUALRUN=1`
* If not, run this command `export MANUALRUN=1`
* Open a shell terminal (again, different from the one running the API)
* Run the following command sequence in the terminal

```bash
make frontend.cli
npm run develop
```

&nbsp;

## **Choose an Edition Module**

---

By default, sextant builds in 'Dev-Mode'. If you'd like to build other editions,
you'll need to edit the volume being copied in `docker-compose.yml`
Instead of `dev.js`, copy the edition you want to build.

```yaml
volumes:
  - ../sextant-api/src:/app/api/src
  - ../sextant-api/test:/app/api/test
  - ../sextant-api/migrations:/app/api/migrations
  - ../sextant-api/editions/dev.js:/app/api/src/edition.js
```

&nbsp;

## **Support for NodeJS GRPC smoke testing**

---

Please refer to instruction [./daml-node-grpc-smoke/README.md](./daml-node-grpc-smoke/README.md)

&nbsp;

## **Boot a Kubernetes cluster locally**

---

It can be useful to have a Kubernetes cluster running locally on your
laptop to test against.

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

You might have done some work in the meantime involving connecting
to a different cluster. If you want to re-point your kubeconfig at
the local cluster:

```bash
kind export kubeconfig
```

Before we can use this cluster with sextant - we need to connect
the api container to the kind-control-plane.

```bashclear
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
