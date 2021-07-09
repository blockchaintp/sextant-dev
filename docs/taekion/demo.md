## taekion demo cluster

 * url: http://a2e52f94a103a40e6acff3485160eb02-1935463034.us-east-1.elb.amazonaws.com/
 * username: admin
 * password: Aqma2Ynul8L9P5oY7im5hKlr02

To deploy new TFS deployment:

 * first undeploy & delete the existing deployment
 * create new TFS deployment with these options:

  * name: tfs
  * namespace: tfs (important the image pull secrets are here)
  * consensus: PBFT
  * image pull secrets:
    * btp-dev
    * btl-lic
  * advanced -> advanced deployment customization:

```
images:
  taekion_rest_api: dev.catenasys.com:8083/test/taekion/middleware/middleware:6c6a0d8bdfe283cde9e36d465286713d89cb8627
  taekion_tp: dev.catenasys.com:8083/test/taekion/fs-tp/fs-tp:6c6a0d8bdfe283cde9e36d465286713d89cb8627
```



----------------------------------------------------

 
 * get access to image repo
 
 * eksctl create cluster
 * create docker registry secret with nexus credentials - see Slack for deets
 * use helm to deploy sextant (sextant-enterprise chart)
   * override images from branch
   * give secret name to docker pull secret
 * kubectl port-forward to get the UI
 * run the create-credentials.sh scripts to get deets
 * use the UI to add the cluster using deets
 * use sextant to deploy a "nginx ingress controller"
 * UI will show the load balancer URL
 * then helm apply with the `ingress.enabled=true` + `ingress.hosts = [<url from UI>]`


images:

 * dev.catenasys.com:8083/base/catenasys/sextant/sextant-sft:2.1.1-156-g512c409
 * dev.catenasys.com:8083/base/catenasys/sextant-api/sextant-api-sft:2.2.1-102-g27e9366
 * dev.catenasys.com:8083/test/taekion/middleware/middleware:6c6a0d8bdfe283cde9e36d465286713d89cb8627
 * dev.catenasys.com:8083/test/taekion/fs-tp/fs-tp:6c6a0d8bdfe283cde9e36d465286713d89cb8627


```bash
eksctl create cluster \
  --name taekion-demo \
  --region us-east-1 \
  --node-type m5.large \
  --nodes 4
```

```bash
export NEXUS_EMAIL=XXX
export NEXUS_USERNAME=XXX
export NEXUS_PASSWORD=XXX
export DOCKERHUB_EMAIL=XXX
export DOCKERHUB_USERNAME=XXX
export DOCKERHUB_PASSWORD=XXX
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add sextant https://btp-charts-stable.s3.amazonaws.com/charts/
kubectl create secret docker-registry btp-dev \
  --docker-server=https://dev.catenasys.com:8083/ \
  --docker-username=$NEXUS_USERNAME \
  --docker-password=$NEXUS_PASSWORD \
  --docker-email=$NEXUS_EMAIL
kubectl create secret docker-registry btp-lic \
  --docker-server=https://dev.catenasys.com:8084/ \
  --docker-username=$NEXUS_USERNAME \
  --docker-password=$NEXUS_PASSWORD \
  --docker-email=$NEXUS_EMAIL
helm install taekion-demo sextant/sextant-enterprise -f values.yaml
```

Then we install nginx-ingress controller from sextant - get the URL - update `values.yaml`

```
hosts:
      - host: "a3ba3881a46804a8cb7e905f9eb22631-735769867.us-east-1.elb.amazonaws.com"
        paths:
          - "/"
```

```bash
helm upgrade -f values.yaml taekion-demo sextant/sextant-enterprise
```


Now create a new namespace with secrets:

```bash
kubectl create ns tfs
kubectl create -n tfs secret docker-registry btp-dev \
  --docker-server=https://dev.catenasys.com:8083/ \
  --docker-username=$NEXUS_USERNAME \
  --docker-password=$NEXUS_PASSWORD \
  --docker-email=$NEXUS_EMAIL
kubectl create -n tfs secret docker-registry btp-lic \
  --docker-server=https://dev.catenasys.com:8084/ \
  --docker-username=$NEXUS_USERNAME \
  --docker-password=$NEXUS_PASSWORD \
  --docker-email=$NEXUS_EMAIL
kubectl create -n tfs secret docker-registry dockerhub \
  --docker-username=$DOCKERHUB_USERNAME \
  --docker-password=$DOCKERHUB_PASSWORD \
  --docker-email=$DOCKERHUB_EMAIL
```

Now use sextant to deploy TFS - put this in YAML override:

```
images:
  taekion_rest_api: dev.catenasys.com:8083/test/taekion/middleware/middleware:6c6a0d8bdfe283cde9e36d465286713d89cb8627
  taekion_tp: dev.catenasys.com:8083/test/taekion/fs-tp/fs-tp:6c6a0d8bdfe283cde9e36d465286713d89cb8627
```

```bash
helm uninstall taekion-demo
```

