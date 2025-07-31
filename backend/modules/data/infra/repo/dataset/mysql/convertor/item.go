// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"github.com/bytedance/sonic"
	"github.com/pkg/errors"

	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/coze-loop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
)

func ItemDO2PO(s *entity.Item) (*model.DatasetItem, error) {
	if s == nil {
		return nil, nil
	}

	t := &model.DatasetItem{
		ID:        s.ID,
		AppID:     s.AppID,
		SpaceID:   s.SpaceID,
		DatasetID: s.DatasetID,
		SchemaID:  s.SchemaID,
		ItemID:    s.ItemID,
		ItemKey:   s.ItemKey,
		AddVn:     s.AddVN,
		DelVn:     s.DelVN,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy,
		UpdatedAt: s.UpdatedAt,
	}
	if len(s.Data) != 0 { // save empty slice
		data, err := sonic.Marshal(s.Data)
		if err != nil {
			return nil, errors.WithMessage(err, "marshal data")
		}
		t.Data = data
	}
	if len(s.RepeatedData) != 0 { // save empty slice
		data, err := sonic.Marshal(s.RepeatedData)
		if err != nil {
			return nil, errors.WithMessage(err, "marshal repeated data")
		}
		t.RepeatedData = data
	}
	if s.DataProperties != nil {
		data, err := sonic.Marshal(s.DataProperties)
		if err != nil {
			return nil, errors.WithMessage(err, "marshal data properties")
		}
		t.DataProperties = data
	}
	return t, nil
}

func ItemPO2DO(s *model.DatasetItem) (*entity.Item, error) {
	if s == nil {
		return nil, nil
	}
	t := &entity.Item{
		ID:        s.ID,
		AppID:     s.AppID,
		SpaceID:   s.SpaceID,
		DatasetID: s.DatasetID,
		SchemaID:  s.SchemaID,
		ItemID:    s.ItemID,
		ItemKey:   s.ItemKey,
		AddVN:     s.AddVn,
		DelVN:     s.DelVn,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy,
		UpdatedAt: s.UpdatedAt,
	}
	if len(s.Data) > 0 {
		if err := sonic.Unmarshal(s.Data, &t.Data); err != nil {
			return nil, errors.WithMessage(err, "unmarshal data")
		}
	}
	if len(s.RepeatedData) > 0 {
		if err := sonic.Unmarshal(s.RepeatedData, &t.RepeatedData); err != nil {
			return nil, errors.WithMessage(err, "unmarshal repeated data")
		}
	}
	if len(s.DataProperties) > 0 {
		if err := sonic.Unmarshal(s.DataProperties, &t.DataProperties); err != nil {
			return nil, errors.WithMessage(err, "unmarshal data properties")
		}
	}
	return t, nil
}
