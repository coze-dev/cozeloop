// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package encoding

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	type args struct {
		ctx context.Context
		val interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
	}{
		{"success_string", args{context.Background(), "123"}, "202cb962ac59075b964b07152d234b70"},
		{"success_int", args{context.Background(), 123}, "202cb962ac59075b964b07152d234b70"},
		{"success_struct", args{context.Background(), struct{ Num int }{123}}, ""},
		// 模拟 sonic.Marshal 失败的情况
		{"marshal_failure", args{context.Background(), func() {}}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Encode(tt.args.ctx, tt.args.val)
			if tt.name == "marshal_failure" {
				// 验证在 marshal 失败时返回的是 UUID 格式的字符串
				_, err := uuid.Parse(res)
				assert.NoError(t, err)
			} else {
				assert.NotEmpty(t, res)
			}
		})
	}
}
