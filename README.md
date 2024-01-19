### Source

- https://github.com/honwen/aliyun-ddns-cli

### Thanks (package alidns)

- https://github.com/denverdino/aliyungo

### Docker

- https://hub.docker.com/r/chenhw2/aliyun-ddns-cli/

### Usage

```shell
$ docker pull chenhw2/aliyun-ddns-cli

$ docker run -d \
    -e "AKID=[ALIYUN's AccessKey-ID]" \
    -e "AKSCT=[ALIYUN's AccessKey-Secret]" \
    -e "DOMAIN=ddns.aliyun.win" \
    -e "REDO=600" \
    -e "TTL=600" \
    chenhw2/aliyun-ddns-cli
```

### Example (for Synology)

- https://github.com/honwen/aliyun-ddns-cli/tree/master/example

### Help

```shell
$ docker run --rm chenhw2/aliyun-ddns-cli -h
NAME:
   aliddns - aliyun-ddns-cli

USAGE:
   aliyun-ddns-cli [global options] command [command options] [arguments...]

VERSION:
   Git:[MISSING BUILD VERSION [GIT HASH]] (go1.21)

COMMANDS:
   help, h  Shows a list of commands or help for one command

   DDNS:
     list         List AliYun's DNS DomainRecords Record
     delete       Delete AliYun's DNS DomainRecords Record
     update       Update AliYun's DNS DomainRecords Record, Create Record if not exist
     auto-update  Auto-Update AliYun's DNS DomainRecords Record, Get IP using its getip

   GET-IP:
     getip          Get IP Combine 10+ different Web-API
     resolve        Get DNS-IPv4 Combine 4+ DNS Upstream

GLOBAL OPTIONS:
   --access-key-id value, --id value          AliYun's Access Key ID
   --access-key-secret value, --secret value  AliYun's Access Key Secret
   --ipapi value, --api value                 Web-API to Get IP, like: http://v6r.ipip.net
   --ipv6, -6                                 IPv6
   --help, -h                                 show help
   --version, -v                              print the version
```

### CLI Example:

```shell
aliddns --id ${AccessKeyID} --secret ${AccessKeySecret} \
    auto-update --domain ddns.example.win

aliddns --id ${AccessKeyID} --secret ${AccessKeySecret} \
    update --domain ddns.example.win \
    --ipaddr $(ifconfig pppoe-wan | sed -n '2{s/[^0-9]*://;s/[^0-9.].*//p}')
```
