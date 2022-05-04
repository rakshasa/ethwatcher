FROM golang:1.17

WORKDIR /go/src

ENV GOCACHE="${GOPATH}/cache"

COPY . /go/src 

RUN --mount=type=cache,target=/go/cache go test ./...
RUN --mount=type=cache,target=/go/cache go build -o bin/ethereum-watcher cli/main.go


FROM alpine

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN --mount=type=cache,target=/var/cache/apt apk add ca-certificates

COPY --from=0 /go/src/bin/* /bin/

CMD ["/bin/ethereum-watcher"]
