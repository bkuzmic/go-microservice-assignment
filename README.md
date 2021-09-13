# Microservices assignment in Go
 
Solution to this assignment consists of two services:

- [Person Service](docs/person-service.md)

- [Backup Service](https://github.com/bkuzmic/go-backup-service) - in different repository

## Deployment architecture

Services are deployed to Kubernetes cluster together with Redis cluster (primary-secondary).
Person! Service has 2 instances, while Backup Service only 1.

![Deployment on Kubernetes](docs/deployment.png)

## Local deployment

### Prerequisite
Install Minikube by following instructions on: https://minikube.sigs.k8s.io/docs/start/

Start Minikube and setup kubectl utility and Helm:
```bash
minikube start
minikube kubectl -- get po -A
# create alias for kubectl
alias kubectl="minikube kubectl --"
# install Helm using installer script
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
chmod +x get_helm.sh
./get_helm.sh
# validate Helm installation
helm
```

### Building service

Using _make_ build Person service docker image:
```bash
make docker
```

TODO: describe how to deploy Redis, Person Service and then point to Backup Service repository

Steps:
```bash
# create namespace
kubectl create ns assignment
# create secret
kubectl apply -f deployment/app-secret.yaml -n assignment
kubectl -n assignment describe secret app-secret
# install Redis cluster
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install -n assignment rds bitnami/redis --values deployment/redis-values.yaml
# create config map
kubectl apply -n assignment -f deployment/app-configmap.yaml
kubectl describe -n assignment configmap/app-configmap

# NOTE!: complex setup Minikube to support local Docker registry
# Follow tutorial on: https://developers.redhat.com/blog/2019/07/11/deploying-an-internal-container-registry-with-minikube-add-ons
# Beware that tutorial's source code is not complete so one needs to manually execute these
# steps to successfully complete tutorial:
# instead of patching coredns using script one has to execute (from https://github.com/kameshsampath/minikube-helpers/tree/master/registry repo)
git clone https://github.com/kameshsampath/minikube-helpers
cd registry
kubectl apply -f registry-aliases-sa.yaml
kubectl apply -f registry-aliases-sa-crb.yaml
kubectl apply -n kube-system -f patch-coredns-job.yaml

# install Nginx Ingress
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install -n assignment ingress-nginx ingress-nginx/ingress-nginx

```

Tasks:

- [ ] Write service description and documentation in README
- [ ] Sketch deployment diagram in README
- [x] Implement update optimistic method
- [x] Implement update pessimistic method
- [ ] Implement new service in Go that will clean up expired keys
- [x] Write unit tests
- [x] Write integration tests for storage (PARTIALLY)
- [x] Create Make file for building and testing main service
- [ ] Create Make file for building and testing clean up service
- [ ] Create deployment files for deployment on Minikube
- [ ] Create manual E2E tests in Insomnia to test the service
- [x] Create Dockerfile and local running/testing using docker-compose
- [x] Add basic tests for CRU methods using Insomnia - REST API test tool
- [x] Implement connection to Redis database
- [x] Implement routes for Create, Read and Update (CRU) API calls
- [x] Create initial project on GitHub