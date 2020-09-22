export VERSION := $(shell git show -s --format=%h)
export BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
export GOPATH := $(shell pwd)
export PATH := $(shell pwd)/bin:$(PATH)
import = projecto

all: .generate
	cd src && go build -v -ldflags "-X $(import)/app.version=$(VERSION)-$(BRANCH) -X $(import)/app.name=$(import)" -o ../bin ./...

test: .generate
	cd src && go test -coverprofile=coverage.txt -covermode=atomic ./...

lint: .generate
	golangci-lint run -v ./...

clean:
	rm -rf pkg/* bin/*
	cd src && find ./ -name '*.gen.go' -delete

.generate:
	cd src && go generate ./...

