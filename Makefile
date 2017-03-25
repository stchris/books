
all: deps fmt build test vet golint

deps:
	go get -u github.com/golang/lint/golint

build:
	go get
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
