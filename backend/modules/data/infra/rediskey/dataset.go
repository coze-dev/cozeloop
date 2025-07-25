// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package rediskey

import "fmt"

const (
	DatasetItemCountKey        = `dataset:%d:item_count`         // dataset:dataset_id:item_count, int, 数据集item 数量
	DatasetVersionItemCountKey = `dataset_version:%d:item_count` // dataset_version:dataset_version_id:item_count, int, 数据集版本 item 数量
	DatasetOperationKey        = `dataset:%d:op:%s`              // dataset:dataset_id:op:operation_type, set[op_id, op_entity], 数据集操作
)

func FormatDatasetItemCountKey(datasetID int64) string {
	return fmt.Sprintf(DatasetItemCountKey, datasetID)
}

func FormatDatasetVersionItemCountKey(versionID int64) string {
	return fmt.Sprintf(DatasetVersionItemCountKey, versionID)
}

func FormatDatasetOperationKey(datasetID int64, opType string) string {
	return fmt.Sprintf(DatasetOperationKey, datasetID, opType)
}
