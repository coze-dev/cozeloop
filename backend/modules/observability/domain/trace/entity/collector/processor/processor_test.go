// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package processor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/component"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/consumer"
	consumermocks "github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/consumer/mocks"
	processormock "github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/processor/mocks"
)

func getComponentID(s string) component.ID {
	id := new(component.ID)
	_ = id.UnmarshalText([]byte(s))
	return *id
}

func TestNewFactory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cfgType := component.Type("test")
	createFunc := func(ctx context.Context, set CreateSettings, cfg component.Config, next consumer.Consumer) (Processor, error) {
		return processormock.NewMockProcessor(ctrl), nil
	}
	factory := NewFactory(cfgType, nil, createFunc)
	assert.Equal(t, cfgType, factory.Type())
}

func TestMakeFactoryMap(t *testing.T) {
	t.Run("duplicate_factory_type", func(t *testing.T) {
		f1 := NewFactory("test", nil, nil)
		f2 := NewFactory("test", nil, nil)
		_, err := MakeFactoryMap(f1, f2)
		assert.ErrorContains(t, err, "duplicate processor factory")
	})
}

func TestBuilder_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful_creation", func(t *testing.T) {
		cfg := processormock.NewMockProcessor(ctrl)
		factory := NewFactory("test", func() component.Config { return cfg },
			func(ctx context.Context, settings CreateSettings, config component.Config, c consumer.Consumer) (Processor, error) {
				return nil, nil
			})
		builder := NewBuilder(
			map[component.ID]component.Config{getComponentID("test"): cfg},
			map[component.Type]Factory{factory.Type(): factory},
		)

		_, err := builder.Create(context.Background(), CreateSettings{
			ID: getComponentID("test"),
		}, consumermocks.NewMockConsumer(ctrl))
		assert.NoError(t, err)
	})
}

func TestFactory_CreateTracesProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil_create_func_error", func(t *testing.T) {
		factory := &factory{
			cfgType:             component.Type("test"),
			CreateProcessorFunc: nil,
		}

		_, err := factory.CreateTracesProcessor(
			context.Background(),
			CreateSettings{},
			nil,
			consumermocks.NewMockConsumer(ctrl),
		)
		assert.ErrorContains(t, err, "nil create trace exporter function")
	})
}
