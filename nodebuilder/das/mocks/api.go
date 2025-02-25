// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sunriselayer/sunrise-da/nodebuilder/das (interfaces: Module)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	das "github.com/sunriselayer/sunrise-da/das"
)

// MockModule is a mock of Module interface.
type MockModule struct {
	ctrl     *gomock.Controller
	recorder *MockModuleMockRecorder
}

// MockModuleMockRecorder is the mock recorder for MockModule.
type MockModuleMockRecorder struct {
	mock *MockModule
}

// NewMockModule creates a new mock instance.
func NewMockModule(ctrl *gomock.Controller) *MockModule {
	mock := &MockModule{ctrl: ctrl}
	mock.recorder = &MockModuleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModule) EXPECT() *MockModuleMockRecorder {
	return m.recorder
}

// SamplingStats mocks base method.
func (m *MockModule) SamplingStats(arg0 context.Context) (das.SamplingStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SamplingStats", arg0)
	ret0, _ := ret[0].(das.SamplingStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SamplingStats indicates an expected call of SamplingStats.
func (mr *MockModuleMockRecorder) SamplingStats(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SamplingStats", reflect.TypeOf((*MockModule)(nil).SamplingStats), arg0)
}

// WaitCatchUp mocks base method.
func (m *MockModule) WaitCatchUp(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitCatchUp", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitCatchUp indicates an expected call of WaitCatchUp.
func (mr *MockModuleMockRecorder) WaitCatchUp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitCatchUp", reflect.TypeOf((*MockModule)(nil).WaitCatchUp), arg0)
}
