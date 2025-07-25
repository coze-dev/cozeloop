// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"context"
)

//go:generate mockgen -destination=mocks/config_provider.go -package=mocks . IConfigProvider
type IConfigProvider interface {
	GetPromptHubMaxQPSBySpace(ctx context.Context, spaceID int64) (maxQPS int, err error)
}
