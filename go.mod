module github.com/honwen/aliyun-ddns-cli

go 1.20

require (
	// locked before tracing/logging https://github.com/denverdino/aliyungo/commits/master/go.mod
	github.com/denverdino/aliyungo v0.0.0-20220321085828-46dabbd9e212
	github.com/honwen/golibs v0.4.5
	github.com/honwen/ip2loc v0.2.8
	github.com/urfave/cli v1.22.12
	github.com/ysmood/got v0.33.0
)

require (
	github.com/AdguardTeam/golibs v0.12.0 // indirect
	github.com/PuerkitoBio/goquery v1.8.1 // indirect
	github.com/Workiva/go-datastructures v1.0.53 // indirect
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/aead/poly1305 v0.0.0-20180717145839-3fee0db0b635 // indirect
	github.com/ameshkov/dnscrypt/v2 v2.2.5 // indirect
	github.com/ameshkov/dnsstamps v1.0.3 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/pprof v0.0.0-20221219190121-3cb0bae90811 // indirect
	github.com/miekg/dns v1.1.51 // indirect
	github.com/mr-karan/doggo v0.5.5 // indirect
	github.com/onsi/ginkgo/v2 v2.6.1 // indirect
	github.com/quic-go/qtls-go1-18 v0.2.0 // indirect
	github.com/quic-go/qtls-go1-19 v0.2.0 // indirect
	github.com/quic-go/qtls-go1-20 v0.1.0 // indirect
	github.com/quic-go/quic-go v0.32.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/ysmood/gop v0.0.1 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	golang.org/x/exp v0.0.0-20221230185412-738e83a70c30 // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/tools v0.4.0 // indirect
)

replace github.com/mr-karan/doggo => github.com/honwen/doggo v0.0.0-20230203023054-7db5c2144fa4
