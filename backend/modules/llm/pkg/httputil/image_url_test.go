// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package httputil

import (
	"fmt"
	"testing"
)

func TestImageURLToBase64(t *testing.T) {
	b64, mt, err := ImageURLToBase64("https://p6-devops-evaluation-sign.byteimg.com/tos-cn-i-r1igazux8i/8697052347b94717b8cff92b947607ea.jpeg~tplv-r1igazux8i-image.jpeg?rk3s=5ff5f0d4&x-expires=1749018364&x-signature=M71EPRcO%2B7RPV8yXAY8HtiYHuWY%3D")
	fmt.Println(b64, "\n", mt, "\n", err)
}
