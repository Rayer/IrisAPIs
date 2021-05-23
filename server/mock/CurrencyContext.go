// Code generated by MockGen. DO NOT EDIT.
// Source: ../CurrencyContext.go

// Package mock_IrisAPIs is a generated GoMock package.
package mock_IrisAPIs

import (
	IrisAPIs "IrisAPIs"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCurrencyService is a mock of CurrencyService interface
type MockCurrencyService struct {
	ctrl     *gomock.Controller
	recorder *MockCurrencyServiceMockRecorder
}

// MockCurrencyServiceMockRecorder is the mock recorder for MockCurrencyService
type MockCurrencyServiceMockRecorder struct {
	mock *MockCurrencyService
}

// NewMockCurrencyService creates a new mock instance
func NewMockCurrencyService(ctrl *gomock.Controller) *MockCurrencyService {
	mock := &MockCurrencyService{ctrl: ctrl}
	mock.recorder = &MockCurrencyServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCurrencyService) EXPECT() *MockCurrencyServiceMockRecorder {
	return m.recorder
}

// Convert mocks base method
func (m *MockCurrencyService) Convert(from, to string, amount float64) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Convert", from, to, amount)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Convert indicates an expected call of Convert
func (mr *MockCurrencyServiceMockRecorder) Convert(from, to, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Convert", reflect.TypeOf((*MockCurrencyService)(nil).Convert), from, to, amount)
}

// GetMostRecentCurrencyDataRaw mocks base method
func (m *MockCurrencyService) GetMostRecentCurrencyDataRaw() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMostRecentCurrencyDataRaw")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMostRecentCurrencyDataRaw indicates an expected call of GetMostRecentCurrencyDataRaw
func (mr *MockCurrencyServiceMockRecorder) GetMostRecentCurrencyDataRaw() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMostRecentCurrencyDataRaw", reflect.TypeOf((*MockCurrencyService)(nil).GetMostRecentCurrencyDataRaw))
}

// SyncToDb mocks base method
func (m *MockCurrencyService) SyncToDb() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncToDb")
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncToDb indicates an expected call of SyncToDb
func (mr *MockCurrencyServiceMockRecorder) SyncToDb() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncToDb", reflect.TypeOf((*MockCurrencyService)(nil).SyncToDb))
}

// CurrencySyncRoutine mocks base method
func (m *MockCurrencyService) CurrencySyncRoutine() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CurrencySyncRoutine")
}

// CurrencySyncRoutine indicates an expected call of CurrencySyncRoutine
func (mr *MockCurrencyServiceMockRecorder) CurrencySyncRoutine() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrencySyncRoutine", reflect.TypeOf((*MockCurrencyService)(nil).CurrencySyncRoutine))
}

// CurrencySyncWorker mocks base method
func (m *MockCurrencyService) CurrencySyncWorker() (*IrisAPIs.CurrencySyncResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrencySyncWorker")
	ret0, _ := ret[0].(*IrisAPIs.CurrencySyncResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CurrencySyncWorker indicates an expected call of CurrencySyncWorker
func (mr *MockCurrencyServiceMockRecorder) CurrencySyncWorker() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrencySyncWorker", reflect.TypeOf((*MockCurrencyService)(nil).CurrencySyncWorker))
}
