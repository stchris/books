
all: deps fmt build test vet golint

deps:
	go get github.com/golang/lint/golint
	go get -u github.com/kardianos/govendor

build:
	govendor sync
	go build

test:
	go test

vet:
	go vet

golint:
	golint

fmt:
	go fmt

clean:
	rm -f books
