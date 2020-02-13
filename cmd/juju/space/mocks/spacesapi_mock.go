// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/cmd/juju/space (interfaces: SpaceAPI)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	params "github.com/juju/juju/apiserver/params"
	network "github.com/juju/juju/core/network"
	reflect "reflect"
)

// MockSpaceAPI is a mock of SpaceAPI interface
type MockSpaceAPI struct {
	ctrl     *gomock.Controller
	recorder *MockSpaceAPIMockRecorder
}

// MockSpaceAPIMockRecorder is the mock recorder for MockSpaceAPI
type MockSpaceAPIMockRecorder struct {
	mock *MockSpaceAPI
}

// NewMockSpaceAPI creates a new mock instance
func NewMockSpaceAPI(ctrl *gomock.Controller) *MockSpaceAPI {
	mock := &MockSpaceAPI{ctrl: ctrl}
	mock.recorder = &MockSpaceAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSpaceAPI) EXPECT() *MockSpaceAPIMockRecorder {
	return m.recorder
}

// AddSpace mocks base method
func (m *MockSpaceAPI) AddSpace(arg0 string, arg1 []string, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSpace", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSpace indicates an expected call of AddSpace
func (mr *MockSpaceAPIMockRecorder) AddSpace(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSpace", reflect.TypeOf((*MockSpaceAPI)(nil).AddSpace), arg0, arg1, arg2)
}

// Close mocks base method
func (m *MockSpaceAPI) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockSpaceAPIMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSpaceAPI)(nil).Close))
}

// ListSpaces mocks base method
func (m *MockSpaceAPI) ListSpaces() ([]params.Space, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSpaces")
	ret0, _ := ret[0].([]params.Space)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSpaces indicates an expected call of ListSpaces
func (mr *MockSpaceAPIMockRecorder) ListSpaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSpaces", reflect.TypeOf((*MockSpaceAPI)(nil).ListSpaces))
}

// ReloadSpaces mocks base method
func (m *MockSpaceAPI) ReloadSpaces() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReloadSpaces")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReloadSpaces indicates an expected call of ReloadSpaces
func (mr *MockSpaceAPIMockRecorder) ReloadSpaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReloadSpaces", reflect.TypeOf((*MockSpaceAPI)(nil).ReloadSpaces))
}

// RemoveSpace mocks base method
func (m *MockSpaceAPI) RemoveSpace(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSpace", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSpace indicates an expected call of RemoveSpace
func (mr *MockSpaceAPIMockRecorder) RemoveSpace(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSpace", reflect.TypeOf((*MockSpaceAPI)(nil).RemoveSpace), arg0)
}

// RenameSpace mocks base method
func (m *MockSpaceAPI) RenameSpace(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenameSpace", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RenameSpace indicates an expected call of RenameSpace
func (mr *MockSpaceAPIMockRecorder) RenameSpace(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenameSpace", reflect.TypeOf((*MockSpaceAPI)(nil).RenameSpace), arg0, arg1)
}

// ShowSpace mocks base method
func (m *MockSpaceAPI) ShowSpace(arg0 string) (network.ShowSpace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShowSpace", arg0)
	ret0, _ := ret[0].(network.ShowSpace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShowSpace indicates an expected call of ShowSpace
func (mr *MockSpaceAPIMockRecorder) ShowSpace(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowSpace", reflect.TypeOf((*MockSpaceAPI)(nil).ShowSpace), arg0)
}

// MoveToSpace mocks base method
func (m *MockSpaceAPI) MoveToSpace(arg0 string, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MoveToSpace", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// MoveToSpace indicates an expected call of MoveToSpace
func (mr *MockSpaceAPIMockRecorder) UpdateSpace(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MoveToSpace", reflect.TypeOf((*MockSpaceAPI)(nil).MoveToSpace), arg0, arg1)
}
