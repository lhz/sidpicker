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

binaries: cmd/sidpicker.go config/*.go ui/*.go hvsc/*.go player/*.go csdb/*.go
	for os in linux darwin windows ; do \
	    for arch in 386 amd64 ; do \
	        mkdir -p binaries/$$os/$$arch ; \
		if [ "$$os" = "windows" ] ; then \
		    filename=sidpicker.exe ; \
		else \
		    filename=sidpicker ; \
		fi ; \
		echo "Building for $$os/$$arch" ; \
	        GOOS=$$os GOARCH=$$arch go build -o binaries/$$os/$$arch/$$filename $< ; \
	    done ; \
	done

godeps:
	go get -d ./...

.PHONY: binaries
