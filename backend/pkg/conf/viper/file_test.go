// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package viper

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/cozeloop/backend/pkg/conf"
)

type TestConfig struct {
	TestStruct *TestStruct          `mapstructure:"test_struct"`
	TestString string               `mapstructure:"test_string"`
	TaskBool   bool                 `mapstructure:"test_bool"`
	TestMap1   map[string]string    `mapstructure:"test_map_1"`
	TestMap2   map[int64]TestStruct `mapstructure:"test_map_2"`
}

type TestStruct struct {
	FieldI64    int64  `mapstructure:"field_i64"`
	FieldString string `mapstructure:"field_string"`
	FieldBool   bool   `mapstructure:"field_bool"`
	FieldInt    int    `mapstructure:"field_int"`
}

type TestTagConfig struct {
	TestStruct *TestStruct          `json:"test_struct"`
	TestString string               `json:"test_string"`
	TaskBool   bool                 `json:"test_bool"`
	TestMap1   map[string]string    `json:"test_map_1"`
	TestMap2   map[int64]TestStruct `json:"test_map_2"`
}

type TestTagStruct struct {
	FieldI64    int64  `json:"field_i64"`
	FieldString string `json:"field_string"`
	FieldBool   bool   `json:"field_bool"`
	FieldInt    int    `json:"field_int"`
}

func TestFileConfLoader(t *testing.T) {
	tests := []struct {
		name          string
		configFile    string
		configPath    string
		searchPathDir bool
		expectedError bool
	}{
		{
			name:          "current directory",
			configFile:    "config.test.yaml",
			configPath:    os.Getenv("PWD"),
			searchPathDir: false,
			expectedError: false,
		},
		{
			name:          "parent directory",
			configFile:    "config.test.yaml",
			configPath:    os.Getenv("PWD") + "/../",
			searchPathDir: true,
			expectedError: false,
		},
		{
			name:          "sub directory",
			configFile:    "config.sub.test.yaml",
			configPath:    "",
			searchPathDir: true,
			expectedError: false,
		},
		{
			name:          "sub directory2",
			configFile:    "/sub/config.sub.test.yaml",
			configPath:    "",
			searchPathDir: true,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			var loader conf.IConfigLoader
			var err error
			var opts []FileConfLoaderOpt
			opts = append(opts, WithSearchPathDir(tt.searchPathDir))
			if tt.configPath != "" {
				opts = append(opts, WithConfigPath(tt.configPath))
			}

			loader, err = NewFileConfLoader(tt.configFile, opts...)

			if tt.expectedError {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)

			ts := loader.Get(ctx, "test_struct")
			assert.NotNil(t, ts)

			tm2 := make(map[int64]TestStruct)
			err = loader.UnmarshalKey(ctx, "test_map_2", &tm2)
			assert.Nil(t, err)
			assert.NotNil(t, tm2)
			assert.Equal(t, 2, len(tm2))

			ttm2 := make(map[int64]TestTagStruct)
			err = loader.UnmarshalKey(ctx, "test_map_2", &ttm2, conf.WithTagName("json"))
			assert.Nil(t, err)
			assert.NotNil(t, ttm2)
			assert.Equal(t, 2, len(ttm2))

			config := &TestConfig{}
			err = loader.Unmarshal(ctx, config)
			assert.Nil(t, err)
			assert.NotNil(t, config)
			assert.Equal(t, 2, len(config.TestMap2))

			tconfig := &TestTagConfig{}
			err = loader.Unmarshal(ctx, tconfig, conf.WithTagName("json"))
			assert.Nil(t, err)
			assert.NotNil(t, tconfig)
			assert.Equal(t, 2, len(tconfig.TestMap2))
		})
	}
}
