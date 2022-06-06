# 编译镜像
FROM songangweb/durl-build:v1.0.0 as builder

WORKDIR /go/src/durl-jump

COPY ./ ./

RUN cd app/exec/jump \
    && bee pack

# 运行镜像
FROM songangweb/durl-run:v1.0.0 as run

LABEL description="durl-jump"

ARG ENV=prod

ENV RUN_MODE=$ENV APP_NAME="durl-jump"

WORKDIR /durl/durl-jump

COPY --from=builder /go/src/durl-jump/app/exec/jump/jump.tar.gz .

RUN tar -zxvf jump.tar.gz \
    && rm -f jump.tar.gz

EXPOSE 8082
EXPOSE 9082
CMD ["/durl/durl-jump/jump"]

## 在根目录执行
## docker build -f build/durl-jump/Dockerfile  . -t songangweb/durl-jump:v1.0.4
## or 使用 buildx 构建多平台 Docker 镜像 https://blog.csdn.net/alex_yangchuansheng/article/details/103343697/
## docker buildx build -t songangweb/durl-jump:v1.0.4 --platform=linux/arm,linux/arm64,linux/amd64 -f build/durl-jump/Dockerfile . --push