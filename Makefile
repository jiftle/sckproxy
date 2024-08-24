OutDir=build

all:
	make clean
	make build-linux-amd64
	make build-linux-arm64
	make compress
	tree ${OutDir} -lh

clean:
	rm -rf build

build-linux-amd64:
	CGO_ENABLED=0
	GOOS=linux \
	GOARCH=amd64 \
	go build \
	-ldflags=' -s -w' \
	-o ${OutDir}/linux/amd64/skpy

build-linux-arm64:
	CGO_ENABLED=0
	GOOS=linux \
	GOARCH=arm64 \
	go build \
	-ldflags=' -s -w' \
	-o ${OutDir}/linux/arm64/skpy

compress:
	upx ${OutDir}/linux/amd64/skpy
	upx ${OutDir}/linux/arm64/skpy