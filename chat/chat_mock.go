// Code generated by MockGen. DO NOT EDIT.
// Source: chat.go

// Package chat is a generated GoMock package.
package chat

import (
	gomock "github.com/golang/mock/gomock"
	sqlx "github.com/jmoiron/sqlx"
	reflect "reflect"
	time "time"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Insert mocks base method
func (m *MockService) Insert(tx *sqlx.Tx, senderID, recipientID int64, messageContent Content) (int64, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", tx, senderID, recipientID, messageContent)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(time.Time)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Insert indicates an expected call of Insert
func (mr *MockServiceMockRecorder) Insert(tx, senderID, recipientID, messageContent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockService)(nil).Insert), tx, senderID, recipientID, messageContent)
}

// GetMessagesByRecipient mocks base method
func (m *MockService) GetMessagesByRecipient(tx *sqlx.Tx, senderID, recipientID, start int64, limit int) ([]*Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessagesByRecipient", tx, senderID, recipientID, start, limit)
	ret0, _ := ret[0].([]*Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessagesByRecipient indicates an expected call of GetMessagesByRecipient
func (mr *MockServiceMockRecorder) GetMessagesByRecipient(tx, senderID, recipientID, start, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessagesByRecipient", reflect.TypeOf((*MockService)(nil).GetMessagesByRecipient), tx, senderID, recipientID, start, limit)
}
