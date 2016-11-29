consid=$(GOPATH)/bin/consid

default: all

all: $(consid)

run: $(consid)
	$^

godeps:
	go get -d ./...

$(consid): command/consid.go config/*.go ui/*.go hvsc/*.go player/*.go
	go build -o $@ $<
