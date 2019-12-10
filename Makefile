# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
SERVICE = ledger
BUILD_TIME=`date +%FT%T%z`
VERSION := $(shell sh -c 'git describe --always --tags')
BRANCH := $(shell sh -c 'git rev-parse --abbrev-ref HEAD')
COMMIT := $(shell sh -c 'git rev-parse --short HEAD')
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.branch=$(BRANCH) -X main.buildDate=$(BUILD_TIME)"
BUILD_TAGS=-tags aws_lambda
LINT_TOOL=$(shell go env GOPATH)/bin/golangci-lint
GO_PKGS=$(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/)
GO_FILES=$(shell find . -type f -name '*.go' -not -path './vendor/*')

setup_dev:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get -u github.com/golang/dep/cmd/dep	
	sudo curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/download/v1.7.0/dbmate-linux-amd64
	sudo chmod +x /usr/local/bin/dbmate

setup_deploy:
	yarn install --frozen-lockfile

setup: setup_dev setup_deploy

clean:
	rm -rf ./gen ./bin

validate:
	swagger validate swagger/$(SERVICE).yaml

swagger: clean
	mkdir gen
	swagger -q generate server -t gen -f swagger/$(SERVICE).yaml --exclude-main -A $(SERVICE) --keep-spec-order

up:
	dbmate up
	
test: up 
	go test -p 1 -race -cover ./...

run:
	go run main.go

deps:
	dep ensure -v

build: swagger deps lint
	env GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(SERVICE)
	chmod +x bin/$(SERVICE)

$(LINT_TOOL):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.19.0

qc: $(LINT_TOOL)
	$(LINT_TOOL) run --config=.golangci.yaml ./...

lint: qc

rebuild: swagger qc build run

.PHONY: setup clean qc swagger up build test run setup_deploy setup_dev

deploy: clean build
	sls deploy --verbose
