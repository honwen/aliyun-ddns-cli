### Source
- https://github.com/chenhw2/aliyun-ddns-cli
  
### Thanks to
- https://github.com/denverdino/aliyungo
  
### Docker
- https://hub.docker.com/r/chenhw2/aliyun-ddns-cli/
  
### Usage
```
$ docker pull chenhw2/aliyun-ddns-cli

$ docker run -d \
    -e "AccessKeyID=1234567890" \
    -e "AccessKeySecret=abcdefghijklmn" \
    -e "AcceDomain=ddns.example.winssKeyID" \
    -e "Redo=0" \
    chenhw2/aliyun-ddns-cli

```
### Help
```
$ docker run --rm chenhw2/aliyun-ddns-cli -h
NAME:
   aliddns - aliyun-ddns-cli

USAGE:
   aliyun-ddns-cli [global options] command [command options] [arguments...]

VERSION:
   XXXXXX

COMMANDS:
     help, h  Shows a list of commands or help for one command
   DDNS:
     list         List AliYun's DNS DomainRecords Record
     update       Update AliYun's DNS DomainRecords Record
     auto-update  Auto-Update AliYun's DNS DomainRecords Record, Get IP using its getip
   GET-IP:
     getip        Get IP Combine 10 different Web-API

GLOBAL OPTIONS:
   --access-key-id value, --id value          AliYun's Access Key ID
   --access-key-secret value, --secret value  AliYun's Access Key Secret
   --help, -h                                 show help
   --version, -v                              print the version

```
### CLI Example:
```
aliddns --id ${AccessKeyID} --secret ${AccessKeySecret} \
    auto-update --domain ddns.example.win

aliddns --id ${AccessKeyID} --secret ${AccessKeySecret} \
    update --domain ddns.example.win \
    --ipaddr $(ifconfig pppoe-wan | sed -n '2{s/[^0-9]*://;s/[^0-9.].*//p}')
```