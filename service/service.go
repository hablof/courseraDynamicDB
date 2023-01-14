package service

import (
	"hw6coursera/dbexplorer"
	"hw6coursera/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type RecordService interface {
	GetAllTables() (data []byte, err error)
	GetAllRecords(tableName string, limit int, offset int) (data []byte, err error)
	GetById(tableName string, id int) (data []byte, err error)
	Create(tableName string, data map[string]string) (lastInsertedId int, err error)
	UpdateById(tableName string, id int, data map[string]string) (err error)
	DeleteById(tableName string, id int) (err error)
	InitSchema() error
}

type Service struct {
	RecordService
}

func NewService(r *repository.Repository, dbe *dbexplorer.DBexplorer) *Service {
	return &Service{
		RecordService: newRecordService(r, dbe),
	}
}
