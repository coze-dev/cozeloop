// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/gg/gptr"
	"github.com/samber/lo"

	"github.com/coze-dev/cozeloop/backend/infra/middleware/session"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/authn"
	authn2 "github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/domain/authn"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/authn/entity"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/authn/repo"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
)

const (
	PageSizeDefault   = 10
	PageNumberDefault = 1
)

type AuthNApplicationImpl struct {
	authNRepo repo.IAuthNRepo
}

func NewAuthNApplication(authNRepo repo.IAuthNRepo) authn.AuthNService {
	return &AuthNApplicationImpl{
		authNRepo: authNRepo,
	}
}

func (a AuthNApplicationImpl) CreatePersonalAccessToken(ctx context.Context, req *authn.CreatePersonalAccessTokenRequest) (r *authn.CreatePersonalAccessTokenResponse, err error) {
	userIDStr, ok := session.UserIDInCtx(ctx)
	if !ok {
		return nil, errorx.NewByCode(errno.CommonNoPermissionCode)
	}
	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, err
	}
	if req.DurationDay == nil && req.ExpireAt == nil {
		return nil, errorx.WrapByCode(errorx.New("param expire_at and duration_day is empty"), errno.CommonInvalidParamCode)
	}

	now := time.Now()
	var expireAt int64
	if req.ExpireAt != nil {
		expireAt = *req.ExpireAt
	}
	if req.DurationDay != nil {
		if *req.DurationDay == authn.DurationDayPermanent {
			expireAt = time.Now().Add(time.Duration(99) * time.Hour * 24 * 365).Unix()
		} else {
			expireDay, err := strconv.ParseInt(*req.DurationDay, 10, 64)
			if err != nil {
				return nil, errorx.WrapByCode(err, errno.CommonInvalidParamCode)
			}
			expireAt = time.Now().Add(time.Duration(expireDay) * time.Hour * 24).Unix()
		}
	}

	apiKeyID, apiKey, err := a.authNRepo.CreateAPIKey(ctx, &entity.APIKey{
		Name:       req.Name,
		Status:     entity.APIKeyStatusNormal,
		UserID:     userIDInt,
		ExpiredAt:  expireAt,
		CreatedAt:  now,
		UpdatedAt:  now,
		DeletedAt:  0,
		LastUsedAt: -1,
	})
	if err != nil {
		return nil, err
	}

	return &authn.CreatePersonalAccessTokenResponse{
		PersonalAccessToken: &authn2.PersonalAccessToken{
			ID:         strconv.FormatInt(apiKeyID, 10),
			Name:       req.Name,
			CreatedAt:  now.Unix(),
			UpdatedAt:  now.Unix(),
			LastUsedAt: 0,
			ExpireAt:   expireAt,
		},
		Token: lo.ToPtr(apiKey),
	}, nil
}

func (a AuthNApplicationImpl) DeletePersonalAccessToken(ctx context.Context, req *authn.DeletePersonalAccessTokenRequest) (r *authn.DeletePersonalAccessTokenResponse, err error) {
	userIDStr, ok := session.UserIDInCtx(ctx)
	if !ok {
		return nil, errorx.Wrapf(err, "user id not found")
	}
	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonInvalidParamCode, errorx.WithExtraMsg("invalid user id"))
	}

	ds, err := a.authNRepo.GetAPIKeyByIDs(ctx, []int64{req.ID})
	if err != nil {
		return nil, err
	}
	if len(ds) == 0 {
		return &authn.DeletePersonalAccessTokenResponse{}, nil
	}

	if ds[0].UserID != userIDInt {
		return nil, errorx.NewByCode(errno.CommonNoPermissionCode)
	}

	if err = a.authNRepo.DeleteAPIKey(ctx, req.ID); err != nil {
		return nil, err
	}

	return &authn.DeletePersonalAccessTokenResponse{}, nil
}

func (a AuthNApplicationImpl) UpdatePersonalAccessToken(ctx context.Context, req *authn.UpdatePersonalAccessTokenRequest) (r *authn.UpdatePersonalAccessTokenResponse, err error) {
	userIDStr, ok := session.UserIDInCtx(ctx)
	if !ok {
		return nil, errorx.Wrapf(err, "user id not found")
	}
	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, err
	}

	ds, err := a.authNRepo.GetAPIKeyByIDs(ctx, []int64{req.ID})
	if err != nil {
		return nil, err
	}
	if len(ds) == 0 {
		return &authn.UpdatePersonalAccessTokenResponse{}, nil
	}

	if ds[0].UserID != userIDInt {
		return nil, errorx.NewByCode(errno.CommonNoPermissionCode)
	}

	err = a.authNRepo.UpdateAPIKeyName(ctx, req.ID, req.Name)
	if err != nil {
		return nil, err
	}

	return &authn.UpdatePersonalAccessTokenResponse{}, nil
}

