// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package metrics

import "time"

//go:generate mockgen -destination=mocks/evaluator.go -package=mocks . EvaluatorExecMetrics
type EvaluatorExecMetrics interface {
	EmitRun(spaceID int64, err error, start time.Time, modelName string)
	EmitCreate(spaceID int64, err error)
}
