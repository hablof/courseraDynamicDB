package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

const (
	defaultLimit  = 5
	defaultOffset = 0
)

type table struct {
	name    string
	columns []column
}

type column struct {
	name         string
	columnType   reflect.Type
	nullable     bool
	isPrimaryKey bool
}

type databaseHandler struct {
	schema            map[string]table
	db                *sql.DB
	tableAndIdPattern *regexp.Regexp
	tablePattern      *regexp.Regexp
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	schema, err := parseSchema(db)
	if err != nil {
		return nil, err
	}

	h := initRoutes(db, schema)
	return h, nil
}

// Получаем мапу[название таблицы]её структура ()
func parseSchema(db *sql.DB) (map[string]table, error) {
	log.Println("getting tables...")
	tableRecords, err := db.Query(`SHOW TABLES`)

	if err != nil {
		return nil, err
	}

	tables := make(map[string]table)

	for tableRecords.Next() {
		t := table{}
		err = tableRecords.Scan(&t.name)
		if err != nil {
			return nil, err
		}
		tables[t.name] = t
	}
	err = tableRecords.Close()
	if err != nil {
		return nil, err
	}

	log.Println("parsing colunms...")
	for name := range tables {
		log.Printf("in table: %s...", name)
		rows, err := db.Query(fmt.Sprintf("SELECT * FROM `%s`", name))

		if err != nil {
			return nil, err
		}

		sqlColumns, err := rows.ColumnTypes()
		if err != nil {
			return nil, err
		}
		err = rows.Close()
		if err != nil {
			return nil, err
		}

		//	Получаем информацию о столбце, который язляется 'PRIMARY KEY'
		//	Например
		//	+-------+------+------+-----+---------+----------------+
		//	| Field | Type | Null | Key | Default | Extra          |
		//	+-------+------+------+-----+---------+----------------+
		//	| id    | int  | NO   | PRI | NULL    | auto_increment |
		//	+-------+------+------+-----+---------+----------------+
		row := db.QueryRow(fmt.Sprintf("SHOW COLUMNS FROM %s WHERE `Key` = 'PRI';", name))
		err = row.Err()
		if err != nil {
			return nil, err
		}
		buf := []byte{}
		waste := new(interface{})
		err = row.Scan(&buf, waste, waste, waste, waste, waste)
		if err != nil {
			return nil, err
		}
		PrimaryKey := string(buf)

		cols := make([]column, 0)
		for _, ct := range sqlColumns {
			c := column{}
			c.name = ct.Name()
			if c.name == PrimaryKey {
				c.isPrimaryKey = true
			}
			nullable, ok := ct.Nullable()
			if !ok {
				return nil, errors.New("shitty driver")
			}
			c.nullable = nullable
			c.columnType = ct.ScanType()
			cols = append(cols, c)
		}
		t := tables[name]
		t.columns = cols
		tables[name] = t
	}
	return tables, nil
}

func initRoutes(db *sql.DB, schema map[string]table) http.Handler {
	mymux := http.NewServeMux()
	tableAndIdPattern := regexp.MustCompile(`\A\/\w+\/\d+\/?\z`) //	"/букво-цифры/цифры и, может быть, слэш"
	tablePattern := regexp.MustCompile(`\A\/\w+\/?\z`)           // "/букво-цифры и, может быть, слэш"

	dbHandler := databaseHandler{
		schema:            schema,
		db:                db,
		tableAndIdPattern: tableAndIdPattern,
		tablePattern:      tablePattern,
	}
	mymux.Handle("/", dbHandler)
	return mymux
}

