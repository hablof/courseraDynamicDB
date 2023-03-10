// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRecordService is a mock of RecordService interface.
type MockRecordService struct {
	ctrl     *gomock.Controller
	recorder *MockRecordServiceMockRecorder
}

// MockRecordServiceMockRecorder is the mock recorder for MockRecordService.
type MockRecordServiceMockRecorder struct {
	mock *MockRecordService
}

// NewMockRecordService creates a new mock instance.
func NewMockRecordService(ctrl *gomock.Controller) *MockRecordService {
	mock := &MockRecordService{ctrl: ctrl}
	mock.recorder = &MockRecordServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRecordService) EXPECT() *MockRecordServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRecordService) Create(tableName string, data map[string]string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", tableName, data)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRecordServiceMockRecorder) Create(tableName, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRecordService)(nil).Create), tableName, data)
}

// DeleteById mocks base method.
func (m *MockRecordService) DeleteById(tableName string, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", tableName, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockRecordServiceMockRecorder) DeleteById(tableName, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockRecordService)(nil).DeleteById), tableName, id)
}

// GetAllRecords mocks base method.
func (m *MockRecordService) GetAllRecords(tableName string, limit, offset int) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllRecords", tableName, limit, offset)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllRecords indicates an expected call of GetAllRecords.
func (mr *MockRecordServiceMockRecorder) GetAllRecords(tableName, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllRecords", reflect.TypeOf((*MockRecordService)(nil).GetAllRecords), tableName, limit, offset)
}

// GetAllTables mocks base method.
func (m *MockRecordService) GetAllTables() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTables")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTables indicates an expected call of GetAllTables.
func (mr *MockRecordServiceMockRecorder) GetAllTables() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTables", reflect.TypeOf((*MockRecordService)(nil).GetAllTables))
}

// GetById mocks base method.
func (m *MockRecordService) GetById(tableName string, id int) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", tableName, id)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockRecordServiceMockRecorder) GetById(tableName, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockRecordService)(nil).GetById), tableName, id)
}

// InitSchema mocks base method.
func (m *MockRecordService) InitSchema() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitSchema")
	ret0, _ := ret[0].(error)
	return ret0
}

// InitSchema indicates an expected call of InitSchema.
func (mr *MockRecordServiceMockRecorder) InitSchema() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitSchema", reflect.TypeOf((*MockRecordService)(nil).InitSchema))
}

// UpdateById mocks base method.
func (m *MockRecordService) UpdateById(tableName string, id int, data map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateById", tableName, id, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateById indicates an expected call of UpdateById.
func (mr *MockRecordServiceMockRecorder) UpdateById(tableName, id, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateById", reflect.TypeOf((*MockRecordService)(nil).UpdateById), tableName, id, data)
}
