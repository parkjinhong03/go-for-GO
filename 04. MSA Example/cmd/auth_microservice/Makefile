buildServer:
	rm -rf ./main
	GOOS=linux GOARCH=amd64 go build ./main.go

buildImage:
	docker rmi auth_microservice_image
	docker build --tag auth_microservice_image .

runContainer:
	docker rm -f auth_microservice_dev
	docker-compose up

upload:
	docker restart auth_microservice_dev

attach:
	docker attach --sig-proxy=false auth_microservice_dev
