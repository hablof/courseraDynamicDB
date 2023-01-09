package router

import (
	"bytes"
	"fmt"
	"hw6coursera/service"
	mock_service "hw6coursera/service/mocks"
	"net/http/httptest"
	"net/url"
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
				ms.EXPECT().DeleteById(tableName, id).Return(fmt.Errorf("some service error"))
			},
		},
		{
			name:              "not found (table)",
			urlPath:           "/table/1",
			tableName:         "table",
			id:                1,
			expectedSatusCode: 404,
			expectedBody:      "",
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, id int) {
				ms.EXPECT().DeleteById(tableName, id).Return(service.ErrTableNotFound)
			},
		},
		{
			name:              "not found (record)",
			urlPath:           "/table/1",
			tableName:         "table",
			id:                1,
			expectedSatusCode: 404,
			expectedBody:      "",
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, id int) {
				ms.EXPECT().DeleteById(tableName, id).Return(service.ErrRecordNotFound)
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
			name:              "service error",
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

func TestRouter_getRecords(t *testing.T) {
	testCases := []struct {
		name              string
		urlPath           string
		expectedSatusCode int
		expectedBody      string
		tableName         string
		limit             int
		offset            int
		mockBehaviour     func(ms *mock_service.MockRecordService, tableName string, limit int, offset int)
	}{
		{
			name:              "OK",
			urlPath:           "/table",
			expectedSatusCode: 200,
			expectedBody:      bigJSON,
			tableName:         "table",
			limit:             5,
			offset:            0,
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, limit int, offset int) {
				ms.EXPECT().GetAllRecords(tableName, limit, offset).Return([]byte(bigJSON), nil)
			},
		},
		{
			name:              "limit & offset check",
			urlPath:           "/table?limit=1&offset=2",
			expectedSatusCode: 200,
			expectedBody:      smallJSON,
			tableName:         "table",
			limit:             1,
			offset:            2,
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, limit int, offset int) {
				ms.EXPECT().GetAllRecords(tableName, limit, offset).Return([]byte(smallJSON), nil)
			},
		},
		{
			name:              "service error",
			urlPath:           "/table",
			expectedSatusCode: 500,
			expectedBody:      "unable to get records",
			tableName:         "table",
			limit:             5,
			offset:            0,
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, limit int, offset int) {
				ms.EXPECT().GetAllRecords(tableName, limit, offset).Return(nil, fmt.Errorf("some service error"))
			},
		},
		{
			name:              "not found (table)",
			urlPath:           "/table",
			expectedSatusCode: 404,
			expectedBody:      "",
			tableName:         "table",
			limit:             5,
			offset:            0,
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, limit int, offset int) {
				ms.EXPECT().GetAllRecords(tableName, limit, offset).Return(nil, service.ErrTableNotFound)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			recordService := mock_service.NewMockRecordService(c)
			tc.mockBehaviour(recordService, tc.tableName, tc.limit, tc.offset)

			servicies := &service.Service{
				RecordService: recordService,
			}

			router := NewRouter(servicies)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", tc.urlPath, bytes.NewBufferString(""))

			router.getRecords(w, r)

			assert.Equal(t, tc.expectedSatusCode, w.Result().StatusCode)
			assert.Equal(t, tc.expectedBody, w.Body.String())

		})
	}
}

