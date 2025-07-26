// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package clickhouseexporter

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/component"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/exporter"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/repo"
)

const (
	exporterType = "clickhouse"
)

func createDefaultConfig() component.Config {
	return &Config{}
}

func NewFactory(traceRepo repo.ITraceRepo) exporter.Factory {
	return exporter.NewFactory(
		exporterType,
		createDefaultConfig,
		func(ctx context.Context, params exporter.CreateSettings, cfg component.Config) (exporter.Exporter, error) {
			return &ckExporter{
				config:    cfg.(*Config),
				traceRepo: traceRepo,
			}, nil
		},
	)
}
