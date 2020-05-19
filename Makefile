SOURCE=./cmd/parser-bank
APP=parser-bank
VERSION=1.0
ARCH= $(shell uname -m)
GOBASE=$(shell pwd)
RELEASE_DIR=$(GOBASE)/bin

.DEFAULT_GOAL = build 

GO_SRC_DIRS := $(shell \
	find . -name "*.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)
GO_TEST_DIRS := $(shell \
	find . -name "*_test.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)	

build: 
	@go build -v -o ${APP} ${SOURCE}

lint:
	@goimports -w ${GO_SRC_DIRS}	
	@gofmt -s -w ${GO_SRC_DIRS}
	@golint ${GO_SRC_DIRS}
	@golangci-lint run

test:
	go test -v ${GO_TEST_DIRS}

mod:
	go mod verify
	go mod vendor
	go mod tidy

run:
	@go run ${SOURCE} --config=configs/prod.yaml

.PHONY: dockerbuild
dockerbuild: 
	docker build -t puzanovma/${APP} -f ./deployments/parser-bank.Dockerfile .

.PHONY: up
up: # 
	docker-compose -f ./deployments/docker-compose.yaml up --build --remove-orphans

.PHONY: down
down: 
	docker-compose -f ./deployments/docker-compose.yaml down -v
	
.PHONY: push
push:
	docker push puzanovma/${APP}

.PHONY: build run release lint test mod