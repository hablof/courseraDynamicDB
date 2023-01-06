package service

import (
	"hw6coursera/internal"
	mock_repository "hw6coursera/repository/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	testingSchema internal.Schema = map[string]internal.Table{
		"example_table_1": {
			Name: "example_table_1",
			Columns: []internal.Column{
				{
					Name:         "primary_key",
					ColumnType:   nil,
					Nullable:     false,
					IsPrimaryKey: true,
				},
				{
					Name:         "name",
					ColumnType:   nil,
					Nullable:     false,
					IsPrimaryKey: false,
				},
				{
					Name:         "nullable_field",
					ColumnType:   nil,
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
					ColumnType:   nil,
					Nullable:     false,
					IsPrimaryKey: true,
				},
				{
					Name:         "field",
					ColumnType:   nil,
					Nullable:     false,
					IsPrimaryKey: false,
				},
				{
					Name:         "additional_field",
					ColumnType:   nil,
					Nullable:     true,
					IsPrimaryKey: false,
				},
			},
		}}

	jsonTables string = "[\n    \"example_table_1\",\n    \"example_table_2\"\n]"
)

func TestService_GetAllTables(t *testing.T) {
	testCases := []struct {
		name          string
		schema        internal.Schema
		expectedData  string
		expectedErr   error
		mockBehaviour func(mr *mock_repository.MockRecordManager)
	}{
		{
			name:         "OK",
			schema:       testingSchema,
			expectedData: jsonTables,
			expectedErr:  nil,
			mockBehaviour: func(mr *mock_repository.MockRecordManager) {
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepo := mock_repository.NewMockRecordManager(c)
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
