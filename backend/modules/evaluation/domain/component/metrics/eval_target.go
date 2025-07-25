// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package metrics

import "time"

//go:generate mockgen -destination=mocks/eval_target.go -package=mocks . EvalTargetMetrics
type EvalTargetMetrics interface {
	EmitRun(spaceID int64, err error, start time.Time)
	EmitCreate(spaceID int64, err error)
}
