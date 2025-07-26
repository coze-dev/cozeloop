// Original Files: open-telemetry/opentelemetry-collector
// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// This file may have been modified by ByteDance Ltd.

package processor

import (
	"context"
	"fmt"

	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/component"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/consumer"
)

//go:generate mockgen -destination=mocks/processor.go -package=mocks . Processor
type Processor interface {
	component.Component
	consumer.Consumer
}

type Factory interface {
	component.Factory
	CreateTracesProcessor(ctx context.Context, set CreateSettings, cfg component.Config, nextConsumer consumer.Consumer) (Processor, error)
}

func NewFactory(cfgType component.Type, createDefaultConfig component.CreateDefaultConfigFunc, createProcessorFunc CreateProcessorFunc) Factory {
	f := &factory{
		cfgType:                 cfgType,
		CreateDefaultConfigFunc: createDefaultConfig,
		CreateProcessorFunc:     createProcessorFunc,
	}
	return f
}

type CreateProcessorFunc func(context.Context, CreateSettings, component.Config, consumer.Consumer) (Processor, error)

type factory struct {
	cfgType component.Type
	component.CreateDefaultConfigFunc
	CreateProcessorFunc
}

func (f factory) Type() component.Type {
	return f.cfgType
}

func (f CreateProcessorFunc) CreateTracesProcessor(
	ctx context.Context,
	set CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Consumer,
) (Processor, error) {
	if f == nil {
		return nil, fmt.Errorf("nil create trace exporter function")
	}
	t, err := f(ctx, set, cfg, nextConsumer)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func MakeFactoryMap(factories ...Factory) (map[component.Type]Factory, error) {
	fMap := map[component.Type]Factory{}
	for _, f := range factories {
		if _, ok := fMap[f.Type()]; ok {
			return nil, fmt.Errorf("duplicate processor factory %q", f.Type())
		}
		fMap[f.Type()] = f
	}
	return fMap, nil
}

type CreateSettings struct {
	ID component.ID
}

type Builder struct {
	cfgs      map[component.ID]component.Config
	factories map[component.Type]Factory
}

func (b *Builder) Create(ctx context.Context, set CreateSettings, next consumer.Consumer) (Processor, error) {
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("processor %q is not configured", set.ID)
	}
	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("processor factory not available for: %q", set.ID)
	}
	t, err := f.CreateTracesProcessor(ctx, set, cfg, next)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func NewBuilder(cfgs map[component.ID]component.Config, factories map[component.Type]Factory) *Builder {
	return &Builder{cfgs: cfgs, factories: factories}
}
