# sextant-dev

Local development of the [sextant](https://github.com/catenasys/sextant) frontend and [sextant-api](https://github.com/catenasys/sextant-api) api stack.

## install

First install: 

 * [docker](https://docs.docker.com/install/)
 * [docker-compose](https://docs.docker.com/compose/install/)

Then - clone the following repos in the same folder as this repo:

 * [sextant](https://github.com/catenasys/sextant)
 * [sextant-api](https://github.com/catenasys/sextant-api)

PLEASE NOTE: If you are running this under WSL with Docker-Win the folder needs to be addressable from windows under the same basic name, and the drive needs to be shared to Docker-Win

For example:

```bash
git clone git@github.com:catenasys/sextant-dev.git
git clone git@github.com:catenasys/sextant.git
git clone git@github.com:catenasys/sextant-api.git
```

## running locally

#### AWS credentials

You will need a `credentials.env` file inside the `sextant-dev` folder that has the following variables:

```
AWS_ACCESS_KEY_ID=XXX
AWS_SECRET_ACCESS_KEY=XXX
```

These credentials will be used by the api container when connecting to AWS.

#### boot stack

From within the `sextant-dev` folder:

```bash
make dev
```

This is the equivalent of doing `docker-compose up`

Then you can view the app in your browser:

```bash
open http://localhost
```

#### running locally with manual restarts

Sometimes - it's useful to have a command line inside the api & frontend containers for quick restarts.

To do this - export the `MANUALRUN` variable before running `make dev`:

```bash
export MANUALRUN=1
make dev
```

Then - in two seperate terminals you can manually start the frontend and api:


**frontend:**

```bash
make frontend.cli
yarn run develop
```

**api:**

```bash
make api.cli
node src/index.js
```


