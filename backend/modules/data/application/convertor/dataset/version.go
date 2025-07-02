// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/sonic"

	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/data/domain/dataset"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
)

func VersionDO2DTO(v *entity.DatasetVersion) (*dataset.DatasetVersion, error) {
	dto := &dataset.DatasetVersion{
		ID:             v.ID,
		AppID:          gptr.Of(v.AppID),
		SpaceID:        v.SpaceID,
		DatasetID:      v.DatasetID,
		SchemaID:       v.SchemaID,
		Version:        gptr.Of(v.Version),
		VersionNum:     gptr.Of(v.VersionNum),
		Description:    v.Description,
		ItemCount:      gptr.Of(v.ItemCount),
		SnapshotStatus: gptr.Of(SnapshotStatusDO2DTO(v.SnapshotStatus)),
		CreatedBy:      gptr.Of(v.CreatedBy),
		CreatedAt:      gptr.Of(v.CreatedAt.UnixMilli()),
	}
	if v.DisabledAt != nil {
		dto.DisabledAt = gptr.Of(v.DisabledAt.UnixMilli())
	}
	if v.DatasetBrief != nil {
		dsBriefStr, err := sonic.MarshalString(v.DatasetBrief)
		if err != nil {
			return nil, errno.JSONErr(err, "marshal dataset_version.dataset_brief failed, data=%v", v.DatasetBrief)
		}
		dto.DatasetBrief = gptr.Of(dsBriefStr)
	}
	return dto, nil
}

func SnapshotStatusDO2DTO(ss entity.SnapshotStatus) dataset.SnapshotStatus {
	switch ss {
	case entity.SnapshotStatusUnstarted:
		return dataset.SnapshotStatus_Unstarted
	case entity.SnapshotStatusInProgress:
		return dataset.SnapshotStatus_InProgress
	case entity.SnapshotStatusCompleted:
		return dataset.SnapshotStatus_Completed
	case entity.SnapshotStatusFailed:
		return dataset.SnapshotStatus_Failed
	default:
		return dataset.DatasetVersion_SnapshotStatus_DEFAULT
	}
}
