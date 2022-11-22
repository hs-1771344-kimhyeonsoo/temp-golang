BINARY=movie-backend

GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
GOBUILD = env GOOS=$(GOOS) GOARCH=$(GOARCH) go build

test:
	go test -v -cover -covermode=atomic ./...

install:
	$(GOBUILD) -o ${BINARY} app/*.go

linux-amd64:
	@docker build \
		--target linux-amd64 \
		--output . \
		--build-arg TARGETOS=linux \
		--build-arg TARGETARCH=amd64 \
		.

darwin-amd64:
	@docker build \
		--target darwin-amd64 \
		--output . \
		--build-arg TARGETOS=darwin \
		--build-arg TARGETARCH=amd64 \
		.

windows-amd64:
	@docker build \
		--target windows-amd64 \
		--output . \
		--build-arg TARGETOS=windows \
		--build-arg TARGETARCH=amd64 \
		.

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: test install linux-amd64  darwin-amd64 windows-amd64 clean