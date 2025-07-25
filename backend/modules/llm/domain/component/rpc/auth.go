// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package rpc

import "context"

//go:generate mockgen -destination=mocks/auth_provider.go -package=mocks . IAuthProvider
type IAuthProvider interface {
	CheckSpacePermission(ctx context.Context, spaceID int64, action string) error
}
