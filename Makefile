sidpicker=       $(GOPATH)/bin/sidpicker
release_check=   $(GOPATH)/bin/sidpicker_release_check
release_extract= $(GOPATH)/bin/sidpicker_release_extract

default: all

all: $(sidpicker) $(release_check) $(release_extract)

$(sidpicker): cmd/sidpicker.go config/*.go ui/*.go hvsc/*.go player/*.go csdb/*.go
	go build -o $@ $<

$(release_check): cmd/release_check.go config/*.go
	go build -o $@ $<

$(release_extract): cmd/release_extract.go config/*.go
	go build -o $@ $<

godeps:
	go get -d ./...
