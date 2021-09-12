build:
	go build

test:
	go test ./...

integration-test:
	export REDIS_PASSWORD=test && docker-compose up redis > /dev/null &
	sleep 10 && echo "Waiting for containers to start"
	go test ./... -tags=integration
	docker-compose down

docker:
	docker build -t bkuzmic/go-person-service .

run-docker:
	export REDIS_PASSWORD=test && docker-compose up