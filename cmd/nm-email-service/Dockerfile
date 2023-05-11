FROM golang:1.17.0-alpine as builder

MAINTAINER Bogdan Ungureanu <ubogdan@gmail.com>

RUN apk update \
    && apk add --no-cache ca-certificates git make

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -v --ldflags '-s -w' -o ./build/service github.com/ubogdan/network-manager-api/cmd/nm-email-service

FROM alpine:3.17

RUN apk update \
    && apk add --no-cache ca-certificates tzdata \
    && addgroup --gid 1000 service && adduser --disabled-password --gecos "" --ingroup service --no-create-home --uid 1000 service

WORKDIR /app

COPY --from=builder /app/build/service /app

USER service

CMD /app/service