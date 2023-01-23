package service

import (
	"fmt"
	"hw6coursera/dto"
	"hw6coursera/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	testingSchema dto.Schema = map[string]dto.Table{
		"example_table_1": {
			Name: "example_table_1",
			Columns: []dto.Column{
				{
					Name:         "primary_key",
					ColumnType:   dto.IntType,
					Nullable:     false,
					IsPrimaryKey: true,
				},
				{
					Name:         "name",
					ColumnType:   dto.StringType,
					Nullable:     false,
					IsPrimaryKey: false,
				},
				{
					Name:         "nullable_field",
					ColumnType:   dto.StringType,
					Nullable:     true,
					IsPrimaryKey: false,
				},
			},
		},
		"example_table_2": {
			Name: "example_table_2",
			Columns: []dto.Column{
				{
					Name:         "primary_column",
					ColumnType:   dto.IntType,
					Nullable:     false,
					IsPrimaryKey: true,
				},
				{
					Name:         "field",
					ColumnType:   dto.StringType,
					Nullable:     false,
					IsPrimaryKey: false,
				},
				{
					Name:         "additional_field",
					ColumnType:   dto.StringType,
					Nullable:     true,
					IsPrimaryKey: false,
				},
			},
		}}

	jsonTables string = "[\n    \"example_table_1\",\n    \"example_table_2\"\n]"

	exampleData []map[string]interface{} = []map[string]interface{}{
		{
			"primary_column":   3,
			"field":            "value",
			"additional_field": "additional value",
		},
		{
			"primary_column":   4,
			"field":            "another value",
			"additional_field": "another additional value",
		},
	}

	serializedExampleData       string = "[\n    {\n        \"additional_field\": \"additional value\",\n        \"field\": \"value\",\n        \"primary_column\": 3\n    },\n    {\n        \"additional_field\": \"another additional value\",\n        \"field\": \"another value\",\n        \"primary_column\": 4\n    }\n]"
	serializedExampleSingleData string = "{\n    \"additional_field\": \"additional value\",\n    \"field\": \"value\",\n    \"primary_column\": 3\n}"

	exampleDataWithNull []map[string]interface{} = []map[string]interface{}{
		{
			"primary_column":   3,
			"field":            "value",
			"additional_field": nil,
		},
		{
			"primary_column":   4,
			"field":            "another value",
			"additional_field": nil,
		},
	}

	serializedExampleDataWithNull string = "[\n    {\n        \"field\": \"value\",\n        \"primary_column\": 3\n    },\n    {\n        \"field\": \"another value\",\n        \"primary_column\": 4\n    }\n]"
)

