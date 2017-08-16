FROM chenhw2/alpine:base
MAINTAINER CHENHW2 <https://github.com/chenhw2>

ARG VER=20170816
ARG URL=https://github.com/chenhw2/aliyun-ddns-cli/releases/download/v$VER/aliddns_linux-amd64-$VER.tar.gz

RUN mkdir -p /usr/bin \
    && cd /usr/bin \
    && wget -qO- ${URL} | tar xz \
    && mv aliddns_* aliddns

ENV AKID=1234567890 \
    AKSCT=abcdefghijklmn \
    DOMAIN=ddns.example.win \
    REDO=0

CMD aliddns \
    --id ${AKID} \
    --secret ${AKSCT} \
    auto-update \
    --domain ${DOMAIN} \
    --redo ${REDO}
