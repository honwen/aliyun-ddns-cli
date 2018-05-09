FROM golang:alpine as builder
RUN apk add --update git curl
RUN go get -u -v \
        -ldflags "-X main.version=$(curl -sSL https://api.github.com/repos/chenhw2/aliyun-ddns-cli/commits/master | \
            sed -n '1,9{/sha/p; /date/p}' | sed 's/.* \"//g' | cut -c1-10 | tr [a-z] [A-Z] | sed 'N;s/\n/@/g')" \
        github.com/chenhw2/aliyun-ddns-cli

FROM chenhw2/alpine:base
LABEL MAINTAINER CHENHW2 <https://github.com/chenhw2>

# /usr/bin/aliyun-ddns-cli
COPY --from=builder /go/bin /usr/bin

ENV AKID=1234567890 \
    AKSCT=abcdefghijklmn \
    DOMAIN=ddns.example.win \
    IPAPI=[IPAPI-GROUP] \
    REDO=0

CMD aliyun-ddns-cli \
    --id ${AKID} \
    --secret ${AKSCT} \
    --ipapi ${IPAPI} \
    auto-update \
    --domain ${DOMAIN} \
    --redo ${REDO}
