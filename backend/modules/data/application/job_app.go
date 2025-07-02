// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gslice"

	"github.com/coze-dev/cozeloop/backend/infra/middleware/session"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/data/dataset"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/data/domain/dataset_job"
	convertor "github.com/coze-dev/cozeloop/backend/modules/data/application/convertor/dataset"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/component/rpc"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/service"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
)

func (d *DatasetApplicationImpl) ImportDataset(ctx context.Context, req *dataset.ImportDatasetRequest) (r *dataset.ImportDatasetResponse, err error) {
	// 鉴权
	err = d.authByDatasetID(ctx, req.GetWorkspaceID(), req.GetDatasetID(), rpc.CommonActionEdit)
	if err != nil {
		return nil, err
	}
	ds, err := d.checkImportDatasetReq(ctx, req)
	if err != nil {
		return nil, err
	}

	job := d.buildJob(ctx, req, ds)
	do := convertor.IOJobDTO2DO(job)
	if err := d.svc.CreateIOJob(ctx, do); err != nil {
		return nil, err
	}

	return &dataset.ImportDatasetResponse{JobID: gptr.Of(do.ID)}, nil
}

func (h *DatasetApplicationImpl) checkImportDatasetReq(ctx context.Context, req *dataset.ImportDatasetRequest) (*service.DatasetWithSchema, error) {
	ds, err := h.svc.GetDataset(ctx, req.GetWorkspaceID(), req.GetDatasetID())
	if err != nil {
		return nil, err
	}
	// todo: check dataset status

	allFields := gslice.Map(ds.Schema.AvailableFields(), func(f *entity.FieldSchema) string { return f.Name })
	mappingFields := gslice.Map(req.GetFieldMappings(), func(m *dataset_job.FieldMapping) string { return m.GetTarget() })
	if diff := gslice.Diff(mappingFields, allFields); len(diff) > 0 {
		return nil, errno.BadReqErrorf("target field %s not found in dataset", strings.Join(diff, ", "))
	}

	// check file
	provider := convertor.StorageProviderDTO2DO(req.GetFile().GetProvider())
	stat, err := h.svc.StatFile(ctx, provider, req.GetFile().GetPath())
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, errno.BadReqErrorf("file is a directory")
	}
	ext := strings.ToLower(filepath.Ext(req.GetFile().GetPath()))
	format := req.GetFile().GetFormat()
	if req.GetFile().IsSetCompressFormat() {
		format = req.GetFile().GetCompressFormat()
	}
	if "."+strings.ToLower(format.String()) != ext {
		return nil, errno.BadReqErrorf("file format mismatch, want %s, got %s", format, ext)
	}

	return ds, nil
}

func (h *DatasetApplicationImpl) buildJob(ctx context.Context, req *dataset.ImportDatasetRequest, ds *service.DatasetWithSchema) *dataset_job.DatasetIOJob {
	userID := session.UserIDInCtxOrEmpty(ctx)
	j := &dataset_job.DatasetIOJob{
		AppID:         gptr.Of(ds.AppID),
		SpaceID:       ds.SpaceID,
		DatasetID:     ds.ID,
		JobType:       dataset_job.JobType_ImportFromFile,
		Source:        &dataset_job.DatasetIOEndpoint{File: req.GetFile()},
		Target:        &dataset_job.DatasetIOEndpoint{Dataset: &dataset_job.DatasetIODataset{DatasetID: ds.ID}},
		FieldMappings: req.GetFieldMappings(),
		Option:        req.GetOption(),
		Status:        gptr.Of(dataset_job.JobStatus_Pending),
		CreatedBy:     gptr.Of(userID),
		UpdatedBy:     gptr.Of(userID),
	}

	return j
}

func (d *DatasetApplicationImpl) GetDatasetIOJob(ctx context.Context, req *dataset.GetDatasetIOJobRequest) (r *dataset.GetDatasetIOJobResponse, err error) {
	// 鉴权
	err = d.authByJobID(ctx, req.GetWorkspaceID(), req.GetJobID(), rpc.CommonActionRead)
	if err != nil {
		return nil, err
	}
	job, err := d.svc.GetIOJob(ctx, req.GetJobID())
	if err != nil {
		return nil, err
	}
	return &dataset.GetDatasetIOJobResponse{Job: convertor.IOJobDO2DTO(job)}, nil
}

func (d *DatasetApplicationImpl) ListDatasetIOJobs(ctx context.Context, req *dataset.ListDatasetIOJobsRequest) (r *dataset.ListDatasetIOJobsResponse, err error) {
	// 鉴权
	err = d.auth.Authorization(ctx, &rpc.AuthorizationParam{
		ObjectID:      strconv.FormatInt(req.GetWorkspaceID(), 10),
		SpaceID:       req.GetWorkspaceID(),
		ActionObjects: []*rpc.ActionObject{{Action: gptr.Of(rpc.CozeActionListLoopEvaluationSet), EntityType: gptr.Of(rpc.AuthEntityType_Space)}},
	})
	if err != nil {
		return nil, err
	}
	jobs, err := d.repo.ListIOJobs(ctx, &repo.ListIOJobsParams{
		SpaceID:   req.GetWorkspaceID(),
		DatasetID: req.GetDatasetID(),
		Types:     gslice.Map(req.GetTypes(), func(t dataset_job.JobType) entity.JobType { return entity.JobType(t) }),
		Statuses:  gslice.Map(req.GetStatuses(), func(t dataset_job.JobStatus) entity.JobStatus { return entity.JobStatus(t) }),
	})
	if err != nil {
		return nil, err
	}
	return &dataset.ListDatasetIOJobsResponse{Jobs: gslice.Map(jobs, convertor.IOJobDO2DTO)}, nil
}
