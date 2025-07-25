// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"time"
)

type SpaceType int32

const (
	SpaceTypePersonal SpaceType = 1
	SpaceTypeTeam     SpaceType = 2
)

type Space struct {
	ID          int64
	Name        string
	Description string
	IconURI     string
	IconURL     string
	SpaceType   SpaceType
	OwnerID     int64
	CreatorID   int64
	CreatedAt   time.Time // 创建时间
	UpdatedAt   time.Time // 更新时间
}

type SpaceUserType int32

const (
	SpaceUserTypeOwner  SpaceUserType = 1
	SpaceUserTypeAdmin  SpaceUserType = 2
	SpaceUserTypeMember SpaceUserType = 3
)
