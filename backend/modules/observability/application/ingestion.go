// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/service"
)

type ITraceIngestionApplication interface {
	RunSync(context.Context) error
	RunAsync(context.Context)
}

func NewIngestionApplication(svc service.IngestionService) ITraceIngestionApplication {
	impl := &IngestionApplicationImpl{
		ingestionService: svc,
	}
	return impl
}

type IngestionApplicationImpl struct {
	ingestionService service.IngestionService
}

func (i *IngestionApplicationImpl) RunAsync(ctx context.Context) {
	i.ingestionService.RunAsync(ctx)
}

func (i *IngestionApplicationImpl) RunSync(ctx context.Context) error {
	return i.ingestionService.RunSync(ctx)
}
