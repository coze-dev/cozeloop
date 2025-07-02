// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package vfs

import (
	"io/fs"

	"github.com/parquet-go/parquet-go"
	"github.com/pkg/errors"
)

func NewReader(r Reader, info fs.FileInfo) (pr *parquet.GenericReader[map[string]any], err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("new parquet reader panic %v", r)
		}
	}()
	rr := &pReader{Reader: r, info: info}
	pf := parquet.NewGenericReader[map[string]any](rr)
	return pf, nil
}

type pReader struct {
	Reader
	info fs.FileInfo
}

func (r *pReader) Size() int64 {
	return r.info.Size()
}
