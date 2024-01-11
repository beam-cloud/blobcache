FROM --platform=linux/amd64 golang:1.19 as base

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
FROM base AS release

COPY --from=build /usr/local/bin/blobcache /usr/local/bin/

CMD ["blobcache"]
