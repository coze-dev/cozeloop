// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package consumer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockConsumer struct{}

func (m *mockConsumer) ConsumeTraces(ctx context.Context, tds Traces) error {
	return nil
}

func TestComsumer_Fanout(t *testing.T) {
	fNode := &fanoutConsumer{
		traces: make([]Consumer, 4),
	}
	for i := 0; i < len(fNode.traces); i++ {
		fNode.traces[i] = &mockConsumer{}
	}
	assert.Nil(t, fNode.ConsumeTraces(context.Background(), Traces{}))
}
