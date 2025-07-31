// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/coze-dev/coze-loop/backend/api/handler/coze/loop/apis"
	router "github.com/coze-dev/coze-loop/backend/api/router"
)

// register registers all routers.
func register(r *server.Hertz, handler *apis.APIHandler) {
	router.GeneratedRegister(r, handler)

	customizedRegister(r)
}
