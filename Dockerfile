FROM golang:1.17

WORKDIR /go/src/github.com/rakshasa/ethwatcher

ARG ETHEREUM_RPC_ENDPOINT
ARG SKIP_RPC_TESTS

ENV GOCACHE="${GOPATH}/cache"
ENV ETHEREUM_RPC_ENDPOINT="${ETHEREUM_RPC_ENDPOINT}"
ENV SKIP_RPC_TESTS="${SKIP_RPC_TESTS}"

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/cache go mod download

COPY . ./

RUN --mount=type=cache,target=/go/cache go test ./...
RUN --mount=type=cache,target=/go/cache go build -o bin/ethereum-watcher cli/main.go


FROM alpine

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN --mount=type=cache,target=/var/cache/apt apk add ca-certificates

COPY --from=0 /go/src/github.com/rakshasa/ethwatcher/bin/ethereum-watcher /

CMD ["/ethereum-watcher"]
