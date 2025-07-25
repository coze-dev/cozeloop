// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package metrics

type MetricType string

const (
	MetricTypeCounter     MetricType = "counter"
	MetricTypeRateCounter MetricType = "rate_counter"
	MetricTypeStore       MetricType = "store"
	MetricTypeTimer       MetricType = "timer"
	MetricTypeHistogram   MetricType = "histogram"
)

//go:generate mockgen -destination ./mocks/provider.go  --package mocks . Meter
type Meter interface {
	NewMetric(name string, types []MetricType, tagNames []string) (Metric, error)
}

var provider Meter = noopMeter{}

// GetMeter Get the metric provider. Must call InitMeter first.
func GetMeter() Meter {
	return provider
}

// InitMeter Init the metric provider. Must call before GetMeter.
func InitMeter(p Meter) {
	provider = p
}
