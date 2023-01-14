package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"hw6coursera/internal"
	"log"
	"strings"
)

type dbExplorer struct {
	db *sql.DB
}

// GetColumns implements Explorer
func (e *dbExplorer) GetColumns(tableName string) ([]internal.Column, error) {
	rows, err := e.db.Query(fmt.Sprintf("SELECT * FROM `%s` LIMIT 1", tableName))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	sqlColumns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	primaryKey, err := e.getPrimaryKeyFieldName(tableName)
	if err != nil {
		return nil, err
	}

	cols := make([]internal.Column, 0, len(sqlColumns))
	// лучше бы всё это (то, что ниже) делать не в репозитории?
	for _, ct := range sqlColumns {
		c := internal.Column{}
		c.Name = ct.Name()
		if c.Name == primaryKey {
			c.IsPrimaryKey = true
		}
		nullable, ok := ct.Nullable()
		if !ok {
			return nil, errors.New("shitty driver")
		}
		c.Nullable = nullable
		typeName := ct.DatabaseTypeName()
		switch {
		case strings.Contains(typeName, "INT"):
			c.ColumnType = internal.IntType
		case strings.Contains(typeName, "FLOAT") || strings.Contains(typeName, "DECIMAL") || strings.Contains(typeName, "DOUBLE"):
			c.ColumnType = internal.FloatType
		case strings.Contains(typeName, "TEXT") || strings.Contains(typeName, "CHAR"):
			c.ColumnType = internal.StringType
		default:
			c.ColumnType = internal.UnknownType
		}
		cols = append(cols, c)
	}
	return cols, nil
}

// GetTables implements Explorer
func (e *dbExplorer) GetTableNames() ([]string, error) {
	tableRecords, err := e.db.Query(`SHOW TABLES`)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = tableRecords.Close(); err != nil {
			log.Println(err)
		}
	}()

	tableNames := make([]string, 0)

	for tableRecords.Next() {
		n := new(string)
		if err = tableRecords.Scan(n); err != nil {
			return nil, err
		}
		tableNames = append(tableNames, *n)
	}
	return tableNames, nil
}

func newExplorer(db *sql.DB) *dbExplorer {
	return &dbExplorer{
		db: db,
	}
}

func (e *dbExplorer) getPrimaryKeyFieldName(tableName string) (string, error) {

	//	Получаем информацию о столбце, который язляется 'PRIMARY KEY'
	//	Например
	//	+-------+------+------+-----+---------+----------------+
	//	| Field | Type | Null | Key | Default | Extra          |
	//	+-------+------+------+-----+---------+----------------+
	//	| id    | int  | NO   | PRI | NULL    | auto_increment |
	//	+-------+------+------+-----+---------+----------------+
	row := e.db.QueryRow(fmt.Sprintf("SHOW COLUMNS FROM %s WHERE `Key` = 'PRI';", tableName))
	if err := row.Err(); err != nil {
		return "", err
	}
	primaryKeyByteSlice := []byte{}
	waste := new(interface{})
	if err := row.Scan(&primaryKeyByteSlice, waste, waste, waste, waste, waste); err != nil {
		return "", err
	}
	return string(primaryKeyByteSlice), nil
}
