package repository

import (
	"fmt"
	"hw6coursera/internal"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	testingSchema internal.Schema = map[string]internal.Table{
		"example_table_1": {
			Name: "example_table_1",
			Columns: []internal.Column{
				{
					Name:         "primary_key",
					ColumnType:   internal.IntType,
					Nullable:     false,
					IsPrimaryKey: true,
				},
				{
					Name:         "name",
					ColumnType:   internal.StringType,
					Nullable:     false,
					IsPrimaryKey: false,
				},
				{
					Name:         "nullable_field",
					ColumnType:   internal.StringType,
					Nullable:     true,
					IsPrimaryKey: false,
				},
			},
		},
		"example_table_2": {
			Name: "example_table_2",
			Columns: []internal.Column{
				{
					Name:         "primary_column",
					ColumnType:   internal.IntType,
					Nullable:     false,
					IsPrimaryKey: true,
				},
				{
					Name:         "field",
					ColumnType:   internal.StringType,
					Nullable:     false,
					IsPrimaryKey: false,
				},
				{
					Name:         "additional_field",
					ColumnType:   internal.StringType,
					Nullable:     true,
					IsPrimaryKey: false,
				},
			},
		}}
)

func TestRecordManageer_Create(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp)) //не требует полного совпадения запроса
	if err != nil {
		log.Fatalf("unable to mock db: %v", err)
	}
	defer db.Close()

	testCases := []struct {
		name          string
		tableStruct   internal.Table
		data          map[string]interface{}
		expectedQuery string
		mockBehaviour func(query string)
		expectedId    int
		expectedError error
	}{
		{
			name:          "OK",
			tableStruct:   testingSchema["example_table_1"],
			data:          map[string]interface{}{"name": "name value"},
			expectedQuery: "INSERT INTO example_table_1",
			mockBehaviour: func(query string) {
				mock.ExpectExec(query).WithArgs("name value").WillReturnResult(sqlmock.NewResult(5, 1))
			},
			expectedId:    5,
			expectedError: nil,
		},
		{
			name:          "db error",
			tableStruct:   testingSchema["example_table_1"],
			data:          map[string]interface{}{"name": "name value"},
			expectedQuery: "INSERT INTO example_table_1",
			mockBehaviour: func(query string) {
				mock.ExpectExec(query).WithArgs("name value").WillReturnError(fmt.Errorf("db error"))
			},
			expectedId:    0,
			expectedError: fmt.Errorf("error on inserting values: %v", fmt.Errorf("db error")),
		},
	}

	for _, tc := range testCases {
		rm := newRecordManager(db)
		tc.mockBehaviour(tc.expectedQuery)

		lastInsertedId, err := rm.Create(tc.tableStruct, tc.data)

		assert.Equal(t, tc.expectedId, lastInsertedId)
		assert.Equal(t, tc.expectedError, err)
	}
}

func TestRecordManageer_DeleteById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp)) //не требует полного совпадения запроса
	if err != nil {
		log.Fatalf("unable to mock db: %v", err)
	}
	defer db.Close()

	testCases := []struct {
		name          string
		tableStruct   internal.Table
		primaryKey    string
		expectedQuery string
		id            int
		mockBehaviour func(query string, id int)
		expectedError error
	}{
		{
			name:          "OK",
			tableStruct:   testingSchema["example_table_1"],
			primaryKey:    "primary_key",
			expectedQuery: "DELETE FROM example_table_1 WHERE primary_key",
			id:            6,
			mockBehaviour: func(query string, id int) {
				mock.ExpectExec(query).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			name:          "row not found",
			tableStruct:   testingSchema["example_table_1"],
			primaryKey:    "primary_key",
			expectedQuery: "DELETE FROM example_table_1 WHERE primary_key",
			id:            6,
			mockBehaviour: func(query string, id int) {
				mock.ExpectExec(query).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: ErrRowNotFound,
		},
		{
			name:          "db error",
			tableStruct:   testingSchema["example_table_1"],
			primaryKey:    "primary_key",
			expectedQuery: "DELETE FROM example_table_1 WHERE primary_key",
			id:            6,
			mockBehaviour: func(query string, id int) {
				mock.ExpectExec(query).WithArgs(id).WillReturnError(fmt.Errorf("db error"))
			},
			expectedError: fmt.Errorf("error on deleting values: %v", fmt.Errorf("db error")),
		},
	}

	for _, tc := range testCases {
		rm := newRecordManager(db)
		tc.mockBehaviour(tc.expectedQuery, tc.id)

		err := rm.DeleteById(tc.tableStruct, tc.primaryKey, tc.id)

		assert.Equal(t, tc.expectedError, err)
	}
}

func TestRecordManageer_GetAllRecords(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp)) //не требует полного совпадения запроса
	if err != nil {
		log.Fatalf("unable to mock db: %v", err)
	}
	defer db.Close()

	testCases := []struct {
		name          string
		tableStruct   internal.Table
		limit         int
		offset        int
		expectedQuery string
		mockBehaviour func(query string, limit int, offset int)
		expectedData  []map[string]interface{}
		expectedError error
	}{
		{
			name:          "OK",
			tableStruct:   testingSchema["example_table_1"],
			limit:         2,
			offset:        2,
			expectedQuery: "SELECT",
			mockBehaviour: func(query string, limit int, offset int) {
				rows := sqlmock.NewRows([]string{"primary_key", "name", "nullable_field"}).AddRow(3, "name 3", nil).AddRow(4, "name 4", "not null")
				mock.ExpectQuery(query).WithArgs(limit, offset).WillReturnRows(rows)
			},
			expectedData: []map[string]interface{}{
				{
					"primary_key":    int64(3),
					"name":           "name 3",
					"nullable_field": nil,
				},
				{
					"primary_key":    int64(4),
					"name":           "name 4",
					"nullable_field": "not null",
				},
			},
			expectedError: nil,
		},
		{
			name:          "db error",
			tableStruct:   testingSchema["example_table_1"],
			limit:         2,
			offset:        2,
			expectedQuery: "SELECT",
			mockBehaviour: func(query string, limit int, offset int) {
				mock.ExpectQuery(query).WithArgs(limit, offset).WillReturnError(fmt.Errorf("db error"))
			},
			expectedError: fmt.Errorf("unable to get records due to error: %+v", fmt.Errorf("db error")),
		},
	}

	for _, tc := range testCases {
		rm := newRecordManager(db)
		tc.mockBehaviour(tc.expectedQuery, tc.limit, tc.offset)

		data, err := rm.GetAllRecords(tc.tableStruct, tc.limit, tc.offset)

		assert.Equal(t, tc.expectedData, data)
		assert.Equal(t, tc.expectedError, err)
	}
}