func TestRouter_getSingleRecord(t *testing.T) {
	testCases := []struct {
		name              string
		urlPath           string
		expectedSatusCode int
		expectedBody      string
		tableName         string
		id                int
		mockBehaviour     func(ms *mock_service.MockRecordService, tableName string, id int)
	}{
		{
			name:              "OK",
			urlPath:           "/table/3",
			expectedSatusCode: 200,
			expectedBody:      smallJSON,
			tableName:         "table",
			id:                3,
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, id int) {
				ms.EXPECT().GetById(tableName, id).Return([]byte(smallJSON), nil)
			},
		},
		{
			name:              "service error",
			urlPath:           "/table/3",
			expectedSatusCode: 500,
			expectedBody:      "unable to service",
			tableName:         "table",
			id:                3,
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, id int) {
				ms.EXPECT().GetById(tableName, id).Return(nil, fmt.Errorf("some service error"))
			},
		},
		{
			name:              "not found (table)",
			urlPath:           "/table/3",
			expectedSatusCode: 404,
			expectedBody:      "",
			tableName:         "table",
			id:                3,
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, id int) {
				ms.EXPECT().GetById(tableName, id).Return(nil, service.ErrTableNotFound)
			},
		},
		{
			name:              "not found (record)",
			urlPath:           "/table/3",
			expectedSatusCode: 404,
			expectedBody:      "",
			tableName:         "table",
			id:                3,
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, id int) {
				ms.EXPECT().GetById(tableName, id).Return(nil, service.ErrRecordNotFound)
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
			r := httptest.NewRequest("GET", tc.urlPath, bytes.NewBufferString(""))

			router.getSingleRecord(w, r)

			assert.Equal(t, tc.expectedSatusCode, w.Result().StatusCode)
			assert.Equal(t, tc.expectedBody, w.Body.String())

		})
	}
}

func TestRouter_updateRecord(t *testing.T) {
	testCases := []struct {
		name               string
		urlPath            string
		expectedSatusCode  int
		expectedBody       string
		tableName          string
		id                 int
		requestData        map[string]string
		updateDataToExpect map[string]interface{}
		mockBehaviour      func(ms *mock_service.MockRecordService, tableName string, id int, data map[string]interface{})
	}{
		{
			name:               "OK",
			urlPath:            "/table/3",
			expectedSatusCode:  200,
			expectedBody:       "updated record id 3",
			tableName:          "table",
			id:                 3,
			requestData:        map[string]string{"updating field": "new data"},
			updateDataToExpect: map[string]interface{}{"updating field": "new data"},
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, id int, data map[string]interface{}) {
				ms.EXPECT().UpdateById(tableName, id, data).Return(nil)
			},
		},
		{
			name:               "missing form body",
			urlPath:            "/table/3",
			expectedSatusCode:  500,
			expectedBody:       "unable to update record",
			tableName:          "table",
			id:                 3,
			requestData:        map[string]string{},
			updateDataToExpect: map[string]interface{}{},
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, id int, data map[string]interface{}) {
				ms.EXPECT().UpdateById(tableName, id, data).Return(fmt.Errorf("missing data to update"))
			},
		},
		{
			name:               "bad id",
			urlPath:            "/table/bad_id",
			expectedSatusCode:  500,
			expectedBody:       "",
			tableName:          "",
			id:                 0,
			requestData:        map[string]string{},
			updateDataToExpect: map[string]interface{}{},
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, id int, data map[string]interface{}) {
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			recordService := mock_service.NewMockRecordService(c)
			tc.mockBehaviour(recordService, tc.tableName, tc.id, tc.updateDataToExpect)

			servicies := &service.Service{
				RecordService: recordService,
			}

			router := NewRouter(servicies)
			w := httptest.NewRecorder()
			params := url.Values{}
			for k, v := range tc.requestData {
				params.Add(k, v)
			}

			r := httptest.NewRequest("POST", tc.urlPath, bytes.NewBufferString(params.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			router.updateRecord(w, r)

			assert.Equal(t, tc.expectedSatusCode, w.Result().StatusCode)
			assert.Equal(t, tc.expectedBody, w.Body.String())

		})
	}
}

