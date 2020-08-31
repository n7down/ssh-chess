VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
MAKEFLAGS += --silent
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

.PHONY: get 
get:
	echo "getting go dependencies..."
	@go get ./...
	echo "done"

.PHONY: generate
generate:
	echo "generating dependency files..."
	go generate ./...
	echo "done"

.PHONY: test-unit
test-unit:
	echo "running unit tests..."
	go test --tags unit -v ./...
	echo "done"

.PHONY: cover-unit
cover-unit:
	go test --tags unit -v ./... -coverprofile c.out; go tool cover -func c.out

.PHONY: cover-unit-html
cover-unit-html:
	go test --tags unit -v ./... -coverprofile c.out; go tool cover -html c.out

.PHONY: lint
lint:
	golint ./...

.PHONY: stop
stop:
	echo "stopping docker containers..."
	docker-compose stop
	echo "done"

.PHONY: rm
rm:
	echo "removing docker containers..."
	docker-compose rm
	echo "done"

.PHONY: clean
clean: stop rm

.PHONY: build-ssh-chess
build-ssh-chess:
	echo "building ssh-chess..."
	docker build -t "$(PROJECTNAME)"/ssh-chess:"$(VERSION)" --label "version"="$(VERSION)" --label "build"="$(BUILD)" -f build/dockerfiles/ssh-chess/Dockerfile .
	echo "done"

.PHONY: up
up: 
	docker-compose up -d "$(PROJECTNAME)"/ssh-chess:"$(VERSION)"

.PHONY: help
help:
	echo "Choose a command run in $(PROJECTNAME):"
	echo " - get: get all dependencies"
	echo " - geneate: generate dependencies"
	echo " - test-unit: run unit tests"
	echo " - cover-unit: run code coverage of unit tests"
	echo " - cover-unit-html: show html document for code coverage of unit tests"
	echo " - lint: run go lint"
	echo " - stop: stop ssh-chess docker containers"
	echo " - rm: remove ssh-chess docker containers"
	echo " - clean: stop and remove ssh-chess docker containers"
	echo " - build-ssh-chess: build ssh-chess docker container"
	echo " - up: runs ssh-chess"
