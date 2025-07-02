// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package mysql

import (
	"context"
	"time"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/authn/entity"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql/gorm_gen/query"
)

type IAuthNDAO interface {
	CreateAPIKey(ctx context.Context, apiKeyEntity *model.APIKey) (err error)
	DeleteAPIKey(ctx context.Context, apiKeyID int64) (err error)
	GetAPIKeyByIDs(ctx context.Context, apiKeyIDs []int64) (apiKeys []*model.APIKey, err error)
	GetAPIKeyByUser(ctx context.Context, userID int64, pageNumber, pageSize int) (apiKeys []*model.APIKey, err error)
	GetAPIKeyByKey(ctx context.Context, key string) (apiKey *model.APIKey, err error)
	UpdateAPIKeyName(ctx context.Context, apiKeyID int64, name string) (err error)
	FlushAPIKeyUsedTime(ctx context.Context, apiKeyID int64) (err error)
}

type AuthNDAOImpl struct {
	db    db.Provider
	query *query.Query
}

func NewAuthNDAOImpl(db db.Provider) IAuthNDAO {
	return &AuthNDAOImpl{
		db:    db,
		query: query.Use(db.NewSession(context.Background())),
	}
}

func (dao *AuthNDAOImpl) CreateAPIKey(ctx context.Context, apiKeyPO *model.APIKey) (err error) {
	return dao.query.APIKey.WithContext(ctx).Create(apiKeyPO)
}

func (dao *AuthNDAOImpl) DeleteAPIKey(ctx context.Context, apiKeyID int64) (err error) {
	_, err = dao.query.APIKey.WithContext(ctx).Where(
		dao.query.APIKey.ID.Eq(apiKeyID),
	).Updates(map[string]interface{}{
		"status":     entity.APIKeyStatusDeleted,
		"deleted_at": time.Now().Unix(),
	})

	return err
}

func (dao *AuthNDAOImpl) GetAPIKeyByIDs(ctx context.Context, apiKeyIDs []int64) (apiKeys []*model.APIKey, err error) {
	return dao.query.APIKey.WithContext(ctx).
		Where(dao.query.APIKey.ID.In(apiKeyIDs...)).
		Where(dao.query.APIKey.Status.Eq(entity.APIKeyStatusNormal)).
		Find()
}

func (dao *AuthNDAOImpl) UpdateAPIKeyName(ctx context.Context, apiKeyID int64, name string) (err error) {
	_, err = dao.query.APIKey.WithContext(ctx).
		Where(dao.query.APIKey.ID.Eq(apiKeyID)).
		Updates(map[string]interface{}{
			"name": name,
		})

	return err
}

func (dao *AuthNDAOImpl) GetAPIKeyByUser(ctx context.Context, userID int64, pageNumber, pageSize int) (apiKeys []*model.APIKey, err error) {
	return dao.query.APIKey.WithContext(ctx).
		Where(dao.query.APIKey.UserID.Eq(userID)).
		Where(dao.query.APIKey.Status.Eq(entity.APIKeyStatusNormal)).
		Offset((pageNumber - 1) * pageSize).
		Limit(pageSize).
		Find()
}

func (dao *AuthNDAOImpl) GetAPIKeyByKey(ctx context.Context, key string) (apiKey *model.APIKey, err error) {
	apiKeys, err := dao.query.APIKey.WithContext(ctx).
		Where(dao.query.APIKey.Key.Eq(key)).
		Where(dao.query.APIKey.Status.Eq(entity.APIKeyStatusNormal)).
		Find()
	if err != nil {
		return nil, err
	}

	if len(apiKeys) == 0 {
		return nil, nil
	}
	return apiKeys[0], nil
}

func (dao *AuthNDAOImpl) FlushAPIKeyUsedTime(ctx context.Context, apiKeyID int64) (err error) {
	_, err = dao.query.APIKey.WithContext(ctx).
		Where(dao.query.APIKey.ID.Eq(apiKeyID)).
		Updates(map[string]interface{}{
			"last_used_at": time.Now().Unix(),
		})

	return err
}
