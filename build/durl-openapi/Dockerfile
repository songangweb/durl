# 编译镜像
FROM songangweb/durl-build:v1.0.0 as builder

WORKDIR /go/src/durl-openapi

COPY ./ ./

RUN cd app/exec/openapi \
    && bee pack

# 运行镜像
FROM songangweb/durl-run:v1.0.0 as run

LABEL description="durl-openapi"

ARG ENV=prod

ENV RUN_MODE=$ENV APP_NAME="durl-openapi"

WORKDIR /durl/durl-openapi

COPY --from=builder /go/src/durl-openapi/app/exec/openapi/openapi.tar.gz .

RUN tar -zxvf openapi.tar.gz \
    && rm -f openapi.tar.gz

EXPOSE 8081
CMD ["./openapi"]

# 在根目录执行
# docker build -f build/durl-openapi/Dockerfile . -t durl-openapi:v1.0.0
