// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/repository_interfaces/worker.go

// Package mock_repository_interfaces is a generated GoMock package.
package mock_repository_interfaces

import (
	reflect "reflect"
	models "teamdev/internal/models"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockIWorkerRepository is a mock of IWorkerRepository interface.
type MockIWorkerRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIWorkerRepositoryMockRecorder
}

// MockIWorkerRepositoryMockRecorder is the mock recorder for MockIWorkerRepository.
type MockIWorkerRepositoryMockRecorder struct {
	mock *MockIWorkerRepository
}

// NewMockIWorkerRepository creates a new mock instance.
func NewMockIWorkerRepository(ctrl *gomock.Controller) *MockIWorkerRepository {
	mock := &MockIWorkerRepository{ctrl: ctrl}
	mock.recorder = &MockIWorkerRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIWorkerRepository) EXPECT() *MockIWorkerRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIWorkerRepository) Create(worker *models.Worker) (*models.Worker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", worker)
	ret0, _ := ret[0].(*models.Worker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIWorkerRepositoryMockRecorder) Create(worker interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIWorkerRepository)(nil).Create), worker)
}

// Delete mocks base method.
func (m *MockIWorkerRepository) Delete(id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockIWorkerRepositoryMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIWorkerRepository)(nil).Delete), id)
}

// GetAllWorkers mocks base method.
func (m *MockIWorkerRepository) GetAllWorkers() ([]models.Worker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllWorkers")
	ret0, _ := ret[0].([]models.Worker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllWorkers indicates an expected call of GetAllWorkers.
func (mr *MockIWorkerRepositoryMockRecorder) GetAllWorkers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllWorkers", reflect.TypeOf((*MockIWorkerRepository)(nil).GetAllWorkers))
}

// GetAverageOrderRate mocks base method.
func (m *MockIWorkerRepository) GetAverageOrderRate(worker *models.Worker) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAverageOrderRate", worker)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAverageOrderRate indicates an expected call of GetAverageOrderRate.
func (mr *MockIWorkerRepositoryMockRecorder) GetAverageOrderRate(worker interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAverageOrderRate", reflect.TypeOf((*MockIWorkerRepository)(nil).GetAverageOrderRate), worker)
}

// GetWorkerByEmail mocks base method.
func (m *MockIWorkerRepository) GetWorkerByEmail(email string) (*models.Worker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkerByEmail", email)
	ret0, _ := ret[0].(*models.Worker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkerByEmail indicates an expected call of GetWorkerByEmail.
func (mr *MockIWorkerRepositoryMockRecorder) GetWorkerByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkerByEmail", reflect.TypeOf((*MockIWorkerRepository)(nil).GetWorkerByEmail), email)
}

// GetWorkerByID mocks base method.
func (m *MockIWorkerRepository) GetWorkerByID(id uuid.UUID) (*models.Worker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkerByID", id)
	ret0, _ := ret[0].(*models.Worker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkerByID indicates an expected call of GetWorkerByID.
func (mr *MockIWorkerRepositoryMockRecorder) GetWorkerByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkerByID", reflect.TypeOf((*MockIWorkerRepository)(nil).GetWorkerByID), id)
}

// GetWorkersByRole mocks base method.
func (m *MockIWorkerRepository) GetWorkersByRole(role int) ([]models.Worker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkersByRole", role)
	ret0, _ := ret[0].([]models.Worker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkersByRole indicates an expected call of GetWorkersByRole.
func (mr *MockIWorkerRepositoryMockRecorder) GetWorkersByRole(role interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkersByRole", reflect.TypeOf((*MockIWorkerRepository)(nil).GetWorkersByRole), role)
}

// Update mocks base method.
func (m *MockIWorkerRepository) Update(worker *models.Worker) (*models.Worker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", worker)
	ret0, _ := ret[0].(*models.Worker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockIWorkerRepositoryMockRecorder) Update(worker interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIWorkerRepository)(nil).Update), worker)
}
