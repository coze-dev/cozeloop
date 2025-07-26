// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mysql

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/infra/db"
	"github.com/coze-dev/coze-loop/backend/modules/foundation/infra/repo/mysql/gorm_gen/model"
	"github.com/coze-dev/coze-loop/backend/modules/foundation/infra/repo/mysql/gorm_gen/query"
)

//go:generate mockgen -destination=mocks/user_dao.go -package=mocks . IUserDAO
type IUserDAO interface {
	Create(ctx context.Context, user *model.User, opts ...db.Option) error
	Save(ctx context.Context, user *model.User, opts ...db.Option) error

	GetByID(ctx context.Context, userID int64, opts ...db.Option) (*model.User, error)
	MGetByIDs(ctx context.Context, userIDs []int64, opts ...db.Option) ([]*model.User, error)
	FindByEmail(ctx context.Context, email string, opts ...db.Option) (*model.User, error)
	FindByUniqueName(ctx context.Context, uniqueName string, opts ...db.Option) (*model.User, error)

	Update(ctx context.Context, userID int64, updates map[string]interface{}, opts ...db.Option) error
}

type UserDAOImpl struct {
	db    db.Provider
	query *query.Query
}

func NewUserDAOImpl(db db.Provider) IUserDAO {
	return &UserDAOImpl{
		db:    db,
		query: query.Use(db.NewSession(context.Background())),
	}
}

func (dao *UserDAOImpl) Create(ctx context.Context, user *model.User, opts ...db.Option) error {
	dao.query = query.Use(dao.db.NewSession(ctx, opts...))
	return dao.query.User.WithContext(ctx).Create(user)
}

func (dao *UserDAOImpl) Save(ctx context.Context, user *model.User, opts ...db.Option) error {
	dao.query = query.Use(dao.db.NewSession(ctx, opts...))
	return dao.query.User.WithContext(ctx).Save(user)
}

func (dao *UserDAOImpl) FindByUniqueName(ctx context.Context, uniqueName string, opts ...db.Option) (*model.User, error) {
	dao.query = query.Use(dao.db.NewSession(ctx, opts...))
	return dao.query.User.WithContext(ctx).Where(
		dao.query.User.UniqueName.Eq(uniqueName),
	).First()
}

func (dao *UserDAOImpl) FindByEmail(ctx context.Context, email string, opts ...db.Option) (*model.User, error) {
	dao.query = query.Use(dao.db.NewSession(ctx, opts...))

	return dao.query.User.WithContext(ctx).Where(
		dao.query.User.Email.Eq(email),
	).First()
}

func (dao *UserDAOImpl) Update(ctx context.Context, userID int64, updates map[string]interface{}, opts ...db.Option) error {
	dao.query = query.Use(dao.db.NewSession(ctx, opts...))

	_, err := dao.query.User.WithContext(ctx).Where(
		dao.query.User.ID.Eq(userID),
	).Updates(updates)

	return err
}

func (dao *UserDAOImpl) GetByID(ctx context.Context, userID int64, opts ...db.Option) (*model.User, error) {
	dao.query = query.Use(dao.db.NewSession(ctx, opts...))

	return dao.query.User.WithContext(ctx).Where(
		dao.query.User.ID.Eq(userID),
	).First()
}

func (dao *UserDAOImpl) MGetByIDs(ctx context.Context, userIDs []int64, opts ...db.Option) ([]*model.User, error) {
	dao.query = query.Use(dao.db.NewSession(ctx, opts...))
	return dao.query.User.WithContext(ctx).Where(
		dao.query.User.ID.In(userIDs...),
	).Find()
}
