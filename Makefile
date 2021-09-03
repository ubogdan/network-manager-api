
GOPATH=$(shell go env GOPATH)
GOFLAGS=--trimpath --tags "netgo" -ldflags "-s -w"

test:
	go test -mod=readonly -v ./...

build:
	go build $(GOFLAGS)  -o bin/api github.com/ubogdan/network-manager-api/cmd/api

fmt:
	go fmt ./...

swagger:
	swag init -g cmd/api/main.go

version=latest

%-service:
	docker build -t 747256501865.dkr.ecr.eu-central-1.amazonaws.com/$@:$(version) -f cmd/$@/Dockerfile .
	docker push 747256501865.dkr.ecr.eu-central-1.amazonaws.com/$@:$(version)
