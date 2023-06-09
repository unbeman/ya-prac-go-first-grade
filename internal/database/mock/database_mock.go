// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/unbeman/ya-prac-go-first-grade/internal/database (interfaces: Database)

// Package mock_database is a generated GoMock package.
package mock_database

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/unbeman/ya-prac-go-first-grade/internal/model"
)

// MockDatabase is a mock of Database interface.
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase.
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance.
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// CreateNewSession mocks base method.
func (m *MockDatabase) CreateNewSession(arg0 context.Context, arg1 *model.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewSession", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNewSession indicates an expected call of CreateNewSession.
func (mr *MockDatabaseMockRecorder) CreateNewSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewSession", reflect.TypeOf((*MockDatabase)(nil).CreateNewSession), arg0, arg1)
}

// CreateNewUser mocks base method.
func (m *MockDatabase) CreateNewUser(arg0 context.Context, arg1 *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNewUser indicates an expected call of CreateNewUser.
func (mr *MockDatabaseMockRecorder) CreateNewUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewUser", reflect.TypeOf((*MockDatabase)(nil).CreateNewUser), arg0, arg1)
}

// CreateNewUserOrder mocks base method.
func (m *MockDatabase) CreateNewUserOrder(arg0 context.Context, arg1 uint, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewUserOrder", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNewUserOrder indicates an expected call of CreateNewUserOrder.
func (mr *MockDatabaseMockRecorder) CreateNewUserOrder(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewUserOrder", reflect.TypeOf((*MockDatabase)(nil).CreateNewUserOrder), arg0, arg1, arg2)
}

// CreateWithdraw mocks base method.
func (m *MockDatabase) CreateWithdraw(arg0 context.Context, arg1 *model.User, arg2 model.WithdrawnInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWithdraw", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateWithdraw indicates an expected call of CreateWithdraw.
func (mr *MockDatabaseMockRecorder) CreateWithdraw(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWithdraw", reflect.TypeOf((*MockDatabase)(nil).CreateWithdraw), arg0, arg1, arg2)
}

// GetNotReadyUserOrders mocks base method.
func (m *MockDatabase) GetNotReadyUserOrders(arg0 context.Context, arg1 uint) ([]model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotReadyUserOrders", arg0, arg1)
	ret0, _ := ret[0].([]model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNotReadyUserOrders indicates an expected call of GetNotReadyUserOrders.
func (mr *MockDatabaseMockRecorder) GetNotReadyUserOrders(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotReadyUserOrders", reflect.TypeOf((*MockDatabase)(nil).GetNotReadyUserOrders), arg0, arg1)
}

// GetOrderByNumber mocks base method.
func (m *MockDatabase) GetOrderByNumber(arg0 context.Context, arg1 string) (*model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByNumber", arg0, arg1)
	ret0, _ := ret[0].(*model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByNumber indicates an expected call of GetOrderByNumber.
func (mr *MockDatabaseMockRecorder) GetOrderByNumber(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByNumber", reflect.TypeOf((*MockDatabase)(nil).GetOrderByNumber), arg0, arg1)
}

// GetUserByID mocks base method.
func (m *MockDatabase) GetUserByID(arg0 context.Context, arg1 uint) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", arg0, arg1)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockDatabaseMockRecorder) GetUserByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockDatabase)(nil).GetUserByID), arg0, arg1)
}

// GetUserByLogin mocks base method.
func (m *MockDatabase) GetUserByLogin(arg0 context.Context, arg1 string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", arg0, arg1)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockDatabaseMockRecorder) GetUserByLogin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockDatabase)(nil).GetUserByLogin), arg0, arg1)
}

// GetUserByToken mocks base method.
func (m *MockDatabase) GetUserByToken(arg0 context.Context, arg1 string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByToken", arg0, arg1)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByToken indicates an expected call of GetUserByToken.
func (mr *MockDatabaseMockRecorder) GetUserByToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByToken", reflect.TypeOf((*MockDatabase)(nil).GetUserByToken), arg0, arg1)
}

// GetUserOrders mocks base method.
func (m *MockDatabase) GetUserOrders(arg0 context.Context, arg1 uint) ([]model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserOrders", arg0, arg1)
	ret0, _ := ret[0].([]model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserOrders indicates an expected call of GetUserOrders.
func (mr *MockDatabaseMockRecorder) GetUserOrders(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserOrders", reflect.TypeOf((*MockDatabase)(nil).GetUserOrders), arg0, arg1)
}

// GetUserWithdrawals mocks base method.
func (m *MockDatabase) GetUserWithdrawals(arg0 context.Context, arg1 uint) ([]model.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserWithdrawals", arg0, arg1)
	ret0, _ := ret[0].([]model.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserWithdrawals indicates an expected call of GetUserWithdrawals.
func (mr *MockDatabaseMockRecorder) GetUserWithdrawals(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserWithdrawals", reflect.TypeOf((*MockDatabase)(nil).GetUserWithdrawals), arg0, arg1)
}

// UpdateUserBalanceAndOrder mocks base method.
func (m *MockDatabase) UpdateUserBalanceAndOrder(arg0 *model.Order, arg1 model.OrderAccrualInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserBalanceAndOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserBalanceAndOrder indicates an expected call of UpdateUserBalanceAndOrder.
func (mr *MockDatabaseMockRecorder) UpdateUserBalanceAndOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserBalanceAndOrder", reflect.TypeOf((*MockDatabase)(nil).UpdateUserBalanceAndOrder), arg0, arg1)
}
