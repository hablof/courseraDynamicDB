package repository

import (
	"database/sql"
	"hw6coursera/internal"
)

type Explorer interface {
	//Connect(driverName string, dataSourceName string) (sql.DB, error)
	GetTableNames() ([]string, error)
	GetColumns(tableName string) ([]internal.Column, error)
}

type RecordManager interface {
	GetAllRecords(table internal.Table, limit int, offset int) (data []map[string]interface{}, err error)
	GetById(table internal.Table, id int) (data map[string]interface{}, err error)
	Create(table internal.Table, data map[string]interface{}) (lastInsertedId int, err error)
	UpdateById(table internal.Table, primaryKey string, id int, data map[string]interface{}) (err error)
	DeleteById(table internal.Table, primaryKey string, id int) (err error)
}

type Repository struct {
	Explorer
	RecordManager
}

func NewRepository(db sql.DB) *Repository {
	return &Repository{

		Explorer:      newExplorer(db),
		RecordManager: newRecordManager(db),
	}
}
