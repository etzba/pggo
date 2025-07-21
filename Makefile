TAG ?= latest
REPO ?= etzba/pggo

all: test up exec down

# unit tests
test:
	go test -v ./...

# TODO: set test from golang client
test-e2e:
	echo e2e
	
run:
	go run main.go

lint:
	golangci-lint run ./...

# local development
start:
	sh scripts/start.sh

# test with etzba 
exec:
	etz --config=etzba/config.yaml
	etz api --exec=etzba/locations.yaml -d=3s -w=2
	etz api --exec=etzba/locations.yaml -d=3s -w=4 -r=12 --output=etzba/results/$$(date +%Y%m%d_%H%M%S)_result.json
	etz api --exec=etzba/locations.yaml -d=3s -w=6 -r=24 --output=etzba/results/$$(date +%Y%m%d_%H%M%S)_result.json

# docker
up:
	docker-compose down
	sleep 3
	docker-compose up -d pg
	sleep 8
	docker-compose up -d pggo 

down:
	docker-compose down

# cleanup running docker containers
cleanup: cleanup-api cleanup-pg

cleanup-api:
	docker rm $$(docker stop $$(docker ps -a -q --filter ancestor=etzba/pggo:latest --format="{{.ID}}"))

cleanup-pg:
	docker rm $$(docker stop $$(docker ps -a -q --filter ancestor=postgres:14 --format="{{.ID}}"))

# build image and push to dockerhub
.PHONY: docker-build
docker-build:
	docker build -t ${REPO}:${TAG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${REPO}:${TAG}

helm:
	helm install pggo chart/ -n pggo --create-namespace
