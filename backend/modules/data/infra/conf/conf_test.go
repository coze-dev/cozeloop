// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	dataset_conf "github.com/coze-dev/cozeloop/backend/modules/data/domain/component/conf"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/consts"
	mock_conf "github.com/coze-dev/cozeloop/backend/pkg/conf/mocks" // 假设 mock 文件在此路径
)

func TestNewConfiger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfigFactory := mock_conf.NewMockIConfigLoaderFactory(ctrl)
	mockConfigLoader := mock_conf.NewMockIConfigLoader(ctrl)

	tests := []struct {
		name              string
		mockSetup         func()
		expectedConfiger  dataset_conf.IConfig
		expectedErr       error
		checkConfigerFunc func(t *testing.T, c dataset_conf.IConfig)
	}{
		{
			name: "正常创建 Configer",
			mockSetup: func() {
				mockConfigFactory.EXPECT().NewConfigLoader(consts.DataConfigFileName).Return(mockConfigLoader, nil)
			},
			expectedConfiger: &configer{
				loader: mockConfigLoader,
			},
			expectedErr: nil,
			checkConfigerFunc: func(t *testing.T, c dataset_conf.IConfig) {
				assert.NotNil(t, c)
				// 可以进一步检查 configer 内部的 loader 是否被正确设置
				actualConf, ok := c.(*configer)
				assert.True(t, ok)
				assert.Equal(t, mockConfigLoader, actualConf.loader)
			},
		},
		{
			name: "异常场景 - NewConfigLoader 失败",
			mockSetup: func() {
				mockConfigFactory.EXPECT().NewConfigLoader(consts.DataConfigFileName).Return(nil, errors.New("failed to create loader"))
			},
			expectedConfiger: nil,
			expectedErr:      errors.New("failed to create loader"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			configerInstance, err := NewConfiger(mockConfigFactory)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Nil(t, configerInstance)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, configerInstance)
				// 由于我们比较的是接口和具体类型的实例，直接 assert.Equal 可能不够准确
				// 最好是类型断言后比较内部字段，或者像 expectedConfiger 一样创建一个期望的实例进行比较
				if tt.checkConfigerFunc != nil {
					tt.checkConfigerFunc(t, configerInstance)
				} else {
					assert.Equal(t, tt.expectedConfiger, configerInstance)
				}
			}
		})
	}
}

func TestConfiger_GetConsumerConfigs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoader := mock_conf.NewMockIConfigLoader(ctrl)
	c := &configer{loader: mockLoader}

	// 定义测试用例
	tests := []struct {
		name           string
		mockSetup      func()
		expectedConfig *dataset_conf.ConsumerConfig
	}{
		{
			name: "正常获取 ConsumerConfigs",
			mockSetup: func() {
				// 模拟 UnmarshalKey 成功并返回预期配置
				mockLoader.EXPECT().UnmarshalKey(context.Background(), "consumer_configs", gomock.Any()).Return(nil)
			},
			expectedConfig: nil,
		},
		{
			name: "异常场景 - UnmarshalKey 失败",
			mockSetup: func() {
				// 模拟 UnmarshalKey 失败
				mockLoader.EXPECT().UnmarshalKey(context.Background(), "consumer_configs", gomock.Any()).
					Return(errors.New("配置解析失败"))
			},
			expectedConfig: &dataset_conf.ConsumerConfig{}, // 失败时返回空结构体
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			result := c.GetConsumerConfigs()
			assert.Equal(t, tt.expectedConfig, result)
		})
	}
}

func TestConfiger_GetSnapshotRetry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoader := mock_conf.NewMockIConfigLoader(ctrl)
	c := &configer{loader: mockLoader}

	tests := []struct {
		name           string
		mockSetup      func()
		expectedConfig *dataset_conf.SnapshotRetry
	}{
		{
			name: "正常获取 SnapshotRetry",
			mockSetup: func() {
				mockLoader.EXPECT().UnmarshalKey(context.Background(), "snapshot_retry", gomock.Any()).Return(nil)
			},
			expectedConfig: nil,
		},
		{
			name: "异常场景 - UnmarshalKey 失败",
			mockSetup: func() {
				mockLoader.EXPECT().UnmarshalKey(context.Background(), "snapshot_retry", gomock.Any()).
					Return(errors.New("配置解析失败"))
			},
			expectedConfig: &dataset_conf.SnapshotRetry{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			result := c.GetSnapshotRetry()
			assert.Equal(t, tt.expectedConfig, result)
		})
	}
}