func TestRecordManageer_GetById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp)) //не требует полного совпадения запроса
	if err != nil {
		log.Fatalf("unable to mock db: %+v", err)
	}
	defer db.Close()

	testCases := []struct {
		name          string
		tableStruct   internal.Table
		primaryKey    string
		expectedQuery string
		id            int
		mockBehaviour func(query string, id int)
		expectedData  map[string]interface{}
		expectedError error
	}{
		{
			name:          "OK",
			tableStruct:   testingSchema["example_table_1"],
			primaryKey:    "primary_key",
			expectedQuery: "SELECT",
			id:            3,
			mockBehaviour: func(query string, id int) {
				rows := sqlmock.NewRows([]string{"primary_key", "name", "nullable_field"}).AddRow(3, "name 3", nil)
				mock.ExpectQuery(query).WithArgs(id).WillReturnRows(rows)
			},
			expectedData: map[string]interface{}{
				"primary_key":    int64(3),
				"name":           "name 3",
				"nullable_field": nil,
			},
			expectedError: nil,
		},
		{
			name:          "row not found",
			tableStruct:   testingSchema["example_table_1"],
			primaryKey:    "primary_key",
			expectedQuery: "SELECT",
			id:            6,
			mockBehaviour: func(query string, id int) {
				rows := sqlmock.NewRows([]string{"primary_key", "name", "nullable_field"})
				mock.ExpectQuery(query).WithArgs(id).WillReturnRows(rows)
			},
			expectedError: ErrRowNotFound,
		},
		{
			name:          "db error",
			tableStruct:   testingSchema["example_table_1"],
			primaryKey:    "primary_key",
			expectedQuery: "SELECT",
			id:            6,
			mockBehaviour: func(query string, id int) {
				mock.ExpectQuery(query).WithArgs(id).WillReturnError(fmt.Errorf("db error"))
			},
			expectedError: fmt.Errorf("unable to get records due to error: %+v", fmt.Errorf("db error")),
		},
	}

	for _, tc := range testCases {
		rm := newRecordManager(db)
		tc.mockBehaviour(tc.expectedQuery, tc.id)

		data, err := rm.GetById(tc.tableStruct, tc.primaryKey, tc.id)

		assert.Equal(t, tc.expectedData, data)
		assert.Equal(t, tc.expectedError, err)
	}
}

func TestRecordManageer_UpdateById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp)) //не требует полного совпадения запроса
	if err != nil {
		log.Fatalf("unable to mock db: %v", err)
	}
	defer db.Close()

	testCases := []struct {
		name          string
		tableStruct   internal.Table
		primaryKey    string
		id            int
		data          map[string]interface{}
		expectedQuery string
		mockBehaviour func(query string)
		expectedError error
	}{
		{
			name:          "OK",
			tableStruct:   testingSchema["example_table_1"],
			primaryKey:    "primary_key",
			id:            3,
			data:          map[string]interface{}{"name": "new name"},
			expectedQuery: "UPDATE",
			mockBehaviour: func(query string) {
				mock.ExpectExec(query).WithArgs("new name", 3).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			name:          "row not found",
			tableStruct:   testingSchema["example_table_1"],
			primaryKey:    "primary_key",
			id:            100500,
			data:          map[string]interface{}{"name": "new name"},
			expectedQuery: "UPDATE",
			mockBehaviour: func(query string) {
				mock.ExpectExec(query).WithArgs("new name", 100500).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: ErrRowNotFound,
		},
		{
			name:          "db error",
			tableStruct:   testingSchema["example_table_1"],
			primaryKey:    "primary_key",
			id:            3,
			data:          map[string]interface{}{"name": "new name"},
			expectedQuery: "UPDATE",
			mockBehaviour: func(query string) {
				mock.ExpectExec(query).WithArgs("new name", 3).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: ErrRowNotFound,
		},
	}

	for _, tc := range testCases {
		rm := newRecordManager(db)
		tc.mockBehaviour(tc.expectedQuery)

		err := rm.UpdateById(tc.tableStruct, tc.primaryKey, tc.id, tc.data)

		assert.Equal(t, tc.expectedError, err)
	}
}
