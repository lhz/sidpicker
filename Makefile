sidpicker=$(GOPATH)/bin/sidpicker
extract=$(GOPATH)/bin/sid_release_extract

default: all

all: $(sidpicker) $(extract)

$(sidpicker): cmd/sidpicker.go config/*.go ui/*.go hvsc/*.go player/*.go csdb/*.go
	go build -o $@ $<

$(extract): cmd/sid_release_extract.go config/*.go
	go build -o $@ $<

godeps:
	go get -d ./...
