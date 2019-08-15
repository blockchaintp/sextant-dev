# sextant-dev

This is a supporting repository to help developers work on sextant related codes, which includes the:

* [sextant](https://github.com/catenasys/sextant) frontend and;
* [sextant-api](https://github.com/catenasys/sextant-api) api stack.

## Pre-requisite

Install:

 * [docker](https://docs.docker.com/install/)
 * [docker-compose](https://docs.docker.com/compose/install/)

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

STEP 1: Set up your AWS credentials
STEP 2: Build artefacts and prep docker containers for execution
STEP 3: Start up api and sextant frontend

#### STEP 1 - AWS credentials

You will need a `credentials.env` file inside the `sextant-dev` folder that has the following variables:

```
AWS_ACCESS_KEY_ID=XXX
AWS_SECRET_ACCESS_KEY=XXX
```

These credentials will be used by the api container when connecting to AWS.

#### STEP 2 - Build executable artefacts

From within the `sextant-dev` folder:

STEP 2.1: Set environmental variable `MANUALRUN` to 1
STEP 2.2: Run the command `make dev`

For example, open one bash terminal and run the following sequence of commands:

```bash
export MANUALRUN=1
make dev
```

#### STEP 3 - Start up api and sextant frontend

All Sextant artefacts run from within docker containers. To start and stop sextant, you will need to manipulate the running containers accordingly.

**Docker operations**

You should consult docker documentation to ensure that you find execution sequence appropriate for your needs.

The following is an example where you wish to ensure a **comnpletely** clean state with no containers and images in your system. 

```bash
docker rm -f $(docker ps -a)
docker rmi -f $(docker images -q)
```
**Starting api:**

STEP 3.1.1: Open a shell terminal.
STEP 3.1.2: Access the internals of the sextant-api container by running the `make api.cli` script.
STEP 3.1.3: Assuming that you have a completely clean sextant state, run the preserve script. If you are merely re-starting stopped containers skip this step.
STEP 3.1.4: Activate the api code.

Assuming a completely clean state, execute the following sequence of commands:

```bash
make api.cli
npm run preserve 
node src/index.js
```
Running the command `npm run preserve` populate the postgres with appropriate schema.

Assuming you are working from a shutdown state, execute the following sequence of commands:

```bash
make api.cli
node src/index.js
```

**Starting frontend:**

STEP 3.2.1: Open a shell terminal (different from the one for the API).
STEP 3.2.2: Run the following command sequence in the terminal

```bash
make frontend.cli
yarn run develop
```
