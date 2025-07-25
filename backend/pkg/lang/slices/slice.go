// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package slices

import (
	"github.com/bytedance/gg/gslice"
	"github.com/samber/lo"
)

func ToMap[T any, K comparable, V any](s []T, t func(e T) (K, V)) map[K]V {
	return gslice.ToMap(s, t)
}

func Transform[T any, R any](s []T, iteratee func(e T, idx int) R) []R {
	return lo.Map(s, iteratee)
}

func Uniq[T comparable, Slice ~[]T](s Slice) Slice {
	return gslice.Uniq(s)
}

func Map[F, T any](s []F, f func(F) T) []T {
	return gslice.Map(s, f)
}

func Contains[T comparable](s []T, v T) bool {
	return gslice.Contains(s, v)
}
