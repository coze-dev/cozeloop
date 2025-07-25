// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDeepCopy(t *testing.T) {
	a := A{
		NumPtr: func() *int64 {
			a := int64(9)
			return &a
		}(),
		NumLen: 10,
		StrPtr: func() *string {
			a := "nihao"
			return &a
		}(),
		Str: "Hello",
		BPtr: &B{
			NumPtr: func() *int64 {
				a := int64(9)
				return &a
			}(),
			NumLen: 20,
			StrPtr: func() *string {
				a := "nihao"
				return &a
			}(),
			Str: "World",
		},
		B: B{
			NumPtr: func() *int64 {
				b := int64(8)
				return &b
			}(),
			NumLen: 15,
			StrPtr: func() *string {
				b := "world"
				return &b
			}(),
			Str: "Golang",
		},
	}

	cp_a := A{}
	_ = DeepCopy(&a, &cp_a)
	fmt.Printf("a is:\n%s\na copy is:\n%s\n", Jsonfy(a), Jsonfy(cp_a))

	fmt.Printf("a = %p, cp_a = %p, eq = %t\n", &a, &cp_a, &a == &cp_a)
	fmt.Printf("&a.NumPtr = %p, &cp_a.NumPtr = %p, eq = %t\n", &a.NumPtr, &cp_a.NumPtr, &a.NumPtr == &cp_a.NumPtr)
	fmt.Printf("a.NumPtr = %p, cp_a.NumPtr = %p, eq = %t\n", a.NumPtr, cp_a.NumPtr, a.NumPtr == cp_a.NumPtr)
	fmt.Printf("&a.StrPtr = %p, &cp_a.StrPtr = %p, eq = %t\n", &a.StrPtr, &cp_a.StrPtr, &a.StrPtr == &cp_a.StrPtr)
	fmt.Printf("a.StrPtr = %p, cp_a.StrPtr = %p, eq = %t\n", a.StrPtr, cp_a.StrPtr, a.StrPtr == cp_a.StrPtr)
	fmt.Printf("&a.NumLen = %p, &cp_a.NumLen = %p, eq = %t\n", &a.NumLen, &cp_a.NumLen, &a.NumLen == &cp_a.NumLen)
	fmt.Printf("&a.Str = %p, &cp_a.Str = %p, eq = %t\n", &a.Str, &cp_a.Str, &a.Str == &cp_a.Str)

	fmt.Printf("&a.BPtr = %p, &cp_a.BPtr = %p, eq = %t\n", &a.BPtr, &cp_a.BPtr, &a.BPtr == &cp_a.BPtr)
	fmt.Printf("a.BPtr = %p, cp_a.BPtr = %p, eq = %t\n", a.BPtr, cp_a.BPtr, a.BPtr == cp_a.BPtr)
	fmt.Printf("&a.BPtr.NumPtr = %p, &cp_a.BPtr.NumPtr = %p, eq = %t\n", &a.BPtr.NumPtr, &cp_a.BPtr.NumPtr, &a.BPtr.NumPtr == &cp_a.BPtr.NumPtr)
	fmt.Printf("a.BPtr.NumPtr = %p, cp_a.BPtr.NumPtr = %p, eq = %t\n", a.BPtr.NumPtr, cp_a.BPtr.NumPtr, a.BPtr.NumPtr == cp_a.BPtr.NumPtr)
	fmt.Printf("&a.BPtr.StrPtr = %p, &cp_a.BPtr.StrPtr = %p, eq = %t\n", &a.BPtr.StrPtr, &cp_a.BPtr.StrPtr, &a.BPtr.StrPtr == &cp_a.BPtr.StrPtr)
	fmt.Printf("a.BPtr.StrPtr = %p, cp_a.BPtr.StrPtr = %p, eq = %t\n", a.BPtr.StrPtr, cp_a.BPtr.StrPtr, a.BPtr.StrPtr == cp_a.BPtr.StrPtr)
	fmt.Printf("&a.BPtr.NumLen = %p, &cp_a.BPtr.NumLen = %p, eq = %t\n", &a.BPtr.NumLen, &cp_a.BPtr.NumLen, &a.BPtr.NumLen == &cp_a.BPtr.NumLen)
	fmt.Printf("&a.BPtr.Str = %p, &cp_a.BPtr.Str = %p, eq = %t\n", &a.BPtr.Str, &cp_a.BPtr.Str, &a.BPtr.Str == &cp_a.BPtr.Str)

	fmt.Printf("&a.B = %p, &cp_a.B = %p, eq = %t\n", &a.B, &cp_a.B, &a.B == &cp_a.B)
	fmt.Printf("&a.B.NumPtr = %p, &cp_a.B.NumPtr = %p, eq = %t\n", &a.B.NumPtr, &cp_a.B.NumPtr, &a.B.NumPtr == &cp_a.B.NumPtr)
	fmt.Printf("a.B.NumPtr = %p, cp_a.B.NumPtr = %p, eq = %t\n", a.B.NumPtr, cp_a.B.NumPtr, a.B.NumPtr == cp_a.B.NumPtr)
	fmt.Printf("&a.B.StrPtr = %p, &cp_a.B.StrPtr = %p, eq = %t\n", &a.B.StrPtr, &cp_a.B.StrPtr, &a.B.StrPtr == &cp_a.B.StrPtr)
	fmt.Printf("a.B.StrPtr = %p, cp_a.B.StrPtr = %p, eq = %t\n", a.B.StrPtr, cp_a.B.StrPtr, a.B.StrPtr == cp_a.B.StrPtr)
	fmt.Printf("&a.B.NumLen = %p, &cp_a.B.NumLen = %p, eq = %t\n", &a.B.NumLen, &cp_a.B.NumLen, &a.B.NumLen == &cp_a.B.NumLen)
	fmt.Printf("&a.B.Str = %p, &cp_a.B.Str = %p, eq = %t\n", &a.B.Str, &cp_a.B.Str, &a.B.Str == &cp_a.B.Str)
}

func Jsonfy(i interface{}) string {
	js, _ := json.MarshalIndent(i, "", "  ")
	return string(js)
}

type A struct {
	NumPtr *int64
	NumLen int
	StrPtr *string
	Str    string

	BPtr *B
	B    B
}

type B struct {
	NumPtr *int64
	NumLen int
	StrPtr *string
	Str    string
}
