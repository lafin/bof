SHA := $(shell git rev-parse --short HEAD)
DATE := $(shell TZ=UTC date +%FT%T)Z
VERSION := $(shell cat VERSION)-$(DATE)-$(SHA)

RELEASES=release/imgdiff-darwin-amd64 \
	 release/imgdiff-darwin-386 \
	 release/imgdiff-linux-amd64 \
	 release/imgdiff-linux-386 \
	 release/imgdiff-windows-amd64.exe \
	 release/imgdiff-windows-386.exe

SRCS=./cmd/imgdiff/*.go
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

imgdiff: $(SRCS)
	go build -o imgdiff $(LDFLAGS) $(SRCS)

prof:
	go test -bench . -cpuprofile cpu.prof
	go tool pprof -top imgdiff.test cpu.prof

release: $(RELEASES)

release/imgdiff-darwin-amd64: $(SRCS)
	GOOS=darwin GOARCH=amd64 go build -o $@ $(LDFLAGS) $(SRCS)

release/imgdiff-darwin-386: $(SRCS)
	GOOS=darwin GOARCH=386 go build -o $@ $(LDFLAGS) $(SRCS)

release/imgdiff-linux-amd64: $(SRCS)
	GOOS=linux GOARCH=amd64 go build -o $@ $(LDFLAGS) $(SRCS)

release/imgdiff-linux-386: $(SRCS)
	GOOS=linux GOARCH=386 go build -o $@ $(LDFLAGS) $(SRCS)

release/imgdiff-windows-amd64.exe: $(SRCS)
	GOOS=windows GOARCH=amd64 go build -o $@ $(LDFLAGS) $(SRCS)

release/imgdiff-windows-386.exe: $(SRCS)
	GOOS=windows GOARCH=386 go build -o $@ $(LDFLAGS) $(SRCS)

clean:
	rm -rf release imgdiff *.prof

default: imgdiff
