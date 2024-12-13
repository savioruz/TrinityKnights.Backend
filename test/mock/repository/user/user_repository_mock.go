package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of UserRepository.
type MockUserRepository struct {
	mock.Mock
}

// GetFirst mocks the GetFirst method of UserRepository.
func (m *MockUserRepository) GetFirst(db *gorm.DB, user *entity.User) error {
	args := m.Called(db, user)
	return args.Error(0)
}

// GetByID mocks the GetByID method of UserRepository.
func (m *MockUserRepository) GetByID(db *gorm.DB, user *entity.User, id string) error {
	args := m.Called(db, user, id)
	return args.Error(0)
}

// GetByEmail mocks the GetByEmail method of UserRepository.
func (m *MockUserRepository) GetByEmail(db *gorm.DB, user *entity.User, email string) error {
	args := m.Called(db, user, email)
	return args.Error(0)
}

// CountByRole mocks the CountByRole method of UserRepository.
func (m *MockUserRepository) CountByRole(db *gorm.DB, role string) (int64, error) {
	args := m.Called(db, role)
	return args.Get(0).(int64), args.Error(1)
}

// Create mocks the Create method of the generic Repository.
func (m *MockUserRepository) Create(db *gorm.DB, entity *entity.User) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

// Update mocks the Update method of the generic Repository.
func (m *MockUserRepository) Update(db *gorm.DB, entity *entity.User) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

// Delete mocks the Delete method of the generic Repository.
func (m *MockUserRepository) Delete(db *gorm.DB, entity *entity.User) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

// FindAll mocks the FindAll method of the generic Repository.
func (m *MockUserRepository) FindAll(db *gorm.DB, entities *[]entity.User) error {
	args := m.Called(db, entities)
	return args.Error(0)
}
