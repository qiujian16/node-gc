.PHONY: build doc fmt lint run test vendor_clean vendor_get vendor_update vet

default: build

build:
	go build -v -i -o bin/node-gc github.com/qiujian16/node-gc/cmd

docker-binary:
	CGO_ENABLED=0 go build -a -installsuffix cgo -v -i -o bin/node-gc github.com/qiujian16/node-gc/cmd
	strip bin/node-gc

image: docker-binary
	docker build -t node-gc .
