// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"github.com/samber/lo"
)

type PageParam struct {
	PageType        int
	ZipPageParam    ZipPageParam
	NormalPageParam NormalPageParam
	PageSize        int
	OrderDesc       bool
	Count           bool
}

type PageResult[T any] struct {
	Total      *int64
	NextCursor *int64
	Items      []*T
}

const (
	PageTypeNormal = 1
	PageTypeZip    = 2
)

type ZipPageParam struct {
	Cursor int64
}

type NormalPageParam struct {
	PageNo  int
	OrderBy int
}

func (pp PageParam) GetCursor() int64 {
	return pp.ZipPageParam.Cursor
}

func (pp PageParam) GetOffset() int {
	return (pp.NormalPageParam.PageNo - 1) * pp.PageSize
}

func (pp PageParam) GetLimit() int {
	return lo.Ternary(pp.PageType == PageTypeNormal, pp.PageSize, pp.PageSize+1)
}
