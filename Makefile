
all: build test vet golint

build:
	go build

test:
	go test

vet:
	go vet

golint:
	golint
