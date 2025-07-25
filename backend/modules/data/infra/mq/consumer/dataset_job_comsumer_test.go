// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package consumer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/infra/mq"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/component/conf"
	"github.com/coze-dev/cozeloop/backend/pkg/conf/mocks"
)

func TestDatasetJobConsumer_ConsumerCfg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfigLoader := mocks.NewMockIConfigLoader(ctrl)
	consumer := &DatasetJobConsumer{IConfigLoader: mockConfigLoader}

	// 测试用例定义
	tests := []struct {
		name          string
		mockSetup     func()
		expectedCfg   *mq.ConsumerConfig
		expectedError error
	}{{
		name: "正常场景: 配置加载成功且ConsumeGoroutineNums为正数",
		mockSetup: func() {
			cfg := &conf.ConsumerConfig{
				Topic:                "test_topic",
				ConsumerGroup:        "test_group",
				Orderly:              true,
				ConsumeTimeout:       5000,
				TagExpression:        "tagA",
				ConsumeGoroutineNums: 20,
			}
			mockConfigLoader.EXPECT().UnmarshalKey(gomock.Any(), "consumer_configs", gomock.Any()).
				SetArg(2, *cfg).Return(nil)
		},
		expectedCfg: &mq.ConsumerConfig{
			Topic:                "test_topic",
			ConsumerGroup:        "test_group",
			Orderly:              true,
			ConsumeTimeout:       5000,
			TagExpression:        "tagA",
			ConsumeGoroutineNums: 20,
		},
		expectedError: nil,
	}, {
		name: "边界场景: ConsumeGoroutineNums为0时使用默认值10",
		mockSetup: func() {
			cfg := &conf.ConsumerConfig{
				ConsumeGoroutineNums: 0,
			}
			mockConfigLoader.EXPECT().UnmarshalKey(gomock.Any(), "consumer_configs", gomock.Any()).
				SetArg(2, *cfg).Return(nil)
		},
		expectedCfg: &mq.ConsumerConfig{
			ConsumeGoroutineNums: 10,
		},
		expectedError: nil,
	}, {
		name: "边界场景: ConsumeGoroutineNums为负数时使用默认值10",
		mockSetup: func() {
			cfg := &conf.ConsumerConfig{
				ConsumeGoroutineNums: -5,
			}
			mockConfigLoader.EXPECT().UnmarshalKey(gomock.Any(), "consumer_configs", gomock.Any()).
				SetArg(2, *cfg).Return(nil)
		},
		expectedCfg: &mq.ConsumerConfig{
			ConsumeGoroutineNums: 10,
		},
		expectedError: nil,
	}}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			cfg, err := consumer.ConsumerCfg(context.Background())

			assert.Equal(t, tt.expectedError, err)
			if tt.expectedError == nil {
				assert.Equal(t, tt.expectedCfg.Addr, cfg.Addr)
				assert.Equal(t, tt.expectedCfg.Topic, cfg.Topic)
				assert.Equal(t, tt.expectedCfg.ConsumerGroup, cfg.ConsumerGroup)
				assert.Equal(t, tt.expectedCfg.Orderly, cfg.Orderly)
				assert.Equal(t, tt.expectedCfg.ConsumeTimeout, cfg.ConsumeTimeout)
				assert.Equal(t, tt.expectedCfg.TagExpression, cfg.TagExpression)
				assert.Equal(t, tt.expectedCfg.ConsumeGoroutineNums, cfg.ConsumeGoroutineNums)
			}
		})
	}
}
