// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"fmt"
	"go/token"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/samber/lo"
)

type ImportRewriteFunc func(importPath string) string

type GoIdent struct {
	GoName       string
	GoImportPath string
}

// GoImportHelper helps build go file imports.
type GoImportHelper struct {
	importRewriteFunc ImportRewriteFunc
	packageNames      map[string] /*GoImportPath*/ string    /*GoPackageName*/
	pkgNames          map[string] /*GoPackageName*/ struct{} // pkg name set
}

func NewGoImportHelper(fn ImportRewriteFunc) *GoImportHelper {
	return &GoImportHelper{
		importRewriteFunc: fn,
		packageNames:      make(map[string]string),
		pkgNames:          make(map[string]struct{}),
	}
}

func (g *GoImportHelper) QualifiedGoIdent(ident GoIdent) string {
	if packageName, ok := g.packageNames[ident.GoImportPath]; ok {
		return string(packageName) + "." + ident.GoName
	}

	base := path.Base(ident.GoImportPath)
	name := g.addPkgName(base)

	g.packageNames[ident.GoImportPath] = name
	return string(name) + "." + ident.GoName
}

func (g *GoImportHelper) QualifiedGoIdentWithAlias(ident GoIdent, alias string) string {
	if packageName, ok := g.packageNames[ident.GoImportPath]; ok {
		// should check packageName != alias?
		return string(packageName) + "." + ident.GoName
	}

	name := g.addPkgName(alias)
	g.packageNames[ident.GoImportPath] = name
	return string(name) + "." + ident.GoName
}

func (g *GoImportHelper) AddImportAlias(importPath string, alias string) {
	if pkg, ok := g.packageNames[importPath]; ok && alias != pkg {
		panic(fmt.Sprintf("package name of %s already set to %s", importPath, pkg))
	}

	name := g.addPkgName(alias)
	g.packageNames[importPath] = name
}

func (g *GoImportHelper) AddImport(importPath string) (pkg string) {
	if p, ok := g.packageNames[importPath]; ok {
		return p
	}

	base := path.Base(importPath)
	name := g.addPkgName(base)
	g.packageNames[importPath] = name
	return name
}

// Imports returns all import statements(quoted).
func (g *GoImportHelper) Imports() []string {
	imports := make([]string, 0, len(g.packageNames))

	m := g.packageNames
	if fn := g.importRewriteFunc; fn != nil {
		// rewrite import path
		m = make(map[string]string, len(g.packageNames))
		for pkg, name := range g.packageNames {
			pkg = fn(pkg)
			m[pkg] = name
		}
	}

	pkgs := lo.Keys(m)
	sort.Strings(pkgs)
	for _, p := range pkgs {
		n := m[p]
		if filepath.Base(p) == n {
			imports = append(imports, strconv.Quote(p))
		} else {
			imports = append(imports, n+" "+strconv.Quote(p))
		}
	}
	return imports
}

func (g *GoImportHelper) addPkgName(name string) string {
	n := goSanitized(name)

	suffix := 0
	for {
		if _, ok := g.pkgNames[name]; !ok {
			break
		}
		suffix++
		name = fmt.Sprintf("%s%d", n, suffix)
	}

	g.pkgNames[name] = struct{}{}
	return name
}

// Copyright 2019 The Go Authors. All rights reserved.
//
// Code forked from https://github.com/protocolbuffers/protobuf-go/blob/master/internal/strs/strings.go
//
// Original code was released under BSD 3-Clause, with the full license text
// available at https://github.com/protocolbuffers/protobuf-go/blob/master/LICENSE.
func goSanitized(s string) string {
	// Sanitize the input to the set of valid characters,
	// which must be '_' or be in the Unicode L or N categories.
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '_'
	}, s)

	// Prepend '_' in the event of a Go keyword conflict or if
	// the identifier is invalid (does not start in the Unicode L category).
	r, _ := utf8.DecodeRuneInString(s)
	if token.Lookup(s).IsKeyword() || !unicode.IsLetter(r) {
		return "_" + s
	}
	return s
}
