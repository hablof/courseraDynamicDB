run: dep build start

dep:
	go mod download
	go mod vendor
build:
	docker compose build main-application
start: 
	docker compose up -d
stop:
	docker compose down


DIR := ${CURDIR}
up-test-db:
	docker run -p 3366:3306 -v $(DIR)/init-db:/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=golang -d --name testing-mysql --rm mysql
down-test-db:
	docker stop testing-mysql

start-local:
	go run . local 3366

test:
	go test ./...