func TestService_GetAllTables(t *testing.T) {
	testCases := []struct {
		name          string
		schema        dto.Schema
		expectedData  string
		expectedErr   error
		mockBehaviour func(mr *repository.MockRecordManager)
	}{
		{
			name:         "OK",
			schema:       testingSchema,
			expectedData: jsonTables,
			expectedErr:  nil,
			mockBehaviour: func(mr *repository.MockRecordManager) {
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepo := repository.NewMockRecordManager(c)
			recordManager := &RecordManager{
				repo:   mockRepo,
				dbe:    nil, ///??????????
				Schema: tc.schema,
			}

			service := Service{
				RecordService: recordManager,
			}

			data, err := service.GetAllTables()

			assert.Equal(t, tc.expectedData, string(data))
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestService_Create(t *testing.T) {
	testCases := []struct {
		name             string
		schema           dto.Schema
		tableName        string
		inputData        map[string]string
		dataToExpect     map[string]interface{}
		expectedInsertId int
		expectedErr      error
		mockBehaviour    func(mr *repository.MockRecordManager, schema dto.Schema, table string, validatedData map[string]interface{})
	}{
		{
			name:             "OK",
			schema:           testingSchema,
			tableName:        "example_table_1",
			inputData:        map[string]string{"name": "name", "nullable_field": "not null"},
			dataToExpect:     map[string]interface{}{"name": "name", "nullable_field": "not null"},
			expectedInsertId: 10,
			expectedErr:      nil,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, table string, validatedData map[string]interface{}) {
				mr.EXPECT().Create(schema[table], validatedData).Return(10, nil)
			},
		},
		{
			name:             "unknown fields",
			schema:           testingSchema,
			tableName:        "example_table_1",
			inputData:        map[string]string{"name": "name", "nullable_field": "not null", "unknown_field": "literal", "unknown_field_2": "literal_2"},
			dataToExpect:     map[string]interface{}{"name": "name", "nullable_field": "not null"},
			expectedInsertId: 20,
			expectedErr:      nil,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, table string, validatedData map[string]interface{}) {
				mr.EXPECT().Create(schema[table], validatedData).Return(20, nil)
			},
		},
		{
			name:             "okay to skip nullable field",
			schema:           testingSchema,
			tableName:        "example_table_1",
			inputData:        map[string]string{"name": "name"},
			dataToExpect:     map[string]interface{}{"name": "name"},
			expectedInsertId: 30,
			expectedErr:      nil,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, table string, validatedData map[string]interface{}) {
				mr.EXPECT().Create(schema[table], validatedData).Return(30, nil)
			},
		},
		{
			name:             "get primary key field no effect",
			schema:           testingSchema,
			tableName:        "example_table_1",
			inputData:        map[string]string{"name": "name", "primary_key": "11"},
			dataToExpect:     map[string]interface{}{"name": "name"},
			expectedInsertId: 40,
			expectedErr:      nil,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, table string, validatedData map[string]interface{}) {
				mr.EXPECT().Create(schema[table], validatedData).Return(40, nil)
			},
		},
		{
			name:             "get primary key field no effect",
			schema:           testingSchema,
			tableName:        "example_table_1",
			inputData:        map[string]string{"name": "name", "primary_key": "11"},
			dataToExpect:     map[string]interface{}{"name": "name"},
			expectedInsertId: 40,
			expectedErr:      nil,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, table string, validatedData map[string]interface{}) {
				mr.EXPECT().Create(schema[table], validatedData).Return(40, nil)
			},
		},
		{
			name:             "not found (table)",
			schema:           testingSchema,
			tableName:        "unknown_table_1",
			inputData:        map[string]string{"name": "name", "nullable_field": "not null"},
			dataToExpect:     map[string]interface{}{},
			expectedInsertId: 0,
			expectedErr:      ErrTableNotFound,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, table string, validatedData map[string]interface{}) {
			},
		},
		{
			name:             "missing non-nullable field",
			schema:           testingSchema,
			tableName:        "example_table_1",
			inputData:        map[string]string{"nullable_field": "not null"},
			dataToExpect:     map[string]interface{}{},
			expectedInsertId: 0,
			expectedErr:      ErrCannotBeNull{"name"},
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, table string, validatedData map[string]interface{}) {
			},
		},
		{
			name:             "repository error",
			schema:           testingSchema,
			tableName:        "example_table_1",
			inputData:        map[string]string{"name": "name", "nullable_field": "not null"},
			dataToExpect:     map[string]interface{}{"name": "name", "nullable_field": "not null"},
			expectedInsertId: 0,
			expectedErr:      fmt.Errorf("repository error"),
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, table string, validatedData map[string]interface{}) {
				mr.EXPECT().Create(schema[table], validatedData).Return(0, fmt.Errorf("repository error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepo := repository.NewMockRecordManager(c)
			recordManager := &RecordManager{
				repo:   mockRepo,
				dbe:    nil, ///??????????
				Schema: tc.schema,
			}

			tc.mockBehaviour(mockRepo, tc.schema, tc.tableName, tc.dataToExpect)
			service := Service{
				RecordService: recordManager,
			}

			insertedId, err := service.Create(tc.tableName, tc.inputData)

			assert.Equal(t, tc.expectedInsertId, insertedId)
			assert.Equal(t, tc.expectedErr, err)
		})
	}

}

func TestService_DeleteById(t *testing.T) {
	testCases := []struct {
		name          string
		schema        dto.Schema
		tableName     string
		primaryKey    string
		idToDelete    int
		expectedErr   error
		mockBehaviour func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int)
	}{
		{
			name:        "OK",
			schema:      testingSchema,
			tableName:   "example_table_1",
			primaryKey:  "primary_key",
			idToDelete:  5,
			expectedErr: nil,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int) {
				mr.EXPECT().DeleteById(schema[tableName], primaryKey, id).Return(nil)
			},
		},
		{
			name:        "not found (table)",
			schema:      testingSchema,
			tableName:   "unknown_table",
			primaryKey:  "",
			idToDelete:  0,
			expectedErr: ErrTableNotFound,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int) {
			},
		},
		{
			name:        "not found (record)",
			schema:      testingSchema,
			tableName:   "example_table_1",
			primaryKey:  "primary_key",
			idToDelete:  5,
			expectedErr: ErrRecordNotFound,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int) {
				mr.EXPECT().DeleteById(schema[tableName], primaryKey, id).Return(repository.ErrRowNotFound)
			},
		},
		{
			name:        "repository error",
			schema:      testingSchema,
			tableName:   "example_table_1",
			primaryKey:  "primary_key",
			idToDelete:  5,
			expectedErr: fmt.Errorf("repository error"),
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int) {
				mr.EXPECT().DeleteById(schema[tableName], primaryKey, id).Return(fmt.Errorf("repository error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepo := repository.NewMockRecordManager(c)
			recordManager := &RecordManager{
				repo:   mockRepo,
				dbe:    nil, ///??????????
				Schema: tc.schema,
			}

			tc.mockBehaviour(mockRepo, tc.schema, tc.tableName, tc.primaryKey, tc.idToDelete)
			service := Service{
				RecordService: recordManager,
			}

			err := service.DeleteById(tc.tableName, tc.idToDelete)

			assert.Equal(t, tc.expectedErr, err)
		})
	}

}

