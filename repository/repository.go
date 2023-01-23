package repository

import (
	"database/sql"
	"hw6coursera/dto"
)

//go:generate mockgen -source=repository.go -destination=mock.go

type Explorer interface {
	GetTableNames() ([]string, error)
	GetColumns(tableName string) ([]dto.Column, error)
}

type RecordManager interface {
	GetAllRecords(table dto.Table, limit int, offset int) (data []map[string]interface{}, err error)
	GetById(table dto.Table, primaryKey string, id int) (data map[string]interface{}, err error)
	Create(table dto.Table, data map[string]interface{}) (lastInsertedId int, err error)
	UpdateById(table dto.Table, primaryKey string, id int, data map[string]interface{}) (err error)
	DeleteById(table dto.Table, primaryKey string, id int) (err error)
}

type Repository struct {
	Explorer
	RecordManager
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Explorer:      newExplorer(db),
		RecordManager: newRecordManager(db),
	}
}
