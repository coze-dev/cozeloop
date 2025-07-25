// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mysql

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/bytedance/gg/gslice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/infra/redis"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/convertor"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/consts"
)

func TestItemDAOImpl_ClearDataset(t *testing.T) {
	r := newTestItemDAO(t)
	ctx := context.TODO()
	items := make([]*model.DatasetItem, 0, 10)
	spaceID := int64(100)
	datasetID := int64(101)
	currentVN := int64(3)

	var wantedItems []*entity.Item         // 预期被删除的 item
	var wantedItemsToUpdate []*entity.Item // 预期被更新的 item
	for i := 0; i < 2000; i++ {
		item := newTestItem(t, func(item *entity.Item) {
			item.ItemKey = fmt.Sprintf("key-%d", i)
			item.SpaceID = spaceID
			item.DatasetID = datasetID
			item.AddVN = 1
			switch i % 3 {
			case 0:
				item.DelVN = 2
			case 1:
				wantedItems = append(wantedItems, item)
				wantedItemsToUpdate = append(wantedItemsToUpdate, item)
			case 2:
				item.AddVN = currentVN
				wantedItems = append(wantedItems, item)
			}
		})
		items = append(items, item)
	}
	_, err := r.MCreateItems(ctx, items)
	require.NoError(t, err)

	// clear
	{
		got, err := r.ClearDataset(ctx, spaceID, datasetID, currentVN)
		require.NoError(t, err)
		assert.Len(t, got, len(wantedItems))

		gotIDs := gslice.Map(got, func(i *entity.ItemIdentity) int64 { return i.ID })
		assert.ElementsMatch(t, gslice.Map(wantedItems, func(i *entity.Item) int64 { return i.ID }), gotIDs)
	}

	// check after clear
	{
		got, _, err := r.ListItems(ctx, &ListItemsParams{
			SpaceID:   spaceID,
			DatasetID: datasetID,
			ItemIDs:   gslice.Map(wantedItems, func(i *entity.Item) int64 { return i.ID }),
		})
		require.NoError(t, err)
		assert.Len(t, got, len(wantedItemsToUpdate))
		assert.ElementsMatch(t, gslice.Map(wantedItemsToUpdate, func(i *entity.Item) int64 { return i.ID }), gslice.Map(got, func(i *model.DatasetItem) int64 { return i.ID }))
	}
}

func newTestItemDAO(t *testing.T) IItemDAO {
	testDB := db.NewTestDB(t, &model.DatasetItem{})
	testRedis := redis.NewTestRedis(t)
	return NewDatasetItemDAO(testDB, testRedis)
}

func newTestItem(t *testing.T, taps ...func(item *entity.Item)) *model.DatasetItem {
	now := time.Now().Truncate(time.Second)
	id := int64(rand.Int())
	item := &entity.Item{
		ID:        id,
		AppID:     1,
		SpaceID:   2,
		DatasetID: 3,
		SchemaID:  4,
		ItemID:    id,
		ItemKey:   fmt.Sprintf("%d", id),
		Data: []*entity.FieldData{
			{Key: "name", Content: "someone"},
			{Key: "age", Content: "19"},
		},
		DataProperties: &entity.ItemDataProperties{Bytes: 21, Runes: 21},
		AddVN:          1,
		DelVN:          consts.MaxVersionNum,
		CreatedBy:      "someone",
		CreatedAt:      now,
		UpdatedBy:      "someone",
		UpdatedAt:      now,
	}
	for _, tap := range taps {
		tap(item)
	}

	po, err := convertor.ItemDO2PO(item)
	require.NoError(t, err)
	return po
}
