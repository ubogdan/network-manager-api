
GOPATH=$(shell go env GOPATH)
GOFLAGS=--trimpath --tags "netgo" -ldflags "-s -w"

test:
	go test -mod=readonly -v ./...

build:
	go build $(GOFLAGS)  -o bin/api github.com/ubogdan/network-manager/cmd/api

fmt:
	go fmt ./...

assets:
	go install github.com/rakyll/statik
	${GOPATH}/bin/statik -f -src=web/dist -include=*.html,*.css,*.js,*.ico,*.eot,*.ttf,*woff2,*.svg,*.jpg,*.gif,*.png -p="dist" -dest=web/

swagger:
	swag init -g cmd/api/main.go

