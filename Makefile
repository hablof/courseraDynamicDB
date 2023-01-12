dep:
	go mod download
	go mod vendor
build:
	docker compose build main-application
run:
	docker compose up -d
stop:
	docker compose down
up-test-db:
	docker run -p 3366:3306 -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=golang -d --name testing-mysql --rm mysql
down-test-db:
	docker stop testing-mysql
test:
	go test ./...