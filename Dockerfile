# syntax=docker/dockerfile:1.6

FROM golang:1.22-bookworm AS build

WORKDIR /workspace

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -o blobcache /workspace/cmd/main.go

CMD ["blobcache"]


FROM gcr.io/distroless/static-debian12 AS release

COPY --from=build /workspace/blobcache /usr/local/bin/

CMD ["blobcache"]
