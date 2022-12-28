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
	GetAllRecords(tableName string, limit int, offset int) (data []map[string]interface{}, err error)
	GetById(tableName string, id int) (data map[string]interface{}, err error)
	Create(tableName string, data map[string]interface{}) (lastInsertedId int, err error)
	UpdateById(tableName string, id int, data map[string]interface{}) (err error)
	DeleteById(tableName string, id int) (err error)
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
