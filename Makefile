
all: deps fmt test build

deps:
	go get github.com/golang/lint/golint
	go get honnef.co/go/tools/cmd/megacheck
	go get -u github.com/kardianos/govendor
	govendor sync

build:
	go build

test: lint
	go test

lint:
	golint
	megacheck

fmt:
	go fmt

clean:
	rm -f books
