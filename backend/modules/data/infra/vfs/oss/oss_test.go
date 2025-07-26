// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package oss

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-loop/backend/infra/fileserver"
	mock_fileserver "github.com/coze-dev/coze-loop/backend/infra/fileserver/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/data/infra/vfs"
)

func TestClient_Stat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOS := mock_fileserver.NewMockObjectStorage(ctrl)
	client := NewClient(mockOS)

	// 测试用例定义
	tests := []struct {
		name         string
		path         string
		mockSetup    func()
		expectedInfo *vfs.FSInformation
		expectedErr  error
	}{{
		name: "正常场景: 获取文件信息成功",
		path: "valid/path/file.txt",
		mockSetup: func() {
			mockInfo := &fileserver.ObjectInfo{}
			mockOS.EXPECT().Stat(gomock.Any(), "valid/path/file.txt").Return(mockInfo, nil)
		},
		expectedInfo: &vfs.FSInformation{},
		expectedErr:  nil,
	}, {
		name: "正常场景: 获取目录信息成功",
		path: "valid/directory/",
		mockSetup: func() {
			mockInfo := &fileserver.ObjectInfo{}
			mockOS.EXPECT().Stat(gomock.Any(), "valid/directory/").Return(mockInfo, nil)
		},
		expectedInfo: &vfs.FSInformation{},
		expectedErr:  nil,
	}, {
		name: "边界场景: 空路径",
		path: "",
		mockSetup: func() {
			mockOS.EXPECT().Stat(gomock.Any(), "").Return(nil, os.ErrNotExist)
		},
		expectedInfo: nil,
		expectedErr:  os.ErrNotExist,
	}, {
		name: "异常场景: 文件不存在",
		path: "invalid/path/nonexistent.txt",
		mockSetup: func() {
			mockOS.EXPECT().Stat(gomock.Any(), "invalid/path/nonexistent.txt").Return(nil, os.ErrNotExist)
		},
		expectedInfo: nil,
		expectedErr:  os.ErrNotExist,
	}, {
		name: "异常场景: 权限不足",
		path: "restricted/path/file.txt",
		mockSetup: func() {
			mockOS.EXPECT().Stat(gomock.Any(), "restricted/path/file.txt").Return(nil, os.ErrPermission)
		},
		expectedInfo: nil,
		expectedErr:  os.ErrPermission,
	}}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			info, err := client.Stat(context.Background(), tt.path)

			// 验证错误
			assert.ErrorIs(t, err, tt.expectedErr)

			// 验证返回信息
			if tt.expectedErr == nil {
				assert.Equal(t, tt.expectedInfo.FName, info.(*vfs.FSInformation).FName)
				assert.Equal(t, tt.expectedInfo.FSize, info.(*vfs.FSInformation).FSize)
				assert.Equal(t, tt.expectedInfo.FModTime, info.(*vfs.FSInformation).FModTime)
				assert.Equal(t, tt.expectedInfo.FIsDir, info.(*vfs.FSInformation).FIsDir)
			} else {
				assert.Nil(t, info)
			}
		})
	}
}
