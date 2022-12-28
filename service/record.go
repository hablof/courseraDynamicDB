package service

import (
	"encoding/json"
	"hw6coursera/dbexplorer"
	"hw6coursera/internal"
	"hw6coursera/repository"
	"log"
)

const (
	DefaultLimit  = 5
	DefaultOffset = 0
)

type RecordManager struct {
	repo   repository.RecordManager
	dbe    dbexplorer.SchemeParser
	Schema internal.Schema
}

// GetAllTables implements RecordService
func (r *RecordManager) GetAllTables() (data []byte, err error) {
	log.Println("getting all tables...")

	tablesList := make([]string, 0, len(r.Schema))
	for n := range r.Schema {
		tablesList = append(tablesList, n)
	}

	b, err := json.MarshalIndent(tablesList, "", "    ")
	return b, nil
}

// Create implements RecordService
func (r *RecordManager) Create(tableName string, data map[string]string) (lastInsertedId int, err error) {
	panic("unimplemented")
}

// DeleteById implements RecordService
func (r *RecordManager) DeleteById(tableName string, id int) (err error) {
	panic("unimplemented")
}

// GetAllRecords implements RecordService
func (r *RecordManager) GetAllRecords(tableName string, limit int, offset int) (data []byte, err error) {
	panic("unimplemented")
}

// GetById implements RecordService
func (r *RecordManager) GetById(tableName string, id int) (data []byte, err error) {
	panic("unimplemented")
}

// UpdateById implements RecordService
func (r *RecordManager) UpdateById(tableName string, id int, data map[string]string) (err error) {
	panic("unimplemented")
}

func (r *RecordManager) InitSchema() error {
	s, err := r.dbe.ParseSchema()
	if err != nil {
		return err
	}
	r.Schema = s
	return nil
}

func newRecordService(repo *repository.Repository, dbe dbexplorer.SchemeParser) RecordService {
	return &RecordManager{
		repo:   repo.RecordManager,
		dbe:    dbe,
		Schema: map[string]internal.Table{},
	}
}
