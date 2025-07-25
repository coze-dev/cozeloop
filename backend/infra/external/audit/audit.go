// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package audit

import (
	"context"
)

//go:generate mockgen -destination mocks/audit_mock.go -package mocks . IAuditService
type IAuditService interface {
	Audit(ctx context.Context, param AuditParam) (AuditRecord, error)
}

type AuditParam struct {
	ObjectID     int64
	AuditType    AuditType
	ReqID        string
	AuditContext *string
	AuditData    map[string]string
}

type AuditRecord struct {
	ObjectID     int64
	AuditType    AuditType
	AuditStatus  AuditStatus
	FailedReason *string
	ReqID        *string
}

type AuditType int64

const (
	AuditType_CozeLoopPEModify        AuditType = 107
	AuditType_CozeLoopExptModify      AuditType = 108
	AuditType_CozeLoopDatasetModify   AuditType = 109
	AuditType_CozeLoopEvaluatorModify AuditType = 110
)

type AuditStatus int64

const (
	AuditStatus_Default   AuditStatus = 0
	AuditStatus_Pending   AuditStatus = 1
	AuditStatus_Approved  AuditStatus = 2
	AuditStatus_Rejected  AuditStatus = 3
	AuditStatus_Abandoned AuditStatus = 4
	AuditStatus_Failed    AuditStatus = 99
)
