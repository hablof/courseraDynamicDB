// тут лежит тестовый код
// менять вам может потребоваться только коннект к базе
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// dsn это соединение с базой
	// вы можете изменить этот на тот который вам нужен
	// docker run -p 3366:3306 -v $(PWD):/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=golang -d mysql
	dsn = "root:1234@tcp(localhost:3366)/golang?charset=utf8"
	// DSN = "coursera:5QPbAUufx7@tcp(localhost:3306)/coursera?charset=utf8"
)

func main() {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		return
	}

	defer db.Close()
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Println(err)
		return
	}

	handler, err := NewDbExplorer(db)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("starting server at :8082")
	http.ListenAndServe(":8082", handler)
}
