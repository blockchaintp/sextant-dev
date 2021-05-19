## taekion-ui

How to spin up the dev environment for developing on the Taekion UI.

```bash
export CODE=~/projects/blockchaintp
cd $CODE
# git clone sextant, sextant-dev, sextant-api, taekion-fs
```

## start sextant

```bash
cd $CODE/sextant-dev
export MANUALRUN=1
make dev
```

get a few terminals open with tmux - start the frontend in one:

```bash
make frontend.cli
yarn develop
```

start the api in another:

```bash
make api.cli
npm run serve
```

## boot k8s

get a `kind` cluster in another terminal and setup:

```bash
kind create cluster
docker network connect sextant-dev_default kind-control-plane
bash $CODE/sextant-api/scripts/create-remote-credentials.sh
bash $CODE/sextant-api/scripts/get-remote-credentials.sh
```

`https://kind-control-plane:6443` is the api server

## build taekion images

Now you can open `http://localhost` and create a cluster using the above

Build the taekion images:

```bash
cd $CODE/taekion-fs
make docker-tp
make docker-middleware
make docker-client
kind load docker-image taekion/taekion-fs-tp:latest
kind load docker-image taekion/taekion-fs-middleware:latest
```

## update taekion chart

If the chart needs updating - you need to stop the api server and `export USE_LOCAL_CHARTS=1` then restart

on host:

```bash
docker cp $CODE/helm-charts/charts/tfs-on-sawtooth sextant-dev_api_1:/app/api/helmCharts/tfs-on-sawtooth/0.1
```

in api terminal:

```bash
ctrl+c
export USE_LOCAL_CHART=/app/api/helmCharts/tfs-on-sawtooth/0.1/tfs-on-sawtooth
node src/index.js
```

## deploy taekion

Now you can use the UI to deploy a taekion cluster.

(give taekion a few mins to spin up otherwise you will get `socket closed` errors)

## use CLI

Now - create a new volume called `test` - then use the `CLI` interface to start a CLI

```bash
cd $HOME
mkdir test
tfs-fuse -v test -m $HOME/test
cd test
echo hello > file.txt
```