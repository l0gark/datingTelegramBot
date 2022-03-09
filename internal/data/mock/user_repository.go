// Code generated by MockGen. DO NOT EDIT.
// Source: user_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/Eretic431/datingTelegramBot/internal/data/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockUserRepository) Add(ctx context.Context, user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockUserRepositoryMockRecorder) Add(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockUserRepository)(nil).Add), ctx, user)
}

// DeleteByUserId mocks base method.
func (m *MockUserRepository) DeleteByUserId(ctx context.Context, userId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByUserId", ctx, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByUserId indicates an expected call of DeleteByUserId.
func (mr *MockUserRepositoryMockRecorder) DeleteByUserId(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByUserId", reflect.TypeOf((*MockUserRepository)(nil).DeleteByUserId), ctx, userId)
}

// GetByUserId mocks base method.
func (m *MockUserRepository) GetByUserId(ctx context.Context, userId string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserId", ctx, userId)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserId indicates an expected call of GetByUserId.
func (mr *MockUserRepositoryMockRecorder) GetByUserId(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserId", reflect.TypeOf((*MockUserRepository)(nil).GetByUserId), ctx, userId)
}

// GetNextUser mocks base method.
func (m *MockUserRepository) GetNextUser(ctx context.Context, userId string, sex bool) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextUser", ctx, userId, sex)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNextUser indicates an expected call of GetNextUser.
func (mr *MockUserRepositoryMockRecorder) GetNextUser(ctx, userId, sex interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextUser", reflect.TypeOf((*MockUserRepository)(nil).GetNextUser), ctx, userId, sex)
}

// UpdateByUserId mocks base method.
func (m *MockUserRepository) UpdateByUserId(ctx context.Context, user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateByUserId", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateByUserId indicates an expected call of UpdateByUserId.
func (mr *MockUserRepositoryMockRecorder) UpdateByUserId(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateByUserId", reflect.TypeOf((*MockUserRepository)(nil).UpdateByUserId), ctx, user)
}
