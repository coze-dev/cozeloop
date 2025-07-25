// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Timeout           time.Duration `yaml:"timeout"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	User              string        `yaml:"user"`
	Password          string        `yaml:"password"`
	Loc               string        `yaml:"loc"`
	DBName            string        `yaml:"db_name"`
	DBCharset         string        `yaml:"db_charset"`
	DBHostname        string        `yaml:"db_hostname"`
	DBPort            string        `yaml:"db_port"`
	InterpolateParams bool          `yaml:"interpolate_params"`
	DSNParams         url.Values    `yaml:"dsn_params"`
	WithReturning     bool          `yaml:"with_returning"`

	// TODO: support read replica
}

func (cfg *Config) buildDSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.User, cfg.Password, cfg.DBHostname, cfg.DBPort, cfg.DBName)

	charset := cfg.DBCharset
	if charset == "" {
		charset = "utf8mb4"
	}
	args := []string{
		"charset=" + charset,
		"parseTime=True",
		"loc=" + cfg.Loc,
		"timeout=" + cfg.Timeout.String(),
		"readTimeout=" + cfg.ReadTimeout.String(),
		"writeTimeout=" + cfg.WriteTimeout.String(),
		"interpolateParams=" + strconv.FormatBool(cfg.InterpolateParams),
	}

	for key := range cfg.DSNParams {
		args = append(args, key+"="+cfg.DSNParams.Get(key))
	}

	return fmt.Sprintf("%s?%s", dsn, strings.Join(args, "&"))
}
