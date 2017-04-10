FROM alpine:edge
MAINTAINER CHENHW2 <https://github.com/chenhw2>

ARG BIN_URL=https://github.com/chenhw2/aliyun-ddns-cli/releases/download/v20170410/aliddns_linux-amd64-20170410.tar.gz

RUN apk add --update --no-cache wget supervisor ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

RUN mkdir -p /opt \
    && cd /opt \
    && wget -qO- ${BIN_URL} | tar xz \
    && mv aliddns_* aliddns

ADD Docker_entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
