FROM golang:1.15 AS build

COPY . .

RUN make build

FROM ubuntu:18.04

RUN apt-get update && \
    apt-get install -y ca-certificates && update-ca-certificates && \
    apt-get -s dist-upgrade | grep '^Insta' | grep -i 'security' | awk -F ' ' '{print $2}' | xargs apt-get install -y

EXPOSE 8080

RUN useradd -m -u 1000 service
USER service

COPY --from=build /go/bin/api /bin/api

CMD /bin/api