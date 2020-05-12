FROM golang:alpine
MAINTAINER Ben Whitley <dev@purplg.com>

# Need 'make' to build
RUN apk add build-base

WORKDIR /build

# Pull deps
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy source files
COPY . .

# Do build
RUN make build

WORKDIR /dist

RUN cp /build/bin/main .

CMD ["./main"]
