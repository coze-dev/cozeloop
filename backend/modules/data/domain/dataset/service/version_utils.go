// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"github.com/Masterminds/semver/v3"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
)

func validateVersion(preVersion, newVersion string) error {
	newV, err := semver.StrictNewVersion(newVersion)
	if err != nil {
		return errno.InvalidParamErr(err, "version '%s' not a valid semantic version", newVersion)
	}

	if preVersion == "" { // 无历史版本，直接返回
		return nil
	}

	preV, err := semver.StrictNewVersion(preVersion)
	if err != nil {
		return errno.InternalErr(err, "previous version '%s' not a valid semantic version", preVersion)
	}
	if newV.LessThanEqual(preV) {
		return errno.InvalidParamErrorf("new version '%s' should be greater than '%s'", newVersion, preVersion)
	}
	return nil
}

func patchVersionWithDataset(d *entity.Dataset, v *entity.DatasetVersion) *entity.DatasetVersion {
	v.AppID = d.AppID
	v.SpaceID = d.SpaceID
	v.DatasetID = d.ID
	v.SchemaID = d.SchemaID
	v.VersionNum = d.NextVersionNum
	v.DatasetBrief = d
	return v
}
