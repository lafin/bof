FROM lafin/alpine-golang:platform

ADD . /go/src/github.com/lafin/bof
RUN go install github.com/lafin/bof \
  && rm -rf /go/{src,pkg} /usr/local/go

ENTRYPOINT /go/bin/bof