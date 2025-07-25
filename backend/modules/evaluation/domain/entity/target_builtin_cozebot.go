// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

type CozeBot struct {
	BotID          int64
	BotVersion     string
	BotInfoType    CozeBotInfoType
	PublishVersion *string

	BotName     string    `json:"-"`
	AvatarURL   string    `json:"-"`
	Description string    `json:"-"`
	BaseInfo    *BaseInfo `json:"-"`
}

type CozeBotInfoType int64

const (
	// 草稿 bot
	CozeBotInfoTypeDraftBot CozeBotInfoType = 1
	// 商店 bot
	CozeBotInfoTypeProductBot CozeBotInfoType = 2
)
