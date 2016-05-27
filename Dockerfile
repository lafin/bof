FROM golang

ADD . /go/src/github.com/lafin/bof
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN go install github.com/lafin/bof
ENTRYPOINT /go/bin/bof
