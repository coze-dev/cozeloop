// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package metrics

type noopMeter struct{}

func (n noopMeter) NewMetric(name string, types []MetricType, tagNames []string) (Metric, error) {
	return noopMetric{}, nil
}

type noopMetric struct{}

func (n noopMetric) Emit(tags []T, values ...*Value) {}
