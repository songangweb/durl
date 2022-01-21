# 编译镜像
FROM songangweb/durl-build:v1.0.0 as builder

WORKDIR /go/src/durl-portal

COPY ./ ./

RUN cd app/exec/portal \
    && bee pack

# 运行镜像
FROM songangweb/durl-run:v1.0.0 as run

LABEL description="durl-portal"

ARG ENV=prod

ENV RUN_MODE=$ENV APP_NAME="durl-portal"

WORKDIR /durl/durl-portal

COPY --from=builder /go/src/durl-portal/app/exec/portal/portal.tar.gz .

RUN tar -zxvf portal.tar.gz \
    && rm -f portal.tar.gz

EXPOSE 8080
EXPOSE 9080
CMD ["/durl/durl-portal/portal"]

## 在根目录执行
## docker build -f build/durl-portal/Dockerfile  . -t durl-portal:v1.0.2
## or 使用 buildx 构建多平台 Docker 镜像 https://blog.csdn.net/alex_yangchuansheng/article/details/103343697/
## docker buildx build -t songangweb/durl-portal:v1.0.2 --platform=linux/arm,linux/arm64,linux/amd64 -f build/durl-portal/Dockerfile . --push
