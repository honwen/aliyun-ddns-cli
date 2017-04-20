FROM alpine:edge
MAINTAINER CHENHW2 <https://github.com/chenhw2>

ARG BIN_URL=https://github.com/chenhw2/aliyun-ddns-cli/releases/download/v20170420/aliddns_linux-amd64-20170420.tar.gz
ARG TZ=Asia/Hong_Kong

RUN apk add --update --no-cache wget supervisor ca-certificates tzdata \
    && update-ca-certificates \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime \
    && rm -rf /var/cache/apk/*

RUN mkdir -p /opt \
    && cd /opt \
    && wget -qO- ${BIN_URL} | tar xz \
    && mv aliddns_* aliddns

ENV AccessKeyID=1234567890 \
    AccessKeySecret=abcdefghijklmn \
    Domain=ddns.example.win \
    Redo=0

ADD Docker_entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
