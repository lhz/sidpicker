binary=bin/consid

default: run

godeps:
	(cd src && go get -d ./...)

$(binary): src/*.go
	mkdir -p bin
	go build -o $@ $^

run: $(binary)
	@$(binary)
