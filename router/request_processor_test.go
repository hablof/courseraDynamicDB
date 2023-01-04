package router

import (
	"bytes"
	"fmt"
	"hw6coursera/service"
	mock_service "hw6coursera/service/mocks"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRouter_deleteRecord(t *testing.T) {
	//type mock_service

	testCases := []struct {
		name              string
		urlPath           string
		tableName         string
		id                int
		expectedSatusCode int
		expectedBody      string
		mockBehaviour     func(ms *mock_service.MockRecordService, tableName string, id int)
	}{
		{
			name:              "OK",
			urlPath:           "/table/1",
			tableName:         "table",
			id:                1,
			expectedSatusCode: 200,
			expectedBody:      "deleted record id 1",
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, id int) {
				ms.EXPECT().DeleteById(tableName, id).Return(nil)
			},
		},
		{
			name:              "id is not integer",
			urlPath:           "/table/i",
			tableName:         "table",
			id:                0,
			expectedSatusCode: 500,
			expectedBody:      "",
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, id int) {
			},
		},
		{
			name:              "service error",
			urlPath:           "/table/1",
			tableName:         "table",
			id:                1,
			expectedSatusCode: 500,
			expectedBody:      "",
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, id int) {
				ms.EXPECT().DeleteById(tableName, id).Return(fmt.Errorf("table %s not found", tableName))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			recordService := mock_service.NewMockRecordService(c)
			tc.mockBehaviour(recordService, tc.tableName, tc.id)

			servicies := &service.Service{
				RecordService: recordService,
			}
			router := NewRouter(servicies)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", tc.urlPath, bytes.NewBufferString(""))

			router.deleteRecord(w, r)

			assert.Equal(t, tc.expectedSatusCode, w.Result().StatusCode)
			assert.Equal(t, tc.expectedBody, w.Body.String())

		})
	}
}

func TestRouter_getAllTables(t *testing.T) {
	testCases := []struct {
		name              string
		urlPath           string
		expectedSatusCode int
		expectedBody      string
		mockBehaviour     func(ms *mock_service.MockRecordService)
	}{
		{
			name:              "OK",
			urlPath:           "/",
			expectedSatusCode: 200,
			expectedBody:      "[\n\t\"haha\",\n\t\"hoho\"\n]",
			mockBehaviour: func(ms *mock_service.MockRecordService) {
				ms.EXPECT().GetAllTables().Return([]byte("[\n\t\"haha\",\n\t\"hoho\"\n]"), nil)
			},
		},
		{
			name:              "not OK",
			urlPath:           "/",
			expectedSatusCode: 500,
			expectedBody:      "unable to get tables",
			mockBehaviour: func(ms *mock_service.MockRecordService) {
				ms.EXPECT().GetAllTables().Return([]byte("[\n\t\"haha\",\n\t\"hoho\"\n]"), fmt.Errorf("unable to serialize data"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			recordService := mock_service.NewMockRecordService(c)
			tc.mockBehaviour(recordService)

			servicies := &service.Service{
				RecordService: recordService,
			}
			router := NewRouter(servicies)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", tc.urlPath, bytes.NewBufferString(""))

			router.getAllTables(w, r)

			assert.Equal(t, tc.expectedSatusCode, w.Result().StatusCode)
			assert.Equal(t, tc.expectedBody, w.Body.String())

		})
	}
}