func (dbh databaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch {
	case dbh.tablePattern.MatchString(r.RequestURI):
		switch r.Method {
		case "GET":
			dbh.GetRecords(w, r)
		case "PUT":
			dbh.InsertRecord(w, r)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	case dbh.tableAndIdPattern.MatchString(r.RequestURI):
		switch r.Method {
		case "GET":
			dbh.GetSingleRecord(w, r)
		case "PUT":
			dbh.UpdateRecord(w, r)
		case "DELETE":
			dbh.DeleteRecord(w, r)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		dbh.GetAllTables(w, r)
	}
}

func (dbh databaseHandler) GetAllTables(w http.ResponseWriter, r *http.Request) {
	log.Println("get all tables request")
	tablesList := make([]string, 0)
	for n := range dbh.schema {
		tablesList = append(tablesList, n)
	}

	b, err := json.Marshal(tablesList)
	if err != nil {
		log.Printf("got error marshaling tables list: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application-json")
	w.Write(b)
}

func (dbh databaseHandler) GetRecords(w http.ResponseWriter, r *http.Request) {
	tableName := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), "/")
	log.Printf("getting records from table %s", tableName)
	tableStruct, ok := dbh.schema[tableName]
	if !ok {
		log.Printf("table %s not found", tableName)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("table %s not found", tableName)))
		return
	}

	//  uri: /$table?limit=5&offset=7
	limit := defaultLimit
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			// Просто логгируемся, не крашимся
			log.Printf("got error parsing to int limit value (%s): %+v\n", limitStr, err)
		} else {
			limit = l
		}
	}

	offset := defaultOffset
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err != nil {
			log.Printf("got error parsing to int limit value (%s): %+v\n", offsetStr, err)
		} else {
			offset = o
		}
	}

	fields := getQueryFields(tableStruct)
	q := `SELECT %s
	FROM %s
	LIMIT ?
	OFFSET ?;`
	query := fmt.Sprintf(q, fields, tableName)

	rows, err := dbh.db.Query(query, limit, offset)
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		log.Printf("unable to get records from sql due to error: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to get records"))
		return
	}

	content := make([]map[string]interface{}, 0)
	dest := initScanDestination(tableStruct)
	for rows.Next() {
		rows.Scan(dest...)
		unit := extractSqlVals(tableStruct, dest)
		content = append(content, unit)
	}
	b, err := json.MarshalIndent(content, "", "    ")
	if err != nil {
		log.Printf("got error marshaling %s table content: %+v", tableName, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	log.Printf("successfully...")
}

func (dbh databaseHandler) InsertRecord(w http.ResponseWriter, r *http.Request) {
	tableName := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), "/")
	log.Printf("inserting record to table %s\n", tableName)
	tableStruct, ok := dbh.schema[tableName]
	if !ok {
		log.Printf("table %s not found", tableName)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("table %s not found", tableName)))
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	urlVals := r.PostForm
	unit := make(map[string]interface{})

	for _, c := range tableStruct.columns {
		if c.isPrimaryKey {
			continue
		}
		if urlVals.Has(c.name) {
			unit[c.name] = urlVals.Get(c.name)
		} else {
			if !c.nullable {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	fields, placehoders, sqlVals := getInsertParams(unit)
	q := `
	INSERT INTO %s (%s)
	VALUES (%s);`
	query := fmt.Sprintf(q, tableName, fields, placehoders)
	res, err := dbh.db.Exec(query, sqlVals...)
	if err != nil {
		log.Printf("error on inserting values: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	i, err := res.RowsAffected()
	if err != nil {
		log.Printf("error on RowsAffected(): %v", err)
	}
	if i != 1 {
		log.Printf("wrong behaviour, affectet rows: %d\n", i)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("error on LastInsertId(): %v", err)
	}
	log.Printf("last insert id %d", id)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("last insert id %d", id)))
}

func (dbh databaseHandler) GetSingleRecord(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	tableName := path[0]
	tableStruct, ok := dbh.schema[tableName]
	if !ok {
		log.Printf("table not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("table not found"))
		return
	}
	id, err := strconv.Atoi(path[1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("getting record (id=%d) from table %s...", id, tableName)

	fields := getQueryFields(tableStruct)
	PK, err := getPKColumnName(tableStruct)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	q := `SELECT %s
	FROM %s
	WHERE %s = ?;`
	query := fmt.Sprintf(q, fields, tableName, PK)

	row := dbh.db.QueryRow(query, id)
	err = row.Err()
	if err != nil {
		log.Printf("unable to get records from sql due to error: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to get records"))
		return
	}
	dest := initScanDestination(tableStruct)
	err = row.Scan(dest...)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("record (id=%d) not found", id)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("record (id=%d) not found", id)))
			return
		}
		log.Printf("unable to scan sql row due to error: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to get records"))
		return
	}

	unit := extractSqlVals(tableStruct, dest)

	b, err := json.MarshalIndent(unit, "", "    ")
	if err != nil {
		log.Printf("got error marshaling %s table content: %+v", tableName, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	log.Printf("successfully...")
}

func (dbh databaseHandler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	tableName := path[0]
	tableStruct, ok := dbh.schema[tableName]
	if !ok {
		log.Printf("table not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("table not found"))
		return
	}
	id, err := strconv.Atoi(path[1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("updating record (id=%d) from table %s", id, tableName)

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	urlVals := r.PostForm
	unit := make(map[string]interface{})
	PK := ""
	for _, c := range tableStruct.columns {
		if c.isPrimaryKey {
			PK = c.name
			continue
		}
		if urlVals.Has(c.name) {
			unit[c.name] = urlVals.Get(c.name)
		}
	}
	keyValues, sqlVals := getUpdateParams(unit)
	q := `
	UPDATE %s
	SET %s
	WHERE %s = ?;`
	query := fmt.Sprintf(q, tableName, keyValues, PK)
	sqlVals = append(sqlVals, id)
	res, err := dbh.db.Exec(query, sqlVals...)
	if err != nil {
		log.Printf("error on updating values: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	i, err := res.RowsAffected()
	if err != nil {
		log.Printf("error on rowsaffected(): %v", err)
	}
	if i == 0 {
		log.Println("no rows affected")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("no records updated"))
		return
	} else if i > 1 {
		log.Println("affected more then 1 row")
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("updated record id %d", id)))
}

func (dbh databaseHandler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	tableName := path[0]
	tableStruct, ok := dbh.schema[tableName]
	if !ok {
		log.Printf("table not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("table not found"))
		return
	}
	id, err := strconv.Atoi(path[1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("deleting record (id=%d) from table %s", id, tableName)

	PK := ""
	for _, c := range tableStruct.columns {
		if c.isPrimaryKey {
			PK = c.name
			break
		}
	}
	q := `
	DELETE FROM %s
	WHERE %s = ?;`
	query := fmt.Sprintf(q, tableName, PK)

	res, err := dbh.db.Exec(query, id)
	if err != nil {
		log.Printf("error on deleting values: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	i, err := res.RowsAffected()
	if err != nil {
		log.Printf("error on rowsaffected(): %v", err)
	}
	if i == 0 {
		log.Println("no rows affected")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("no records deleted"))
		return
	} else if i > 1 {
		log.Println("affected more then 1 row")
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("deleted record id %d", id)))
}

func initScanDestination(t table) []interface{} {
	a := make([]interface{}, len(t.columns))
	for i := range a {
		a[i] = new(interface{})
	}
	return a
}

func extractSqlVals(tableStruct table, dest []interface{}) map[string]interface{} {
	unit := make(map[string]interface{})
	for i, c := range tableStruct.columns {
		v1 := reflect.ValueOf(dest[i])
		v2 := reflect.Indirect(v1)
		v3 := v2.Interface()
		switch v3 := v3.(type) {
		case []byte:
			unit[c.name] = string(v3)
		default:
			unit[c.name] = v3
		}
	}
	return unit
}

func getQueryFields(t table) string {
	switch len(t.columns) {
	case 0:
		return ""
	case 1:
		return t.columns[0].name
	}
	var sb strings.Builder
	sb.WriteString(t.columns[0].name)
	for _, col := range t.columns[1:] {
		sb.WriteString(", ")
		sb.WriteString(col.name)
	}
	return sb.String()
}

func getInsertParams(unit map[string]interface{}) (fieldNames string, placehoderStr string, data []interface{}) {
	l := len(unit)
	names := make([]string, 0, l)
	placehoders := make([]string, 0, l)
	output := make([]interface{}, 0, l)
	for k, v := range unit {
		names = append(names, k)
		placehoders = append(placehoders, "?")
		output = append(output, v)
	}
	return strings.Join(names, ", "), strings.Join(placehoders, ", "), output
}

func getUpdateParams(unit map[string]interface{}) (fieldPlacehoderStr string, data []interface{}) {
	l := len(unit)
	placehoders := make([]string, 0, l)
	output := make([]interface{}, 0, l)
	for k, v := range unit {
		placehoders = append(placehoders, fmt.Sprintf("%s = ?", k))
		output = append(output, v)
	}
	return strings.Join(placehoders, ", "), output
}

func getPKColumnName(t table) (string, error) {
	for _, c := range t.columns {
		if c.isPrimaryKey {
			return c.name, nil
		}
	}
	return "", errors.New("there is no primary key column")
}
