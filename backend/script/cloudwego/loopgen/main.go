// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/coze-dev/coze-loop/backend/script/kitex/loopgen/internal"
)

const (
	Version = "v0.0.1"
)

var rootCmd = &cobra.Command{
	Use:     "loopgen",
	Short:   "generate localrpc code",
	Version: Version,
}

var cfg Config

type Config struct {
	GoMod                 string
	KitexPrefix           string
	PackagePrefix         string
	ImportPrefix          string
	IDLDir                string
	IDLFile               string
	ScanIDL               bool
	OutputDir             string
	LocalStreamImportPath string
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`{{.Name}} version {{.Version}}` + "\n")

	fs := rootCmd.Flags()

	fs.StringVarP(&cfg.GoMod, "gomod", "g", "github.com/coze-dev/coze-loop/backend/loop_gen", "go mod")
	fs.StringVarP(&cfg.KitexPrefix, "kitex-prefix", "k", "kitex_gen", "kitex prefix")
	fs.StringVarP(&cfg.PackagePrefix, "package-prefix", "p", "lo", "package prefix")
	fs.StringVarP(&cfg.ImportPrefix, "import-prefix", "i", "loop_gen", "import prefix")
	fs.StringVarP(&cfg.OutputDir, "output-dir", "o", "output", "output dir")
	fs.StringVar(&cfg.LocalStreamImportPath, "local-stream-import-path", "", "local stream import path")
	fs.StringVarP(&cfg.IDLDir, "idl-dir", "d", "", "idl dir")
	fs.StringVarP(&cfg.IDLFile, "idl-file", "f", "", "idl file")
	fs.BoolVarP(&cfg.ScanIDL, "scan-idl", "s", true, "scan idl")

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		idls := []string{cfg.IDLFile}
		if cfg.ScanIDL {
			idls = scanIDLFiles(cfg.IDLDir)
		}

		for _, idl := range idls {
			if idl == "" {
				continue
			}
			g, err := internal.NewGenerator(&internal.ThriftFile{
				IncludeDirs:           []string{cfg.IDLDir},
				IDL:                   idl,
				GoMod:                 cfg.GoMod,
				KitexPrefix:           cfg.KitexPrefix,
				PackagePrefix:         cfg.PackagePrefix,
				ImportPrefix:          cfg.ImportPrefix,
				OutputDir:             cfg.OutputDir,
				LocalStreamImportPath: cfg.LocalStreamImportPath,
			})
			if err != nil {
				log.Printf("[Error] new generator %s failed, err=%+v", idl, err)
				continue
			}
			if err := g.Generate(); err != nil {
				log.Printf("[Error] generate %s failed, err=%+v", idl, err)
				continue
			}
		}
	}
}

func scanIDLFiles(dir string) []string {
	// walk all files in dir recursively
	files := make([]string, 0)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".thrift") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files
}
