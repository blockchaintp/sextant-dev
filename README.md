# sextant-dev

Local development of the [sextant](https://github.com/catenasys/sextant) frontend and [sextant-api](https://github.com/catenasys/sextant-api) api stack.

## install

First install: 

 * [docker](https://docs.docker.com/install/)
 * [docker-compose](https://docs.docker.com/compose/install/)

Then - clone the following repos in the same folder as this repo:

 * [sextant](https://github.com/catenasys/sextant)
 * [sextant-api](https://github.com/catenasys/sextant-api)


For example:

```bash
git clone git@github.com:catenasys/sextant-dev.git
git clone git@github.com:catenasys/sextant.git
git clone git@github.com:catenasys/sextant-api.git
```

## running locally

From within the `sextant-dev` folder:

```bash
make dev
```

This is the equivalent of doing `docker-compose up`

Then you can view the app in your browser:

```bash
open http://localhost
```

## running locally with manual restarts

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
yarn run watch
```

**api:**

```bash
make api.cli
node src/index.js
```


