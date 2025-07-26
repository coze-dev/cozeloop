// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package producer

import (
	"context"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	infra_mq "github.com/coze-dev/coze-loop/backend/infra/mq"
	mock_infra_mq "github.com/coze-dev/coze-loop/backend/infra/mq/mocks" // 使用这个路径
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/component/conf"
	mock_config "github.com/coze-dev/coze-loop/backend/modules/data/domain/component/conf/mocks"
	component_mq "github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/component/mq"
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
)

func TestDatasetJobPublisherSend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProducer := mock_infra_mq.NewMockIProducer(ctrl) // 使用 gomock 生成的 mock

	// 准备一些通用的测试数据
	ctx := context.Background()
	testTopic := "test-topic"
	testTag := "test-tag"
	baseMsg := &entity.JobRunMessage{
		Type:     "test_type",
		SpaceID:  1,
		TaskID:   2,
		RunID:    3,
		JobID:    4,
		Extra:    map[string]string{"key": "value"},
		Operator: "tester",
	}

	// 定义测试用例
	tests := []struct {
		name        string
		publisher   *DatasetJobPublisher
		msg         *entity.JobRunMessage
		opts        []component_mq.MessageOpt
		setupMock   func(msgBody []byte, expectedMqMsg *infra_mq.Message) // 用于设置 mock producer 的期望
		expectedErr error
	}{
		{
			name: "成功发送消息 - 无选项",
			publisher: &DatasetJobPublisher{
				Topic:    testTopic,
				Tag:      testTag,
				producer: mockProducer,
			},
			msg:  baseMsg,
			opts: []component_mq.MessageOpt{},
			setupMock: func(msgBody []byte, expectedMqMsg *infra_mq.Message) {
				// 预期 mq.Message
				// 注意：这里直接比较 expectedMqMsg 可能因为 Body 是 []byte 而导致指针不同而不匹配。
				// gomock.Eq 在比较结构体时，如果包含 slice/map 等，会进行深比较。
				// 或者使用 gomock.Cond 来进行更灵活的匹配。
				// 这里我们假设 gomock.Eq 能够处理好 *infra_mq.Message 的比较，或者使用更具体的匹配器。
				// 为了简单，我们这里直接构造一个期望的 mq.Message
				// 实际项目中，可能需要更细致地匹配 mqMsg 的字段。
				mockProducer.EXPECT().Send(gomock.Any(), gomock.Cond(func(x interface{}) bool {
					actualMsg, ok := x.(*infra_mq.Message)
					if !ok {
						return false
					}
					return actualMsg.Topic == testTopic &&
						actualMsg.Tag == testTag &&
						string(actualMsg.Body) == string(msgBody) && // 比较 body 内容
						actualMsg.PartitionKey == "" &&
						actualMsg.DeferDuration == 0
				})).Return(infra_mq.SendResponse{MessageID: "mock-id", Offset: 123}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "成功发送消息 - 带选项",
			publisher: &DatasetJobPublisher{
				Topic:    testTopic,
				Tag:      testTag,
				producer: mockProducer,
			},
			msg: baseMsg,
			opts: []component_mq.MessageOpt{
				func(opt *component_mq.MessageOption) {
					opt.Key = "partition-key-1"
					opt.RetryInterval = 5 * time.Second
				},
			},
			setupMock: func(msgBody []byte, expectedMqMsg *infra_mq.Message) {
				mockProducer.EXPECT().Send(gomock.Any(), gomock.Cond(func(x interface{}) bool {
					actualMsg, ok := x.(*infra_mq.Message)
					if !ok {
						return false
					}
					return actualMsg.Topic == testTopic &&
						actualMsg.Tag == testTag &&
						string(actualMsg.Body) == string(msgBody) &&
						actualMsg.PartitionKey == "partition-key-1" &&
						actualMsg.DeferDuration == 5*time.Second
				})).Return(infra_mq.SendResponse{MessageID: "mock-id-2", Offset: 124}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重新计算当前测试用例的 msgBody，因为 msg 可能在不同用例中不同
			currentMsgBody, marshalErr := sonic.Marshal(tt.msg)
			if tt.name == "序列化消息失败" { // 特殊处理序列化失败的场景
				assert.Error(t, marshalErr) // 确认这里确实应该序列化失败
			}

			// 构造预期的 mq.Message，用于 setupMock
			// 注意：这个 expectedMqMsg 主要是为了在 setupMock 中方便地引用其字段，
			// 实际的 gomock 匹配逻辑在 setupMock 内部实现。
			opt := &component_mq.MessageOption{}
			for _, o := range tt.opts {
				o(opt)
			}
			expectedMqMsg := &infra_mq.Message{
				Topic:         tt.publisher.Topic,
				Body:          currentMsgBody, // 使用当前用例的 body
				Tag:           tt.publisher.Tag,
				PartitionKey:  opt.Key,
				DeferDuration: opt.RetryInterval,
			}

			if tt.setupMock != nil {
				tt.setupMock(currentMsgBody, expectedMqMsg)
			}

			err := tt.publisher.Send(ctx, tt.msg, tt.opts...)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				// 比较错误信息可能需要更灵活的方式，例如 errors.Is 或者字符串包含
				// assert.EqualError(t, err, tt.expectedErr.Error()) // 直接比较 error string 可能不稳定
				assert.Contains(t, err.Error(), tt.expectedErr.Error()) // 改为包含判断，因为包装后的错误信息可能更复杂
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewDatasetJobPublisher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks instances, created once and configured per test case
	mockConfiger := mock_config.NewMockIConfig(ctrl)
	mockMqFactory := mock_infra_mq.NewMockIFactory(ctrl)
	// mockProducer will be configured to be returned by mockMqFactory
	mockActualProducer := mock_infra_mq.NewMockIProducer(ctrl)

	topic := "test-topic"
	tag := "test-tag"
	// Common configuration data
	commonDomainProducerConf := &conf.ProducerConfig{
		Topic:          topic,
		Tag:            tag,
		Addr:           []string{"localhost:9876"},
		ProduceTimeout: 5 * time.Second,
		ProducerGroup:  "test-group",
	}

	expectedInfraMqProducerConf := infra_mq.ProducerConfig{
		Addr:           commonDomainProducerConf.Addr,
		ProduceTimeout: commonDomainProducerConf.ProduceTimeout,
		ProducerGroup:  &commonDomainProducerConf.ProducerGroup, // Important: pointer to the string
	}

	// Test cases
	testCases := []struct {
		name           string
		setupMocks     func()
		expectedPanic  bool
		expectedError  error
		validateResult func(t *testing.T, publisher component_mq.IDatasetJobPublisher, err error)
	}{
		{
			name: "成功创建Publisher",
			setupMocks: func() {
				mockConfiger.EXPECT().GetProducerConfig().Return(commonDomainProducerConf)
				mockMqFactory.EXPECT().NewProducer(gomock.Eq(expectedInfraMqProducerConf)).Return(mockActualProducer, nil)
				mockActualProducer.EXPECT().Start().Return(nil)
			},
			expectedPanic: false,
			expectedError: nil,
			validateResult: func(t *testing.T, publisher component_mq.IDatasetJobPublisher, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, publisher)
				// Check if the returned publisher is of the concrete type and has correct fields
				p, ok := publisher.(*DatasetJobPublisher)
				assert.True(t, ok, "Publisher is not of type *DatasetJobPublisher")
				if ok {
					assert.Equal(t, topic, p.Topic)
					assert.Equal(t, tag, p.Tag)
					assert.Equal(t, mockActualProducer, p.producer)
				}
			},
		},
		{
			name: "失败 - GetProducerConfig返回nil (导致panic)",
			setupMocks: func() {
				mockConfiger.EXPECT().GetProducerConfig().Return(nil)
				// NewProducer and Start should not be called
				mockMqFactory.EXPECT().NewProducer(gomock.Any()).Times(0)
				mockActualProducer.EXPECT().Start().Times(0)
			},
			expectedPanic: true, // Expecting a panic due to nil pointer dereference
			expectedError: nil,  // No error returned, but panic
			validateResult: func(t *testing.T, publisher component_mq.IDatasetJobPublisher, err error) {
				// Validated by assert.Panics in the loop
				assert.Nil(t, publisher) // Publisher should be nil if panic occurred before return
			},
		},
		{
			name: "失败 - mqFactory.NewProducer返回错误",
			setupMocks: func() {
				mockConfiger.EXPECT().GetProducerConfig().Return(commonDomainProducerConf)
				mockMqFactory.EXPECT().NewProducer(gomock.Eq(expectedInfraMqProducerConf)).Return(nil, errors.New("new producer failed"))
				// Start should not be called
				mockActualProducer.EXPECT().Start().Times(0)
			},
			expectedPanic: false,
			expectedError: errors.New("new producer failed"),
			validateResult: func(t *testing.T, publisher component_mq.IDatasetJobPublisher, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "new producer failed")
				assert.Nil(t, publisher)
			},
		},
		{
			name: "失败 - producer.Start返回错误",
			setupMocks: func() {
				mockConfiger.EXPECT().GetProducerConfig().Return(commonDomainProducerConf)
				mockMqFactory.EXPECT().NewProducer(gomock.Eq(expectedInfraMqProducerConf)).Return(mockActualProducer, nil)
				mockActualProducer.EXPECT().Start().Return(errors.New("producer start failed"))
			},
			expectedPanic: false,
			expectedError: errors.New("producer start failed"),
			validateResult: func(t *testing.T, publisher component_mq.IDatasetJobPublisher, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "producer start failed")
				assert.Nil(t, publisher)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks for the current test case
			tc.setupMocks()

			var publisher component_mq.IDatasetJobPublisher
			var err error

			if tc.expectedPanic {
				// Defer a function to recover from panic and check if it occurred
				// This is necessary because assert.Panics executes the function itself
				// and we want to capture the return values if no panic occurs (though not expected here)
				assert.Panics(t, func() {
					// Act: Call the function under test
					// We don't assign to publisher and err here directly if we expect a panic
					// because the assignment might not happen.
					// The purpose of calling it inside assert.Panics is just to check for panic.
					_, _ = NewDatasetJobPublisher(mockConfiger, mockMqFactory)
				}, "Expected NewDatasetJobPublisher to panic")
				// For panic cases, the publisher and err might not be meaningfully set if panic happens early.
				// The validateResult can check for nil publisher.
				if tc.validateResult != nil {
					// Pass nil explicitly as publisher and err are not reliably set if panic occurs.
					tc.validateResult(t, nil, nil)
				}

			} else {
				// Act: Call the function under test
				publisher, err = NewDatasetJobPublisher(mockConfiger, mockMqFactory)

				// Assert: Check results
				if tc.validateResult != nil {
					tc.validateResult(t, publisher, err)
				}
			}
		})
	}
}
