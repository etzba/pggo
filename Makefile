TAG ?= latest
REPO ?= etzba/pggo

all: test up 

# unit tests
test:
	go test -v ./...

run:
	go run main.go

start:
	sh scripts/start.sh

up:
	docker-compose down
	sleep 3
	docker-compose up -d pg
	sleep 8
	docker-compose up -d pggo 

# prepare docker for cli tests
docker-cleanup: cleanup-api cleanup-pg

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
