FROM golang:alpine as builder

WORKDIR /go/src/durl-portal

COPY ./ ./

RUN cd app/exec/portal \
    && go build .

EXPOSE 8080
CMD ["./portal"]

## 在根目录执行
## docker build -f build/portal/Dockerfile . -t durl-portal:v1
