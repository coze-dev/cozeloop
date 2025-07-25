// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package experiment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
	"github.com/coze-dev/cozeloop/backend/pkg/json"
)

func TestEvalConfConvert_ConvertEntityToDTO(t *testing.T) {
	raw := `{
    "ConnectorConf":
    {
        "TargetConf":
        {
            "TargetVersionID": 7486074365205872641,
            "IngressConf":
            {
                "EvalSetAdapter":
                {
                    "FieldConfs":
                    [
                        {
                            "FieldName": "role",
                            "FromField": "role",
                            "Value": ""
                        },
                        {
                            "FieldName": "question",
                            "FromField": "input",
                            "Value": ""
                        }
                    ]
                },
                "CustomConf": null
            }
        },
        "EvaluatorsConf":
        {
            "EvaluatorConcurNum": null,
            "EvaluatorConf":
            [
                {
                    "EvaluatorVersionID": 7486074365205823489,
                    "IngressConf":
                    {
                        "EvalSetAdapter":
                        {
                            "FieldConfs":
                            [
                                {
                                    "FieldName": "input",
                                    "FromField": "input",
                                    "Value": ""
                                },
                                {
                                    "FieldName": "reference_output",
                                    "FromField": "reference_output",
                                    "Value": ""
                                }
                            ]
                        },
                        "TargetAdapter":
                        {
                            "FieldConfs":
                            [
                                {
                                    "FieldName": "output",
                                    "FromField": "actual_output",
                                    "Value": ""
                                }
                            ]
                        },
                        "CustomConf": null
                    }
                }
            ]
        }
    },
    "ItemConcurNum": null
}`
	conf := &entity.EvaluationConfiguration{}
	err := json.Unmarshal([]byte(raw), &conf)
	assert.Nil(t, err)

	target, evaluators := NewEvalConfConvert().ConvertEntityToDTO(conf)
	t.Logf("target: %v", json.Jsonify(target))
	t.Logf("evaluators: %v", json.Jsonify(evaluators))
}
