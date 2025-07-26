// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
)

func TestScenarioValue(t *testing.T) {
	type args struct {
		scenario *Scenario
	}
	tests := []struct {
		name string
		args args
		want Scenario
	}{
		{
			name: "scenario nil",
			args: args{
				scenario: nil,
			},
			want: ScenarioDefault,
		},
		{
			name: "scenario prompt debug",
			args: args{
				scenario: ptr.Of(ScenarioPromptDebug),
			},
			want: ScenarioPromptDebug,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ScenarioValue(tt.args.scenario), "ScenarioValue(%v)", tt.args.scenario)
		})
	}
}
