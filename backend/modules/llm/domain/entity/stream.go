// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"fmt"

	"github.com/cloudwego/eino/schema"

	llm_errorx "github.com/coze-dev/cozeloop/backend/modules/llm/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
)

//go:generate mockgen -destination=mocks/stream.go -package=mocks . IStreamReader
type IStreamReader interface {
	Recv() (*Message, error)
}

type StreamReader struct {
	frame      Frame
	einoReader *schema.StreamReader[*Message]
}

func NewStreamReader(frame Frame, einoReader *schema.StreamReader[*schema.Message]) IStreamReader {
	return &StreamReader{
		frame:      frame,
		einoReader: schema.StreamReaderWithConvert(einoReader, ToDOMessage),
	}
}

func (sr *StreamReader) Recv() (message *Message, err error) {
	switch sr.frame {
	case FrameDefault:
		return sr.einoReader.Recv()
	case FrameEino:
		return sr.einoReader.Recv()
	default:
		return nil, errorx.NewByCode(llm_errorx.ModelInvalidCode, errorx.WithExtraMsg(fmt.Sprintf("frame:%s is not valid", sr.frame)))
	}
}
