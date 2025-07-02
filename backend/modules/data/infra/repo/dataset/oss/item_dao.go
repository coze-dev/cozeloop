// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package oss

import (
	"bytes"
	"context"
	"io"

	"github.com/bytedance/gg/gslice"
	"github.com/bytedance/sonic"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/coze-dev/cozeloop/backend/infra/fileserver"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	entity2 "github.com/coze-dev/cozeloop/backend/modules/data/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/item_dao"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/oss/model"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
)

type ItemDAOImpl struct {
	batchObjectStorage fileserver.BatchObjectStorage
}

func NewDatasetItemDAO(batchObjectStorage fileserver.BatchObjectStorage) item_dao.ItemDAO {
	return &ItemDAOImpl{
		batchObjectStorage: batchObjectStorage,
	}
}

func (r *ItemDAOImpl) MSetItemData(ctx context.Context, items []*entity.Item) (int, error) {
	items = r.filterAbaseItems(items)
	if len(items) == 0 {
		return 0, nil
	}
	var (
		merr  = &multierror.Error{}
		count int
	)
	keys := make([]string, 0)
	readers := make([]io.Reader, 0)
	for _, item := range items {
		value, err := sonic.Marshal(&model.ItemDataPO{Data: item.Data, RepeatedData: item.RepeatedData})
		if err != nil {
			merr = multierror.Append(merr, errors.Wrapf(err, "marshal item data, id=%d", item.ID))
			continue
		}
		key := item.GetOrBuildProperties().StorageKey
		keys = append(keys, key)
		readers = append(readers, bytes.NewReader(value))
		count++
	}
	err := r.batchObjectStorage.BatchUpload(ctx, keys, readers)
	if err != nil {
		return 0, err
	}
	return count, merr.ErrorOrNil()
}

func (r *ItemDAOImpl) MGetItemData(ctx context.Context, items []*entity.Item) error {
	items = r.filterAbaseItems(items)
	if len(items) == 0 {
		return nil
	}
	merr := &multierror.Error{}
	keys := gslice.Map(items, func(item *entity.Item) string { return item.GetOrBuildProperties().StorageKey })
	values, err := r.batchObjectStorage.BatchRead(ctx, keys)
	if err != nil {
		return err
	}
	if len(values) != len(items) {
		return errno.Errorf(errno.CommonInternalErrorCode, "get %d item data from os, got %d", len(values), len(items))
	}
	for i, value := range values {
		item := items[i]
		if value == nil {
			continue
		}
		s, err := io.ReadAll(value)
		if err != nil {
			merr = multierror.Append(merr, errors.Wrapf(err, "expect string of item data, but got %T, id=%d", value, item.ID))
			continue
		}
		po := &model.ItemDataPO{}
		if err = sonic.Unmarshal(s, &po); err != nil {
			merr = multierror.Append(merr, errors.Wrapf(err, "unmarshal item data, id=%d", item.ID))
			continue
		}
		item.Data = po.Data
		item.RepeatedData = po.RepeatedData
	}
	return merr.ErrorOrNil()
}

func (r *ItemDAOImpl) filterAbaseItems(items []*entity.Item) []*entity.Item {
	return gslice.Filter(items, func(i *entity.Item) bool {
		return i.DataProperties != nil &&
			i.DataProperties.Storage == entity2.ProviderS3 &&
			i.DataProperties.StorageKey != ""
	})
}
