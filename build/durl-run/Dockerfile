FROM alpine:3.14.1

LABEL description="运行镜像" maintainer="songangweb@foxmail.com"

# 设置代理
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk --no-cache add ca-certificates tzdata openssl nginx wget bash libstdc++ logrotate

# 设置语言
ENV LANG=en_US.UTF-8 LANGUAGE=en_US.UTF-8 TZ=Asia/Shanghai

# 安装supervisor
COPY --from=ochinchina/supervisord:latest /usr/local/bin/supervisord /usr/local/bin/supervisord