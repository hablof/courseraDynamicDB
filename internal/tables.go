package internal

const (
	StringType  = "string"
	IntType     = "int"
	FloatType   = "float"
	BoolType    = "bool"
	UnknownType = "unknown"
)

type Schema map[string]Table

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	Name         string
	ColumnType   string //"VARCHAR", "TEXT", "NVARCHAR", "DECIMAL", "BOOL", "INT", "BIGINT" ....
	Nullable     bool
	IsPrimaryKey bool
}
