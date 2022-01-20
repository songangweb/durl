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

EXPOSE [8082,9082]
CMD ["./jump"]

# 在根目录执行
# docker build -f build/durl-jump/Dockerfile . -t durl-jump:v1.0.0
