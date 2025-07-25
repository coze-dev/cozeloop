// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package rpc

import "context"

//go:generate mockgen -destination=mocks/auth.go -package=mocks . IAuthProvider
type IAuthProvider interface {
	Authorization(ctx context.Context, param *AuthorizationParam) (err error)
	AuthorizationWithoutSPI(ctx context.Context, param *AuthorizationWithoutSPIParam) (err error)
}

type AuthorizationParam struct {
	ObjectID      string
	SpaceID       int64
	ActionObjects []*ActionObject
}

type AuthorizationWithoutSPIParam struct {
	ObjectID      string
	SpaceID       int64
	ActionObjects []*ActionObject

	OwnerID         *string
	ResourceSpaceID int64
}

type ActionObject struct {
	Action     *string
	EntityType *AuthEntityType
}

type AuthEntityType = string

const (
	AuthEntityType_Space         = "Space"
	AuthEntityType_EvaluationSet = "EvaluationSet"
)

const (
	CommonActionRead = "read"
	CommonActionEdit = "edit"

	CozeActionListLoopEvaluationSet   = "listLoopEvaluationSet"   // 数据集列表读权限
	CozeActionCreateLoopEvaluationSet = "createLoopEvaluationSet" // 数据集创建权限
)
