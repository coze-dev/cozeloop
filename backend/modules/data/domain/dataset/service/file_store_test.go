// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_vfs "github.com/coze-dev/coze-loop/backend/modules/data/domain/component/vfs/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/entity"
)

func TestDatasetServiceImpl_StatFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFS := mock_vfs.NewMockIUnionFS(ctrl)
	service := &DatasetServiceImpl{
		fsUnion: mockFS,
	}

	tests := []struct {
		name     string
		provider entity.Provider
		path     string
		mockFS   func()
		want     fs.FileInfo
		wantErr  bool
	}{
		{
			name: "成功获取文件信息",
			mockFS: func() {
				mockFS.EXPECT().StatFile(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, nil)
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFS()

			got, err := service.StatFile(context.Background(), tt.provider, tt.path)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
