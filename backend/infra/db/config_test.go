// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_buildDSN(t *testing.T) {
	defaultArgs := `?charset=utf8mb4&parseTime=True&loc=&timeout=0s&readTimeout=0s&writeTimeout=0s&interpolateParams=false`

	for _, tt := range []struct {
		name string
		cfg  *Config
		want string
	}{
		{
			name: "tcp",
			cfg: &Config{
				User:       "user",
				Password:   "pass",
				DBHostname: "localhost",
				DBPort:     "3306",
				DBName:     "test",
			},
			want: "user:pass@tcp(localhost:3306)/test" + defaultArgs,
		},
		{
			name: "with params",
			cfg: &Config{
				User:       "user",
				Password:   "pass",
				DBHostname: "localhost",
				DBPort:     "3306",
				DBName:     "test",
				DSNParams:  url.Values{"multiStatements": []string{"true"}},
			},
			want: "user:pass@tcp(localhost:3306)/test" + defaultArgs + "&multiStatements=true",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.buildDSN()
			assert.Equal(t, tt.want, got)
		})
	}
}
