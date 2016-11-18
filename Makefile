consid=$(GOPATH)/bin/consid

default: all

all: $(consid)

run: $(consid)
	$^

godeps:
	(cd src && go get -d ./...)

$(consid): command/consid.go config/config.go hvsc/hvsc.go
	go build -o $@ $<
