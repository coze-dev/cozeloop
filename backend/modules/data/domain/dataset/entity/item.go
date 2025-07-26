// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gslice"

	"github.com/coze-dev/coze-loop/backend/modules/data/domain/entity"
)

type Item struct {
	ID        int64
	AppID     int32
	SpaceID   int64
	DatasetID int64
	SchemaID  int64

	ItemID         int64
	ItemKey        string              // 幂等 key
	Data           []*FieldData        // 数据内容
	RepeatedData   []*ItemData         // 多轮数据内容，与 Data 互斥
	DataProperties *ItemDataProperties // 内容属性

	AddVN int64
	DelVN int64

	CreatedBy string
	CreatedAt time.Time
	UpdatedBy string
	UpdatedAt time.Time
}

func (i *Item) SetID(id int64) { i.ID = id }

func (i *Item) GetID() int64 { return i.ID }

func (i *Item) GetOrBuildProperties() *ItemDataProperties {
	if i.DataProperties == nil {
		i.BuildProperties()
	}
	return i.DataProperties
}

func (i *Item) BuildProperties() {
	data := i.AllData()
	i.DataProperties = &ItemDataProperties{}
	i.DataProperties.Bytes = int64(gslice.SumBy(gslice.Flatten(data), func(f *FieldData) int { return f.DataBytes() }))
	i.DataProperties.Runes = int64(gslice.SumBy(gslice.Flatten(data), func(f *FieldData) int { return f.DataRunes() }))
}

func (i *Item) AllData() [][]*FieldData {
	if len(i.RepeatedData) != 0 {
		return gslice.Map(i.RepeatedData, func(d *ItemData) []*FieldData { return d.Data })
	}
	return [][]*FieldData{i.Data}
}

func (i *Item) ClearData() {
	i.Data = nil
	i.RepeatedData = nil
}

type ItemIdentity struct {
	SpaceID   int64
	DatasetID int64
	ID        int64
	ItemID    int64
	AddVN     int64
}

type FieldData struct {
	Key         string             `json:"key,omitempty"`
	Name        string             `json:"-"` // DB 中不存储
	ContentType ContentType        `json:"content_type,omitempty"`
	Format      FieldDisplayFormat `json:"format,omitempty"`
	Content     string             `json:"content,omitempty"`
	Attachments []*ObjectStorage   `json:"attachments,omitempty"`
	Parts       []*FieldData       `json:"parts,omitempty"`
}

type ItemData struct {
	ID   int64        `json:"id,omitempty"`
	Data []*FieldData `json:"data,omitempty"`
}

type ItemDataProperties struct {
	Storage        entity.Provider `json:"storage,omitempty"`         // 外部存储，为空表示存在 RDS 中
	StorageKey     string          `json:"storage_key,omitempty"`     // 外部存储 key
	CompressFormat string          `json:"compress_format,omitempty"` // 压缩格式, 为空表示未压缩
	Bytes          int64           `json:"bytes,omitempty"`           // 字节数
	Runes          int64           `json:"characters,omitempty"`      // 字符数
}

type ObjectStorage struct {
	Provider entity.Provider `json:"provider,omitempty"`
	Name     string          `json:"name,omitempty"`
	URI      string          `json:"uri,omitempty"`
	URL      string          `json:"-"`
	ThumbURL string          `json:"-"`
}

func (f *FieldData) DataBytes() int {
	var n int
	n += len(f.Content)
	for _, att := range f.Attachments {
		n += len(att.Name)
		n += len(att.URI)
	}
	for _, part := range f.Parts {
		n += part.DataBytes()
	}
	return n
}

func (f *FieldData) DataRunes() int {
	var n int
	n += utf8.RuneCountInString(f.Content)
	for _, att := range f.Attachments {
		n += utf8.RuneCountInString(att.Name)
		n += utf8.RuneCountInString(att.URI)
	}
	for _, part := range f.Parts {
		n += part.DataRunes()
	}
	return n
}

type ItemErrorGroup struct {
	Type    *ItemErrorType
	Summary *string
	// 错误条数
	ErrorCount *int32
	// 批量写入时，每类错误至多提供 5 个错误详情；导入任务，至多提供 10 个错误详情
	Details []*ItemErrorDetail
}

type ItemErrorDetail struct {
	Message *string
	// 单条错误数据在输入数据中的索引。从 0 开始，下同
	Index *int32
	// [startIndex, endIndex] 表示区间错误范围, 如 ExceedDatasetCapacity 错误时
	StartIndex *int32
	EndIndex   *int32
}

type ItemErrorType int64

const (
	// schema 不匹配
	ItemErrorType_MismatchSchema ItemErrorType = 1
	// 空数据
	ItemErrorType_EmptyData ItemErrorType = 2
	// 单条数据大小超限
	ItemErrorType_ExceedMaxItemSize ItemErrorType = 3
	// 数据集容量超限
	ItemErrorType_ExceedDatasetCapacity ItemErrorType = 4
	// 文件格式错误
	ItemErrorType_MalformedFile ItemErrorType = 5
	// 包含非法内容
	ItemErrorType_IllegalContent ItemErrorType = 6
	/* system error*/
	ItemErrorType_InternalError ItemErrorType = 100
)

func ItemErrorDetailToString(d *ItemErrorDetail) string {
	if gptr.Indirect(d.Index) > 0 || gptr.Indirect(d.EndIndex) > 0 {
		return fmt.Sprintf(`%s, range=%d-%d`, gptr.Indirect(d.Message), gptr.Indirect(d.StartIndex), gptr.Indirect(d.EndIndex))
	}
	return fmt.Sprintf(`%s, index=%d`, gptr.Indirect(d.Message), gptr.Indirect(d.Index))
}

func SanitizeItemErrorGroup(eg *ItemErrorGroup, maxDetailCnt int) {
	// count errors
	var count int32
	for _, d := range eg.Details {
		if n := gptr.Indirect(d.EndIndex) - gptr.Indirect(d.StartIndex) + 1; n > 1 {
			count += n
		} else {
			count++
		}
	}
	eg.ErrorCount = gptr.Of(count)

	if len(eg.Details) <= maxDetailCnt {
		return
	}

	// summarize details
	indices := gslice.Map(eg.Details, func(d *ItemErrorDetail) string {
		if gptr.Indirect(d.StartIndex) > 0 || gptr.Indirect(d.EndIndex) > 0 {
			return fmt.Sprintf(`%d-%d`, gptr.Indirect(d.StartIndex), gptr.Indirect(d.EndIndex))
		}
		return fmt.Sprintf(`%d`, gptr.Indirect(d.Index))
	})
	gslice.Sort(indices)
	eg.Summary = gptr.Of(fmt.Sprintf(`%d errors happened, indices=%v`, eg.ErrorCount, indices))

	// prune details
	details := gslice.Filter(eg.Details, func(d *ItemErrorDetail) bool { return gptr.Indirect(d.Message) != "" })
	if len(details) > maxDetailCnt {
		details = details[:maxDetailCnt]
	}
	eg.Details = details
}
