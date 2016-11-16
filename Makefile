default: all

all: consid

run: consid
	$(GOPATH)/bin/consid

godeps:
	(cd src && go get -d ./...)

consid: $(GOPATH)/bin/consid

$(GOPATH)/bin/consid: consid/consid.go
	go build -o $@ $^
