
all: fmt build test vet golint

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
