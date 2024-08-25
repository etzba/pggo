TAG ?= latest
REPO ?= etzba/pggo

all: go-test 

# unit tests
test:
	go test -v ./...

run:
	go run main.go

# build image and push to dockerhub
.PHONY: docker-build
docker-build:
	docker build -t ${REPO}:${TAG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${REPO}:${TAG}