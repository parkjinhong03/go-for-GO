buildServer:
	rm -rf ./main
	GOOS=linux GOARCH=amd64 go build ./main.go

buildImage:
	docker rmi api_gateway_image
	docker build --tag api_gateway_image .

runContainer:
	docker rm -f api_gateway_dev
	docker-compose up

upload:
	docker restart api_gateway_dev

attach:
	docker attach --sig-proxy=false api_gateway_dev

log:
	docker logs api_gateway_dev