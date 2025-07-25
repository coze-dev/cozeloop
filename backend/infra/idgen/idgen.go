// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package idgen

import (
	"context"
)

//go:generate mockgen -destination=mocks/idgen.go -package=mocks . IIDGenerator
type IIDGenerator interface {
	GenID(ctx context.Context) (int64, error)
	GenMultiIDs(ctx context.Context, counts int) ([]int64, error)
}
