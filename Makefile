# meta
NAME := tempshelf
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X github.com/aspick/tempshelf/cmd.version=$(VERSION) \
		   -X github.com/aspick/tempshelf/cmd.revision=$(REVISION)
COMPILE := go build -ldflags "$(LDFLAGS)" -o bin/tempshelf main.go

# setup
## setup
setup:
	go get github.com/Masterminds/glide
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/Songmu/make2help/cmd/make2help

## test
test: deps
	go test &&(glide novendor)

## install dependences
deps: setup
	glide install

## lint
lint: setup
	go vet $$(glide novendor)
	for pkg in $$(glide novendor -x); do\
		golint -set_exit_status $$pkg || exit $$?; \
	done

## format source code
fmt: setup
	goimports -w $$(glide nv -x)

## build binaries
make: main.go cmd/*.go deps
	$(COMPILE)

## setup closs compile
setupx:
	go get github.com/mitchellh/gox

## build all platform
all: make setupx
	gox -ldflags="$(LDFLAGS)" -output="bin/{{.OS}}_{{.Arch}}/{{.Dir}}" -os="linux darwin windows" -arch="386 amd64"


## cleanup binaries
clean:
	rm -rf bin


## show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps update test lint help clean
