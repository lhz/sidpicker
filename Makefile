consid=$(GOPATH)/bin/consid
usedex=$(GOPATH)/bin/consid_used_extract

default: all

all: $(consid) $(usedex)

run: $(consid)
	$^

godeps:
	go get -d ./...

$(consid): command/consid.go config/*.go ui/*.go hvsc/*.go player/*.go
	go build -o $@ $<

$(usedex): command/usedex.go config/*.go
	go build -o $@ $<