func TestService_GetAllRecords(t *testing.T) {
	testCases := []struct {
		name          string
		schema        dto.Schema
		tableName     string
		limit         int
		offset        int
		dataToReturn  []map[string]interface{}
		errorToReturn error
		expectedErr   error
		expectedData  string
		mockBehaviour func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, limit int, offset int, data []map[string]interface{}, errorToReturn error)
	}{
		{
			name:          "OK",
			schema:        testingSchema,
			tableName:     "example_table_1",
			limit:         2,
			offset:        2,
			dataToReturn:  exampleData,
			errorToReturn: nil,
			expectedErr:   nil,
			expectedData:  serializedExampleData,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, limit int, offset int, data []map[string]interface{}, errorToReturn error) {
				mr.EXPECT().GetAllRecords(schema[tableName], limit, offset).Return(data, errorToReturn)
			},
		},
		{
			name:          "data with null",
			schema:        testingSchema,
			tableName:     "example_table_1",
			limit:         2,
			offset:        2,
			dataToReturn:  exampleDataWithNull,
			errorToReturn: nil,
			expectedErr:   nil,
			expectedData:  serializedExampleDataWithNull,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, limit int, offset int, data []map[string]interface{}, errorToReturn error) {
				mr.EXPECT().GetAllRecords(schema[tableName], limit, offset).Return(data, errorToReturn)
			},
		},
		{
			name:         "not found (table)",
			schema:       testingSchema,
			tableName:    "unknown_table_1",
			limit:        5,
			offset:       0,
			expectedErr:  ErrTableNotFound,
			expectedData: "",
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, limit int, offset int, data []map[string]interface{}, errorToReturn error) {
			},
		},
		{
			name:          "repository error",
			schema:        testingSchema,
			tableName:     "example_table_1",
			limit:         2,
			offset:        2,
			dataToReturn:  nil,
			errorToReturn: fmt.Errorf("repository error"),
			expectedErr:   fmt.Errorf("repository error"),
			expectedData:  "",
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, limit int, offset int, data []map[string]interface{}, errorToReturn error) {
				mr.EXPECT().GetAllRecords(schema[tableName], limit, offset).Return(data, errorToReturn)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepo := repository.NewMockRecordManager(c)
			recordManager := &RecordManager{
				repo:   mockRepo,
				dbe:    nil,
				Schema: tc.schema,
			}

			tc.mockBehaviour(mockRepo, tc.schema, tc.tableName, tc.limit, tc.offset, tc.dataToReturn, tc.errorToReturn)
			service := Service{
				RecordService: recordManager,
			}

			data, err := service.GetAllRecords(tc.tableName, tc.limit, tc.offset)

			assert.Equal(t, tc.expectedData, string(data))
			assert.Equal(t, tc.expectedErr, err)
		})
	}

}

