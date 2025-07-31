// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package conf_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/coze-loop/backend/modules/data/domain/component/conf"
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
)

func TestDatasetSpec_GetSpecByCategory(t *testing.T) {
	// 定义测试用的 DatasetCategory
	category1 := entity.DatasetCategoryEvaluation
	category2 := entity.DatasetCategoryUnknown

	// 定义测试用的 DatasetSpec
	spec1 := &entity.DatasetSpec{}

	// 创建 DatasetSpec 实例
	datasetSpec := &conf.DatasetSpec{
		Spec: &entity.DatasetSpec{},
		SpecsByCategory: map[entity.DatasetCategory]*entity.DatasetSpec{
			category1: spec1,
		},
	}

	// 定义测试用例
	tests := []struct {
		name     string
		spec     *conf.DatasetSpec
		category entity.DatasetCategory
		want     *entity.DatasetSpec
	}{
		{
			name:     "从 SpecsByCategory 获取 Spec",
			spec:     datasetSpec,
			category: category1,
			want:     spec1,
		},
		{
			name:     "从 Spec 获取 Spec",
			spec:     datasetSpec,
			category: category2,
			want:     datasetSpec.Spec,
		},
		{
			name:     "Spec 为 nil",
			spec:     nil,
			category: category1,
			want:     nil,
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spec.GetSpecByCategory(tt.category); got != tt.want {
				t.Errorf("GetSpecByCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatasetFeature_GetFeatureByCategory(t *testing.T) {
	tests := []struct {
		name     string
		feature  *conf.DatasetFeature
		category entity.DatasetCategory
		want     *entity.DatasetFeatures
	}{
		{
			name:     "nil feature",
			feature:  nil,
			category: entity.DatasetCategoryGeneral,
			want:     nil,
		},
		{
			name: "category exists in map",
			feature: &conf.DatasetFeature{
				Feature: &entity.DatasetFeatures{
					EditSchema: false,
					MultiModal: false,
				},
				FeatureByCategory: map[entity.DatasetCategory]*entity.DatasetFeatures{
					entity.DatasetCategoryGeneral: {
						EditSchema: true,
						MultiModal: true,
					},
				},
			},
			category: entity.DatasetCategoryGeneral,
			want: &entity.DatasetFeatures{
				EditSchema: true,
				MultiModal: true,
			},
		},
		{
			name: "category not exists in map - fallback to default",
			feature: &conf.DatasetFeature{
				Feature: &entity.DatasetFeatures{
					EditSchema: false,
					MultiModal: false,
				},
				FeatureByCategory: map[entity.DatasetCategory]*entity.DatasetFeatures{
					entity.DatasetCategoryTraining: {
						EditSchema: true,
						MultiModal: true,
					},
				},
			},
			category: entity.DatasetCategoryGeneral,
			want: &entity.DatasetFeatures{
				EditSchema: false,
				MultiModal: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.feature.GetFeatureByCategory(tt.category)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSnapshotRetry_GetRetryInterval(t *testing.T) {
	tests := []struct {
		name string
		sr   *conf.SnapshotRetry
		want time.Duration
	}{
		{
			name: "nil snapshot retry",
			sr:   nil,
			want: 5 * time.Second,
		},
		{
			name: "custom retry interval",
			sr: &conf.SnapshotRetry{
				RetryIntervalMS: 10000, // 10 seconds
			},
			want: 10 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sr.GetRetryInterval()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSnapshotRetry_GetMaxProcessingTime(t *testing.T) {
	defaultTTL := 5 * time.Minute

	tests := []struct {
		name string
		sr   *conf.SnapshotRetry
		want time.Duration
	}{
		{
			name: "zero max processing time",
			sr: &conf.SnapshotRetry{
				MaxProcessingTimeS: 0,
			},
			want: defaultTTL,
		},
		{
			name: "custom max processing time",
			sr: &conf.SnapshotRetry{
				MaxProcessingTimeS: 600, // 10 minutes
			},
			want: 10 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sr.GetMaxProcessingTime()
			assert.Equal(t, tt.want, got)
		})
	}
}
