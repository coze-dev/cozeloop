// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"testing"

	"github.com/bytedance/gg/gptr"
	"github.com/stretchr/testify/assert"

	kitexDataset "github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/data/domain/dataset"
	jobPKG "github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/data/domain/dataset_job"
	datasetEntity "github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	domainEntity "github.com/coze-dev/cozeloop/backend/modules/data/domain/entity"
)

// TestDatasetIOEndpointDO2DTO 用于测试 DatasetIOEndpointDO2DTO 函数的转换逻辑
func TestDatasetIOEndpointDO2DTO(t *testing.T) {
	var jobFormatZero jobPKG.FileFormat = 0 // FileFormat 的零值

	// 定义测试用例
	testCases := []struct {
		name     string
		input    *datasetEntity.DatasetIOEndpoint // 输入的 DO
		expected *jobPKG.DatasetIOEndpoint        // 期望输出的 DTO
	}{
		{
			name:     "输入为nil",
			input:    nil,
			expected: nil,
		},
		{
			name:  "空的DatasetIOEndpoint, File和Dataset都为nil",
			input: &datasetEntity.DatasetIOEndpoint{}, // File 和 Dataset 字段默认为 nil
			expected: &jobPKG.DatasetIOEndpoint{ // 对应的 DTO File 和 Dataset 也应为 nil
				File:    nil,
				Dataset: nil,
			},
		},
		{
			name: "只有Dataset部分, File为nil",
			input: &datasetEntity.DatasetIOEndpoint{
				File: nil,
				Dataset: &datasetEntity.DatasetIODataset{
					SpaceID:   gptr.Of(int64(123)),
					DatasetID: 456,
					VersionID: gptr.Of(int64(789)),
				},
			},
			expected: &jobPKG.DatasetIOEndpoint{
				File: nil,
				Dataset: &jobPKG.DatasetIODataset{
					SpaceID:   gptr.Of(int64(123)),
					DatasetID: 456,
					VersionID: gptr.Of(int64(789)),
				},
			},
		},
		{
			name: "File部分Provider未知, Format和CompressFormat为nil",
			input: &datasetEntity.DatasetIOEndpoint{
				File: &datasetEntity.DatasetIOFile{
					Provider:       domainEntity.Provider("unknown_provider"), // 一个无法识别的 provider
					Path:           "/path/to/unknown",
					Format:         nil, // Format 为 nil
					CompressFormat: nil, // CompressFormat 也为 nil
				},
				Dataset: nil,
			},
			expected: &jobPKG.DatasetIOEndpoint{
				File: &jobPKG.DatasetIOFile{
					// StorageProviderFromString 对于未知字符串通常返回枚举的零值
					Provider:       kitexDataset.StorageProvider(0),
					Path:           "/path/to/unknown",
					Format:         gptr.Of(jobFormatZero), // nil Format 转换为指向零值的指针
					CompressFormat: gptr.Of(jobFormatZero), // nil CompressFormat 转换为指向零值的指针
					Files:          nil,                    // 假设 Files 为 nil 如果未提供
				},
				Dataset: nil,
			},
		},
	}

	// 遍历所有测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 调用被测函数
			actual := DatasetIOEndpointDO2DTO(tc.input)

			// 使用 testify/assert 进行断言，比较期望结果和实际结果
			assert.Equal(t, tc.expected, actual)
		})
	}
}