func (a AuthNApplicationImpl) GetPersonalAccessToken(ctx context.Context, req *authn.GetPersonalAccessTokenRequest) (r *authn.GetPersonalAccessTokenResponse, err error) {
	userIDStr, ok := session.UserIDInCtx(ctx)
	if !ok {
		return nil, errorx.Wrapf(err, "user id not found")
	}
	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, err
	}

	ds, err := a.authNRepo.GetAPIKeyByIDs(ctx, []int64{req.ID})
	if err != nil {
		return nil, err
	}
	if len(ds) == 0 {
		return &authn.GetPersonalAccessTokenResponse{}, nil
	}

	if ds[0].UserID != userIDInt {
		return nil, errorx.NewByCode(errno.CommonNoPermissionCode)
	}

	return &authn.GetPersonalAccessTokenResponse{
		PersonalAccessToken: &authn2.PersonalAccessToken{
			ID:         strconv.FormatInt(ds[0].ID, 10),
			Name:       ds[0].Name,
			CreatedAt:  ds[0].CreatedAt.Unix(),
			UpdatedAt:  ds[0].UpdatedAt.Unix(),
			LastUsedAt: ds[0].LastUsedAt,
			ExpireAt:   ds[0].ExpiredAt,
		},
	}, nil
}

func (a AuthNApplicationImpl) ListPersonalAccessToken(ctx context.Context, req *authn.ListPersonalAccessTokenRequest) (r *authn.ListPersonalAccessTokenResponse, err error) {
	userIDStr, ok := session.UserIDInCtx(ctx)
	if !ok {
		return nil, errorx.Wrapf(err, "user id not found")
	}
	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonInvalidParamCode, errorx.WithExtraMsg("invalid user id"))
	}

	pageNumber := PageNumberDefault
	pageSize := PageSizeDefault
	if req.PageNumber != nil {
		pageNumber = int(*req.PageNumber)
	}
	if req.PageSize != nil {
		pageSize = int(*req.PageSize)
	}
	apiKeys, err := a.authNRepo.GetAPIKeyByUser(ctx, userIDInt, pageNumber, pageSize)
	if err != nil {
		return nil, err
	}

	return &authn.ListPersonalAccessTokenResponse{
		PersonalAccessTokens: lo.Map(apiKeys, func(item *entity.APIKey, index int) *authn2.PersonalAccessToken {
			return &authn2.PersonalAccessToken{
				ID:         strconv.FormatInt(item.ID, 10),
				Name:       item.Name,
				CreatedAt:  item.CreatedAt.Unix(),
				UpdatedAt:  item.UpdatedAt.Unix(),
				LastUsedAt: item.LastUsedAt,
				ExpireAt:   item.ExpiredAt,
			}
		}),
	}, nil
}

func (a AuthNApplicationImpl) VerifyToken(ctx context.Context, req *authn.VerifyTokenRequest) (r *authn.VerifyTokenResponse, err error) {
	varifyFail := &authn.VerifyTokenResponse{Valid: lo.ToPtr(false)}
	varifyPass := &authn.VerifyTokenResponse{Valid: lo.ToPtr(true)}

	ds, err := a.authNRepo.GetAPIKeyByKey(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	if ds == nil {
		return varifyFail, errorx.NewByCode(errno.CommonNoPermissionCode)
	}

	if ds.Status == entity.APIKeyStatusDeleted {
		return varifyFail, errorx.NewByCode(errno.CommonNoPermissionCode)
	}

	if time.Now().Unix() > ds.ExpiredAt {
		return varifyFail, errorx.NewByCode(errno.CommonNoPermissionCode)
	}

	err = a.authNRepo.FlushAPIKeyUsedTime(ctx, ds.ID)
	if err != nil {
		return nil, err
	}

	varifyPass.UserID = gptr.Of(strconv.FormatInt(ds.UserID, 10))
	return varifyPass, nil
}
