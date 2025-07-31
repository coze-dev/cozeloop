// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package viper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/viper"

	"github.com/coze-dev/coze-loop/backend/pkg/conf"
	"github.com/coze-dev/coze-loop/backend/pkg/logs"
)

var errStopWalk = fmt.Errorf("stop walking")

func NewFileConfLoader(file string, opts ...FileConfLoaderOpt) (conf.IConfigLoader, error) {
	opt := &fileConfLoaderOpt{}
	for _, fn := range opts {
		fn(opt)
	}

	var (
		abs  = ""
		fb   = filepath.Base(file)
		path = lo.Ternary(len(opt.path) > 0, opt.path, os.Getenv("PWD"))
		err  error
	)

	if opt.search && file == fb {
		if err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && d.Name() == fb {
				abs, err = filepath.Abs(filepath.Join(path, abs))
				if err != nil {
					return err
				}
				return errStopWalk
			}
			return nil
		}); err != nil && !errors.Is(err, errStopWalk) {
			return nil, err
		}

		if abs == "" {
			return nil, fmt.Errorf("file %s not found", file)
		}
	} else {
		abs, err = filepath.Abs(filepath.Join(path, file))
		if err != nil {
			return nil, err
		}
	}

	v := viper.New()
	v.SetConfigFile(abs)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read file config fail, file=%s, err=%v", path+file, err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		logs.Info("on changed file: %s", e.Name)
	})

	return &fileConfLoader{
		v: v,
	}, nil
}

func NewFileConfigLoaderFactory(opts ...FileConfLoaderFactoryOpt) conf.IConfigLoaderFactory {
	opt := &fileConfLoaderFactoryOpt{}
	for _, fn := range opts {
		fn(opt)
	}
	return &fileConfLoaderFactory{
		configPath: opt.path,
	}
}
