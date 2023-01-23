package main

import (
	"database/sql"
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

var (
	client = &http.Client{Timeout: time.Second}
)

func prepareTestApis(db *sql.DB) error {

	qs := []string{
		`DROP TABLE IF EXISTS items_test;`,

		`CREATE TABLE items_test (
  		id int(11) NOT NULL AUTO_INCREMENT,
  		title varchar(255) NOT NULL,
  		description text NOT NULL,
  		updated varchar(255) DEFAULT NULL,
		rating decimal(5,2),
		level int,
  		PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`INSERT INTO items_test (title, description, updated, rating, level) VALUES
		('database/sql', 'Рассказать про базы данных', 'rvasily', '2.71828182', '15'),
		('memcache', 'Рассказать про мемкеш с примером использования', NULL, '0.0', '80');`,

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
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func cleanupTestApis(db *sql.DB) error {
	qs := []string{
		"DROP TABLE IF EXISTS items_test;",
		"DROP TABLE IF EXISTS users_test;",
	}
	for _, q := range qs {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func TestApis(t *testing.T) {
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3366)/integration_testing")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := prepareTestApis(db); err != nil {
		log.Printf("failed to prepare tests: %v", err)
		assert.Equal(t, nil, err)
		return
	}

	defer cleanupTestApis(db)

	repo := repository.NewRepository(db)
	explorer := dbexplorer.NewDbExplorer(repo)
	service := service.NewService(repo, explorer)
	if err := service.InitSchema(); err != nil {
		log.Printf("failed to init database shcema: %v", err)
		assert.Equal(t, nil, err)
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
			"level":       15,
			"rating":      2.72,
		},
		{
			"id":          2,
			"title":       "memcache",
			"description": "Рассказать про мемкеш с примером использования",
			"level":       80,
			"rating":      0,
		},
	}
	jsonItems, _ := json.MarshalIndent(tableItemsContent, "", "    ")
	itemsData := string(jsonItems)

	jsonSlice1stItem, _ := json.MarshalIndent(tableItemsContent[:1], "", "    ")
	slice1stItem := string(jsonSlice1stItem)

	jsonSlice2stItem, _ := json.MarshalIndent(tableItemsContent[1:], "", "    ")
	slice2stItem := string(jsonSlice2stItem)

	json1stItem, _ := json.MarshalIndent(tableItemsContent[0], "", "    ")
	items1stRecord := string(json1stItem)

	newItem := map[string]string{"id": "42", "title": "db_crud", "description": ""}
	newItemCorrectId := map[string]interface{}{"id": 3, "title": "db_crud", "description": ""}
	newItemBytes, _ := json.MarshalIndent(newItemCorrectId, "", "    ")
	newItemString := string(newItemBytes)

	updatedItem := map[string]interface{}{"id": 3, "title": "db_crud", "description": "Написать программу db_crud"}
	updatedItemBytes, _ := json.MarshalIndent(updatedItem, "", "    ")
	updatedItemString := string(updatedItemBytes)

	updatedItem2 := map[string]interface{}{"id": 3, "title": "db_crud", "updated": "autotests", "description": "Написать программу db_crud"}
	updatedItemBytes2, _ := json.MarshalIndent(updatedItem2, "", "    ")
	updatedItemString2 := string(updatedItemBytes2)

	finaleItem := map[string]interface{}{"id": 3, "title": "db_crud", "description": "Написать программу db_crud"}
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
		"info":     "try update",
		"updated":  "now",
	}
	updatedVasiliyBytes, _ := json.MarshalIndent(updatedVasiliy, "", "    ")
	updatedVasiliyString := string(updatedVasiliyBytes)

	sqlGuy := map[string]interface{}{
		"user_id":  3,
		"login":    "petya",
		"password": "pass",
		"email":    "pochta@yandex.ru",
		"info":     "info); DELETE FROM table WHERE 1=1; now(",
	}
	sqlGuyBytes, _ := json.MarshalIndent(sqlGuy, "", "    ")
	sqlGuyString := string(sqlGuyBytes)

	// users := []map[string]interface{}{
	// 	{
	// 		"user_id":  1,
	// 		"login":    "rvasily",
	// 		"password": "love",
	// 		"email":    "rvasily@example.com",
	// 		"info":     "try update",
	// 		"updated":  "now",
	// 	},
	// 	{
	// 		"user_id":  2,
	// 		"login":    "qwerty'",
	// 		"password": "love\"",
	// 		"email":    "",
	// 		"info":     "",
	// 	},
	// }
	// usersBytes, _ := json.MarshalIndent(users, "", "    ")
	// usersString := string(usersBytes)

	testCases := []struct {
		name                   string
		method                 string // GET по-умолчанию в http.NewRequest если передали пустую строку
		path                   string
		queryParams            string
		expectedResponseStatus int
		expectedResponseBody   string
		requestBody            map[string]string
	}{
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
			queryParams:          "?limit=1",
			expectedResponseBody: slice1stItem,
		},
		{
			name:                 "limit 1 offset 1",
			path:                 "/items_test",
			queryParams:          "?limit=1&offset=1",
			expectedResponseBody: slice2stItem,
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
			requestBody:          map[string]string{"description": "Написать программу db_crud"},
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
			requestBody:          map[string]string{"updated": "autotests"},
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
			expectedResponseBody:   "missing data to update",
		},
		{
			name:                   "try update float with int",
			path:                   "/items_test/3",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusOK,
			requestBody:            map[string]string{"rating": "15"}, // int -> float
			expectedResponseBody:   "updated record id 3",
		},
		{
			name:                   "try update float with string",
			path:                   "/items_test/3",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody:            map[string]string{"rating": "string"}, // string -> float
			expectedResponseBody:   "invalid type rating",
		},
		{
			name:                   "try update int with bool",
			path:                   "/items_test/3",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody:            map[string]string{"level": "true"}, // bool -> int
			expectedResponseBody:   "invalid type level",
		},
		{
			name:                   "try set null to not-null field",
			path:                   "/items_test/3",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody:            map[string]string{"title": "%00"},
			expectedResponseBody:   "title cannot be null",
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

		{
			name:                 "get R Vasiliy",
			path:                 "/users_test/1",
			expectedResponseBody: userRVasiliyString,
		},
		{
			name:   "update Vasiliy",
			path:   "/users_test/1",
			method: http.MethodPost,
			requestBody: map[string]string{
				"info":    "try update",
				"updated": "now",
			},
			expectedResponseBody: "updated record id 1",
		},
		{
			name:                 "updated Vasiliy",
			path:                 "/users_test/1",
			expectedResponseBody: updatedVasiliyString,
		},
		{
			name:                   "try update user id",
			path:                   "/users_test/1",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody: map[string]string{
				"user_id": "1",
			},
			expectedResponseBody: "missing data to update",
		},
		{
			name:   "SQL injection",
			path:   "/users_test/",
			method: http.MethodPut,
			requestBody: map[string]string{
				"user_id":    "2",
				"login":      "qwerty'",
				"password":   "love\"",
				"unkn_field": "love",
				"email":      "tosi-bosi@ya.ru",
				"info":       "крокодилы ходят лёжа",
			},
			expectedResponseBody: "last insert id 2",
		},
		{
			name:   "SQL injection pro",
			path:   "/users_test/",
			method: http.MethodPut,
			requestBody: map[string]string{
				"login":    "petya",
				"password": "pass",
				"email":    "pochta@yandex.ru",
				"info":     "info); DELETE FROM table WHERE 1=1; now(",
			},
			expectedResponseBody: "last insert id 3",
		},
		{
			name:                 "check injections",
			path:                 "/users_test/3",
			expectedResponseBody: sqlGuyString,
			// expectedResponseBody: "[]",
		},
		{
			name:                   "user not found",
			path:                   "/users_test/4",
			expectedResponseStatus: http.StatusNotFound,
			expectedResponseBody:   "record not found",
		},
		{
			name:   "insert without email and info",
			path:   "/users_test/",
			method: http.MethodPut,
			requestBody: map[string]string{
				"user_id":    "2",
				"login":      "qwerty'",
				"password":   "love\"",
				"unkn_field": "love",
			},
			expectedResponseStatus: http.StatusBadRequest,
			expectedResponseBody:   "email cannot be null",
		},
		{
			name:                   "SQL injection 2",
			path:                   "/users_test",
			queryParams:            "?limit=1'&offset=1\"",
			expectedResponseStatus: http.StatusNotFound,
			expectedResponseBody:   "page not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// caseName := fmt.Sprintf("case %d: [%s] %s %s", idx, item.method, item.path, item.queryParams)

			if db.Stats().OpenConnections > 2 {
				t.Fatalf("[%s] you have %d open connections, must be 2 or less", tc.name, db.Stats().OpenConnections)
			}

			params := url.Values{}
			for k, v := range tc.requestBody {
				params.Add(k, v)
			}
			req, err := http.NewRequest(tc.method, ts.URL+tc.path+tc.queryParams, bytes.NewBufferString(params.Encode()))
			if err != nil {
				log.Printf("do req err: %v", err)
			}

			if tc.method == http.MethodPut || tc.method == http.MethodPost {
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			}

			resp, err := client.Do(req)
			assert.Equal(t, nil, err)
			if err != nil {
				log.Printf("do req err: %v", err)
			} else {
				defer resp.Body.Close()
			}
			if tc.expectedResponseStatus == 0 {
				tc.expectedResponseStatus = http.StatusOK
			}

			assert.Equal(t, tc.expectedResponseStatus, resp.StatusCode)

			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			assert.Equal(t, tc.expectedResponseBody, buf.String())
		})
	}
}
