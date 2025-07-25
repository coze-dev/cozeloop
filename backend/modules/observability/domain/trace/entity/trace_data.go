// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/loop_span"
)

type TTL string

const (
	TTL3d   TTL = "3d"
	TTL7d   TTL = "7d"
	TTL30d  TTL = "30d"
	TTL90d  TTL = "90d"
	TTL180d TTL = "180d"
	TTL365d TTL = "365d"
)

type TenantInfo struct {
	TTL              TTL            `json:"ttl"`
	WorkspaceId      string         `json:"space_id"`
	CozeAccountID    string         `json:"coze_account_id"`
	WhichIsEnough    int            `json:"which_is_enough"`
	VolcanoAccountID int64          `json:"volcano_account_id"`
	Extra            map[string]any `json:"extra,omitempty"`
}

type TraceData struct {
	Tenant     string             `json:"tenant"`
	TenantInfo TenantInfo         `json:"tenant_info"`
	SpanList   loop_span.SpanList `json:"span_list"`
}

func TTLFromInteger(i int64) TTL {
	switch i {
	case 3:
		return TTL3d
	case 7:
		return TTL7d
	case 30:
		return TTL30d
	case 90:
		return TTL90d
	case 180:
		return TTL180d
	case 360:
		return TTL365d
	default:
		return TTL3d
	}
}
