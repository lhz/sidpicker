sidpicker=$(GOPATH)/bin/sidpicker
sidusedex=$(GOPATH)/bin/sidusedex

default: all

all: $(sidpicker) $(sidusedex)

$(sidpicker): cmd/sidpicker.go config/*.go ui/*.go hvsc/*.go player/*.go csdb/*.go
	go build -o $@ $<

$(sidusedex): cmd/releases.go config/*.go
	go build -o $@ $<

godeps:
	go get -d ./...
