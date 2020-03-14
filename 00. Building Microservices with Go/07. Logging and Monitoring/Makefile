buildserver:
	echo "Build API"
	rm -rf ./helloworldserver
	GOOS=linux GOARCH=amd64 go build -o ./helloworldserver ./main.go

runserver: buildserver
	echo "Run API"
	docker-compose -p helloWorldServer -f ./docker-compose.yml up --build