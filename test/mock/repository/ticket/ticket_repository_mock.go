// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/ticket/ticket_repository.go
//
// Generated by this command:
//
//	mockgen -source=./internal/repository/ticket/ticket_repository.go -destination=test/mock/repository/ticket/ticket_repository_mock.go
//

// Package mock_ticket is a generated GoMock package.
package mock_ticket

import (
	reflect "reflect"

	entity "github.com/TrinityKnights/Backend/internal/domain/entity"
	model "github.com/TrinityKnights/Backend/internal/domain/model"
	gomock "go.uber.org/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockTicketRepository is a mock of TicketRepository interface.
type MockTicketRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTicketRepositoryMockRecorder
	isgomock struct{}
}

// MockTicketRepositoryMockRecorder is the mock recorder for MockTicketRepository.
type MockTicketRepositoryMockRecorder struct {
	mock *MockTicketRepository
}

// NewMockTicketRepository creates a new mock instance.
func NewMockTicketRepository(ctrl *gomock.Controller) *MockTicketRepository {
	mock := &MockTicketRepository{ctrl: ctrl}
	mock.recorder = &MockTicketRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTicketRepository) EXPECT() *MockTicketRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTicketRepository) Create(db *gorm.DB, entity *entity.Ticket) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", db, entity)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockTicketRepositoryMockRecorder) Create(db, entity any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTicketRepository)(nil).Create), db, entity)
}

// CreateBatch mocks base method.
func (m *MockTicketRepository) CreateBatch(db *gorm.DB, tickets []*entity.Ticket) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBatch", db, tickets)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateBatch indicates an expected call of CreateBatch.
func (mr *MockTicketRepositoryMockRecorder) CreateBatch(db, tickets any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBatch", reflect.TypeOf((*MockTicketRepository)(nil).CreateBatch), db, tickets)
}

// Delete mocks base method.
func (m *MockTicketRepository) Delete(db *gorm.DB, entity *entity.Ticket) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", db, entity)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTicketRepositoryMockRecorder) Delete(db, entity any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTicketRepository)(nil).Delete), db, entity)
}

// Find mocks base method.
func (m *MockTicketRepository) Find(db *gorm.DB, filter *model.TicketQueryOptions) ([]*entity.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", db, filter)
	ret0, _ := ret[0].([]*entity.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockTicketRepositoryMockRecorder) Find(db, filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockTicketRepository)(nil).Find), db, filter)
}

// GetLastTicketNumber mocks base method.
func (m *MockTicketRepository) GetLastTicketNumber(db *gorm.DB, eventID uint, ticketType string) (*entity.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastTicketNumber", db, eventID, ticketType)
	ret0, _ := ret[0].(*entity.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastTicketNumber indicates an expected call of GetLastTicketNumber.
func (mr *MockTicketRepositoryMockRecorder) GetLastTicketNumber(db, eventID, ticketType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastTicketNumber", reflect.TypeOf((*MockTicketRepository)(nil).GetLastTicketNumber), db, eventID, ticketType)
}

// Update mocks base method.
func (m *MockTicketRepository) Update(db *gorm.DB, entity *entity.Ticket) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", db, entity)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTicketRepositoryMockRecorder) Update(db, entity any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTicketRepository)(nil).Update), db, entity)
}
