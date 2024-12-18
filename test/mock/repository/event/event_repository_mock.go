package mock

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockEventRepository is a mock implementation of EventRepository.
type MockEventRepository struct {
	mock.Mock
}

// GetByID mocks the GetByID method of EventRepository.
func (m *MockEventRepository) GetByID(db *gorm.DB, event *entity.Event, id uint) error {
	args := m.Called(db, event, id)
	return args.Error(0)
}

// GetPaginated mocks the GetPaginated method of EventRepository.
func (m *MockEventRepository) GetPaginated(db *gorm.DB, events *[]entity.Event, opts *model.EventQueryOptions) (int64, error) {
	args := m.Called(db, events, opts)
	return args.Get(0).(int64), args.Error(1)
}

// Create mocks the Create method of the generic Repository.
func (m *MockEventRepository) Create(db *gorm.DB, entity *entity.Event) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

// Update mocks the Update method of the generic Repository.
func (m *MockEventRepository) Update(db *gorm.DB, entity *entity.Event) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

// Delete mocks the Delete method of the generic Repository.
func (m *MockEventRepository) Delete(db *gorm.DB, entity *entity.Event) error {
	args := m.Called(db, entity)
	return args.Error(0)
}
