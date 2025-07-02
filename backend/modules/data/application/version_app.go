// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"strconv"
	"strings"

	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gslice"

	"github.com/coze-dev/cozeloop/backend/infra/external/audit"
	"github.com/coze-dev/cozeloop/backend/infra/middleware/session"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/data/dataset"
	convertor "github.com/coze-dev/cozeloop/backend/modules/data/application/convertor/dataset"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/component/rpc"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/service"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/pagination"
	"github.com/coze-dev/cozeloop/backend/pkg/encoding"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

func (h *DatasetApplicationImpl) CreateDatasetVersion(ctx context.Context, req *dataset.CreateDatasetVersionRequest) (resp *dataset.CreateDatasetVersionResponse, err error) {
	// 鉴权
	err = h.authByDatasetID(ctx, req.GetWorkspaceID(), req.GetDatasetID(), rpc.CommonActionEdit)
	if err != nil {
		return nil, err
	}
	ds, err := h.svc.GetDataset(ctx, req.GetWorkspaceID(), req.GetDatasetID())
	if err != nil {
		return nil, err
	}
	ds.UpdatedBy = session.UserIDInCtxOrEmpty(ctx)
	version := &entity.DatasetVersion{
		AppID:          session.AppIDInCtxOrEmpty(ctx),
		SpaceID:        req.GetWorkspaceID(),
		DatasetID:      req.GetDatasetID(),
		Version:        req.GetVersion(),
		Description:    gptr.Of(req.GetDesc()),
		SnapshotStatus: entity.SnapshotStatusUnstarted,
		CreatedBy:      session.UserIDInCtxOrEmpty(ctx),
	}
	data := map[string]string{
		"texts": strings.Join([]string{gptr.Indirect(version.Description)}, ","), // ignore_security_alert SQL_INJECTION
	}
	record, err := h.auditClient.Audit(ctx, audit.AuditParam{
		ObjectID:  0,
		AuditType: audit.AuditType_CozeLoopDatasetModify,
		AuditData: data,
		ReqID:     encoding.Encode(ctx, data),
	})
	if err != nil {
		logs.CtxError(ctx, "call BatchAuditDatas failed, err=%v", err)
	}
	if record.AuditStatus != audit.AuditStatus_Approved {
		return nil, errno.GetContentAuditFailedErrorf("reason=%s", gptr.Indirect(record.FailedReason))
	}
	if err := h.svc.CreateVersion(ctx, ds, version); err != nil {
		return nil, err
	}

	return &dataset.CreateDatasetVersionResponse{ID: gptr.Of(version.ID)}, nil
}

func (h *DatasetApplicationImpl) ListDatasetVersions(ctx context.Context, req *dataset.ListDatasetVersionsRequest) (resp *dataset.ListDatasetVersionsResponse, err error) {
	// 鉴权
	err = h.authByDatasetID(ctx, req.GetWorkspaceID(), req.GetDatasetID(), rpc.CommonActionRead)
	if err != nil {
		return nil, err
	}
	orderBy := &service.OrderBy{}
	if len(req.GetOrderBys()) != 0 {
		orderBy = &service.OrderBy{
			Field: gptr.Of(req.GetOrderBys()[0].GetField()),
			IsAsc: gptr.Of(req.GetOrderBys()[0].GetIsAsc()),
		}
	}

	pg := pagination.New(
		repo.DatasetVersionOrderBy(gptr.Indirect(orderBy.Field)),
		pagination.WithOrderByAsc(gptr.Indirect(orderBy.IsAsc)),
		pagination.WithPrePage(req.PageNumber, req.PageSize, req.PageToken),
	)

	param := &repo.ListDatasetVersionsParams{
		SpaceID:     req.GetWorkspaceID(),
		DatasetID:   req.GetDatasetID(),
		VersionLike: req.GetVersionLike(),
		Paginator:   pg,
	}

	versions, pr, err := h.repo.ListVersions(ctx, param)
	if err != nil {
		return nil, err
	}
	dtos, err := gslice.TryMap(versions, convertor.VersionDO2DTO).Get()
	if err != nil {
		return nil, err
	}
	total, err := h.repo.CountVersions(ctx, param)
	if err != nil {
		return nil, err
	}
	return &dataset.ListDatasetVersionsResponse{
		Versions:      dtos,
		NextPageToken: gptr.Of(pr.Cursor),
		Total:         gptr.Of(total),
	}, nil
}

func (h *DatasetApplicationImpl) GetDatasetVersion(ctx context.Context, req *dataset.GetDatasetVersionRequest) (resp *dataset.GetDatasetVersionResponse, err error) {
	// 鉴权
	err = h.authByVersionID(ctx, req.GetWorkspaceID(), req.GetVersionID(), rpc.CommonActionRead)
	if err != nil {
		return nil, err
	}
	version, ds, err := h.svc.GetVersionWithOpt(ctx, req.GetWorkspaceID(), req.GetVersionID(), service.WithDeleted(req.GetWithDeleted()))
	if err != nil {
		return nil, err
	}
	ds.Schema.Fields = ds.Schema.AvailableFields()
	dsDTO, err := convertor.DatasetDO2DTO(ds.Dataset, ds.Schema)
	if err != nil {
		return nil, err
	}
	versionDTO, err := convertor.VersionDO2DTO(version)
	if err != nil {
		return nil, err
	}
	return &dataset.GetDatasetVersionResponse{
		Version: versionDTO,
		Dataset: dsDTO,
	}, nil
}

func (h *DatasetApplicationImpl) BatchGetDatasetVersions(ctx context.Context, req *dataset.BatchGetDatasetVersionsRequest) (resp *dataset.BatchGetDatasetVersionsResponse, err error) {
	if len(req.GetVersionIds()) <= 0 {
		return &dataset.BatchGetDatasetVersionsResponse{}, nil
	}
	// 鉴权
	err = h.auth.Authorization(ctx, &rpc.AuthorizationParam{
		ObjectID:      strconv.FormatInt(req.GetWorkspaceID(), 10),
		SpaceID:       req.GetWorkspaceID(),
		ActionObjects: []*rpc.ActionObject{{Action: gptr.Of(rpc.CozeActionListLoopEvaluationSet), EntityType: gptr.Of(rpc.AuthEntityType_Space)}},
	})
	if err != nil {
		return nil, err
	}
	versionedDS, err := h.svc.BatchGetVersionedDatasetsWithOpt(ctx, req.GetWorkspaceID(), req.GetVersionIds(), &service.GetOpt{WithDeleted: req.GetWithDeleted()})
	if err != nil {
		return nil, err
	}
	dtos, err := gslice.TryMap(versionedDS, func(d *service.VersionedDatasetWithSchema) (*dataset.VersionedDataset, error) {
		if !req.GetWithDeleted() {
			d.Schema.Fields = d.Schema.AvailableFields()
		}
		dsDTO, err := convertor.DatasetDO2DTO(d.Dataset, d.Schema)
		if err != nil {
			return nil, err
		}
		versionDTO, err := convertor.VersionDO2DTO(d.Version)
		if err != nil {
			return nil, err
		}
		return &dataset.VersionedDataset{Dataset: dsDTO, Version: versionDTO}, nil
	}).Get()
	if err != nil {
		return nil, err
	}
	return &dataset.BatchGetDatasetVersionsResponse{VersionedDataset: dtos}, nil
}
