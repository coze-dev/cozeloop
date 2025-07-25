// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package metrics

type T struct {
	Name  string
	Value string
}

//go:generate mockgen -destination ./mocks/metrics.go  --package mocks . Metric
type Metric interface {
	Emit(tags []T, values ...*Value)
}

type Value struct {
	suffix string
	mType  MetricType
	value  *int64
	valuef *float64
}

func Counter(n int64, opts ...ValueOption) *Value {
	return buildValue(n, MetricTypeCounter, opts...)
}

func RateCounter(n int64, opts ...ValueOption) *Value {
	return buildValue(n, MetricTypeRateCounter, opts...)
}

func Store(n int64, opts ...ValueOption) *Value {
	return buildValue(n, MetricTypeStore, opts...)
}

func Timer(n int64, opts ...ValueOption) *Value {
	return buildValue(n, MetricTypeTimer, opts...)
}

func Histogram(n int64, opts ...ValueOption) *Value {
	return buildValue(n, MetricTypeHistogram, opts...)
}

func CounterF(n float64, opts ...ValueOption) *Value {
	return buildValueF(n, MetricTypeCounter, opts...)
}

func RateCounterF(n float64, opts ...ValueOption) *Value {
	return buildValueF(n, MetricTypeRateCounter, opts...)
}

func StoreF(n float64, opts ...ValueOption) *Value {
	return buildValueF(n, MetricTypeStore, opts...)
}

func TimerF(n float64, opts ...ValueOption) *Value {
	return buildValueF(n, MetricTypeTimer, opts...)
}

func HistogramF(n float64, opts ...ValueOption) *Value {
	return buildValueF(n, MetricTypeHistogram, opts...)
}

type ValueOption func(v *Value)

func WithSuffix(suffix string) ValueOption {
	return func(v *Value) {
		v.suffix = suffix
	}
}

func buildValue(n int64, mType MetricType, opts ...ValueOption) *Value {
	value := &Value{
		mType: mType,
		value: &n,
	}
	for _, opt := range opts {
		opt(value)
	}
	return value
}

func buildValueF(n float64, mType MetricType, opts ...ValueOption) *Value {
	value := &Value{
		mType:  mType,
		valuef: &n,
	}
	for _, opt := range opts {
		opt(value)
	}
	return value
}

func (v *Value) GetType() MetricType {
	return v.mType
}

func (v *Value) GetSuffix() string {
	return v.suffix
}

func (v *Value) GetValue() *int64 {
	return v.value
}

func (v *Value) GetValueF() *float64 {
	return v.valuef
}
