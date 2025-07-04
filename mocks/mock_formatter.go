// Code generated by MockGen. DO NOT EDIT.
// Source: formatter.go
//
// Generated by this command:
//
//	mockgen -source=formatter.go -destination=../../../mocks/mock_formatter.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockFormatter is a mock of Formatter interface.
type MockFormatter struct {
	ctrl     *gomock.Controller
	recorder *MockFormatterMockRecorder
}

// MockFormatterMockRecorder is the mock recorder for MockFormatter.
type MockFormatterMockRecorder struct {
	mock *MockFormatter
}

// NewMockFormatter creates a new mock instance.
func NewMockFormatter(ctrl *gomock.Controller) *MockFormatter {
	mock := &MockFormatter{ctrl: ctrl}
	mock.recorder = &MockFormatterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFormatter) EXPECT() *MockFormatterMockRecorder {
	return m.recorder
}

// FormatMarkdown mocks base method.
func (m *MockFormatter) FormatMarkdown(ctx context.Context, text string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FormatMarkdown", ctx, text)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FormatMarkdown indicates an expected call of FormatMarkdown.
func (mr *MockFormatterMockRecorder) FormatMarkdown(ctx, text any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FormatMarkdown", reflect.TypeOf((*MockFormatter)(nil).FormatMarkdown), ctx, text)
}
