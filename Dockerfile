FROM golang:alpine as builder
RUN apk add --update git
RUN go get github.com/chenhw2/aliyun-ddns-cli

FROM chenhw2/alpine:base
MAINTAINER CHENHW2 <https://github.com/chenhw2>

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
