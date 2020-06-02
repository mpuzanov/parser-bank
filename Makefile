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
	@gofmt -s -w -d ${GO_SRC_DIRS}
	@golint ${GO_SRC_DIRS}
	@go vet ${GO_SRC_DIRS}
	@#golangci-lint run

test:
	go test -v ${GO_TEST_DIRS}
	@#go test -race -count 100 ${GO_TEST_DIRS}

bench:
	go test -benchmem -bench=. ${GO_TEST_DIRS}

mod:
	go mod verify
	#go mod vendor
	go mod tidy

run:
	@go run ${SOURCE} web_server --config=configs/prod.yaml

.PHONY: run-cli
run-cli:
	@#go run ${SOURCE} shell --path=test_data
	@go run ${SOURCE} shell --path=g:\\payments_test -d=false
	@#go run ${SOURCE} shell --path='g:/payments_test/new folder/CER18B210420.txt' -d
	@#GOMAXPROCS=1 go run ${SOURCE} shell --path=test_data/0_mupspdu_0925.txt -d

.PHONY: dockerbuild
dockerbuild: 
	docker build -t puzanovma/${APP} -f ./deployments/parser-bank.Dockerfile .

.PHONY: up
up: # 
	docker-compose -f ./deployments/docker-compose.yaml up --build

.PHONY: down
down: 
	docker-compose -f ./deployments/docker-compose.yaml down -v
	
.PHONY: push
push:
	docker push puzanovma/${APP}

.PHONY: build run release lint test mod