VERSION =
LDFLAGS =

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

ifeq ($(VERSION),)
VERSION := $(shell echo $${GITHUB_REF_NAME})
endif

LDFLAGS += -X github.com/infuseai/artivc/cmd.tagVersion=${VERSION}
LDFLAGS += -X github.com/infuseai/artivc/cmd.gitCommit=${GIT_COMMIT}
LDFLAGS += -X github.com/infuseai/artivc/cmd.gitTreeState=${GIT_DIRTY}
LDFLAGS += $(EXT_LDFLAGS)


build:
	mkdir -p bin
	go build -o bin/avc -ldflags '$(LDFLAGS)' main.go

test:
	go test ./...

integration-test:
	go test -v ./internal/repository

.PHONY: doc-server
doc-server:
	make -C docs/ start
