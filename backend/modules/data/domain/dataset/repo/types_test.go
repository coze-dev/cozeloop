// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/coze-loop/backend/modules/data/pkg/consts"
)

func TestNewListItemsParamsFromVersion(t *testing.T) {
	version := &entity.DatasetVersion{
		SpaceID:    1,
		DatasetID:  2,
		VersionNum: 3,
	}

	params := NewListItemsParamsFromVersion(version)

	assert.Equal(t, version.SpaceID, params.SpaceID)
	assert.Equal(t, version.DatasetID, params.DatasetID)
	assert.Equal(t, version.VersionNum, params.DelVNGt)
	assert.Equal(t, version.VersionNum, params.AddVNLte)
}

func TestNewListItemsParamsOfDataset(t *testing.T) {
	spaceID := int64(1)
	datasetID := int64(2)

	params := NewListItemsParamsOfDataset(spaceID, datasetID)

	assert.Equal(t, spaceID, params.SpaceID)
	assert.Equal(t, datasetID, params.DatasetID)
	// 假设 consts.MaxVersionNum 是一个已知的常量，这里简单用一个大值代替
	assert.Equal(t, consts.MaxVersionNum, params.DelVNEq)
}

func TestWithDeleted(t *testing.T) {
	opt := &Opt{}
	WithDeleted()(opt)
	assert.True(t, opt.WithDeleted)
}

func TestWithMaster(t *testing.T) {
	opt := &Opt{}
	WithMaster()(opt)
	assert.True(t, opt.WithMaster)
}
