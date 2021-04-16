## test-stack

### install

```bash
# start kind with ingress enabled & local registry
bash start.sh
# install nginx ingress
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
# build local images for sextant
docker build -t localhost:5000/sextant:local -f ../../../sextant/Dockerfile.multistage ../../../sextant
docker build -t localhost:5000/sextant-api:local ../../../sextant-api
docker push localhost:5000/sextant:local
docker push localhost:5000/sextant-api:local
```

### test

```bash
# deploy
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

