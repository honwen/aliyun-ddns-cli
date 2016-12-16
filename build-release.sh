#!/bin/bash
MD5='md5sum'
unamestr=`uname`
if [[ "$unamestr" == 'Darwin' ]]; then
	MD5='md5'
fi

UPX=false
if hash upx 2>/dev/null; then
	UPX=true
fi

VERSION=`date -u +%Y%m%d`
LDFLAGS="-X main.version=$VERSION -s -w -linkmode external -extldflags -static"
GCFLAGS=""

OSES=(windows linux darwin freebsd)
ARCHS=(amd64 386)
rm -rf ./release
mkdir -p ./release
for os in ${OSES[@]}; do
	for arch in ${ARCHS[@]}; do
		suffix=""
		if [ "$os" == "windows" ]; then
			suffix=".exe"
			LDFLAGS="-X main.version=$VERSION -s -w"
		fi
		env CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ./release/aliddns_${os}_${arch}${suffix} github.com/chenhw2/aliyun-ddns-cli
		if $UPX; then upx -9 ./release/aliddns_${os}_${arch}${suffix};fi
		tar -zcf ./release/aliddns_${os}-${arch}-$VERSION.tar.gz ./release/aliddns_${os}_${arch}${suffix}
		$MD5 ./release/aliddns_${os}-${arch}-$VERSION.tar.gz
	done
done

# ARM
ARMS=(5 6 7)
for v in ${ARMS[@]}; do
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ./release/aliddns_arm$v  github.com/chenhw2/aliyun-ddns-cli
done
if $UPX; then upx -9 ./release/aliddns_arm*;fi
tar -zcf ./release/aliddns_arm-$VERSION.tar.gz ./release/aliddns_arm*
$MD5 ./release/aliddns_arm-$VERSION.tar.gz

#MIPS32LE
LDFLAGS="-X main.version=$VERSION -s -w"
env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ./release/aliddns_mipsle github.com/chenhw2/aliyun-ddns-cli
env CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ./release/aliddns_mips github.com/chenhw2/aliyun-ddns-cli

if $UPX; then upx -9 ./release/aliddns_mips**;fi
tar -zcf ./release/aliddns_mipsle-$VERSION.tar.gz ./release/aliddns_mipsle
tar -zcf ./release/aliddns_mips-$VERSION.tar.gz ./release/aliddns_mips
$MD5 ./release/aliddns_mipsle-$VERSION.tar.gz
