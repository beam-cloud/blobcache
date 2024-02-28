FROM --platform=linux/amd64 golang:1.22 as base

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

RUN apt-get update && apt-get install -y wget git curl \
    libseccomp-dev libsndfile1 libsndfile1-dev \
    libaio-dev libzmq3-dev iptables \
    build-essential git libprotobuf-dev libprotobuf-c-dev \
    protobuf-c-compiler protobuf-compiler \
    pkg-config libbsd-dev iproute2 \
    libnftnl-dev libcap-dev libnet1-dev libnl-3-dev

# Build & install criu
RUN git clone https://github.com/checkpoint-restore/criu.git
RUN cd criu && make criu && make install-criu

COPY --from=build /usr/local/bin/blobcache /usr/local/bin/

ENV PATH "$PATH:/usr/local/sbin"

CMD ["blobcache"]
