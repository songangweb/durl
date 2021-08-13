FROM golang:alpine as builder

WORKDIR /go/src/durl-openApi

COPY ./ ./

RUN cd app/exec/openApi \
    && go build .

EXPOSE 8081
CMD ["./openApi"]

## 在根目录执行
## docker build -f build/openApi/Dockerfile . -t durl-openApi:v1
