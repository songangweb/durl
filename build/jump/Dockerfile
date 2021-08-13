FROM golang:alpine as builder

WORKDIR /go/src/durl-jump

COPY ./ ./

RUN cd app/exec/jump \
    && go build .

EXPOSE 8082
CMD ["./jump"]

## 在根目录执行
## docker build -f build/jump/Dockerfile . -t durl-jump:v1
