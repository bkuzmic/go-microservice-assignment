build:
	go build

test:
	go test -v ./...

integration-test: redis-up
	@sleep 5 && echo "Waiting for container to start"
	go test -v ./... -tags=integration
	@docker-compose down

redis-up:
	@export REDIS_PASSWORD=test && docker-compose up -d redis

docker:
	docker build -t bkuzmic/go-person-service .

run-local:
	export REDIS_PASSWORD=test && docker-compose up