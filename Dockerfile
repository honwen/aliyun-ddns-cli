FROM alpine:3.5
MAINTAINER CHENHW2 <https://github.com/chenhw2>

ARG VER=20170627
ARG URL=https://github.com/chenhw2/aliyun-ddns-cli/releases/download/v$VER/aliddns_linux-amd64-$VER.tar.gz
ARG TZ=Asia/Hong_Kong

RUN apk add --update --no-cache wget ca-certificates tzdata \
    && update-ca-certificates \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime \
    && rm -rf /var/cache/apk/*

RUN mkdir -p /usr/bin \
    && cd /usr/bin \
    && wget -qO- $URL | tar xz \
    && mv aliddns_* aliddns

ENV AccessKeyID=1234567890 \
    AccessKeySecret=abcdefghijklmn \
    Domain=ddns.example.win \
    Redo=0

CMD aliddns \
    --id $AccessKeyID \
    --secret $AccessKeySecret \
    auto-update \
    --domain $Domain \
    --redo $Redo
