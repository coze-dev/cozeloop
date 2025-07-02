// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package vfs

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"unicode/utf8"

	"github.com/dimchansky/utfbom"
	"github.com/pkg/errors"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func NewCSVReader(r io.Reader) (*csv.Reader, error) {
	br, err := rebuildReader(r)
	if err != nil {
		return nil, err
	}
	cr := csv.NewReader(br)
	cr.LazyQuotes = true // allow unquoted quotes
	return cr, nil
}

func rebuildReader(r io.Reader) (io.Reader, error) {
	r = utfbom.SkipOnly(r) // 跳过 BOM, issue: https://github.com/golang/go/issues/33887
	br := bufio.NewReader(r)
	r = br // reset r
	// peek first 1024 bytes to detect charset
	head, err := br.Peek(1024)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, bufio.ErrBufferFull) {
		return nil, err
	}
	if len(head) == 0 || isUTF8(head) {
		return r, nil
	}
	// try gb18030, commonly used in WPS and Excel in China
	hr := transform.NewReader(bytes.NewReader(head), simplifiedchinese.GB18030.NewDecoder())
	if head, err := io.ReadAll(hr); err == nil && isUTF8(head) {
		return transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder()), nil
	}
	// detect charset
	charRes, err := chardet.NewTextDetector().DetectBest(head)
	if err != nil || charRes.Confidence < 100 { // ignore error
		return br, nil
	}
	switch charRes.Charset {
	case "GB-18030":
		r = transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder())
	case "GBK":
		r = transform.NewReader(r, simplifiedchinese.GBK.NewDecoder())
	}
	return r, nil
}

func isUTF8(head []byte) bool {
	// avoid corner case: cut utf-8 in the middle
	n := len(head)
	for i := 0; i < n && i < utf8.UTFMax; i++ {
		if utf8.Valid(head[:n-i]) {
			return true
		}
	}
	return false
}
