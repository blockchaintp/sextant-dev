## test-stack

### install

```bash
# start kind with ingress enabled
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF
# install nginx ingress
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
# build local images for sextant
docker build -t sextant:local -f ../../../sextant/Dockerfile.multistage ../../../sextant
docker build -t sextant-api:local ../../../sextant-api
kind load docker-image sextant:local
kind load docker-image sextant-api:local
```

### test

```bash
kubectl create ns ingress-test
kubectl apply -f ingress-test/echo-deployment.yaml
kubectl apply -f ingress-test/echo-service.yaml
kubectl apply -f ingress-test/nginx-deployment.yaml
kubectl apply -f ingress-test/nginx-service.yaml
kubectl apply -f ingress-test/ingress.yaml
```

This will deploy:
 * `/apples` -> nginx
 * `/apples/api/v1` -> echo

The rewrite rules are:
 * `/apples(/|$)(.*)` -> `/apples/index.css` -> nginx => $2 == `/index.css`
 * `/apples(/|$)(/api/v1/.*)` -> `/apples/api/v1/config` -> echo => $2 == `/api/v1/config`

We use `nginx.ingress.kubernetes.io/rewrite-target: /$2` to remove the `/apples` path

We use `nginx.ingress.kubernetes.io/x-forwarded-prefix` to tell the back we have removed the `/apples` path (it will have a `x-forwarded-prefix` header)

To make changes to the ingress:

```bash
kubectl apply -f ingress-test/ingress.yaml
```

Remove:

```bash
kubectl delete ns ingress-test
```

### sextant

```bash
kubectl create ns sextant
kubectl apply -f sextant/postgres-deployment.yaml
kubectl apply -f sextant/postgres-service.yaml
kubectl apply -f sextant/api-deployment.yaml
kubectl apply -f sextant/api-service.yaml
kubectl apply -f sextant/frontend-deployment.yaml
kubectl apply -f sextant/frontend-service.yaml
kubectl apply -f sextant/ingress.yaml
```

To then quickly reload the frontend:

```bash
function redeploy-frontend() {
  export TAG=$(date +%s)
  docker build -t sextant:$TAG -f ../../../sextant/Dockerfile.multistage ../../../sextant
  kind load docker-image sextant:$TAG
  kubectl -n sextant set image deployment/frontend frontend=sextant:$TAG
}
redeploy-frontend
```