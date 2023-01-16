package repository

import (
	"database/sql"
	"fmt"
	"hw6coursera/internal"
	"log"
	"strconv"
	"strings"
)

var ErrRowNotFound = fmt.Errorf("row not found")

type recordManager struct {
	db *sql.DB
}

// Create implements RecordManager
func (rm *recordManager) Create(table internal.Table, data map[string]interface{}) (lastInsertedId int, err error) {
	fields, placehoders, sqlVals := getInsertParams(data)
	queryTemplate := "INSERT INTO %s (%s) VALUES (%s);"
	queryString := fmt.Sprintf(queryTemplate, table.Name, fields, placehoders)
	res, err := rm.db.Exec(queryString, sqlVals...)
	if err != nil {
		return 0, fmt.Errorf("error on inserting values: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("error on RowsAffected(): %v", err)
	}
	if rowsAffected != 1 {
		log.Printf("wrong behaviour, affectet rows: %d\n", rowsAffected)
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		log.Printf("error on LastInsertId(): %v", err)
	}
	return int(lastInsertId), nil
}

// DeleteById implements RecordManager
func (rm *recordManager) DeleteById(table internal.Table, primaryKey string, id int) (err error) {
	queryTemplate := "DELETE FROM %s WHERE %s = ?;"
	queryString := fmt.Sprintf(queryTemplate, table.Name, primaryKey)
	res, err := rm.db.Exec(queryString, id)
	if err != nil {
		return fmt.Errorf("error on deleting values: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error on rowsaffected(): %v", err)
	}
	if rowsAffected == 0 {
		return ErrRowNotFound
	} else if rowsAffected > 1 {
		log.Println("affected more then 1 row")
	}
	return nil
}

// GetAllRecords implements RecordManager
func (rm *recordManager) GetAllRecords(table internal.Table, limit int, offset int) (data []map[string]interface{}, err error) {
	fields := getQueryFields(table)
	queryTemplate := "SELECT %s FROM %s LIMIT ? OFFSET ?;"
	queryString := fmt.Sprintf(queryTemplate, fields, table.Name)
	rows, err := rm.db.Query(queryString, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("unable to get records due to error: %+v", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	content := make([]map[string]interface{}, 0)
	dest := initScanDestination(table)
	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		unit, err := extractSqlVals(table, dest)
		if err != nil {
			return nil, err
		}
		content = append(content, unit)
	}
	return content, nil
}

// GetById implements RecordManager
func (rm *recordManager) GetById(table internal.Table, primaryKey string, id int) (data map[string]interface{}, err error) {
	fields := getQueryFields(table)
	queryTemplate := "SELECT %s FROM %s WHERE %s = ?;"
	queryString := fmt.Sprintf(queryTemplate, fields, table.Name, primaryKey)
	row := rm.db.QueryRow(queryString, id)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("unable to get records due to error: %+v", err)
	}

	dest := initScanDestination(table)
	if err := row.Scan(dest...); err == sql.ErrNoRows {
		return nil, ErrRowNotFound
	} else if err != nil {
		return nil, err
	}
	unit, err := extractSqlVals(table, dest)
	if err != nil {
		return nil, err
	}
	return unit, nil
}

// UpdateById implements RecordManager
func (rm *recordManager) UpdateById(table internal.Table, primaryKey string, id int, data map[string]interface{}) (err error) {
	palceholders, sqlVals := getUpdateParams(data)
	if len(sqlVals) == 0 {
		return fmt.Errorf("required at least one field to update")
	}

	queryTemplate := "UPDATE %s SET %s WHERE %s = ?;"
	queryString := fmt.Sprintf(queryTemplate, table.Name, palceholders, primaryKey)
	sqlVals = append(sqlVals, id)
	result, err := rm.db.Exec(queryString, sqlVals...)
	if err != nil {
		return fmt.Errorf("error on updating values: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error on rowsaffected(): %v", err)
	}

	if rowsAffected == 0 {
		return ErrRowNotFound
	} else if rowsAffected > 1 {
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

func getUpdateParams(unit map[string]interface{}) (placehoderStr string, data []interface{}) {
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

func extractSqlVals(tableStruct internal.Table, dest []interface{}) (map[string]interface{}, error) {
	unit := make(map[string]interface{})
	for i, c := range tableStruct.Columns {
		ptrToInterface, ok := dest[i].(*interface{}) //т.к. dest[i] это указатель на interface{} (см. func initScanDestination)
		if !ok {
			return nil, fmt.Errorf("interface indirect error")
		}
		scannedValue := *ptrToInterface

		switch value := scannedValue.(type) {
		case []byte: // байты преобразуем в строку...
			str := string(value)
			if c.ColumnType == internal.FloatType { // отдадим в json число, а не строку с числом
				floatValue, err := strconv.ParseFloat(str, 64)
				if err != nil {
					return nil, fmt.Errorf("parse float error: %v", err)
				}
				unit[c.Name] = floatValue
			} else {
				unit[c.Name] = str
			}
		default: /// ...остальное просто отдаём
			unit[c.Name] = value
		}
	}
	return unit, nil
}
