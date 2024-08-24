all:
	make clean
	make build

clean:
	rm -rf build

build-linux-amd64:
	go build \
	-ldflags " \
	-s -w \
	" \
	-o build/linux/amd64/skpy



# install package for os
pack:
	cp -f build/linux/amd64/skpy 1-inst/skpy-inst-linux-amd64/files/skpy/
	tar czf 1-inst/skpy-inst-linux-amd64.tar.gz 1-inst/skpy-inst-linux-amd64/
