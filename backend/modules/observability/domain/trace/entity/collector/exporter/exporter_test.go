// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package exporter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/component"
	exportermock "github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/exporter/mocks"
)

func TestNewFactory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfgType := component.Type("test")
	createFunc := func(ctx context.Context, set CreateSettings, cfg component.Config) (Exporter, error) {
		return exportermock.NewMockExporter(ctrl), nil
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
		factory := NewFactory("test", func() component.Config { return nil },
			func(ctx context.Context, settings CreateSettings, config component.Config) (Exporter, error) {
				return nil, nil
			})
		builder := NewBuilder(
			map[component.ID]component.Config{parseComponentID("test"): nil},
			map[component.Type]Factory{factory.Type(): factory},
		)

		_, err := builder.Create(context.Background(), CreateSettings{
			ID: parseComponentID("test"),
		})
		assert.NoError(t, err)
	})
}

func parseComponentID(s string) component.ID {
	id := component.ID{}
	_ = id.UnmarshalText([]byte(s))
	return id
}
