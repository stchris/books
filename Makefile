
all: deps fmt test build

deps:
	go get github.com/golang/lint/golint
	go get honnef.co/go/tools/cmd/megacheck

build:
	go get
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
