// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mysql

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
)

func wrapDBErr(err error, msgFormat string, args ...any) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errno.NotFoundErrorf(msgFormat, args...)
	}
	return errno.DBErr(err, fmt.Sprintf(msgFormat, args...))
}
