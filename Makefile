
GOPATH=$(shell go env GOPATH)
GOFLAGS=--trimpath --tags "netgo" -ldflags "-s -w"

test:
	go test -mod=readonly -v ./...

fmt:
	go fmt ./...

client=747256501865
region=eu-central-1
version=latest

login:
	aws ecr get-login-password --region $(region) | docker login --username AWS --password-stdin $(client).dkr.ecr.$(region).amazonaws.com

%-service:
	docker build -t $(client).dkr.ecr.$(region).amazonaws.com/$@:$(version)-arm64 --build-arg ARCH=arm64v --build-arg GOARCH=arm64 -f cmd/$@/Dockerfile .
	docker push $(client).dkr.ecr.$(region).amazonaws.com/$@:$(version)-arm64
	docker build -t $(client).dkr.ecr.$(region).amazonaws.com/$@:$(version) --build-arg ARCH=amd64 --build-arg GOARCH=amd64 -f cmd/$@/Dockerfile .
	docker push $(client).dkr.ecr.$(region).amazonaws.com/$@:$(version)
