package dto

const (
	StringType  = "string"
	IntType     = "int"
	FloatType   = "float"
	UnknownType = "unknown"
)

type Schema map[string]Table

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	Name         string
	ColumnType   string // одна из констант
	Nullable     bool
	IsPrimaryKey bool
}
