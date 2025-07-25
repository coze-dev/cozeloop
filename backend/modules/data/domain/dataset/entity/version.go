// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import "time"

type DatasetVersion struct {
	ID        int64
	AppID     int32
	SpaceID   int64
	DatasetID int64
	SchemaID  int64

	DatasetBrief     *Dataset          // 数据集元信息备份
	Version          string            // 版本号，SemVer2 三段式
	VersionNum       int64             // 数字版本号，从1开始递增
	Description      *string           // 版本描述
	ItemCount        int64             // 条数
	SnapshotStatus   SnapshotStatus    // 快照状态
	SnapshotProgress *SnapshotProgress // 快照进度

	UpdateVersion int64 // 更新版本号，用于乐观锁
	CreatedBy     string
	CreatedAt     time.Time
	DisabledAt    *time.Time // 版本禁用
}

type SnapshotStatus string

const (
	SnapshotStatusUnknown    SnapshotStatus = ""
	SnapshotStatusUnstarted  SnapshotStatus = "unstarted"
	SnapshotStatusInProgress SnapshotStatus = "in_progress"
	SnapshotStatusCompleted  SnapshotStatus = "completed"
	SnapshotStatusFailed     SnapshotStatus = "failed"
)

func (ss SnapshotStatus) IsFinished() bool {
	switch ss {
	case SnapshotStatusUnstarted, SnapshotStatusInProgress, SnapshotStatusUnknown:
		return false
	default:
		return true
	}
}

func (s *DatasetVersion) GetID() int64 {
	return s.ID
}

func (s *DatasetVersion) SetID(id int64) {
	s.ID = id
}
