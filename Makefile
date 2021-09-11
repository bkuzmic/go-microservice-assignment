build:
	go build

test:
	go test ./...

integration-test:
	go test ./... -tags=integration

docker:
	docker build -t bkuzmic/go-person-service .

run-docker:
	export REDIS_PASSWORD=test && docker-compose up