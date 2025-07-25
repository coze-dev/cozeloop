// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"github.com/bytedance/sonic"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
)

func VersionDO2PO(version *entity.DatasetVersion) (p *model.DatasetVersion, err error) {
	if version == nil {
		return nil, nil
	}
	p = &model.DatasetVersion{
		ID:             version.ID,
		AppID:          version.AppID,
		SpaceID:        version.SpaceID,
		DatasetID:      version.DatasetID,
		SchemaID:       version.SchemaID,
		Version:        version.Version,
		VersionNum:     version.VersionNum,
		Description:    version.Description,
		ItemCount:      version.ItemCount,
		SnapshotStatus: string(version.SnapshotStatus),
		UpdateVersion:  version.UpdateVersion,
		CreatedBy:      version.CreatedBy,
		CreatedAt:      version.CreatedAt,
		DisabledAt:     version.DisabledAt,
	}
	if version.DatasetBrief != nil {
		p.DatasetBrief, err = sonic.Marshal(version.DatasetBrief)
		if err != nil {
			return nil, errno.JSONErr(err, "marshal version.DatasetBrief failed, data=%v", version.DatasetBrief)
		}
	}
	if version.SnapshotProgress != nil {
		p.SnapshotProgress, err = sonic.Marshal(version.SnapshotProgress)
		if err != nil {
			return nil, errno.JSONErr(err, "marshal version.SnapshotProgress failed, data=%v", version.SnapshotProgress)
		}
	}
	return p, nil
}

func ConvertVersionPOToDO(p *model.DatasetVersion) (*entity.DatasetVersion, error) {
	if p == nil {
		return nil, nil
	}
	m := &entity.DatasetVersion{
		ID:             p.ID,
		AppID:          p.AppID,
		SpaceID:        p.SpaceID,
		DatasetID:      p.DatasetID,
		SchemaID:       p.SchemaID,
		Version:        p.Version,
		VersionNum:     p.VersionNum,
		Description:    p.Description,
		ItemCount:      p.ItemCount,
		SnapshotStatus: entity.SnapshotStatus(p.SnapshotStatus),
		UpdateVersion:  p.UpdateVersion,
		CreatedBy:      p.CreatedBy,
		CreatedAt:      p.CreatedAt,
		DisabledAt:     p.DisabledAt,
	}
	if err := sonic.Unmarshal(p.DatasetBrief, &m.DatasetBrief); err != nil {
		return nil, errno.JSONErr(err, "unmarshal version.DatasetBrief failed, data=%v", p.DatasetBrief)
	}
	if p.SnapshotProgress != nil {
		if err := sonic.Unmarshal(p.SnapshotProgress, &m.SnapshotProgress); err != nil {
			return nil, errno.JSONErr(err, "unmarshal version.SnapshotProgress failed, data=%v", p.SnapshotProgress)
		}
	}
	return m, nil
}
