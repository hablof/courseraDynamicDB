package service

import (
	"encoding/json"
	"fmt"
	"hw6coursera/dbexplorer"
	"hw6coursera/internal"
	"hw6coursera/repository"
	"log"
)

const (
	DefaultLimit  = 5
	DefaultOffset = 0
)

var (
	ErrTableNotFound  = fmt.Errorf("table not found")
	ErrRecordNotFound = fmt.Errorf("record not found")
)

type RecordManager struct {
	repo   repository.RecordManager
	dbe    dbexplorer.SchemeParser
	Schema internal.Schema
}

// GetAllTables implements RecordService
func (r *RecordManager) GetAllTables() ([]byte, error) {
	log.Println("getting all tables...")

	tablesList := make([]string, 0, len(r.Schema))
	for n := range r.Schema {
		tablesList = append(tablesList, n)
	}

	b, err := json.MarshalIndent(tablesList, "", "    ")
	if err != nil {
		log.Printf("unable to serialize data: %+v", err)
		return nil, err
	}
	return b, nil
}

// Create implements RecordService
func (r *RecordManager) Create(tableName string, data map[string]string) (int, error) {
	log.Printf("inserting record to table %s\n", tableName)

	unit, err := validateData(data)
	if err != nil {
		log.Printf("invalid data")
		return 0, err
	}

	tableStruct, ok := r.Schema[tableName]
	if !ok {
		log.Printf("table %s not found", tableName)
		return 0, ErrTableNotFound
	}

	insertedId, err := r.repo.Create(tableStruct, unit)
	if err != nil {
		log.Printf("unable to create record: %+v", err)
		return 0, err
	}
	return insertedId, nil
}

// DeleteById implements RecordService
func (r *RecordManager) DeleteById(tableName string, id int) error {
	log.Printf("deleting record from table %s\n", tableName)

	tableStruct, ok := r.Schema[tableName]
	if !ok {
		log.Printf("table %s not found", tableName)
		return ErrTableNotFound
	}

	primaryKey, err := getPKColumnName(tableStruct)
	if err != nil {
		log.Printf("unable to get primary key name: %+v", err)
		return err
	}

	if err := r.repo.DeleteById(tableStruct, primaryKey, id); err == repository.ErrRowNotFound {
		return ErrRecordNotFound
	} else if err != nil {
		log.Printf("unable to delete record: %+v", err)
		return err
	}
	return nil
}

// GetAllRecords implements RecordService
func (r *RecordManager) GetAllRecords(tableName string, limit int, offset int) ([]byte, error) {
	log.Printf("getting records from table %s", tableName)

	tableStruct, ok := r.Schema[tableName]
	if !ok {
		log.Printf("table %s not found", tableName)
		return nil, ErrTableNotFound
	}

	data, err := r.repo.GetAllRecords(tableStruct, limit, offset)
	if err != nil {
		log.Printf("unable to get all records: %+v", err)
		return nil, err
	}

	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Printf("unable to serialize data: %+v", err)
		return nil, err
	}
	return b, nil
}

// GetById implements RecordService
func (r *RecordManager) GetById(tableName string, id int) ([]byte, error) {
	tableStruct, ok := r.Schema[tableName]
	if !ok {
		log.Printf("table %s not found", tableName)
		return nil, ErrTableNotFound
	}

	log.Printf("getting record (id=%d) from table %s", id, tableName)

	primaryKey, err := getPKColumnName(tableStruct)
	if err != nil {
		log.Printf("unable to get primary key name: %+v", err)
		return nil, err
	}

	data, err := r.repo.GetById(tableStruct, primaryKey, id)
	if err == repository.ErrRowNotFound {
		return nil, ErrRecordNotFound
	} else if err != nil {
		log.Printf("unable to get record dy id: %+v", err)
		return nil, err
	}
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Printf("unable to serialize data: %+v", err)
		return nil, err
	}
	return b, nil
}

// UpdateById implements RecordService
func (r *RecordManager) UpdateById(tableName string, id int, data map[string]string) error {
	log.Printf("updating record (id=%d) from table %s", id, tableName)

	if len(data) == 0 {
		return fmt.Errorf("missing data to update")
	}
	unit, err := validateData(data)
	if err != nil {
		log.Printf("invalid data")
		return err
	}

	tableStruct, ok := r.Schema[tableName]
	if !ok {
		log.Printf("table %s not found", tableName)
		return ErrTableNotFound
	}

	primaryKey, err := getPKColumnName(tableStruct)
	if err != nil {
		log.Printf("unable to get primary key name: %+v", err)
		return err
	}

	if err := r.repo.UpdateById(tableStruct, primaryKey, id, unit); err == repository.ErrRowNotFound {
		return ErrRecordNotFound
	} else if err != nil {
		log.Printf("unable to update record by id: %+v", err)
		return err
	}

	return nil
}

func (r *RecordManager) InitSchema() error {
	s, err := r.dbe.ParseSchema()
	if err != nil {
		return err
	}
	r.Schema = s
	return nil
}

func newRecordService(repo *repository.Repository, dbe dbexplorer.SchemeParser) *RecordManager {
	return &RecordManager{
		repo:   repo.RecordManager,
		dbe:    dbe,
		Schema: map[string]internal.Table{},
	}
}

// reflect Type.ConvertibleTo(u Type) bool ??
func validateData(data map[string]string) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	for k, v := range data {
		output[k] = v
	}
	return output, nil
}

func getPKColumnName(t internal.Table) (string, error) {
	for _, c := range t.Columns {
		if c.IsPrimaryKey {
			return c.Name, nil
		}
	}
	return "", fmt.Errorf("there is no primary key column")
}
