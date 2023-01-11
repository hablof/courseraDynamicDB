// тут лежит тестовый код
// менять вам может потребоваться только коннект к базе
package main

import (
	"database/sql"
	"fmt"
	"hw6coursera/dbexplorer"
	"hw6coursera/repository"
	"hw6coursera/router"
	"hw6coursera/service"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// dsn это соединение с базой
	// вы можете изменить этот на тот который вам нужен
	// docker run -p 3366:3306 -v $(PWD):/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=golang -d mysql
	dsn = "root:1234@tcp(database-mysql:3306)/golang?charset=utf8"
	// DSN = "coursera:5QPbAUufx7@tcp(localhost:3306)/coursera?charset=utf8"
)

func main() {
	log.Printf("startig application")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Println(err)
		return
	}

	repo := repository.NewRepository(db)
	explorer := dbexplorer.NewDbExplorer(repo)
	service := service.NewService(repo, explorer)
	if err := service.InitSchema(); err != nil {
		log.Printf("failed to init database shcema: %v", err)
		return
	}
	router := router.NewRouter(service)

	fmt.Println("starting server at :8082")
	http.ListenAndServe(":8082", router)
}
