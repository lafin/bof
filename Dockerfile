FROM lafin/armhf-alpine-golang:latest

ADD . /go/src/github.com/lafin/bof

RUN apk add --no-cache --virtual .build-deps git \
  && go get gopkg.in/mgo.v2 gopkg.in/mgo.v2/bson \
  && apk del .build-deps \
  && go install github.com/lafin/bof

ENTRYPOINT /go/bin/bof