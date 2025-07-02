// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package receiver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/component"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/consumer"
	receivermock "github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/receiver/mocks"
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
	createFunc := func(ctx context.Context, set CreateSettings, cfg component.Config, next consumer.Consumer) (Receiver, error) {
		return receivermock.NewMockReceiver(ctrl), nil
	}
	factory := NewFactory(cfgType, func() component.Config { return nil }, createFunc)
	assert.Equal(t, cfgType, factory.Type())
}

func TestMakeFactoryMap(t *testing.T) {
	f1 := NewFactory("type1", func() component.Config { return nil }, nil)
	f2 := NewFactory("type2", func() component.Config { return nil }, nil)

	t.Run("success", func(t *testing.T) {
		fMap, err := MakeFactoryMap(f1, f2)
		assert.NoError(t, err)
		assert.Len(t, fMap, 2)
	})

	t.Run("error", func(t *testing.T) {
		_, err := MakeFactoryMap(f1, f1)
		assert.ErrorContains(t, err, "duplicate processor factory")
	})
}

func TestBuilder_CreateTraces(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceiver := receivermock.NewMockReceiver(ctrl)
	factory := NewFactory("test", func() component.Config { return nil },
		func(ctx context.Context, set CreateSettings, cfg component.Config, next consumer.Consumer) (Receiver, error) {
			return mockReceiver, nil
		})

	t.Run("success", func(t *testing.T) {
		builder := NewBuilder(
			map[component.ID]component.Config{getComponentID("test"): nil},
			map[component.Type]Factory{"test": factory},
		)
		rec, err := builder.CreateTraces(context.Background(), CreateSettings{
			ID: getComponentID("test"),
		}, nil)
		assert.NoError(t, err)
		assert.Equal(t, mockReceiver, rec)
	})

	t.Run("error", func(t *testing.T) {
		builder := NewBuilder(
			map[component.ID]component.Config{},
			map[component.Type]Factory{"test": factory},
		)

		_, err := builder.CreateTraces(context.Background(), CreateSettings{
			ID: getComponentID("test"),
		}, nil)
		assert.ErrorContains(t, err, "receiver \"test\" is not configured")
	})
}

func TestFactory_CreateTracesReceiver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error", func(t *testing.T) {
		factory := &factory{
			cfgType:            "test",
			CreateReceiverFunc: nil,
		}

		_, err := factory.CreateTracesReceiver(context.Background(), CreateSettings{}, nil, nil)
		assert.ErrorContains(t, err, "nil create trace exporter function")
	})
}
