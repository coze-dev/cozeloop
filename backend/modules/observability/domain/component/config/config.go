// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"

	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/loop_span"
	"github.com/coze-dev/cozeloop/backend/pkg/conf"
)

type SystemView struct {
	ID       int64  `mapstructure:"id" json:"id"`
	ViewName string `mapstructure:"view_name" json:"view_name"`
	Filters  string `mapstructure:"filters" json:"filters"`
}

type PlatformTenantsCfg struct {
	Config map[string][]string `mapstructure:"config" json:"config"`
}

type SpanTransHandlerConfig struct {
	PlatformCfg map[string]loop_span.SpanTransCfgList `mapstructure:"platform_cfg" json:"platform_cfg"`
}

type TraceMqProducerCfg struct {
	Addr          []string `mapstructure:"addr" json:"addr"`
	Timeout       int      `mapstructure:"timeout" json:"timeout"` // ms
	RetryTimes    int      `mapstructure:"retry_times" json:"retry_times"`
	Topic         string   `mapstructure:"topic" json:"topic"`
	ProducerGroup string   `mapstructure:"producer_group" json:"producer_group"`
}

type TraceCKCfg struct {
	Hosts       []string        `mapstructure:"hosts" json:"hosts"`
	DataBase    string          `mapstructure:"database" json:"database"`
	UserName    string          `mapstructure:"username" json:"username"`
	Password    string          `mapstructure:"password" json:"password"`
	DialTimeout int             `mapstructure:"dial_timeout" json:"dial_timeout"` // seconds
	ReadTimeout int             `mapstructure:"read_timeout" json:"read_timeout"` // seconds
	SuperFields map[string]bool `mapstructure:"super_fields" json:"super_fields"`
}

type TenantCfg struct {
	QueryTables map[string][]string              `mapstructure:"query_tables" json:"query_tables"`
	InsertTable map[string]map[entity.TTL]string `mapstructure:"insert_table" json:"insert_table"`
}

type FieldMeta struct {
	FieldType     loop_span.FieldType       `mapstructure:"field_type" json:"field_type"`
	FilterTypes   []loop_span.QueryTypeEnum `mapstructure:"filter_types" json:"filter_types"`
	FieldOptions  *loop_span.FieldOptions   `mapstructure:"field_options" json:"field_options"`
	SupportCustom bool                      `mapstructure:"support_custom" json:"support_custom"`
}

type TraceAttrTosCfg struct {
	Template   string `mapstructure:"template" json:"template"`
	Format     string `mapstructure:"format" json:"format"`
	Expiration int    `mapstructure:"ttl" json:"ttl"` // seconds
}

// AvailableFields: 配置可查询的Tag
// FieldMetas定义不同场景可使用的Key
type TraceFieldMetaInfoCfg struct {
	AvailableFields map[string]*FieldMeta                                          `mapstructure:"available_fields" json:"available_fields"`
	FieldMetas      map[loop_span.PlatformType]map[loop_span.SpanListType][]string `mapstructure:"field_metas" json:"field_metas"`
}

//go:generate mockgen -destination=mocks/config.go -package=mocks . ITraceConfig
type ITraceConfig interface {
	GetSystemViews(ctx context.Context) ([]*SystemView, error)
	GetPlatformTenants(ctx context.Context) (*PlatformTenantsCfg, error)
	GetPlatformSpansTrans(ctx context.Context) (*SpanTransHandlerConfig, error)
	GetTraceMqProducerCfg(ctx context.Context) (*TraceMqProducerCfg, error)
	GetTraceCkCfg(ctx context.Context) (*TraceCKCfg, error)
	GetTenantConfig(ctx context.Context) (*TenantCfg, error)
	GetTraceFieldMetaInfo(ctx context.Context) (*TraceFieldMetaInfoCfg, error)
	GetTraceAttrTosCfg(ctx context.Context) (*TraceAttrTosCfg, error)
	GetTraceDataMaxDurationDay(ctx context.Context, platformType *string) int64
	GetDefaultTraceTenant(ctx context.Context) string
	conf.IConfigLoader
}
