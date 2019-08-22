FROM golang:1.11

COPY . /app
WORKDIR /app

RUN go mod download
RUN go build -o bin/hydro-dex-ctl -v -ldflags '-s -w' cli/admincli/main.go && \
  go build -o bin/adminapi -v -ldflags '-s -w' cli/adminapi/main.go && \
  go build -o bin/api -v -ldflags '-s -w' cli/api/main.go && \
  go build -o bin/engine -v -ldflags '-s -w' cli/engine/main.go && \
  go build -o bin/launcher -v -ldflags '-s -w' cli/launcher/main.go && \
  go build -o bin/watcher -v -ldflags '-s -w' cli/watcher/main.go && \
  go build -o bin/websocket -v -ldflags '-s -w' cli/websocket/main.go && \
  go build -o bin/maker -v -ldflags '-s -w' cli/maker/main.go

FROM alpine
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN apk update && \
  apk add sqlite ca-certificates wget && \
  rm -rf /var/cache/apk/*

RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
  wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.30-r0/glibc-2.30-r0.apk && \
  apk add glibc-2.30-r0.apk

COPY --from=0 /app/db /db/
COPY --from=0 /app/bin/* /bin/
