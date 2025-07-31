// Original Files: open-telemetry/opentelemetry-collector
// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// This file may have been modified by ByteDance Ltd.

package exporter

import (
	"context"
	"fmt"

	component "github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/component"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/consumer"
)

//go:generate mockgen -destination=mocks/exporter.go -package=mocks . Exporter
type Exporter interface {
	component.Component
	consumer.Consumer
}

type Factory interface {
	component.Factory
	CreateTracesExporter(ctx context.Context, set CreateSettings, cfg component.Config) (Exporter, error)
}

func NewFactory(cfgType component.Type, createDefaultConfig component.CreateDefaultConfigFunc, createExporterFunc CreateExporterFunc) Factory {
	f := &factory{
		cfgType:                 cfgType,
		CreateDefaultConfigFunc: createDefaultConfig,
		CreateExporterFunc:      createExporterFunc,
	}
	return f
}

type CreateExporterFunc func(context.Context, CreateSettings, component.Config) (Exporter, error)

type factory struct {
	cfgType component.Type
	component.CreateDefaultConfigFunc
	CreateExporterFunc
}

func (f factory) Type() component.Type {
	return f.cfgType
}

func (f CreateExporterFunc) CreateTracesExporter(
	ctx context.Context,
	set CreateSettings,
	cfg component.Config,
) (Exporter, error) {
	if f == nil {
		return nil, fmt.Errorf("nil create trace exporter function")
	}
	t, err := f(ctx, set, cfg)
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

func (b *Builder) Create(ctx context.Context, set CreateSettings) (Exporter, error) {
	cfg, existsCfg := b.cfgs[set.ID]
	if !existsCfg {
		return nil, fmt.Errorf("exporter %q is not configured", set.ID)
	}
	f, existsFactory := b.factories[set.ID.Type()]
	if !existsFactory {
		return nil, fmt.Errorf("exporter factory not available for: %q", set.ID)
	}
	t, err := f.CreateTracesExporter(ctx, set, cfg)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func NewBuilder(cfgs map[component.ID]component.Config, factories map[component.Type]Factory) *Builder {
	return &Builder{cfgs: cfgs, factories: factories}
}
