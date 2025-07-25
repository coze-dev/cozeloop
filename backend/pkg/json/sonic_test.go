// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJson(t *testing.T) {
	v := map[string]interface{}{
		"name": "evaluator",
		"age":  123,
	}

	b, err := Marshal(v)
	assert.Nil(t, err)

	tar := map[string]interface{}{}
	err = Unmarshal(b, &tar)
	assert.Nil(t, err)
	assert.Equal(t, "evaluator", tar["name"])

	r := bytes.NewReader(b)

	tar2 := map[string]interface{}{}
	err = Decode(r, &tar2)
	assert.Nil(t, err)
	assert.Equal(t, "evaluator", tar2["name"])

	ok := Valid(b)
	assert.True(t, ok)
}

func TestJsonify(t *testing.T) {
	v := map[string]interface{}{
		"name": "evaluator",
		"age":  123,
	}
	res := Jsonify(v)

	tar := map[string]interface{}{}
	err := Unmarshal([]byte(res), &tar)
	assert.Nil(t, err)
	assert.Equal(t, "evaluator", tar["name"])
}
