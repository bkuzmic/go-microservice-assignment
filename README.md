# Microservices assignment in Go
 
Solution to this assignment consists of two services:

- [Person Service](docs/person-service.md)

- [Backup Service](docs/backup-service.md) - in different repository

## Deployment architecture

TODO

## Local deployment

### Prerequisite
Install Minikube by following instructions on: https://minikube.sigs.k8s.io/docs/start/

Start Minikube and setup kubectl utility:
```bash
minikube start
minikube kubectl -- get po -A
# create alias for kubectl
alias kubectl="minikube kubectl --"
```

### Building service

Using _make_ build Person service docker image:
```bash
make docker
```

TODO: describe how to deploy Redis, Person Service and then point to Backup Service repository

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
