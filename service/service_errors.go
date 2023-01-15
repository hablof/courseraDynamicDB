package service

import "fmt"

var (
	ErrTableNotFound  = fmt.Errorf("table not found")
	ErrRecordNotFound = fmt.Errorf("record not found")
	ErrMissingUpdData = fmt.Errorf("missing data to update")
)

type ErrType struct {
	field string
}

func (te ErrType) Error() string {
	return fmt.Sprintf("invalid type %s", te.field)
}

type ErrCannotBeNull struct {
	field string
}

func (ne ErrCannotBeNull) Error() string {
	return fmt.Sprintf("%s cannot be null", ne.field)
}
