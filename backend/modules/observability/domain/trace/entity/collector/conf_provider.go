// Original Files: open-telemetry/opentelemetry-collector
// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// This file may have been modified by ByteDance Ltd.

package collector

import (
	"context"
	"fmt"

	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/component"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/confmap"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/exporter"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/processor"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/receiver"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/service"
	"github.com/coze-dev/coze-loop/backend/pkg/conf"
)

const traceCollectorCfgKey = "trace_collector_cfg"

//go:generate mockgen -destination=mocks/conf_provider.go -package=mocks . ConfigProvider
type ConfigProvider interface {
	Get(ctx context.Context, factories Factories) (*Config, error)
}

type configProvider struct {
	confP conf.IConfigLoader
}

func NewConfigProvider(confP conf.IConfigLoader) ConfigProvider {
	return &configProvider{
		confP: confP,
	}
}

type configSettings struct {
	Receivers  map[component.ID]map[string]any `mapstructure:"receivers"`
	Processors map[component.ID]map[string]any `mapstructure:"processors"`
	Exporters  map[component.ID]map[string]any `mapstructure:"exporters"`
	Tenants    map[string]*service.Config      `mapstructure:"tenants"`
}

func (cm *configProvider) Get(ctx context.Context, factories Factories) (*Config, error) {
	var cfg Config
	ret, err := cm.retrieveValue(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve the configuration, %v", err)
	}
	rawConf, ok := ret.RawConf().(map[string]any)
	if !ok {
		return nil, fmt.Errorf("cannot convert the configuration %s to a map[string]any", ret.RawConf())
	}
	set := &configSettings{}
	if err = confmap.DecodeConfig(rawConf, set); err != nil {
		return nil, fmt.Errorf("cannot unmarshal the configuration, %v", err)
	}
	cfg.Receivers, err = unmarshalTenantConfig[receiver.Factory](ctx, set.Receivers, factories.Receivers)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal receivers configuration, %v", err)
	}
	cfg.Processors, err = unmarshalTenantConfig[processor.Factory](ctx, set.Processors, factories.Processors)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal processors configuration, %v", err)
	}
	cfg.Exporters, err = unmarshalTenantConfig[exporter.Factory](ctx, set.Exporters, factories.Exporters)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal exporters configuration, %v", err)
	}
	cfg.Tenants = set.Tenants
	return &cfg, nil
}

func unmarshalTenantConfig[F component.Factory](ctx context.Context, rawCfgs map[component.ID]map[string]any, factories map[component.Type]F) (cfgs map[component.ID]component.Config, err error) {
	if rawCfgs == nil {
		return nil, fmt.Errorf("empty rawCfgs")
	}
	// Prepare resulting map.
	cfgs = make(map[component.ID]component.Config)
	// Iterate over raw configs and create a config for each.
	for id, value := range rawCfgs {
		// Find factory based on component kind and type that we read from config source.
		factory, ok := factories[id.Type()]
		if !ok {
			return nil, fmt.Errorf("unknown type: %q for id: %q (valid values: %v)", id.Type(), id, factories)
		}
		// Create the default config for this component.
		cfg := factory.CreateDefaultConfig()
		// Now that the default config struct is created we can Unmarshal into it,
		// and it will apply user-defined config on top of the default.
		if err = confmap.DecodeConfig(value, cfg); err != nil {
			return nil, fmt.Errorf("error reading configuration for %q, %v", id, err)
		}
		cfgs[id] = cfg
	}

	return cfgs, nil
}

func (cm *configProvider) retrieveValue(ctx context.Context) (*confmap.Retrieved, error) {
	val := cm.confP.Get(ctx, traceCollectorCfgKey)
	if val == nil {
		return nil, fmt.Errorf("cannot retrieve the collector configuration")
	}
	return confmap.NewRetrieved(val)
}

func (cm *configProvider) Shutdown(ctx context.Context) error {
	return nil
}
