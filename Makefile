.EXPORT_ALL_VARIABLES:


# Allow environment variable override for running targets within Dockerfile
VERSION ?= $(shell git describe --always --dirty --tags | egrep -o '([0-9]+\.){1,2}[0-9]?[a-zA-Z0-9-]+')
BUILD := $(shell date +%FT%T%z)
# [x] fetching git tag here
LASTTAG := $(shell git fetch --tags && git describe --abbrev=0 --tags)
BUILD_FLAGS :=

LDFLAGS := -ldflags "-w -s -X main.version=${VERSION} -extldflags '-static'"
BINARY := password-${VERSION}

test:
	@go test -race ./...

install:
	@go install ${BUILD_FLAGS} ${LDFLAGS}

binary: out/bin/${BINARY}.gz

out/bin/${BINARY}.gz: out/bin/${BINARY}-amd64 out/bin/${BINARY}-arm64 out/bin/${BINARY}-darwin
	@tar --strip-components 2 -czvf out/bin/${BINARY}.tar.gz out/bin/${BINARY}-amd64 out/bin/${BINARY}-arm64 out/bin/${BINARY}-darwin

out/bin/${BINARY}-amd64:
	@mkdir -p ./out/bin
	@GOOS=linux GOARCH=amd64 go build ${BUILD_FLAGS} ${LDFLAGS} -o out/bin/${BINARY}-amd64

out/bin/${BINARY}-arm64:
	@mkdir -p ./out/bin
	@GOOS=linux GOARCH=arm64 go build ${BUILD_FLAGS} ${LDFLAGS} -o out/bin/${BINARY}-arm64

out/bin/${BINARY}-darwin:
	@mkdir -p ./out/bin
	@GOOS=darwin GOARCH=amd64 go build ${BUILD_FLAGS} ${LDFLAGS} -o out/bin/${BINARY}-darwin