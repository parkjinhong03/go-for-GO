
GOPATH:=$(shell go env GOPATH)

.PHONY: build
build: proto

	GOOS=linux GOARCH=amd64 go build -o api-gateway *.go

.PHONY: upload
upload:
	docker build --build-arg CONSUL_ADDRESS=${CONSUL_ADDRESS} -t api-gateway:${IMAGE_VERSION}.RELEASE .
	docker tag api-gateway:${IMAGE_VERSION}.RELEASE jinhong0719/examples-blog-api-gateway:${IMAGE_VERSION}.RELEASE
	docker push jinhong0719/examples-blog-api-gateway:${IMAGE_VERSION}.RELEASE
