## taekion install guide

This guide will show how to:

 * mount a TFS volume into a docker-compose stack
 * connect to the TFS S3 endpoint within your application

Pre-requisites:

 * a running sextant cluster
 * a TFS deployment running on that sextant cluster
 * the S3 ingress URL for the TFS deployment

### docker-compose

Once you have a running sextant cluster (with a TFS deployment) - you can visit the `TFS CLI` page of the deployment settings.

You will need to copy the `TFS_URL` environment variable and export it into your environment.

```bash
export TFS_URL=<copy_pasted_from_sextant_ui>
```

Then - create a volume that you will mount into your docker-compose stack (e.g. a volume named `test`).

Export a variable with the name of that volume:

```bash
export VOLUME_NAME=test
```

Finally we create a folder on our host system that we will use to mount the TFS volume:

```bash
sudo mkdir /tfs_mount
```

Once we have these variables and a folder to mount our volume into - we can use the following docker-compose file to boot our stack.

In this example - there are two services:

 * `taekion` - this container runs alongside the rest of your stack and will look after connecting to the TFS api and mounting the volume locally
 * `app` - this is your application container that will want to consume the volume as part of it's filesystem

The `app` service is just `ubuntu` - apply the `volumes` section to the containers in your stack that need access to the volume:

```yaml
version: '3.2'
services:
  taekion:
    image: <TFS_IMAGE_HERE>
    privileged: true
    init: true
    entrypoint: bash -c 'mkdir -p /home/user/tfs_mount/data && tfs-fuse --allow-other -v $VOLUME_NAME -m /home/user/tfs_mount/data && tail -f /dev/null'
    devices:
      - "/dev/fuse"
    cap_add:
      - SYS_ADMIN
    environment:
      - TFS_URL
      - VOLUME_NAME
    volumes:
      - type: bind
        source: /tfs_mount
        target: /home/user/tfs_mount
        bind:
          propagation: rshared
  app:
    image: ubuntu
    entrypoint: bash -c 'tail -f /dev/null'
    depends_on: 
      - taekion
    volumes:
      - type: bind
        source: /tfs_mount
        target: /var/lib/tfs
        bind:
          propagation: rshared
```

### S3

To connect to the S3 endpoint - use the S3 ingress URL we copied from the sextant deployment as the S3 hostname.

You can use any string for the AccessToken - it will be ignored for the demo installation.

The operations that the S3 server will support:

 * list buckets (each volume shows up as a bucket)
 * list objects in a bucket
 * download an object
 * upload an object

For example - if we have a volume called `test` and we are using the Minio CLI connected to our endpoint called `tfs` - we can list the objects in that bucket:

```bash
mc ls tfs/test/
data/
snapshots/
```

There are two top level keys:

 * `data` - this is the root of the current volume - i.e. the latest content
 * `snapshots` - this holds each of the snapshots you have taken of the volume

The main difference between the `data` and `snapshots` keys is that you cannot write to a path in the `snapshots` key.