func TestService_GetById(t *testing.T) {
	testCases := []struct {
		name          string
		schema        dto.Schema
		tableName     string
		primaryKey    string
		id            int
		dataToReturn  map[string]interface{}
		errorToReturn error
		expectedErr   error
		expectedData  string
		mockBehaviour func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int, data map[string]interface{}, errorToReturn error)
	}{
		{
			name:          "OK",
			schema:        testingSchema,
			tableName:     "example_table_1",
			primaryKey:    "primary_key",
			id:            3,
			dataToReturn:  exampleData[0],
			errorToReturn: nil,
			expectedErr:   nil,
			expectedData:  serializedExampleSingleData,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int, data map[string]interface{}, errorToReturn error) {
				mr.EXPECT().GetById(schema[tableName], primaryKey, id).Return(data, errorToReturn)
			},
		},
		{
			name:         "not found (table)",
			schema:       testingSchema,
			tableName:    "unknown_table_1",
			id:           3,
			expectedErr:  ErrTableNotFound,
			expectedData: "",
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int, data map[string]interface{}, errorToReturn error) {
			},
		},
		{
			name:          "not found (record)",
			schema:        testingSchema,
			tableName:     "example_table_1",
			primaryKey:    "primary_key",
			id:            100500,
			dataToReturn:  nil,
			errorToReturn: repository.ErrRowNotFound,
			expectedErr:   ErrRecordNotFound,
			expectedData:  "",
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int, data map[string]interface{}, errorToReturn error) {
				mr.EXPECT().GetById(schema[tableName], primaryKey, id).Return(data, errorToReturn)
			},
		},
		{
			name:          "repository error",
			schema:        testingSchema,
			tableName:     "example_table_1",
			primaryKey:    "primary_key",
			id:            100500,
			dataToReturn:  nil,
			errorToReturn: fmt.Errorf("repository error"),
			expectedErr:   fmt.Errorf("repository error"),
			expectedData:  "",
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int, data map[string]interface{}, errorToReturn error) {
				mr.EXPECT().GetById(schema[tableName], primaryKey, id).Return(data, errorToReturn)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepo := repository.NewMockRecordManager(c)
			recordManager := &RecordManager{
				repo:   mockRepo,
				dbe:    nil,
				Schema: tc.schema,
			}

			tc.mockBehaviour(mockRepo, tc.schema, tc.tableName, tc.primaryKey, tc.id, tc.dataToReturn, tc.errorToReturn)
			service := Service{
				RecordService: recordManager,
			}

			data, err := service.GetById(tc.tableName, tc.id)

			assert.Equal(t, tc.expectedData, string(data))
			assert.Equal(t, tc.expectedErr, err)
		})
	}

}

func TestService_UpdateById(t *testing.T) {
	testCases := []struct {
		name          string
		schema        dto.Schema
		tableName     string
		primaryKey    string
		id            int
		inputData     map[string]string
		dataToExpect  map[string]interface{}
		errorToReturn error
		expectedErr   error
		mockBehaviour func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int, data map[string]interface{}, errorToReturn error)
	}{
		{
			name:          "OK",
			schema:        testingSchema,
			tableName:     "example_table_1",
			primaryKey:    "primary_key",
			id:            3,
			inputData:     map[string]string{"name": "updated name"},
			dataToExpect:  map[string]interface{}{"name": "updated name"},
			errorToReturn: nil,
			expectedErr:   nil,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int, data map[string]interface{}, errorToReturn error) {
				mr.EXPECT().UpdateById(schema[tableName], primaryKey, id, data).Return(errorToReturn)
			},
		},
		{
			name:          "missing data to update",
			schema:        testingSchema,
			tableName:     "example_table_1",
			id:            3,
			inputData:     map[string]string{"unknown_field": "value"},
			errorToReturn: nil,
			expectedErr:   fmt.Errorf("missing data to update"),
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int, data map[string]interface{}, errorToReturn error) {
			},
		},
		{
			name:        "not found (table)",
			schema:      testingSchema,
			tableName:   "unknown_table",
			id:          3,
			inputData:   map[string]string{"name": "updated name"},
			expectedErr: ErrTableNotFound,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int, data map[string]interface{}, errorToReturn error) {
			},
		},
		{
			name:          "not found (record)",
			schema:        testingSchema,
			tableName:     "example_table_1",
			primaryKey:    "primary_key",
			id:            100500,
			inputData:     map[string]string{"name": "updated name"},
			dataToExpect:  map[string]interface{}{"name": "updated name"},
			errorToReturn: repository.ErrRowNotFound,
			expectedErr:   ErrRecordNotFound,
			mockBehaviour: func(mr *repository.MockRecordManager, schema dto.Schema, tableName string, primaryKey string, id int, data map[string]interface{}, errorToReturn error) {
				mr.EXPECT().UpdateById(schema[tableName], primaryKey, id, data).Return(errorToReturn)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepo := repository.NewMockRecordManager(c)
			recordManager := &RecordManager{
				repo:   mockRepo,
				dbe:    nil,
				Schema: tc.schema,
			}

			tc.mockBehaviour(mockRepo, tc.schema, tc.tableName, tc.primaryKey, tc.id, tc.dataToExpect, tc.errorToReturn)
			service := Service{
				RecordService: recordManager,
			}

			err := service.UpdateById(tc.tableName, tc.id, tc.inputData)

			assert.Equal(t, tc.expectedErr, err)
		})
	}

}
