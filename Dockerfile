FROM golang:1.15

COPY . .

RUN make build

FROM ubuntu:18.04

RUN apt-get update && \
    apt-get -s dist-upgrade | grep '^Insta' | grep -i 'security' | awk -F ' ' '{print $2}' | xargs apt-get install -y

RUN useradd -m -u 1000 service
USER service

COPY bin/api /bin/api

CMD /bin/api