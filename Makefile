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

# запуск базы для тестов
up-test-db:
	docker run -p 3366:3306 -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=integration_testing -d --name testing-mysql --rm mysql
down-test-db:
	docker stop testing-mysql
test:
	go test ./...

# локальный запуск базы и приложения
DIR := ${CURDIR}
up-local-db:
	docker run -p 3366:3306 -v $(DIR)/init-db:/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=golang -d --name local-mysql --rm mysql
down-local-db:
	docker stop local-mysql
start-local:
	go run . local 3366