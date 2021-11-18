FROM golang:1.15-alpine

LABEL Description="编译镜像(v1.0)"  maintainer="songangweb@foxmail.com"

ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
#    && apk --no-cache add git gcc g++ \
    && go get -u -v github.com/beego/bee/v2