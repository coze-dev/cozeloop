// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/authn/entity"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/authn/repo"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql/convertor"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
)

type AuthNRepoImpl struct {
	idgen    idgen.IIDGenerator
	authNDao mysql.IAuthNDAO
}

func NewAuthNRepo(
	idgen idgen.IIDGenerator,
	authNDao mysql.IAuthNDAO,
) repo.IAuthNRepo {
	return &AuthNRepoImpl{
		idgen:    idgen,
		authNDao: authNDao,
	}
}

func (a AuthNRepoImpl) CreateAPIKey(ctx context.Context, apiKeyEntity *entity.APIKey) (apiKeyID int64, apiKey string, err error) {
	if apiKeyEntity == nil {
		return 0, "", errorx.WrapByCode(err, errno.CommonInvalidParamCode, errorx.WithExtraMsg("AuthNRepoImpl.CreateAPIKey invalid param"))
	}

	apiKeyID, err = a.idgen.GenID(ctx)
	if err != nil {
		return 0, "", err
	}
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d", apiKeyID)))
	apiKey = hex.EncodeToString(hash[:])

	apiKeyPO := convertor.APIKeyDO2PO(apiKeyEntity)
	apiKeyPO.ID = apiKeyID
	apiKeyPO.Key = apiKey

	if err = a.authNDao.CreateAPIKey(ctx, apiKeyPO); err != nil {
		return 0, "", err
	}

	return apiKeyID, apiKey, nil
}

func (a AuthNRepoImpl) DeleteAPIKey(ctx context.Context, apiKeyID int64) (err error) {
	if err = a.authNDao.DeleteAPIKey(ctx, apiKeyID); err != nil {
		return err
	}

	return nil
}

func (a AuthNRepoImpl) GetAPIKeyByIDs(ctx context.Context, apiKeyIDs []int64) (apiKeys []*entity.APIKey, err error) {
	ds, err := a.authNDao.GetAPIKeyByIDs(ctx, apiKeyIDs)
	if err != nil {
		return nil, err
	}

	return convertor.APIKeysPO2DO(ds), nil
}

func (a AuthNRepoImpl) UpdateAPIKeyName(ctx context.Context, apiKeyID int64, name string) (err error) {
	err = a.authNDao.UpdateAPIKeyName(ctx, apiKeyID, name)
	if err != nil {
		return err
	}

	return nil
}

func (a AuthNRepoImpl) GetAPIKeyByUser(ctx context.Context, userID int64, pageNumber, pageSize int) (apiKeys []*entity.APIKey, err error) {
	apiKeysPO, err := a.authNDao.GetAPIKeyByUser(ctx, userID, pageNumber, pageSize)
	if err != nil {
		return nil, err
	}

	return convertor.APIKeysPO2DO(apiKeysPO), nil
}

func (a AuthNRepoImpl) GetAPIKeyByKey(ctx context.Context, key string) (apiKey *entity.APIKey, err error) {
	res, err := a.authNDao.GetAPIKeyByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	return convertor.APIKeyPO2DO(res), nil
}

func (a AuthNRepoImpl) FlushAPIKeyUsedTime(ctx context.Context, apiKeyID int64) (err error) {
	err = a.authNDao.FlushAPIKeyUsedTime(ctx, apiKeyID)
	if err != nil {
		return err
	}

	return nil
}
