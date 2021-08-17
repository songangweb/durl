# 编译镜像
FROM registry.cn-beijing.aliyuncs.com/songang-durl/durl-build:v1 as builder

WORKDIR /go/src/durl-openApi

COPY ./ ./

RUN cd app/exec/openApi \
    && bee pack

# 运行镜像
FROM registry.cn-beijing.aliyuncs.com/songang-durl/durl-run:v1 as run

LABEL description="durl-openApi"

ARG ENV=prod

ENV RUN_MODE=$ENV APP_NAME="durl-openApi"

WORKDIR /durl/durl-openApi

COPY --from=builder /go/src/durl-openApi/app/exec/openApi/openApi.tar.gz .

RUN tar -zxvf openApi.tar.gz \
    && rm -f openApi.tar.gz

EXPOSE 8081
CMD ["./openApi"]

# 在根目录执行
# docker build -f build/openApi/Dockerfile . -t durl-openApi:v1
