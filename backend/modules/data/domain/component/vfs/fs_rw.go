// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package vfs

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"io/fs"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/parquet-go/parquet-go"
	"github.com/pkg/errors"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
)

type FileReader struct {
	name   string
	info   fs.FileInfo
	format entity.FileFormat

	csv     *csv.Reader                            // 仅当 format == CSV 时有效
	fields  []string                               // 仅当 format == CSV 时有效
	parquet *parquet.GenericReader[map[string]any] // 仅当 format == Parquet 时有效
	scanner *bufio.Scanner                         // 仅当 format == JSONL 时有效

	cursor int64
	closer io.Closer
}

func (r *FileReader) SetCursor(cursor int64) {
	r.cursor = cursor
}

func (r *FileReader) SetName(name string) {
	r.name = name
}

func (r *FileReader) GetName() string {
	return r.name
}

func (r *FileReader) GetCursor() int64 {
	return r.cursor
}

func NewFileReader(name string, r Reader, info fs.FileInfo, format entity.FileFormat) (*FileReader, error) {
	fr := &FileReader{name: name, info: info, closer: r, format: format}
	switch format {
	case entity.FileFormat_CSV:
		cr, err := NewCSVReader(r)
		if err != nil {
			return nil, err
		}
		fr.csv = cr

		head, err := fr.csv.Read() // 第一行视作表头
		if err != nil {
			return nil, err
		}
		fr.fields = head
		fr.cursor = 1

	case entity.FileFormat_JSONL:
		fr.scanner = bufio.NewScanner(r)

	case entity.FileFormat_Parquet:
		var err error
		fr.parquet, err = NewReader(r, info)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.Errorf("unknown file format: %s", format)
	}
	return fr, nil
}

func (r *FileReader) SeekToOffset(offset int64) error {
	switch r.format {
	case entity.FileFormat_CSV:
		return r.seekCSV(offset)
	case entity.FileFormat_JSONL:
		return r.seekJSONL(offset)
	case entity.FileFormat_Parquet:
		return r.parquet.SeekToRow(offset)
	default:
		return errors.Errorf("unknown file format: %s", r.format)
	}
}

func (r *FileReader) seekCSV(offset int64) error {
	for {
		if _, err := r.csv.Read(); err != nil {
			return err
		}
		r.cursor += 1
		if r.cursor == offset {
			return nil
		}
	}
}

func (r *FileReader) seekJSONL(line int64) error {
	for r.scanner.Scan() {
		r.cursor += 1
		if r.cursor == line {
			break
		}
	}
	if r.cursor != line {
		return errors.Errorf("seek to line %d out of range", r.cursor)
	}
	return r.scanner.Err()
}

func (r *FileReader) Next() (map[string]any, error) {
	switch r.format {
	case entity.FileFormat_CSV:
		return r.nextInCSV()
	case entity.FileFormat_Parquet:
		return r.nextInParquet()
	case entity.FileFormat_JSONL:
		return r.nextInJSONL()
	default:
		return nil, errors.Errorf("unknown file format: %s", r.format)
	}
}

func (r *FileReader) nextInCSV() (map[string]any, error) {
	record, err := r.csv.Read()
	if err != nil {
		return nil, err
	}
	if len(record) != len(r.fields) {
		return nil, errors.Errorf("record length mismatch, expected=%d, got=%d", len(r.fields), len(record))
	}

	r.cursor += 1
	m := make(map[string]any)
	for i, field := range r.fields {
		m[field] = record[i]
	}
	return m, nil
}

func (r *FileReader) nextInParquet() (map[string]any, error) {
	kvs := make([]map[string]any, 1)
	_, err := r.parquet.Read(kvs)
	if err != nil {
		return nil, err
	}
	r.cursor += 1
	return kvs[0], nil
}

func (r *FileReader) nextInJSONL() (map[string]any, error) {
	for r.scanner.Scan() {
		line := r.scanner.Bytes()
		r.cursor += 1

		line = bytes.TrimSpace(line) // skip empty line
		if len(line) > 0 {
			m := make(map[string]any)
			dc := decoder.NewStreamDecoder(bytes.NewBuffer(line))
			dc.UseInt64()
			if err := sonic.Unmarshal(line, &m); err != nil {
				return nil, errors.WithMessagef(err, "line_num=%d", r.cursor)
			}
			return m, nil
		}
	}

	if err := r.scanner.Err(); err != nil {
		return nil, err
	}
	return nil, io.EOF
}

func (r *FileReader) close() error {
	if r.closer != nil {
		return r.closer.Close()
	}
	return nil
}
