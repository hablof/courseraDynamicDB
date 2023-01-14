package main

import (
	"database/sql"
	"fmt"
	"hw6coursera/dbexplorer"
	"hw6coursera/repository"
	"hw6coursera/router"
	"hw6coursera/service"
	"log"
	"testing"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

// // CaseResponse
// type CR map[string]interface{}

type testCase struct {
	name                   string
	method                 string // GET по-умолчанию в http.NewRequest если передали пустую строку
	path                   string
	queryParams            string
	expectedResponseStatus int
	expectedResponseBody   string
	requestBody            map[string]string
}

var (
	client = &http.Client{Timeout: time.Second}
)

func PrepareTestApis(db *sql.DB) {

	qs := []string{
		`DROP DATABASE IF EXISTS integration_testing;`,

		`CREATE DATABASE integration_testing;`,

		`USE integration_testing;`,

		`DROP TABLE IF EXISTS items_test;`,

		`CREATE TABLE items_test (
  		id int(11) NOT NULL AUTO_INCREMENT,
  		title varchar(255) NOT NULL,
  		description text NOT NULL,
  		updated varchar(255) DEFAULT NULL,
		rating decimal(5,2),
		shadow boolean,
		level int,
  		PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`INSERT INTO items_test (title, description, updated, rating, shadow, level) VALUES
		('database/sql', 'Рассказать про базы данных', 'rvasily', '5.0', TRUE, '15'),
		('memcache', 'Рассказать про мемкеш с примером использования', NULL, '0.0', FALSE, '80');`,

		`DROP TABLE IF EXISTS users_test;`,

		`CREATE TABLE users_test (
		user_id int(11) NOT NULL AUTO_INCREMENT,
  		login varchar(255) NOT NULL,
  		password varchar(255) NOT NULL,
  		email varchar(255) NOT NULL,
  		info text NOT NULL,
  		updated varchar(255) DEFAULT NULL,
  		PRIMARY KEY (user_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`INSERT INTO users_test (login, password, email, info, updated) VALUES
		('rvasily',	'love',	'rvasily@example.com',	'none',	NULL);`,
	}

	for _, q := range qs {
		_, err := db.Exec(q)
		if err != nil {
			panic(err)
		}
	}
}

func CleanupTestApis(db *sql.DB) {
	qs := []string{
		"DROP TABLE IF EXISTS items_test;",
		"DROP TABLE IF EXISTS users_test;",
		"DROP DATABASE IF EXISTS integration_testing;",
	}
	for _, q := range qs {
		_, err := db.Exec(q)
		if err != nil {
			panic(err)
		}
	}
}

func TestApis(t *testing.T) {
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3366)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	PrepareTestApis(db)

	defer CleanupTestApis(db)

	repo := repository.NewRepository(db)
	explorer := dbexplorer.NewDbExplorer(repo)
	service := service.NewService(repo, explorer)
	if err := service.InitSchema(); err != nil {
		log.Printf("failed to init database shcema: %v", err)
		return
	}
	router := router.NewRouter(service)

	ts := httptest.NewServer(router)

	tableItemsContent := []map[string]interface{}{
		{
			"id":          1,
			"title":       "database/sql",
			"description": "Рассказать про базы данных",
			"updated":     "rvasily",
		},
		{
			"id":          2,
			"title":       "memcache",
			"description": "Рассказать про мемкеш с примером использования",
			"updated":     nil,
		},
	}
	jsonItems, _ := json.MarshalIndent(tableItemsContent, "", "    ")
	itemsData := string(jsonItems)

	json1stItem, _ := json.MarshalIndent(tableItemsContent[0], "", "    ")
	items1stRecord := string(json1stItem)

	json2stItem, _ := json.MarshalIndent(tableItemsContent[0], "", "    ")
	items2stRecord := string(json2stItem)

	newItem := map[string]string{"id": "42", "title": "db_crud", "description": ""}
	newItemBytes, _ := json.MarshalIndent(newItem, "", "    ")
	newItemString := string(newItemBytes)

	updatingData := map[string]string{"description": "Написать программу db_crud"}

	updatedItem := map[string]string{"id": "3", "title": "db_crud", "description": "Написать программу db_crud"}
	updatedItemBytes, _ := json.MarshalIndent(updatedItem, "", "    ")
	updatedItemString := string(updatedItemBytes)

	updatingData2 := map[string]string{"updated": "autotests"}

	updatedItem2 := map[string]string{"id": "3", "title": "db_crud", "description": "", "updated": "autotests"}
	updatedItemBytes2, _ := json.MarshalIndent(updatedItem2, "", "    ")
	updatedItemString2 := string(updatedItemBytes2)

	finaleItem := map[string]string{"id": "3", "title": "db_crud"}
	finaleItemBytes, _ := json.MarshalIndent(finaleItem, "", "    ")
	finaleItemString := string(finaleItemBytes)

	userRVasiliy := map[string]interface{}{
		"user_id":  1,
		"login":    "rvasily",
		"password": "love",
		"email":    "rvasily@example.com",
		"info":     "none",
	}
	userRVasiliyBytes, _ := json.MarshalIndent(userRVasiliy, "", "    ")
	userRVasiliyString := string(userRVasiliyBytes)

	updatedVasiliy := map[string]interface{}{
		"user_id":  1,
		"login":    "rvasily",
		"password": "love",
		"email":    "rvasily@example.com",
		"info":     "none",
	}
	updatedVasiliyBytes, _ := json.MarshalIndent(updatedVasiliy, "", "    ")
	updatedVasiliyString := string(updatedVasiliyBytes)

	users := []map[string]interface{}{
		{
			"user_id":  1,
			"login":    "rvasily",
			"password": "love",
			"email":    "rvasily@example.com",
			"info":     "try update",
			"updated":  "now",
		},
		{
			"user_id":  2,
			"login":    "qwerty'",
			"password": "love\"",
			"email":    "",
			"info":     "",
		},
	}
	usersBytes, _ := json.MarshalIndent(users, "", "    ")
	usersString := string(usersBytes)

	cases := []testCase{
		{
			name:                 "tables list",
			path:                 "/",
			expectedResponseBody: "[\n    \"items_test\",\n    \"users_test\"\n]",
		},
		{
			name:                   "unknown_table",
			path:                   "/unknown_table",
			expectedResponseStatus: http.StatusNotFound,
			expectedResponseBody:   "unknown table",
		},
		{
			name:                 "items_test",
			path:                 "/items_test",
			expectedResponseBody: itemsData,
		},
		{
			name:                 "limit 1",
			path:                 "/items_test",
			queryParams:          "limit=1",
			expectedResponseBody: items1stRecord,
		},
		{
			name:                 "limit 1 offset 1",
			path:                 "/items_test",
			queryParams:          "limit=1&offset=1",
			expectedResponseBody: items2stRecord,
		},
		{
			name:                 "id 1",
			path:                 "/items_test/1",
			expectedResponseBody: items1stRecord,
		},
		{
			name:                   "record not found",
			path:                   "/items_test/100500",
			expectedResponseStatus: http.StatusNotFound,
			expectedResponseBody:   "record not found",
		},
		{
			name:                 "new record",
			path:                 "/items_test/",
			method:               http.MethodPut,
			requestBody:          newItem,
			expectedResponseBody: "last insert id 3",
		},
		{
			name:                 "id 3",
			path:                 "/items_test/3",
			expectedResponseBody: newItemString,
		},
		{
			path:                 "/items_test/3",
			method:               http.MethodPost,
			requestBody:          updatingData,
			expectedResponseBody: "updated record id 3",
		},
		{
			name:                 "updated id 3",
			path:                 "/items_test/3",
			expectedResponseBody: updatedItemString,
		},
		{
			name:                 "update null",
			path:                 "/items_test/3",
			method:               http.MethodPost,
			requestBody:          updatingData2,
			expectedResponseBody: "updated record id 3",
		},
		{
			name:                 "updated id 3 second time",
			path:                 "/items_test/3",
			expectedResponseBody: updatedItemString2,
		},
		{
			name:                 "set null",
			path:                 "/items_test/3",
			method:               http.MethodPost,
			requestBody:          map[string]string{"updated": "%00"},
			expectedResponseBody: "updated record id 3",
		},
		{
			name:                 "updated id 3 third time",
			path:                 "/items_test/3",
			expectedResponseBody: finaleItemString,
		},
		{
			name:                   "try update primary key",
			path:                   "/items_test/3",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody:            map[string]string{"id": "4"}, // primary key нельзя обновлять у существующей записи
			expectedResponseBody:   "unable to update record",
		},
		{
			name:                   "try update float with int",
			path:                   "/items/3",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody:            map[string]string{"rating": "15"}, // int -> float
			expectedResponseBody:   "invalid type rating",
		},
		{
			name:                   "try update float with string",
			path:                   "/items/3",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody:            map[string]string{"rating": "string"}, // string -> float
			expectedResponseBody:   "invalid type rating",
		},
		{
			name:                   "try update bool with string",
			path:                   "/items/3",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody:            map[string]string{"shadow": "string"}, // string -> bool
			expectedResponseBody:   "invalid type shadow",
		},
		{
			name:                   "try update int with bool",
			path:                   "/items/3",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody:            map[string]string{"level": "true"}, // bool -> int
			expectedResponseBody:   "invalid type level",
		},
		{
			name:                   "try set null",
			path:                   "/items_test/3",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody:            map[string]string{"title": "%00"},
			expectedResponseBody:   "unable to update record",
		},
		{
			name:                 "delete",
			path:                 "/items_test/3",
			method:               http.MethodDelete,
			expectedResponseBody: "deleted record id 3",
		},
		{
			name:                   "delete deleted",
			path:                   "/items_test/3",
			method:                 http.MethodDelete,
			expectedResponseStatus: http.StatusNotFound,
			expectedResponseBody:   "record not found",
		},
		{
			name:                   "deleted not found",
			path:                   "/items_test/3",
			expectedResponseStatus: http.StatusNotFound,
			expectedResponseBody:   "record not found",
		},

		// и немного по другой таблице
		{
			path:                 "/users_test/1",
			expectedResponseBody: userRVasiliyString,
		},

		{
			path:   "/users_test/1",
			method: http.MethodPost,
			requestBody: map[string]string{
				"info":    "try update",
				"updated": "now",
			},
			expectedResponseBody: "updated record id 1",
		},
		{
			path:                 "/users_test/1",
			expectedResponseBody: updatedVasiliyString,
		},
		// ошибки
		{
			path:                   "/users_test/1",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody: map[string]string{
				"user_id": "1", // primary key нельзя обновлять у существующей записи
			},
			expectedResponseBody: "unable to update record",
		},
		// не забываем про sql-инъекции
		{
			path:   "/users_test/",
			method: http.MethodPut,
			requestBody: map[string]string{
				"user_id":    "2",
				"login":      "qwerty'",
				"password":   "love\"",
				"unkn_field": "love",
			},
			expectedResponseBody: "unable to insert record",
		},
		{
			path:                   "/users_test/2",
			expectedResponseStatus: http.StatusNotFound,
			expectedResponseBody:   "record not found",
		},
		{
			path:   "/users_test/",
			method: http.MethodPut,
			requestBody: map[string]string{
				"user_id":    "2",
				"login":      "qwerty'",
				"password":   "love\"",
				"unkn_field": "love",
				"email":      "",
				"info":       "",
			},
			expectedResponseBody: "unable to insert record",
		},

		// тут тоже возможна sql-инъекция
		// если пришло не число на вход - берём дефолтное значене для лимита-оффсета
		{
			path:                 "/users_test",
			queryParams:          "limit=1'&offset=1\"",
			expectedResponseBody: usersString,
		},
	}

	runCases(t, ts, db, cases)
}

func runCases(t *testing.T, ts *httptest.Server, db *sql.DB, cases []testCase) {
	for idx, item := range cases {

		caseName := fmt.Sprintf("case %d: [%s] %s %s", idx, item.method, item.path, item.queryParams)

		if db.Stats().OpenConnections > 2 {
			t.Fatalf("[%s] you have %d open connections, must be 2 or less", caseName, db.Stats().OpenConnections)
		}

		params := url.Values{}
		for k, v := range item.requestBody {
			params.Add(k, v)
		}
		req := httptest.NewRequest(item.method, ts.URL+item.path+item.queryParams, bytes.NewBufferString(params.Encode()))

		if item.method == http.MethodPut || item.method == http.MethodPost {
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		}

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		if item.expectedResponseStatus == 0 {
			item.expectedResponseStatus = http.StatusOK
		}

		assert.Equal(t, item.expectedResponseStatus, resp.StatusCode)

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		assert.Equal(t, item.expectedResponseBody, buf.String())
	}

}
