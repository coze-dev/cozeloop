// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/modules/foundation/domain/authn/entity"
)

//go:generate mockgen -destination=mocks/authn_repo.go -package=mocks . IAuthNRepo
type IAuthNRepo interface {
	CreateAPIKey(ctx context.Context, apiKeyEntity *entity.APIKey) (apiKeyID int64, apiKey string, err error)
	DeleteAPIKey(ctx context.Context, apiKeyID int64) (err error)
	GetAPIKeyByIDs(ctx context.Context, apiKeyIDs []int64) (apiKeys []*entity.APIKey, err error)
	GetAPIKeyByUser(ctx context.Context, userID int64, pageNumber, pageSize int) (apiKeys []*entity.APIKey, err error)
	GetAPIKeyByKey(ctx context.Context, key string) (apiKey *entity.APIKey, err error)
	UpdateAPIKeyName(ctx context.Context, apiKeyID int64, name string) (err error)
	FlushAPIKeyUsedTime(ctx context.Context, apiKeyID int64) (err error)
}
