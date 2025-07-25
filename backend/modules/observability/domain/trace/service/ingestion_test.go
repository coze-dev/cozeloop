// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	mqmocks "github.com/coze-dev/cozeloop/backend/infra/mq/mocks"
	confmocks "github.com/coze-dev/cozeloop/backend/modules/observability/domain/component/config/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/exporter"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/processor"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/receiver"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/service/collector/exporter/clickhouseexporter"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/service/collector/processor/queueprocessor"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/service/collector/receiver/rmqreceiver"
)

var collectorCfg = map[string]any{
	"receivers": map[string]any{
		"rmq/default": map[string]any{
			"addr":           []string{"a"},
			"consumer_group": "11",
			"topic":          "22",
		},
	},
	"processors": map[string]any{
		"queue/default": map[string]any{
			"pool_name":         "a",
			"max_pool_size":     100,
			"queue_size":        100,
			"max_batch_size":    100,
			"tick_intervals_ms": 1000,
			"shard_count":       4,
		},
	},
	"exporters": map[string]any{
		"clickhouse/default": nil,
	},
	"tenants": map[string]any{
		"cozeloop": map[string]any{
			"receivers":  []string{"rmq/default"},
			"processors": []string{"queue/default"},
			"exporters":  []string{"clickhouse/default"},
		},
	},
}

func TestIngestionServiceImpl_Run(t *testing.T) {
	type fields struct {
		c *collector.Collector
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		wantErr      bool
	}{
		{
			name: "run ingestion successfully",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mqConsumer := mqmocks.NewMockIConsumer(ctrl)
				mqConsumer.EXPECT().Start().Return(nil)
				mqConsumer.EXPECT().Close().Return(nil)
				mqConsumer.EXPECT().RegisterHandler(gomock.Any()).Return()
				mqFactory := mqmocks.NewMockIFactory(ctrl)
				mqFactory.EXPECT().NewConsumer(gomock.Any()).Return(mqConsumer, nil)
				confMocks := confmocks.NewMockITraceConfig(ctrl)
				confMocks.EXPECT().Get(gomock.Any(), gomock.Any()).Return(collectorCfg)
				factories := NewIngestionCollectorFactory(
					[]receiver.Factory{
						rmqreceiver.NewFactory(mqFactory),
					},
					[]processor.Factory{
						queueprocessor.NewFactory(),
					},
					[]exporter.Factory{
						clickhouseexporter.NewFactory(nil),
					},
				)
				c, err := collector.New(collector.Settings{
					Factories:      factories.GetCollectorFactory,
					ConfigProvider: collector.NewConfigProvider(confMocks),
				})
				if err != nil {
					t.Fatal(err)
				}
				return fields{
					c: c,
				}
			},
			wantErr: false,
		},
		{
			name: "run ingestion failed",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mqConsumer := mqmocks.NewMockIConsumer(ctrl)
				mqConsumer.EXPECT().Start().Return(fmt.Errorf("err"))
				mqConsumer.EXPECT().Close().Return(nil)
				mqConsumer.EXPECT().RegisterHandler(gomock.Any()).Return()
				mqFactory := mqmocks.NewMockIFactory(ctrl)
				mqFactory.EXPECT().NewConsumer(gomock.Any()).Return(mqConsumer, nil)
				confMocks := confmocks.NewMockITraceConfig(ctrl)
				confMocks.EXPECT().Get(gomock.Any(), gomock.Any()).Return(collectorCfg)
				factories := NewIngestionCollectorFactory(
					[]receiver.Factory{
						rmqreceiver.NewFactory(mqFactory),
					},
					[]processor.Factory{
						queueprocessor.NewFactory(),
					},
					[]exporter.Factory{
						clickhouseexporter.NewFactory(nil),
					},
				)
				c, err := collector.New(collector.Settings{
					Factories:      factories.GetCollectorFactory,
					ConfigProvider: collector.NewConfigProvider(confMocks),
				})
				if err != nil {
					t.Fatal(err)
				}
				return fields{
					c: c,
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			i := &IngestionServiceImpl{
				c: fields.c,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			if err := i.RunSync(ctx); (err != nil) != tt.wantErr {
				t.Errorf("IngestionServiceImpl.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIngestionServiceImpl_RunAsync(t *testing.T) {
	type fields struct {
		c *collector.Collector
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		wantErr      bool
	}{
		{
			name: "run ingestion successfully",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mqConsumer := mqmocks.NewMockIConsumer(ctrl)
				mqConsumer.EXPECT().Start().Return(nil)
				mqConsumer.EXPECT().Close().Return(nil)
				mqConsumer.EXPECT().RegisterHandler(gomock.Any()).Return()
				mqFactory := mqmocks.NewMockIFactory(ctrl)
				mqFactory.EXPECT().NewConsumer(gomock.Any()).Return(mqConsumer, nil)
				confMocks := confmocks.NewMockITraceConfig(ctrl)
				confMocks.EXPECT().Get(gomock.Any(), gomock.Any()).Return(collectorCfg)
				factories := NewIngestionCollectorFactory(
					[]receiver.Factory{
						rmqreceiver.NewFactory(mqFactory),
					},
					[]processor.Factory{
						queueprocessor.NewFactory(),
					},
					[]exporter.Factory{
						clickhouseexporter.NewFactory(nil),
					},
				)
				c, err := collector.New(collector.Settings{
					Factories:      factories.GetCollectorFactory,
					ConfigProvider: collector.NewConfigProvider(confMocks),
				})
				if err != nil {
					t.Fatal(err)
				}
				return fields{
					c: c,
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			i := &IngestionServiceImpl{
				c: fields.c,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			i.RunAsync(ctx)
			time.Sleep(200 * time.Millisecond)
		})
	}
}
