consid=$(GOPATH)/bin/consid
usedex=$(GOPATH)/bin/consid_used_extract

default: all

all: $(consid) $(usedex)

run: $(consid)
	$^

godeps:
	go get -d ./...

$(consid): cmd/consid.go config/*.go ui/*.go hvsc/*.go player/*.go
	go build -o $@ $<

$(usedex): cmd/usedex.go config/*.go
	go build -o $@ $<
