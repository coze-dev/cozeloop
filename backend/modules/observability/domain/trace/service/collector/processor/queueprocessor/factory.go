// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package queueprocessor

import (
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/processor"
)

const (
	procType = "queue"
)

func NewFactory() processor.Factory {
	return processor.NewFactory(
		procType,
		createDefaultConfig,
		createTracesProcessor,
	)
}
