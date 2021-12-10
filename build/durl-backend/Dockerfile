# 编译镜像
FROM songangweb/durl-build:v1.0.0 as builder

WORKDIR /go/src/durl-backend

COPY ./ ./

RUN cd app/exec/backend \
    && bee pack

# 运行镜像
FROM songangweb/durl-run:v1.0.0 as run

LABEL description="durl-backend"

ARG ENV=prod

ENV RUN_MODE=$ENV APP_NAME="durl-backend"

WORKDIR /durl/durl-backend

COPY --from=builder /go/src/durl-backend/app/exec/backend/backend.tar.gz .

RUN tar -zxvf backend.tar.gz \
    && rm -f backend.tar.gz

EXPOSE 8083
CMD ["/durl/durl-backend/backend"]

## 在根目录执行
## docker build -f build/durl-portal/Dockerfile . -t durl-portal:v1.0.0
