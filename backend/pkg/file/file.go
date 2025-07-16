// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func FindSubDir(root, subDir string) (string, error) {
	var found string
	subDir = filepath.Clean(subDir)
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != root {
			rel, relErr := filepath.Rel(root, path)
			if relErr != nil {
				return relErr
			}
			rel = filepath.Clean(rel)
			if rel == subDir || strings.HasSuffix(rel, string(os.PathSeparator)+subDir) {
				found = path
				return filepath.SkipDir
			}
		}
		return nil
	})
	if err != nil && !errors.Is(err, filepath.SkipDir) {
		return "", err
	}
	if found == "" {
		return "", os.ErrNotExist
	}
	return found, nil
}
