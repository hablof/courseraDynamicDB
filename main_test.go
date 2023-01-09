package main

import (
	"database/sql"
	"fmt"
	"hw6coursera/dbexplorer"
	"hw6coursera/repository"
	"hw6coursera/router"
	"hw6coursera/service"
	"log"
	"reflect"
	"testing"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
  		PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`INSERT INTO items (title, description, updated) VALUES
		('database/sql',	'Рассказать про базы данных',	'rvasily'),
		('memcache',	'Рассказать про мемкеш с примером использования',	NULL);`,

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

		`INSERT INTO users (login, password, email, info, updated) VALUES
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

	updatedItem := map[string]string{"id": "42", "title": "db_crud", "description": ""}
	updatedItemBytes, _ := json.MarshalIndent(updatedItem, "", "    ")
	updatedItemString := string(updatedItemBytes)

	updatingData2 := map[string]string{"updated": "autotests"}

	updatedItem2 := map[string]string{"id": "42", "title": "db_crud", "description": "", "updated": "autotests"}
	updatedItemBytes2, _ := json.MarshalIndent(updatedItem2, "", "    ")
	updatedItemString2 := string(updatedItemBytes2)

	cases := []testCase{
		{
			name:                 "список таблиц",
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
			path:                 "/items",
			queryParams:          "limit=1",
			expectedResponseBody: items1stRecord,
		},
		{
			path:                 "/items",
			queryParams:          "limit=1&offset=1",
			expectedResponseBody: items2stRecord,
		},
		{
			path:                 "/items/1",
			expectedResponseBody: items1stRecord,
		},
		{
			path:                   "/items/100500",
			expectedResponseStatus: http.StatusNotFound,
			expectedResponseBody:   "record not found",
		},

		// тут идёт создание и редактирование
		{
			path:                 "/items/",
			method:               http.MethodPut,
			requestBody:          newItem,
			expectedResponseBody: "last insert id 3",
		},
		// это пример хрупкого теста
		// если много раз вызывать один и тот же тест - записи будут добавляться
		// поэтому придётся сделать сброс базы каждый раз в PrepareTestData
		{
			path:                 "/items/3",
			expectedResponseBody: newItemString,
		},
		{
			path:                 "/items/3",
			method:               http.MethodPost,
			requestBody:          updatingData,
			expectedResponseBody: "updated record id 3",
		},
		{
			path:                 "/items/3",
			expectedResponseBody: updatedItemString,
		},

		// обновление null-поля в таблице
		{
			path:                 "/items/3",
			method:               http.MethodPost,
			requestBody:          updatingData2,
			expectedResponseBody: "updated record id 3",
		},
		{
			path:                 "/items/3",
			expectedResponseBody: updatedItemString2,
		},

		// обновление null-поля в таблице
		{
			path:   "/items/3",
			method: http.MethodPost,
			requestBody: CR{
				"updated": nil,
			},
			expectedResponseBody: "updated record id 3",
		},
		{
			path: "/items/3",
			expectedResponseBody: CR{
				"response": CR{
					"record": CR{
						"id":          3,
						"title":       "db_crud",
						"description": "Написать программу db_crud",
						"updated":     nil,
					},
				},
			},
		},

		// ошибки
		{
			path:                   "/items/3",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody: CR{
				"id": 4, // primary key нельзя обновлять у существующей записи
			},
			expectedResponseBody: CR{
				"error": "field id have invalid type",
			},
		},
		{
			path:                   "/items/3",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody: CR{
				"title": 42,
			},
			expectedResponseBody: CR{
				"error": "field title have invalid type",
			},
		},
		{
			path:                   "/items/3",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody: CR{
				"title": nil,
			},
			expectedResponseBody: CR{
				"error": "field title have invalid type",
			},
		},

		{
			path:                   "/items/3",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody: CR{
				"updated": 42,
			},
			expectedResponseBody: CR{
				"error": "field updated have invalid type",
			},
		},

		// удаление
		{
			path:   "/items/3",
			method: http.MethodDelete,
			expectedResponseBody: CR{
				"response": CR{
					"deleted": 1,
				},
			},
		},
		{
			path:   "/items/3",
			method: http.MethodDelete,
			expectedResponseBody: CR{
				"response": CR{
					"deleted": 0,
				},
			},
		},
		{
			path:                   "/items/3",
			expectedResponseStatus: http.StatusNotFound,
			expectedResponseBody: CR{
				"error": "record not found",
			},
		},

		// и немного по другой таблице
		{
			path: "/users/1",
			expectedResponseBody: CR{
				"response": CR{
					"record": CR{
						"user_id":  1,
						"login":    "rvasily",
						"password": "love",
						"email":    "rvasily@example.com",
						"info":     "none",
						"updated":  nil,
					},
				},
			},
		},

		{
			path:   "/users/1",
			method: http.MethodPost,
			requestBody: CR{
				"info":    "try update",
				"updated": "now",
			},
			expectedResponseBody: CR{
				"response": CR{
					"updated": 1,
				},
			},
		},
		{
			path: "/users/1",
			expectedResponseBody: CR{
				"response": CR{
					"record": CR{
						"user_id":  1,
						"login":    "rvasily",
						"password": "love",
						"email":    "rvasily@example.com",
						"info":     "try update",
						"updated":  "now",
					},
				},
			},
		},
		// ошибки
		{
			path:                   "/users/1",
			method:                 http.MethodPost,
			expectedResponseStatus: http.StatusBadRequest,
			requestBody: CR{
				"user_id": 1, // primary key нельзя обновлять у существующей записи
			},
			expectedResponseBody: CR{
				"error": "field user_id have invalid type",
			},
		},
		// не забываем про sql-инъекции
		{
			path:   "/users/",
			method: http.MethodPut,
			requestBody: CR{
				"user_id":    2,
				"login":      "qwerty'",
				"password":   "love\"",
				"unkn_field": "love",
			},
			expectedResponseBody: CR{
				"response": CR{
					"user_id": 2,
				},
			},
		},
		{
			path: "/users/2",
			expectedResponseBody: CR{
				"response": CR{
					"record": CR{
						"user_id":  2,
						"login":    "qwerty'",
						"password": "love\"",
						"email":    "",
						"info":     "",
						"updated":  nil,
					},
				},
			},
		},
		// тут тоже возможна sql-инъекция
		// если пришло не число на вход - берём дефолтное значене для лимита-оффсета
		{
			path:        "/users",
			queryParams: "limit=1'&offset=1\"",
			expectedResponseBody: CR{
				"response": CR{
					"records": []CR{
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
							"updated":  nil,
						},
					},
				},
			},
		},
	}

	runCases(t, ts, db, cases)
}

func runCases(t *testing.T, ts *httptest.Server, db *sql.DB, cases []Case) {
	for idx, item := range cases {
		var (
			err      error
			result   interface{}
			expected interface{}
			req      *http.Request
		)

		caseName := fmt.Sprintf("case %d: [%s] %s %s", idx, item.Method, item.Path, item.QueryParams)

		// если у вас случилась это ошибка - значит вы не делаете где-то rows.Close и у вас текут соединения с базой
		// если такое случилось на первом тесте - значит вы не закрываете коннект где-то при иницаилизации в NewDbExplorer
		if db.Stats().OpenConnections > 2 {
			t.Fatalf("[%s] you have %d open connections, must be 2 or less", caseName, db.Stats().OpenConnections)
		}

		if item.Method == "" || item.Method == http.MethodGet {
			req, err = http.NewRequest(item.Method, ts.URL+item.Path+item.Query, nil)
		} else {
			data, err := json.Marshal(item.Body)
			if err != nil {
				panic(err)
			}
			reqBody := bytes.NewReader(data)
			req, err = http.NewRequest(item.Method, ts.URL+item.Path, reqBody)
			req.Header.Add("Content-Type", "application/json")
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("[%s] request error: %v", caseName, err)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		// fmt.Printf("[%s] requestBody: %s\n", caseName, string(body))
		if item.Status == 0 {
			item.Status = http.StatusOK
		}

		if resp.StatusCode != item.Status {
			t.Errorf("[%s] expected http status %v, got %v", caseName, item.Status, resp.StatusCode)
			continue
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Fatalf("[%s] cant unpack json: %v", caseName, err)
			continue
		}

		// reflect.DeepEqual не работает если нам приходят разные типы
		// а там приходят разные типы (string VS interface{}) по сравнению с тем что в ожидаемом результате
		// этот маленький грязный хак конвертит данные сначала в json, а потом обратно в interface - получаем совместимые результаты
		// не используйте это в продакшен-коде - надо явно писать что ожидается интерфейс или использовать другой подход с точным форматом ответа
		data, err := json.Marshal(item.Result)
		json.Unmarshal(data, &expected)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("[%s] results not match\nGot : %#v\nWant: %#v", caseName, result, expected)
			continue
		}
	}

}
