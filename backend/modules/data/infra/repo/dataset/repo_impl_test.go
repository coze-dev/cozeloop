// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	idgenmocks "github.com/coze-dev/cozeloop/backend/infra/idgen/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/redis"
)

func TestNewDatasetRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 定义测试用例
	tests := []struct {
		name            string
		idgen           idgen.IIDGenerator
		db              db.Provider
		datasetDAO      mysql.IDatasetDAO
		schemaDAO       mysql.ISchemaDAO
		datasetRedisDAO redis.DatasetDAO
		versionDAO      mysql.IVersionDAO
		versionRedisDAO redis.VersionDAO
		optDAO          redis.OperationDAO
		itemDAO         mysql.IItemDAO
		itemSnapshotDAO mysql.IItemSnapshotDAO
		ioJobDAO        mysql.IIOJobDAO
		expectErr       bool
	}{
		{
			name:            "正常场景",
			idgen:           idgenmocks.NewMockIIDGenerator(ctrl),
			db:              nil,
			datasetDAO:      mocks.NewMockIDatasetDAO(ctrl),
			schemaDAO:       mocks.NewMockISchemaDAO(ctrl),
			datasetRedisDAO: nil,
			versionDAO:      mocks.NewMockIVersionDAO(ctrl),
			versionRedisDAO: nil,
			optDAO:          nil,
			itemDAO:         mocks.NewMockIItemDAO(ctrl),
			itemSnapshotDAO: nil,
			ioJobDAO:        mocks.NewMockIIOJobDAO(ctrl),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectErr {
						t.Errorf("NewDatasetRepo 意外 panic: %v", r)
					}
				} else if tt.expectErr {
					t.Error("NewDatasetRepo 应该 panic 但没有")
				}
			}()

			repo := NewDatasetRepo(
				tt.idgen,
				tt.db,
				tt.datasetDAO,
				tt.schemaDAO,
				tt.datasetRedisDAO,
				tt.versionDAO,
				tt.versionRedisDAO,
				tt.optDAO,
				tt.itemDAO,
				tt.itemSnapshotDAO,
				tt.ioJobDAO,
				nil,
			)

			if !tt.expectErr && repo == nil {
				t.Error("NewDatasetRepo 返回了 nil")
			}
		})
	}
}
