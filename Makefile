GOCMD=go
GOBUILD=${GOCMD} build
GOCLEAN=${GOCMD} clean
GOTEST=${GOCMD} test
GOGET=${GOCMD} get

BINARY_NAME=main
BINARY_DIR=bin
BINARY=$(BINARY_DIR)/$(BINARY_NAME)

BUILD_CMD=$(GOBUILD) -v -o $(BINARY) ./cmd/memberchannels/memberchannels.go

all: test build

build:
	${BUILD_CMD}

test:
	${GOTEST} -v ./...

run:
	${BUILD_CMD}
	${BINARY}
