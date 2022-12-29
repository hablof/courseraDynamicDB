package repository

import (
	"database/sql"
	"fmt"
	"hw6coursera/internal"
	"log"
	"reflect"
	"strings"
)

type recordManager struct {
	db *sql.DB
}

// Create implements RecordManager
func (rm *recordManager) Create(table internal.Table, data map[string]interface{}) (lastInsertedId int, err error) {
	fields, placehoders, sqlVals := getInsertParams(data)
	q := `
	INSERT INTO %s (%s)
	VALUES (%s);`
	query := fmt.Sprintf(q, table.Name, fields, placehoders)
	res, err := rm.db.Exec(query, sqlVals...)
	if err != nil {
		return 0, fmt.Errorf("error on inserting values: %v", err)
	}
	i, err := res.RowsAffected()
	if err != nil {
		log.Printf("error on RowsAffected(): %v", err)
	}
	if i != 1 {
		log.Printf("wrong behaviour, affectet rows: %d\n", i)
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		log.Printf("error on LastInsertId(): %v", err)
	}
	return int(lastInsertId), nil
}

// DeleteById implements RecordManager
func (rm *recordManager) DeleteById(table internal.Table, primaryKey string, id int) (err error) {
	q := `
	DELETE
	FROM %s
	WHERE %s = ?;`
	query := fmt.Sprintf(q, table.Name, primaryKey)

	res, err := rm.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error on deleting values: %v", err)
	}
	i, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error on rowsaffected(): %v", err)
	}

	if i == 0 {
		return fmt.Errorf("no rows affected")
	} else if i > 1 {
		log.Println("affected more then 1 row")
	}
	return nil
}

// GetAllRecords implements RecordManager
func (rm *recordManager) GetAllRecords(table internal.Table, limit int, offset int) (data []map[string]interface{}, err error) {
	fields := getQueryFields(table)
	q := `
	SELECT %s
	FROM %s
	LIMIT ?
	OFFSET ?;`
	query := fmt.Sprintf(q, fields, table.Name)

	rows, err := rm.db.Query(query, limit, offset)
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("unable to get records from sql due to error: %+v", err)
	}

	content := make([]map[string]interface{}, 0)
	dest := initScanDestination(table)
	for rows.Next() {
		rows.Scan(dest...)
		unit := extractSqlVals(table, dest)
		content = append(content, unit)
	}
	return content, nil
}

// GetById implements RecordManager
func (rm *recordManager) GetById(table internal.Table, primaryKey string, id int) (data map[string]interface{}, err error) {
	fields := getQueryFields(table)
	q := `
	SELECT %s
	FROM %s
	WHERE %s = ?;`
	query := fmt.Sprintf(q, fields, table.Name, primaryKey)

	row := rm.db.QueryRow(query, id)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("unable to get records from sql due to error: %+v", err)
	}

	dest := initScanDestination(table)
	if err = row.Scan(dest...); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("record not found")
		}
		return nil, err
	}
	unit := extractSqlVals(table, dest)
	return unit, nil
}

// UpdateById implements RecordManager
func (rm *recordManager) UpdateById(table internal.Table, primaryKey string, id int, data map[string]interface{}) (err error) {
	keyValues, sqlVals := getUpdateParams(data)
	if len(sqlVals) == 0 {
		return fmt.Errorf("required at least one field to update")
	}

	q := `
	UPDATE %s
	SET %s
	WHERE %s = ?;`
	query := fmt.Sprintf(q, table.Name, keyValues, primaryKey)
	sqlVals = append(sqlVals, id)
	result, err := rm.db.Exec(query, sqlVals...)
	if err != nil {
		return fmt.Errorf("error on updating values: %v", err)
	}

	i, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error on rowsaffected(): %v", err)
	}

	if i == 0 {
		return fmt.Errorf("no rows affected")
	} else if i > 1 {
		return fmt.Errorf("affected more then 1 row")
	}
	return nil
}

func newRecordManager(db *sql.DB) *recordManager {
	return &recordManager{
		db: db,
	}
}

func getInsertParams(unit map[string]interface{}) (fieldNames string, placehoderStr string, data []interface{}) {
	length := len(unit)
	names := make([]string, 0, length)
	placehoders := make([]string, 0, length)
	output := make([]interface{}, 0, length)
	for k, v := range unit {
		names = append(names, k)
		placehoders = append(placehoders, "?")
		output = append(output, v)
	}
	return strings.Join(names, ", "), strings.Join(placehoders, ", "), output
}

func getUpdateParams(unit map[string]interface{}) (fieldPlacehoderStr string, data []interface{}) {
	length := len(unit)
	placehoders := make([]string, 0, length)
	output := make([]interface{}, 0, length)
	for k, v := range unit {
		placehoders = append(placehoders, fmt.Sprintf("%s = ?", k))
		output = append(output, v)
	}
	return strings.Join(placehoders, ", "), output
}

// копипаста функции strings.Join только для моей структуры table
func getQueryFields(t internal.Table) string {
	switch len(t.Columns) {
	case 0:
		return ""
	case 1:
		return t.Columns[0].Name
	}
	var sb strings.Builder
	sb.WriteString(t.Columns[0].Name)
	for _, col := range t.Columns[1:] {
		sb.WriteString(", ")
		sb.WriteString(col.Name)
	}
	return sb.String()
}

func initScanDestination(t internal.Table) []interface{} {
	a := make([]interface{}, len(t.Columns))
	for i := range a {
		a[i] = new(interface{})
	}
	return a
}

func extractSqlVals(tableStruct internal.Table, dest []interface{}) map[string]interface{} {
	unit := make(map[string]interface{})
	for i, c := range tableStruct.Columns {
		reflectPointerToInterface := reflect.ValueOf(dest[i])
		reflectInterface := reflect.Indirect(reflectPointerToInterface) //т.к. dest[i] это указатель на interface{} (см. func initScanDestination)
		goInterface := reflectInterface.Interface()
		switch goValue := goInterface.(type) {
		case []byte:
			unit[c.Name] = string(goValue)
		default:
			unit[c.Name] = goValue
		}
	}
	return unit
}

func getPKColumnName(t internal.Table) (string, error) {
	for _, c := range t.Columns {
		if c.IsPrimaryKey {
			return c.Name, nil
		}
	}
	return "", fmt.Errorf("there is no primary key column")
}
