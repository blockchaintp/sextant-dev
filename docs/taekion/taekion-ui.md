## taekion-ui

# branches

 * sextant/SXT-535
 * sextant-api/SXT-535
 * sextant-dev/SXT-535
 * taekion-fs/SXT-535
 * helm-charts/SXT-535-features

# local dev

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
docker tag taekion/taekion-fs-tp:latest taekion/taekion-fs-tp:v0.6.0
docker tag taekion/taekion-fs-middleware:latest taekion/taekion-fs-middleware:v0.6.0
kind load docker-image taekion/taekion-fs-tp:v0.6.0
kind load docker-image taekion/taekion-fs-middleware:v0.6.0
```

## preload other images

If you are on a slow connection - it's quicker to pull the other images to the host and then import them:

```bash
export IMAGE_TAG=BTP2.1.0rc14
for image in sawtooth-devmode-engine-rust sawtooth-pbft-engine sawtooth-poet-cli sawtooth-poet-engine sawtooth-poet-validator-registry-tp sawtooth-raft-engine sawtooth-block-info-tp sawtooth-identity-tp sawtooth-intkey-tp-go sawtooth-settings-tp sawtooth-shell metrics-grafana metrics-influxdb sawtooth-validator sawtooth-rest-api; do
  docker pull blockchaintp/$image:$IMAGE_TAG
  kind load docker-image blockchaintp/$image:$IMAGE_TAG
done
docker pull postgres:11
kind load docker-image postgres:11
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

## rebuild taekion images

Once we have a deloyment - we will want to iterate on the taekion containers.

Here is how to rebuild and re-deploy the middleware:

```bash
cd $CODE/taekion-fs
export TAG=$(date +%s)
export MIDDLEWARE_IMAGE=taekion/taekion-fs-middleware
# set this to the namespace you deployed taekion into
export NAMESPACE=tfs
make docker-middleware
docker tag $MIDDLEWARE_IMAGE:latest $MIDDLEWARE_IMAGE:$TAG
kind load docker-image $MIDDLEWARE_IMAGE:$TAG
```
## quick rebuild

do this once

```bash
cd $CODE/taekion-fs
docker build -f docker/Dockerfile --target client-base -t taekion/client-base:quick .
docker build -f docker/Dockerfile --target build -t taekion/build:quick .
```

do this for each rebuild:

```bash
cd $CODE/taekion-fs
export TAG=$(date +%s)
export MIDDLEWARE_IMAGE=taekion/taekion-fs-middleware
export NAMESPACE=tfs
docker build -f $CODE/sextant-dev/docs/taekion/Dockerfile --target middleware -t $MIDDLEWARE_IMAGE:$TAG .
kind load docker-image $MIDDLEWARE_IMAGE:$TAG
```

## deploy taekion image

Now we need to know the index of the middleware container.

Keep changing the index in this command until we get `taekion/taekion-fs-middleware:XXX`

NOTE: there is probably a better way to update the specific container within the stateful set than manually typing indexes :-)

```bash
kubectl -n $NAMESPACE get statefulset.apps/tfs-validator -o=jsonpath='{.spec.template.spec.containers[7].image}'
export CONTAINER_INDEX=7
```

Now we can update the image using that index:

```bash
kubectl -n $NAMESPACE patch statefulset tfs-validator --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/'$CONTAINER_INDEX'/image", "value":"'"$MIDDLEWARE_IMAGE:$TAG"'"}]'
```

Then check for the rollout:

```bash
kubectl -n $NAMESPACE get po -w
```

## s3 development

To build the S3 container:

```bash
cd $CODE/taekion-fs
make docker-s3
```

Open the CLI page of the taekion deployment ane get the TFS_URL env:

```bash
export TFS_URL=XXX
```

Run the container:

```bash
function run_s3() {
  docker run --rm --network sextant-dev_default  \
    --name tfs-s3 \
    -p 8001:8001 \
    --entrypoint /usr/local/bin/tfs-s3 \
    -e TFS_URL=$TFS_URL \
    taekion/taekion-fs-s3:latest
}
run_s3
```

To rebuild:

```bash
make docker-s3 && run_s3
```

To run s3 commands:

 * https://www.npmjs.com/package/minio
 * https://docs.min.io/docs/minio-client-quickstart-guide.html

To install the minion client CLI:

```bash
wget https://dl.min.io/client/mc/release/linux-amd64/mc
chmod a+x mc
sudo mv mc /usr/local/bin
# these are fake credentials but the "mc" cli requires them
mc alias set tfs http://localhost:8001 "" ""
```

List buckets:

```bash
mc ls tfs
```


#### debug XML

To see how the XML should actually look - we can proxy connections off to play.min.io and dump the XML on the way through.

First - let's add our proxy as an alias:

```bash
mc alias set test http://localhost:5050 Q3AM3UQ867SPQQA43P2F zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG
```

Then let's start the proxy:

```bash
cd $CODE/sextand-dev/docs/taekion/s3-client
node proxy.hs
```

Now - we can run commands and see the XML on the way back through:

```bash
mc ls test
```

#### quick s3 iteration

```bash
cd $CODE/taekion-fs
docker build -f docker/Dockerfile --target client-base -t taekion/client-base:quick .
docker build -f docker/Dockerfile --target build -t taekion/build:quick .
```

do this for each rebuild:

```bash
cd $CODE/taekion-fs
docker build -f $CODE/sextant-dev/docs/taekion/Dockerfile.s3 --target s3 -t taekion/taekion-fs-s3:latest .
```