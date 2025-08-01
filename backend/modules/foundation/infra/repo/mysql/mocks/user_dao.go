// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/coze-dev/coze-loop/backend/modules/foundation/infra/repo/mysql (interfaces: IUserDAO)
//
// Generated by this command:
//
//	mockgen -destination=mocks/user_dao.go -package=mocks . IUserDAO
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	db "github.com/coze-dev/coze-loop/backend/infra/db"
	model "github.com/coze-dev/coze-loop/backend/modules/foundation/infra/repo/mysql/gorm_gen/model"
	gomock "go.uber.org/mock/gomock"
)

// MockIUserDAO is a mock of IUserDAO interface.
type MockIUserDAO struct {
	ctrl     *gomock.Controller
	recorder *MockIUserDAOMockRecorder
	isgomock struct{}
}

// MockIUserDAOMockRecorder is the mock recorder for MockIUserDAO.
type MockIUserDAOMockRecorder struct {
	mock *MockIUserDAO
}

// NewMockIUserDAO creates a new mock instance.
func NewMockIUserDAO(ctrl *gomock.Controller) *MockIUserDAO {
	mock := &MockIUserDAO{ctrl: ctrl}
	mock.recorder = &MockIUserDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUserDAO) EXPECT() *MockIUserDAOMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIUserDAO) Create(ctx context.Context, user *model.User, opts ...db.Option) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, user}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockIUserDAOMockRecorder) Create(ctx, user any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, user}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIUserDAO)(nil).Create), varargs...)
}

// FindByEmail mocks base method.
func (m *MockIUserDAO) FindByEmail(ctx context.Context, email string, opts ...db.Option) (*model.User, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, email}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindByEmail", varargs...)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockIUserDAOMockRecorder) FindByEmail(ctx, email any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, email}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockIUserDAO)(nil).FindByEmail), varargs...)
}

// FindByUniqueName mocks base method.
func (m *MockIUserDAO) FindByUniqueName(ctx context.Context, uniqueName string, opts ...db.Option) (*model.User, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, uniqueName}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindByUniqueName", varargs...)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUniqueName indicates an expected call of FindByUniqueName.
func (mr *MockIUserDAOMockRecorder) FindByUniqueName(ctx, uniqueName any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, uniqueName}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUniqueName", reflect.TypeOf((*MockIUserDAO)(nil).FindByUniqueName), varargs...)
}

// GetByID mocks base method.
func (m *MockIUserDAO) GetByID(ctx context.Context, userID int64, opts ...db.Option) (*model.User, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, userID}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByID", varargs...)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIUserDAOMockRecorder) GetByID(ctx, userID any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, userID}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIUserDAO)(nil).GetByID), varargs...)
}

// MGetByIDs mocks base method.
func (m *MockIUserDAO) MGetByIDs(ctx context.Context, userIDs []int64, opts ...db.Option) ([]*model.User, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, userIDs}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "MGetByIDs", varargs...)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MGetByIDs indicates an expected call of MGetByIDs.
func (mr *MockIUserDAOMockRecorder) MGetByIDs(ctx, userIDs any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, userIDs}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MGetByIDs", reflect.TypeOf((*MockIUserDAO)(nil).MGetByIDs), varargs...)
}

// Save mocks base method.
func (m *MockIUserDAO) Save(ctx context.Context, user *model.User, opts ...db.Option) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, user}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Save", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIUserDAOMockRecorder) Save(ctx, user any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, user}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIUserDAO)(nil).Save), varargs...)
}

// Update mocks base method.
func (m *MockIUserDAO) Update(ctx context.Context, userID int64, updates map[string]any, opts ...db.Option) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, userID, updates}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Update", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIUserDAOMockRecorder) Update(ctx, userID, updates any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, userID, updates}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIUserDAO)(nil).Update), varargs...)
}
