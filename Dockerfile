FROM golang:1.25.5-alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,https://goproxy.io,direct

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && set -ex \
    && apk update --no-cache \
    && apk add --no-cache tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && echo -e "\033[42;37m Build Completed :).\033[0m\n"

WORKDIR /go-project

COPY . /go-project

RUN go mod tidy \
    && go build -a -installsuffix cgo -o momentoApiBinary momentoapi.go

CMD ["./momentoApiBinary", "-f", "etc/momentoapi.yaml"]
