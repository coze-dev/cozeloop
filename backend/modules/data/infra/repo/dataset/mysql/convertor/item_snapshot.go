// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"github.com/bytedance/sonic"
	"github.com/pkg/errors"

	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/coze-loop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
)

func ItemSnapshotDO2PO(m *entity.ItemSnapshot) (*model.ItemSnapshot, error) {
	if m == nil || m.Snapshot == nil {
		return nil, nil
	}
	s := m.Snapshot
	t := &model.ItemSnapshot{
		ID:            m.ID,
		AppID:         s.AppID,
		SpaceID:       s.SpaceID,
		DatasetID:     s.DatasetID,
		SchemaID:      s.SchemaID,
		VersionID:     m.VersionID,
		ItemPrimaryID: s.ID,
		ItemID:        s.ItemID,
		ItemKey:       s.ItemKey,
		AddVn:         s.AddVN,
		DelVn:         s.DelVN,
		CreatedAt:     m.CreatedAt,
		ItemCreatedBy: s.CreatedBy,
		ItemCreatedAt: s.CreatedAt,
		ItemUpdatedBy: s.UpdatedBy,
		ItemUpdatedAt: s.UpdatedAt,
	}
	if len(s.Data) > 0 {
		data, err := sonic.Marshal(s.Data)
		if err != nil {
			return nil, errors.WithMessage(err, "marshal data of snapshot")
		}
		t.Data = data
	}
	if len(s.RepeatedData) > 0 {
		data, err := sonic.Marshal(s.RepeatedData)
		if err != nil {
			return nil, errors.WithMessage(err, "marshal repeated data of snapshot")
		}
		t.RepeatedData = data
	}
	if s.DataProperties != nil {
		data, err := sonic.Marshal(s.DataProperties)
		if err != nil {
			return nil, errors.WithMessage(err, "marshal data properties of snapshot")
		}
		t.DataProperties = data
	}
	return t, nil
}

func ConvertItemSnapshotPO2DO(p *model.ItemSnapshot) (*entity.ItemSnapshot, error) {
	if p == nil {
		return nil, nil
	}
	m := &entity.ItemSnapshot{
		ID:        p.ID,
		VersionID: p.VersionID,
		CreatedAt: p.CreatedAt,
		Snapshot: &entity.Item{
			ID:        p.ItemPrimaryID,
			AppID:     p.AppID,
			SpaceID:   p.SpaceID,
			DatasetID: p.DatasetID,
			SchemaID:  p.SchemaID,
			ItemID:    p.ItemID,
			ItemKey:   p.ItemKey,
			AddVN:     p.AddVn,
			DelVN:     p.DelVn,
			CreatedBy: p.ItemCreatedBy,
			CreatedAt: p.ItemCreatedAt,
			UpdatedBy: p.ItemUpdatedBy,
			UpdatedAt: p.ItemUpdatedAt,
		},
	}
	if len(p.Data) > 0 {
		if err := sonic.Unmarshal(p.Data, &m.Snapshot.Data); err != nil {
			return nil, errors.WithMessage(err, "unmarshal data of snapshot")
		}
	}
	if len(p.RepeatedData) > 0 {
		if err := sonic.Unmarshal(p.RepeatedData, &m.Snapshot.RepeatedData); err != nil {
			return nil, errors.WithMessage(err, "unmarshal repeated data of snapshot")
		}
	}
	if len(p.DataProperties) > 0 {
		if err := sonic.Unmarshal(p.DataProperties, &m.Snapshot.DataProperties); err != nil {
			return nil, errors.WithMessage(err, "unmarshal data properties of snapshot")
		}
	}
	return m, nil
}
