FROM lafin/armhf-alpine-golang

ADD . /go/src/github.com/lafin/bof

RUN apk add --update --no-cache --virtual .build-deps git
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN apk del .build-deps

RUN go install github.com/lafin/bof

ENTRYPOINT /go/bin/bof