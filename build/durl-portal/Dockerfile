# 编译镜像
FROM songangweb/durl-build:v1.0.0 as builder

WORKDIR /go/src/durl-portal

COPY ./ ./

RUN cd app/exec/durl-portal \
    && bee pack

# 运行镜像
FROM songangweb/durl-run:v1.0.0 as run

LABEL description="durl-portal"

ARG ENV=prod

ENV RUN_MODE=$ENV APP_NAME="durl-portal"

WORKDIR /durl/durl-portal

COPY --from=builder /go/src/durl-portal/app/exec/portal/portal.tar.gz .

RUN tar -zxvf durl-portal.tar.gz \
    && rm -f durl-portal.tar.gz

EXPOSE 8080
CMD ["/durl/durl-portal/portal"]

## 在根目录执行
## docker build -f build/durl-portal/Dockerfile . -t durl-durl-portal:v1.0.0
