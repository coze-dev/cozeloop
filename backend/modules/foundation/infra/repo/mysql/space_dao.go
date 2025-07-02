// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package mysql

import (
	"context"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql/gorm_gen/query"
)

//go:generate mockgen -destination=mocks/space_dao.go -package=mocks . ISpaceDAO
type ISpaceDAO interface {
	Create(ctx context.Context, space *model.Space, opts ...db.Option) error
	GetByID(ctx context.Context, spaceID int64, opts ...db.Option) (*model.Space, error)
	MGetByIDs(ctx context.Context, spaceIDs []int64, opts ...db.Option) ([]*model.Space, error)
}

type SpaceDAOImpl struct {
	db    db.Provider
	query *query.Query
}

func NewSpaceDAOImpl(db db.Provider) ISpaceDAO {
	return &SpaceDAOImpl{
		db:    db,
		query: query.Use(db.NewSession(context.Background())),
	}
}

func (dao *SpaceDAOImpl) Create(ctx context.Context, space *model.Space, opts ...db.Option) error {
	dao.query = query.Use(dao.db.NewSession(ctx, opts...))
	return dao.query.Space.WithContext(ctx).Create(space)
}

func (dao *SpaceDAOImpl) GetByID(ctx context.Context, spaceID int64, opts ...db.Option) (*model.Space, error) {
	dao.query = query.Use(dao.db.NewSession(ctx, opts...))
	return dao.query.Space.WithContext(ctx).Where(
		dao.query.Space.ID.Eq(spaceID),
	).First()
}

func (dao *SpaceDAOImpl) MGetByIDs(ctx context.Context, spaceIDs []int64, opts ...db.Option) ([]*model.Space, error) {
	dao.query = query.Use(dao.db.NewSession(ctx, opts...))
	return dao.query.Space.WithContext(ctx).Where(
		dao.query.Space.ID.In(spaceIDs...),
	).Find()
}
