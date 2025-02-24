// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Users\deepak.ag\Desktop\MiniProject1\mymodule\utils\log_reader_interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "Mymodule/mymodule/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLogReader is a mock of LogReader interface.
type MockLogReader struct {
	ctrl     *gomock.Controller
	recorder *MockLogReaderMockRecorder
}

// MockLogReaderMockRecorder is the mock recorder for MockLogReader.
type MockLogReaderMockRecorder struct {
	mock *MockLogReader
}

// NewMockLogReader creates a new mock instance.
func NewMockLogReader(ctrl *gomock.Controller) *MockLogReader {
	mock := &MockLogReader{ctrl: ctrl}
	mock.recorder = &MockLogReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogReader) EXPECT() *MockLogReaderMockRecorder {
	return m.recorder
}

// ReadLogsFromFile mocks base method.
func (m *MockLogReader) ReadLogsFromFile() ([]models.ApiLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadLogsFromFile")
	ret0, _ := ret[0].([]models.ApiLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadLogsFromFile indicates an expected call of ReadLogsFromFile.
func (mr *MockLogReaderMockRecorder) ReadLogsFromFile() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadLogsFromFile", reflect.TypeOf((*MockLogReader)(nil).ReadLogsFromFile))
}
