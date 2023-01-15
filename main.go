package main

import (
	"database/sql"
	"flag"
	"fmt"
	"hw6coursera/dbexplorer"
	"hw6coursera/repository"
	"hw6coursera/router"
	"hw6coursera/service"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	localDSN  = "root:1234@tcp(localhost:%d)/golang?charset=utf8"
	dockerDSN = "root:1234@tcp(database-mysql:3306)/golang?charset=utf8"
)

func main() {
	log.Printf("startig application")

	var db *sql.DB
	var err error
	var port int

	flag.Parse()
	if flag.Arg(0) == "local" {
		port, err = strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.Println(err)
			return
		}
		DSN := fmt.Sprintf(localDSN, port)
		log.Printf("db dsn: %s", DSN)
		db, err = sql.Open("mysql", DSN)
	} else {
		db, err = sql.Open("mysql", dockerDSN)
	}

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
		log.Printf("failed to ping db: %v", err)
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
