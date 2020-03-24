cucumber:
	docker-compose up -d
	cd features && godog ./
	docker-compose stop
	docker-compose rm -f

run:
	docker-compose up -d
	go run main.go
	docker-compose stop
	docker-compose rm -f

unit:
	go test -v -race ./...

test:
	make unit
	make cucumber

package:
	go get -t ./...
	go get github.com/DATA-DOG/godog/cmd/godog