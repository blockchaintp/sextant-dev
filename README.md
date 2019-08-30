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
STEP 3.1.4: Assuming that you have a completely clean sextant state, run the preserve script. If you are merely re-starting stopped containers skip this step.<br>
STEP 3.1.5: Activate the api code.<br>

Assuming a completely clean state, execute the following sequence of commands:

```bash
export MANUALRUN=1
make api.cli
npm run preserve
node src/index.js
```
Running the command `npm run preserve` populates the postgres db with the appropriate schema.

Alternatively, if you have already executed the above sequence previously, all you need to do is to run this sequence:

```bash
export MANUALRUN=1
make api.cli
node src/index.js
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