func TestConfiger_GetProducerConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoader := mock_conf.NewMockIConfigLoader(ctrl)
	c := &configer{loader: mockLoader}

	tests := []struct {
		name           string
		mockSetup      func()
		expectedConfig *dataset_conf.ProducerConfig
	}{
		{
			name: "正常获取 ProducerConfig",
			mockSetup: func() {
				mockLoader.EXPECT().UnmarshalKey(context.Background(), "job_mq_producer", gomock.Any()).Return(nil)
			},
			expectedConfig: nil,
		},
		{
			name: "异常场景 - UnmarshalKey 失败",
			mockSetup: func() {
				mockLoader.EXPECT().UnmarshalKey(context.Background(), "job_mq_producer", gomock.Any()).
					Return(errors.New("配置解析失败"))
			},
			expectedConfig: &dataset_conf.ProducerConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			result := c.GetProducerConfig()
			assert.Equal(t, tt.expectedConfig, result)
		})
	}
}

func TestConfiger_GetDatasetFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoader := mock_conf.NewMockIConfigLoader(ctrl)
	c := &configer{loader: mockLoader}

	tests := []struct {
		name     string
		mock     func()
		expected *dataset_conf.DatasetFeature
	}{{
		name: "正常场景: 配置加载成功",
		mock: func() {
			mockConf := &dataset_conf.DatasetFeature{}
			mockLoader.EXPECT().UnmarshalKey(gomock.Any(), "default_dataset_feature", gomock.Any()).
				SetArg(2, mockConf).Return(nil)
		},
		expected: &dataset_conf.DatasetFeature{},
	}, {
		name: "异常场景: 配置加载失败",
		mock: func() {
			mockLoader.EXPECT().UnmarshalKey(gomock.Any(), "default_dataset_feature", gomock.Any()).
				Return(assert.AnError)
		},
		expected: &dataset_conf.DatasetFeature{},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			assert.Equal(t, tt.expected, c.GetDatasetFeature())
		})
	}
}

func TestConfiger_GetDatasetItemStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoader := mock_conf.NewMockIConfigLoader(ctrl)
	c := &configer{loader: mockLoader}

	tests := []struct {
		name     string
		mock     func()
		expected *dataset_conf.DatasetItemStorage
	}{{
		name: "正常场景: 配置加载成功",
		mock: func() {
			mockConf := &dataset_conf.DatasetItemStorage{}
			mockLoader.EXPECT().UnmarshalKey(gomock.Any(), "dataset_item_storage", gomock.Any()).
				SetArg(2, mockConf).Return(nil)
		},
		expected: &dataset_conf.DatasetItemStorage{},
	}, {
		name: "异常场景: 配置加载失败",
		mock: func() {
			mockLoader.EXPECT().UnmarshalKey(gomock.Any(), "dataset_item_storage", gomock.Any()).
				Return(assert.AnError)
		},
		expected: &dataset_conf.DatasetItemStorage{},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			assert.Equal(t, tt.expected, c.GetDatasetItemStorage())
		})
	}
}

func TestConfiger_GetDatasetSpec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoader := mock_conf.NewMockIConfigLoader(ctrl)
	c := &configer{loader: mockLoader}

	tests := []struct {
		name     string
		mock     func()
		expected *dataset_conf.DatasetSpec
	}{{
		name: "正常场景: 配置加载成功",
		mock: func() {
			mockConf := &dataset_conf.DatasetSpec{}
			mockLoader.EXPECT().UnmarshalKey(gomock.Any(), "default_dataset_spec", gomock.Any()).
				SetArg(2, mockConf).Return(nil)
		},
		expected: &dataset_conf.DatasetSpec{}, // 注意：此处逻辑与其他方法相反
	}, {
		name: "异常场景: 配置加载失败",
		mock: func() {
			mockConf := &dataset_conf.DatasetSpec{}
			mockLoader.EXPECT().UnmarshalKey(gomock.Any(), "default_dataset_spec", gomock.Any()).
				SetArg(2, mockConf).Return(assert.AnError)
		},
		expected: &dataset_conf.DatasetSpec{},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			assert.Equal(t, tt.expected, c.GetDatasetSpec())
		})
	}
}
