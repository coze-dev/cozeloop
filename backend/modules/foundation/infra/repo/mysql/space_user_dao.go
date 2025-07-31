// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mysql

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/infra/db"
	"github.com/coze-dev/coze-loop/backend/modules/foundation/infra/repo/mysql/gorm_gen/model"
	"github.com/coze-dev/coze-loop/backend/modules/foundation/infra/repo/mysql/gorm_gen/query"
)

//go:generate mockgen -destination=mocks/space_user_dao.go -package=mocks . ISpaceUserDAO
type ISpaceUserDAO interface {
	Create(ctx context.Context, spaceUser *model.SpaceUser, opts ...db.Option) error
	List(ctx context.Context, userID int64, pageSize, pageNumber int32, opts ...db.Option) ([]*model.SpaceUser, int32, error)
}

const (
	defaultPageSize = 20
	defaultPageNum  = 1
)

type SpaceUserDAOImpl struct {
	db    db.Provider
	query *query.Query
}

func NewSpaceUserDAOImpl(db db.Provider) ISpaceUserDAO {
	return &SpaceUserDAOImpl{
		db:    db,
		query: query.Use(db.NewSession(context.Background())),
	}
}

func (dao *SpaceUserDAOImpl) Create(ctx context.Context, spaceUser *model.SpaceUser, opts ...db.Option) error {
	dao.query = query.Use(dao.db.NewSession(ctx, opts...))
	return dao.query.SpaceUser.WithContext(ctx).Create(spaceUser)
}

func (dao *SpaceUserDAOImpl) List(ctx context.Context, userID int64, pageSize, pageNumber int32, opts ...db.Option) ([]*model.SpaceUser, int32, error) {
	dao.query = query.Use(dao.db.NewSession(ctx, opts...))

	q := dao.query.SpaceUser.WithContext(ctx).Where(
		dao.query.SpaceUser.UserID.Eq(userID),
	)

	total, err := q.Count()
	if err != nil || total == 0 {
		return nil, 0, err
	}

	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	if pageNumber == 0 {
		pageNumber = defaultPageNum
	}
	spaceUsers, err := q.Offset(int((pageNumber - 1) * pageSize)).Limit(int(pageSize)).Find()
	if err != nil {
		return nil, 0, err
	}
	return spaceUsers, int32(total), nil
}
