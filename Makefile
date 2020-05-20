buildtime=$(shell date)
version=0

GOCMD=go
GOBUILD=${GOCMD} build
GOCLEAN=${GOCMD} clean
GOTEST=${GOCMD} test
GOGET=${GOCMD} get

BINARY_NAME=main
BINARY_DIR=bin
BINARY=${BINARY_DIR}/${BINARY_NAME}

FLAGS="-X 'main.BuildTime=${buildtime}' -X 'main.BuildVersion=${version}'"

BUILD_CMD=${GOBUILD} -v -ldflags=${FLAGS} -o ${BINARY} ./cmd/memberchannels/memberchannels.go

all: test build

build:
	echo ${FLAGS}
	${BUILD_CMD}

test:
	${GOTEST} -v ./...

run:
	${BUILD_CMD}
	${BINARY}
