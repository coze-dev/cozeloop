// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/coze-dev/cozeloop/backend/api/handler/coze/loop/apis"
	router "github.com/coze-dev/cozeloop/backend/api/router"
)

// register registers all routers.
func register(r *server.Hertz, handler *apis.APIHandler) {
	router.GeneratedRegister(r, handler)

	customizedRegister(r)
}
