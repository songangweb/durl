# 编译镜像
FROM songangweb/durl-build:v1.0.0 as builder

WORKDIR /go/src/durl-openApi

COPY ./ ./

RUN cd app/exec/durl-openapi \
    && bee pack

# 运行镜像
FROM songangweb/durl-run:v1.0.0 as run

LABEL description="durl-openApi"

ARG ENV=prod

ENV RUN_MODE=$ENV APP_NAME="durl-openApi"

WORKDIR /durl/durl-openApi

COPY --from=builder /go/src/durl-openApi/app/exec/openApi/openApi.tar.gz .

RUN tar -zxvf durl-openapi.tar.gz \
    && rm -f durl-openapi.tar.gz

EXPOSE 8081
CMD ["./openApi"]

# 在根目录执行
# docker build -f build/durl-openapi/Dockerfile . -t durl-openapi:v1.0.0
