FROM golang:1.22 as base

# used for compiling a build and doing local development
FROM base as build

WORKDIR /workspace

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ENV GIN_MODE=release
RUN go build -v -o /usr/local/bin/blobcache /workspace/cmd/main.go


CMD ["blobcache"]

# used for production release
FROM ubuntu:22.04 AS release

COPY --from=golang:1.22 /usr/local/go/ /usr/local/go/

RUN apt-get update && apt-get install -y software-properties-common
RUN add-apt-repository ppa:criu/ppa
RUN apt-get update && apt-get install -y wget git curl criu=3.19-1ppa1.22.04
COPY --from=build /usr/local/bin/blobcache /usr/local/bin/

ENV PATH "$PATH:/usr/sbin"

CMD ["blobcache"]
