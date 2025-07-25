// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package audit

import (
	"context"
)

// NoopAuditServiceImpl 是 IAuditService 接口的模拟实现结构体
type NoopAuditServiceImpl struct{}

func NewNoopAuditService() IAuditService {
	return &NoopAuditServiceImpl{}
}

// Audit 实现 IAuditService 接口的 Audit 方法
// 该方法接收一个上下文和审计参数，返回一个审计记录
func (a *NoopAuditServiceImpl) Audit(ctx context.Context, param AuditParam) (AuditRecord, error) {
	// MockResponse For OpenSource Version
	return AuditRecord{
		ObjectID:     param.ObjectID,
		AuditType:    param.AuditType,
		AuditStatus:  AuditStatus_Approved,
		FailedReason: nil,
		ReqID:        &param.ReqID,
	}, nil
}