func TestRouter_insertRecord(t *testing.T) {
	testCases := []struct {
		name              string
		urlPath           string
		expectedSatusCode int
		expectedBody      string
		tableName         string
		lastInsertId      int
		requestData       map[string]string
		dataToExpect      map[string]interface{}
		mockBehaviour     func(ms *mock_service.MockRecordService, tableName string, data map[string]interface{})
	}{
		{
			name:              "OK",
			urlPath:           "/table",
			expectedSatusCode: 200,
			expectedBody:      "last insert id 3",
			tableName:         "table",
			lastInsertId:      3,
			requestData:       map[string]string{"updating field": "new data"},
			dataToExpect:      map[string]interface{}{"updating field": "new data"},
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, data map[string]interface{}) {
				ms.EXPECT().Create(tableName, data).Return(3, nil)
			},
		},
		{
			name:              "not found (table)",
			urlPath:           "/table",
			expectedSatusCode: 404,
			expectedBody:      "",
			tableName:         "table",
			lastInsertId:      0,
			requestData:       map[string]string{},
			dataToExpect:      map[string]interface{}{},
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, data map[string]interface{}) {
				ms.EXPECT().Create(tableName, data).Return(0, service.ErrTableNotFound)
			},
		},
		{
			name:              "service error",
			urlPath:           "/table",
			expectedSatusCode: 500,
			expectedBody:      "unable to insert record",
			tableName:         "table",
			lastInsertId:      0,
			requestData:       map[string]string{},
			dataToExpect:      map[string]interface{}{},
			mockBehaviour: func(ms *mock_service.MockRecordService, tableName string, data map[string]interface{}) {
				ms.EXPECT().Create(tableName, data).Return(0, fmt.Errorf("some service error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			recordService := mock_service.NewMockRecordService(c)
			tc.mockBehaviour(recordService, tc.tableName, tc.dataToExpect)

			servicies := &service.Service{
				RecordService: recordService,
			}

			router := NewRouter(servicies)
			w := httptest.NewRecorder()
			params := url.Values{}
			for k, v := range tc.requestData {
				params.Add(k, v)
			}

			r := httptest.NewRequest("PUT", tc.urlPath, bytes.NewBufferString(params.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			router.insertRecord(w, r)

			assert.Equal(t, tc.expectedSatusCode, w.Result().StatusCode)
			assert.Equal(t, tc.expectedBody, w.Body.String())

		})
	}

}

func Test_getIntFieldOrDefault(t *testing.T) {
	testCases := []struct {
		name         string
		params       map[string]string
		defaultValue int
		targetField  string
		exeptedValue int
	}{
		{
			name:         "got 5",
			params:       map[string]string{"field1": "notIntValue", "field2": "5", "field3": "17"},
			defaultValue: 1,
			targetField:  "field2",
			exeptedValue: 5,
		},
		{
			name:         "got default due field missing",
			params:       map[string]string{"field1": "notIntValue", "field2": "5", "field3": "17"},
			defaultValue: 1,
			targetField:  "target field",
			exeptedValue: 1,
		},
		{
			name:         "got default due error",
			params:       map[string]string{"field1": "notIntValue", "field2": "5", "field3": "17"},
			defaultValue: 1,
			targetField:  "field1",
			exeptedValue: 1,
		},
		{
			name:         "got negative 5",
			params:       map[string]string{"field1": "notIntValue", "field2": "5", "field3": "-5"},
			defaultValue: 1,
			targetField:  "field3",
			exeptedValue: -5,
		},
		{
			name:         "int overflow",
			params:       map[string]string{"field1": "12345678901234567890123456789012345678901234567890", "field2": "5", "field3": "-5"},
			defaultValue: 1,
			targetField:  "field1",
			exeptedValue: 1,
		},
		{
			name:         "dont parse float",
			params:       map[string]string{"field1": "5.15", "field2": "5", "field3": "-5"},
			defaultValue: 10,
			targetField:  "field1",
			exeptedValue: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqParams := url.Values{}
			for k, v := range tc.params {
				reqParams.Add(k, v)
			}
			r := httptest.NewRequest("GET", "/", nil)
			r.URL.RawQuery = reqParams.Encode()

			result := getIntFieldOrDefault(r, tc.targetField, tc.defaultValue)

			assert.Equal(t, tc.exeptedValue, result)
		})
	}
}
