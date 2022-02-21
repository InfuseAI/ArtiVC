VERSION ?=
LDFLAGS =

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

LDFLAGS += -X github.com/infuseai/art/cmd.tagVersion=${VERSION}
LDFLAGS += -X github.com/infuseai/art/cmd.gitCommit=${GIT_COMMIT}
LDFLAGS += -X github.com/infuseai/art/cmd.gitTreeState=${GIT_DIRTY}
LDFLAGS += $(EXT_LDFLAGS)


build:
	mkdir -p bin
	GOOS='$(GOOS)' GOARCH='$(GOARCH)' go build -o bin/art -ldflags '$(LDFLAGS)' main.go

test:
	go test ./...

build_for_release:
	mkdir -p release
	GOOS='$(GOOS)' GOARCH='$(GOARCH)' go build -o release/art-$(GOOS)-$(GOARCH) -ldflags '$(LDFLAGS)' main.go

clean:
	rm -rf release

release: clean
	make build_for_release GOOS=linux GOARCH=amd64
	make build_for_release GOOS=linux GOARCH=arm64
	make build_for_release GOOS=darwin GOARCH=amd64
	make build_for_release GOOS=darwin GOARCH=arm64