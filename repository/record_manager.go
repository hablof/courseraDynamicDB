package repository

import (
	"database/sql"
)

type recordManager struct {
	db sql.DB
}

// Create implements RecordManager
func (rm *recordManager) Create(tableName string, data map[string]interface{}) (lastInsertedId int, err error) {
	panic("unimplemented")
}

// DeleteById implements RecordManager
func (rm *recordManager) DeleteById(tableName string, id int) (err error) {
	panic("unimplemented")
}

// GetAllRecords implements RecordManager
func (rm *recordManager) GetAllRecords(tableName string, limit int, offset int) (data []map[string]interface{}, err error) {
	panic("unimplemented")
}

// GetById implements RecordManager
func (rm *recordManager) GetById(tableName string, id int) (data map[string]interface{}, err error) {
	panic("unimplemented")
}

// UpdateById implements RecordManager
func (rm *recordManager) UpdateById(tableName string, id int, data map[string]interface{}) (err error) {
	panic("unimplemented")
}

func newRecordManager(db sql.DB) *recordManager {
	return &recordManager{
		db: db,
	}
}
