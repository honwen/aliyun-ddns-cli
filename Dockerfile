FROM chenhw2/alpine:base
LABEL MAINTAINER="honwen <https://github.com/honwen>"

# /usr/bin/aliyun-ddns-cli
RUN mkdir -p /usr/bin/ \
    && cd /usr/bin/ \
    && curl -skSL $( \
    curl -skSL 'https://api.github.com/repos/honwen/aliyun-ddns-cli/releases/latest' | \
    sed -n '/url.*linux-amd64/{s/.*\(https:.*.gz\).*/\1/p}'\
    ) | tar --strip-components=1 -zx linux-amd64/aliddns \
    && ln -sf aliddns aliyun-ddns-cli \
    && aliyun-ddns-cli -v

ENV AKID=1234567890 \
    AKSCT=abcdefghijklmn \
    DOMAIN=ddns.example.win \
    IPAPI=[IPAPI-GROUP] \
    REDO=555r \
    TTL=600

CMD ["/bin/sh", "-cx", "aliyun-ddns-cli --ipapi ${IPAPI} ${IPV6:+-6} auto-update --domain ${DOMAIN} --redo ${REDO} --ttl ${TTL}"]
