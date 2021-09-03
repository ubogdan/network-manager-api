FROM golang:1.17.0-alpine as builder

MAINTAINER Bogdan Ungureanu <ubogdan@gmail.com>

RUN apk update \
    && apk add --no-cache ca-certificates git make

WORKDIR /lambda

COPY . .

RUN CGO_ENABLED=0 go build -v -o ./build/service github.com/ubogdan/network-manager-api/cmd/nm-api-service

FROM alpine:3.13

RUN apk update \
    && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /lambda/build/service /app

CMD /app/service