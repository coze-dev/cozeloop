// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"time"
)

type User struct {
	SpaceID      int64  // 空间ID
	UserID       int64  // 用户ID
	UniqueName   string // 唯一名称
	NickName     string // 昵称
	Email        string // 邮箱
	HashPassword string // 密码哈希
	Description  string // 用户描述
	IconURI      string // 头像URI
	IconURL      string // 头像URL
	UserVerified bool   // 用户是否已验证
	CountryCode  int64  // 国家代码
	SessionKey   string // 会话密钥

	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}
