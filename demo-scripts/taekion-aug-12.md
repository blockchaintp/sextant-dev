# taekion aug 12th demo

### start stack

Before we start - clean up any existing containers.

```bash
(cd sextant-api && git checkout tfs-client)
(cd sextant && git checkout taekion-tweaks)
kind create cluster
docker pull taekion/taekion-fs-tp:latest
docker pull taekion/taekion-fs-middleware:latest
kind load docker-image taekion/taekion-fs-tp:latest
kind load docker-image taekion/taekion-fs-middleware:latest
```

Now create 4 terminal panes so we can see multiple panes at once.

**pane #1**

```bash
export MANUALRUN=1
cd sextant-dev
make dev
cd ../sextant
make cli.frontend
npm run develop
```

**pane #2**

```bash
cd sextant
make cli.api
npm run serve
```

**pane #3**

```bash
docker network connect sextant-dev_default kind-control-plane
(cd sextant-api && bash scripts/create-remote-credentials.sh)
(cd sextant-api && bash scripts/get-remote-credentials.sh)
```

### create deployment

Now:

 * open `http://localhost`
 * setup admin user
 * login as admin user
 * create cluster
   * name=`kind`
   * api server=`https://kind-control-plane:6443`
   * token=`from output above`
   * ca=`from output above`
 * create taekion deployment
   * namespace=`taekion`
   * consensus=devmode

Then we wait for the taekion pod to sort itself out:

```bash
watch kubectl -n taekion get po
```

### demo

Create a new volume called `apples` (don't use encryption).

Create a snapshot of that volume called `test`

Open the `CLI` tab of the settings.

Paste (and run) the docker command in **pane #3**

```bash
tfs-cli volume list
```

Demonstrate that we are wrapping the CLI:

```bash
DEBUG=1 tfs-cli volume list
```

Mount the volume

```bash
tfs-fuse -v apples -m /mnt/tfs --non-empty --nodaemon --debug
```

Write a file:

**pane #4**

```bash
docker exec -ti tfs-cli bash
echo hello > /mnt/tfs/world.txt
ls -la /mnt/tfs
exit
```

Quit the fuse process in **pane #3** and demonstrate that the container is gone:

```bash
ctrl+c
exit
docker ps -a | grep tfs-cli
```

explain that this means the files are gone (but ahah not actually because blockchain):

Paste (and run) the docker command in **pane #3**

```bash
DEBUG=1 tfs-fuse -v apples -m /mnt/tfs --non-empty --nodaemon --debug
```

**pane #4**

```bash
docker exec -ti tfs-cli bash
ls -la /mnt/tfs
exit
```

We just (re-read) our files from the blockchain.

Now let's create a volume from the cli:

**pane #3**

```bash
ctrl+c
tfs-cli volume create oranges
tfs-cli volume list
```

Login to the web UI and show that oranges is now on the screen.
