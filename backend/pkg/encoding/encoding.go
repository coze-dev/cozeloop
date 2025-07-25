// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package encoding

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"

	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

func Encode(ctx context.Context, val interface{}) (res string) {
	bytes, err := sonic.Marshal(val)
	if err != nil {
		logs.CtxError(ctx, "failed to encode data", err)
		res = uuid.New().String()
		return res
	}
	h := md5.New()
	h.Write(bytes)
	res = hex.EncodeToString(h.Sum(nil))
	return res
}
