// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package js_conv

import (
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/evaluation/expt"
)

func BenchmarkJsonIterExtension_UpdateStructDescriptor(b *testing.B) {
	j := jsoniter.ConfigCompatibleWithStandardLibrary
	j.RegisterExtension(&JsonIterExtension{})

	type MyStruct struct {
		M1 map[int64]int64
		M2 map[string]int64
		M3 map[int64]string `json:"m3,omitempty"`
		M4 map[int64]string `json:"m4,omitempty"`
	}

	type WrappedStruct struct {
		M1 map[int64]MyStruct
	}

	obj := MyStruct{
		M1: map[int64]int64{123: 456},
		M2: map[string]int64{"key": 789},
		M3: map[int64]string{456: "value"},
		M4: map[int64]string{},
	}

	for i := 0; i < b.N; i++ {
		data, _ := j.MarshalToString(obj)

		var newObj MyStruct
		_ = j.UnmarshalFromString(data, &newObj)

		data, _ = j.MarshalToString(WrappedStruct{
			M1: map[int64]MyStruct{123: obj},
		})

		var newWObj WrappedStruct
		_ = j.UnmarshalFromString(data, &newWObj)
	}
}

func Test_jsonIterExtension_UpdateStructDescriptor(t *testing.T) {
	j := jsoniter.ConfigCompatibleWithStandardLibrary
	j.RegisterExtension(&JsonIterExtension{})

	t.Run("map", func(t *testing.T) {
		type MyStruct struct {
			M1 map[int64]int64
			M2 map[string]int64
			M3 map[int64]string `json:"m3,omitempty"`
			M4 map[int64]string `json:"m4,omitempty"`
		}

		obj := MyStruct{
			M1: map[int64]int64{123: 456},
			M2: map[string]int64{"key": 789},
			M3: map[int64]string{456: "value"},
			M4: map[int64]string{},
		}
		data, _ := j.MarshalToString(obj) // {"M1":{"123":"456"},"M2":{"key":"789"},"m3":{"456":"value"}}

		var newObj MyStruct
		err := j.UnmarshalFromString(data, &newObj)
		assert.NoError(t, err)
		fmt.Printf("Deserialized MyStruct M1: %v\n", newObj.M1)
		fmt.Printf("Deserialized MyStruct M2: %v\n", newObj.M2)
		fmt.Printf("Deserialized MyStruct M3: %v\n", newObj.M3)

		type WrappedStruct struct {
			M1 map[int64]MyStruct
		}

		wobj := WrappedStruct{
			M1: map[int64]MyStruct{123: obj},
		}
		data, _ = j.MarshalToString(wobj)
		fmt.Println("Serialized WrappedStruct:", data)

		var newWObj WrappedStruct
		err = j.UnmarshalFromString(data, &newWObj)
		assert.NoError(t, err)
		fmt.Printf("Deserialized WrappedStruct M1: %v\n", newWObj.M1)
	})

	t.Run("slice", func(t *testing.T) {
		type MyStruct struct {
			IDs  []int64  `json:"ids"`
			SIDs []string `json:"sids"`
			IIDs []int64  `json:"iids"`
			IVal int64    `json:"ival"`
			AVal int64    `json:"aval"`
			IPtr *int64   `json:"iptr"`
		}

		raw := `{"ids":["4","5","6"],"sids":["1","2","3"],"iids":[1,2,3],"ival":"1","nkey":"1","iptr":"1"}`

		var s MyStruct
		err := j.Unmarshal([]byte(raw), &s)
		assert.NoError(t, err)

		gb, err := j.Marshal(s)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"ids":["4","5","6"],"sids":["1","2","3"],"iids":["1","2","3"],"ival":"1","aval":"0","iptr":"1"}`, string(gb))
	})
}

func TestSerializeEvaluatorRecord(t *testing.T) {
	type args struct {
		raw  string
		dest any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "hybrid",
			args: args{
				raw:  "{\n    \"code\": 0,\n    \"column_eval_set_fields\":\n    [\n        {\n            \"content_type\": \"Text\",\n            \"description\": \"作为输入投递给评测对象\",\n            \"key\": \"input\",\n            \"name\": \"input\"\n        },\n        {\n            \"content_type\": \"Text\",\n            \"description\": \"预期理想输出，可作为评估时的参考标准\",\n            \"key\": \"reference_output\",\n            \"name\": \"reference_output\"\n        }\n    ],\n    \"column_evaluators\":\n    [\n        {\n            \"description\": \"\",\n            \"evaluator_id\": \"7506423405453770752\",\n            \"evaluator_type\": 1,\n            \"evaluator_version_id\": \"7506423405453787136\",\n            \"name\": \"测试评估器\",\n            \"version\": \"0.0.1\"\n        }\n    ],\n    \"item_results\":\n    [\n        {\n            \"item_id\": \"0\",\n            \"item_index\": \"0\",\n            \"system_info\":\n            {\n                \"error\":\n                {\n                    \"code\": \"0\"\n                },\n                \"run_state\": 2\n            },\n            \"turn_results\":\n            [\n                {\n                    \"experiment_results\":\n                    [\n                        {\n                            \"experiment_id\": \"7506735326975492098\",\n                            \"payload\":\n                            {\n                                \"eval_set\":\n                                {\n                                    \"turn\":\n                                    {\n                                        \"field_data_list\":\n                                        [\n                                            {\n                                                \"content\":\n                                                {\n                                                    \"content_type\": \"Text\",\n                                                    \"format\": 1,\n                                                    \"text\": \"test\"\n                                                },\n                                                \"key\": \"input\",\n                                                \"name\": \"input\"\n                                            },\n                                            {\n                                                \"content\":\n                                                {\n                                                    \"content_type\": \"Text\",\n                                                    \"format\": 1,\n                                                    \"text\": \"test\"\n                                                },\n                                                \"key\": \"reference_output\",\n                                                \"name\": \"reference_output\"\n                                            }\n                                        ],\n                                        \"id\": \"0\"\n                                    }\n                                },\n                                \"evaluator_output\":\n                                {\n                                    \"evaluator_records\":\n                                    {\n                                        \"7506423405453787136\":\n                                        {\n                                            \"base_info\":\n                                            {\n                                                \"created_at\": 1747798124000,\n                                                \"created_by\":\n                                                {\n                                                    \"user_id\": \"\"\n                                                },\n                                                \"deleted_at\": null,\n                                                \"updated_at\": 1747798124000,\n                                                \"updated_by\":\n                                                {\n                                                    \"user_id\": \"\"\n                                                }\n                                            },\n                                            \"evaluator_input_data\":\n                                            {\n                                                \"input_fields\":\n                                                {\n                                                    \"USER_NAME\":\n                                                    {\n                                                        \"content_type\": \"Text\",\n                                                        \"format\": 2,\n                                                        \"text\": \"mock text\"\n                                                    },\n                                                    \"input\":\n                                                    {\n                                                        \"content_type\": \"Text\",\n                                                        \"format\": 1,\n                                                        \"text\": \"test\"\n                                                    },\n                                                    \"output\":\n                                                    {\n                                                        \"content_type\": \"Text\",\n                                                        \"format\": 1,\n                                                        \"text\": \"test\"\n                                                    }\n                                                }\n                                            },\n                                            \"evaluator_output_data\":\n                                            {\n                                                \"evaluator_result\":\n                                                {\n                                                    \"reasoning\": \"模型的输出与输入相同，没有包含任何额外的不必要信息，完全符合完美简洁答案的要求。因此，应该给出的分数是1.0\",\n                                                    \"score\": 1\n                                                },\n                                                \"evaluator_usage\":\n                                                {\n                                                    \"input_tokens\": 586,\n                                                    \"output_tokens\": 113\n                                                },\n                                                \"time_consuming_ms\": 0\n                                            },\n                                            \"evaluator_version_id\": 7506423405453787000,\n                                            \"experiment_id\": 7506735326975492000,\n                                            \"experiment_run_id\": 7506735327076155000,\n                                            \"id\": 7506735782262997000,\n                                            \"item_id\": 0,\n                                            \"status\": 1,\n                                            \"trace_id\": \"2b5ab553de528293195c0edd37fbb57b\",\n                                            \"turn_id\": 0\n                                        }\n                                    }\n                                },\n                                \"system_info\":\n                                {\n                                    \"error\":\n                                    {\n                                        \"code\": \"0\"\n                                    },\n                                    \"log_id\": \"\",\n                                    \"turn_run_state\": 1\n                                },\n                                \"target_output\":\n                                {\n                                    \"eval_target_record\":\n                                    {\n                                        \"base_info\":\n                                        {\n                                            \"created_at\": \"1747798122000\",\n                                            \"deleted_at\": null,\n                                            \"updated_at\": \"1747798122000\"\n                                        },\n                                        \"eval_target_input_data\":\n                                        {},\n                                        \"eval_target_output_data\":\n                                        {\n                                            \"eval_target_usage\":\n                                            {\n                                                \"input_tokens\": \"15\",\n                                                \"output_tokens\": \"2054\"\n                                            },\n                                            \"output_fields\":\n                                            {\n                                                \"actual_output\":\n                                                {\n                                                    \"content_type\": \"Text\",\n                                                    \"format\": 2,\n                                                    \"text\": \"mock test\"\n                                                }\n                                            },\n                                            \"time_consuming_ms\": \"65451\"\n                                        },\n                                        \"experiment_run_id\": \"7506735327076155392\",\n                                        \"id\": \"7506735772620292098\",\n                                        \"item_id\": \"0\",\n                                        \"status\": 1,\n                                        \"target_id\": \"7506472875889524736\",\n                                        \"target_version_id\": \"7506731336871198722\",\n                                        \"trace_id\": \"66d830168b45f8a322520ce708f8ba6e\",\n                                        \"turn_id\": \"0\",\n                                        \"workspace_id\": \"7480438112481656876\"\n                                    }\n                                },\n                                \"turn_id\": \"0\"\n                            }\n                        }\n                    ],\n                    \"turn_id\": \"0\",\n                    \"turn_index\": \"0\"\n                }\n            ]\n        }\n    ],\n    \"msg\": \"\",\n    \"total\": \"1\"\n}",
				dest: &expt.BatchGetExperimentResultResponse{},
			},
		},
		{
			name: "js string",
			args: args{
				raw:  "{\"column_eval_set_fields\":[{\"key\":\"input\",\"name\":\"input\",\"description\":\"作为输入投递给评测对象\",\"content_type\":\"Text\"},{\"key\":\"reference_output\",\"name\":\"reference_output\",\"description\":\"预期理想输出，可作为评估时的参考标准\",\"content_type\":\"Text\"}],\"column_evaluators\":[{\"evaluator_version_id\":\"7506423405453787136\",\"evaluator_id\":\"7506423405453770752\",\"evaluator_type\":1,\"name\":\"测试评估器\",\"version\":\"0.0.1\",\"description\":\"\"}],\"item_results\":[{\"item_id\":\"0\",\"turn_results\":[{\"turn_id\":\"0\",\"experiment_results\":[{\"experiment_id\":\"7506735326975492098\",\"payload\":{\"turn_id\":\"0\",\"eval_set\":{\"turn\":{\"id\":\"0\",\"field_data_list\":[{\"key\":\"input\",\"name\":\"input\",\"content\":{\"content_type\":\"Text\",\"format\":1,\"text\":\"test\"}},{\"key\":\"reference_output\",\"name\":\"reference_output\",\"content\":{\"content_type\":\"Text\",\"format\":1,\"text\":\"test\"}}]}},\"target_output\":{\"eval_target_record\":{\"id\":\"7506735772620292098\",\"workspace_id\":\"7480438112481656876\",\"target_id\":\"7506472875889524736\",\"target_version_id\":\"7506731336871198722\",\"experiment_run_id\":\"7506735327076155392\",\"item_id\":\"0\",\"turn_id\":\"0\",\"trace_id\":\"66d830168b45f8a322520ce708f8ba6e\",\"eval_target_input_data\":{},\"eval_target_output_data\":{\"output_fields\":{\"actual_output\":{\"content_type\":\"Text\",\"format\":2,\"text\":\"mock test\"}},\"eval_target_usage\":{\"input_tokens\":\"15\",\"output_tokens\":\"2054\"},\"time_consuming_ms\":\"65451\"},\"status\":1,\"base_info\":{\"created_at\":\"1747798122000\",\"updated_at\":\"1747798122000\",\"deleted_at\":null}}},\"evaluator_output\":{\"evaluator_records\":{\"7506423405453787136\":{\"id\":\"7506735782262997000\",\"experiment_id\":\"7506735326975492000\",\"experiment_run_id\":\"7506735327076155000\",\"item_id\":\"0\",\"turn_id\":\"0\",\"evaluator_version_id\":\"7506423405453787000\",\"trace_id\":\"2b5ab553de528293195c0edd37fbb57b\",\"evaluator_input_data\":{\"input_fields\":{\"USER_NAME\":{\"content_type\":\"Text\",\"format\":2,\"text\":\"mock text\"},\"input\":{\"content_type\":\"Text\",\"format\":1,\"text\":\"test\"},\"output\":{\"content_type\":\"Text\",\"format\":1,\"text\":\"test\"}}},\"evaluator_output_data\":{\"evaluator_result\":{\"score\":1,\"reasoning\":\"模型的输出与输入相同，没有包含任何额外的不必要信息，完全符合完美简洁答案的要求。因此，应该给出的分数是1.0\"},\"evaluator_usage\":{\"input_tokens\":\"586\",\"output_tokens\":\"113\"},\"time_consuming_ms\":\"0\"},\"status\":1,\"base_info\":{\"created_by\":{\"user_id\":\"\"},\"updated_by\":{\"user_id\":\"\"},\"created_at\":\"1747798124000\",\"updated_at\":\"1747798124000\",\"deleted_at\":null}}}},\"system_info\":{\"turn_run_state\":1,\"log_id\":\"\",\"error\":{\"code\":\"0\"}}}}],\"turn_index\":\"0\"}],\"system_info\":{\"run_state\":2,\"error\":{\"code\":\"0\"}},\"item_index\":\"0\"}],\"total\":\"1\"}",
				dest: &expt.BatchGetExperimentResultResponse{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := jsoniter.ConfigCompatibleWithStandardLibrary
			j.RegisterExtension(NewJSONIterExtension())

			err := j.Unmarshal([]byte(tt.args.raw), tt.args.dest)
			assert.NoError(t, err)

			bytes, err := j.Marshal(tt.args.dest)
			assert.NoError(t, err)
			assert.True(t, len(bytes) > 0)
		})
	}
}
