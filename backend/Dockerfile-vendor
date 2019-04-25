FROM golang:1.11

COPY . /app
WORKDIR /app

#RUN go mod download

#RUN echo "0"
#RUN echo $(ls vendor/golang.org/x/sys/unix)
#RUN go build -mod=vendor -o bin/hydro-dex-ctl -v -ldflags '-s -w' cli/admincli/main.go

RUN go build -mod=vendor -o bin/hydro-dex-ctl -v -ldflags '-s -w' cli/admincli/main.go && \
    go build -mod=vendor -o bin/adminapi -v -ldflags '-s -w' cli/adminapi/main.go && \
    go build -mod=vendor -o bin/api -v -ldflags '-s -w' cli/api/main.go && \
    go build -mod=vendor -o bin/engine -v -ldflags '-s -w' cli/engine/main.go && \
    go build -mod=vendor -o bin/launcher -v -ldflags '-s -w' cli/launcher/main.go && \
    go build -mod=vendor -o bin/watcher -v -ldflags '-s -w' cli/watcher/main.go && \
    go build -mod=vendor -o bin/websocket -v -ldflags '-s -w' cli/websocket/main.go && \
    go build -mod=vendor -o bin/maker -v -ldflags '-s -w' cli/maker/main.go

FROM alpine
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN apk update && \
  apk add sqlite ca-certificates && \
  rm -rf /var/cache/apk/*

COPY --from=0 /app/db /db/
COPY --from=0 /app/bin/* /bin/
