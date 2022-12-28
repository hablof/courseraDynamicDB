package internal

import "reflect"

type Schema map[string]Table

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	Name         string
	ColumnType   reflect.Type
	Nullable     bool
	IsPrimaryKey bool
}
