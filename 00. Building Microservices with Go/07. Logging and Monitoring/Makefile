buildserver:
	echo "Build API"
	rm -rf ./helloworldserver
	GOOS=linux GOARCH=amd64 go build -o ./helloworldserver ./main.go

runlogging:
	docker-compose -f ./docker-compose.yml up --build -d

runserver: buildserver
	echo "Run API"
	LOGSTASH=localhost:5000 go run main.go