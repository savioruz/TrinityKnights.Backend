// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/delivery/http/handler/ticket/ticket_handler.go
//
// Generated by this command:
//
//	mockgen -source=./internal/delivery/http/handler/ticket/ticket_handler.go -destination=test/mock/delivery/http/handler/ticket/ticket_handler_mock.go
//

// Package mock_ticket is a generated GoMock package.
package mock_ticket

import (
	reflect "reflect"

	echo "github.com/labstack/echo/v4"
	gomock "go.uber.org/mock/gomock"
)

// MockTicketHandler is a mock of TicketHandler interface.
type MockTicketHandler struct {
	ctrl     *gomock.Controller
	recorder *MockTicketHandlerMockRecorder
	isgomock struct{}
}

// MockTicketHandlerMockRecorder is the mock recorder for MockTicketHandler.
type MockTicketHandlerMockRecorder struct {
	mock *MockTicketHandler
}

// NewMockTicketHandler creates a new mock instance.
func NewMockTicketHandler(ctrl *gomock.Controller) *MockTicketHandler {
	mock := &MockTicketHandler{ctrl: ctrl}
	mock.recorder = &MockTicketHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTicketHandler) EXPECT() *MockTicketHandlerMockRecorder {
	return m.recorder
}

// CreateTicket mocks base method.
func (m *MockTicketHandler) CreateTicket(ctx echo.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTicket", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTicket indicates an expected call of CreateTicket.
func (mr *MockTicketHandlerMockRecorder) CreateTicket(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTicket", reflect.TypeOf((*MockTicketHandler)(nil).CreateTicket), ctx)
}

// GetAllTickets mocks base method.
func (m *MockTicketHandler) GetAllTickets(ctx echo.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTickets", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetAllTickets indicates an expected call of GetAllTickets.
func (mr *MockTicketHandlerMockRecorder) GetAllTickets(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTickets", reflect.TypeOf((*MockTicketHandler)(nil).GetAllTickets), ctx)
}

// GetTicketByID mocks base method.
func (m *MockTicketHandler) GetTicketByID(ctx echo.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTicketByID", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetTicketByID indicates an expected call of GetTicketByID.
func (mr *MockTicketHandlerMockRecorder) GetTicketByID(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTicketByID", reflect.TypeOf((*MockTicketHandler)(nil).GetTicketByID), ctx)
}

// SearchTickets mocks base method.
func (m *MockTicketHandler) SearchTickets(ctx echo.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchTickets", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// SearchTickets indicates an expected call of SearchTickets.
func (mr *MockTicketHandlerMockRecorder) SearchTickets(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchTickets", reflect.TypeOf((*MockTicketHandler)(nil).SearchTickets), ctx)
}

// UpdateTicket mocks base method.
func (m *MockTicketHandler) UpdateTicket(ctx echo.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTicket", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTicket indicates an expected call of UpdateTicket.
func (mr *MockTicketHandlerMockRecorder) UpdateTicket(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTicket", reflect.TypeOf((*MockTicketHandler)(nil).UpdateTicket), ctx)
}
