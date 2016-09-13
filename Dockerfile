FROM lafin/alpine-golang:armhf-1.7.1

ADD . /go/src/github.com/lafin/bof

RUN apk add --no-cache --virtual .build-deps git \
  && go get gopkg.in/mgo.v2 gopkg.in/mgo.v2/bson \
  && apk del .build-deps \
  && go install github.com/lafin/bof

ENTRYPOINT /go/bin/bof