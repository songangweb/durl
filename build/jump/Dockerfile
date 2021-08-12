# 编译镜像
#FROM golang:alpine as builder
FROM  registry.cn-beijing.aliyuncs.com/durl/golang:alpine as builder

WORKDIR /durl

ADD ../../go.mod .
ADD ../../go.sum .

RUN export GO111MODULE=on \
    && export GOPROXY=https://goproxy.cn \
    && go mod download

COPY ../.. ./

RUN go build .

# 运行镜像
#FROM alpine
#FROM registry.cn-beijing.aliyuncs.com/durl/alpine:latest

#FROM golang:alpine
#FROM  registry.cn-beijing.aliyuncs.com/durl/golang:alpine

#LABEL maintainer="durl" version="1.0"

#WORKDIR /code

#COPY --from=builder /durl/conf    ./conf
#COPY --from=builder /durl/views   ./views
#COPY --from=builder /durl/controllers   ./controllers
#COPY --from=builder /durl/durl    ./

EXPOSE 8080

#CMD ["/code/durl"]

CMD ["/durl/durl"]

