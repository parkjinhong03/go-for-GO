
GOPATH:=$(shell go env GOPATH)

.PHONY: build
build: proto

	GOOS=linux GOARCH=amd64 go build -o api-gateway *.go

.PHONY: docker
docker:
	docker build . -t api-gateway:latest
