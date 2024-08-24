ProjectName         = github.com/jiftle/sckproxy
OutAppName          = skpy
OutDir              = build
BUILD_VERSION      := v1.0.1
BUILD_DATE         := $(shell date "+%y%m%d")
BUILD_TIME         := $(shell date "+%F %T")
BUILD_AUTHOR       := $(shell id -u -n)
Version            := ${BUILD_VERSION}.${BUILD_DATE}
BUILD_HASH         := $(shell git log --pretty=format:"%h" -1)

all:
	make clean
	make build-linux-amd64
	make package

clean:
	rm -rf ${OutDir}

build-linux-amd64:
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build \
	-ldflags " \
	-s -w \
	-X '${ProjectName}/version.BuildTime=${BUILD_TIME}' \
	-X '${ProjectName}/version.Version=${Version}' \
	-X '${ProjectName}/version.Author=${BUILD_AUTHOR}' \
	-X '${ProjectName}/version.Hash=${BUILD_HASH}' \
	" \
	-o ${OutDir}/linux/amd64/${OutAppName}
	upx ${OutDir}/linux/amd64/${OutAppName}


# install package for os
package:
	cp -f build/linux/amd64/skpy 1-inst/skpy-inst-linux-amd64/files/skpy/
	cd 1-inst && tar czf skpy-inst-linux-amd64.tar.gz skpy-inst-linux-amd64/
	ls -lh  1-inst |grep .tar.gz
