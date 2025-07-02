// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"time"

	"github.com/bytedance/gg/gcond"
	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/sonic"
	"github.com/pkg/errors"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
)

func IoJobPO2DO(p *model.DatasetIOJob) (*entity.IOJob, error) {
	m := &entity.IOJob{
		ID:        p.ID,
		AppID:     gptr.Of(p.AppID),
		SpaceID:   p.SpaceID,
		DatasetID: p.DatasetID,
		Source:    &entity.DatasetIOEndpoint{},
		Target:    &entity.DatasetIOEndpoint{},
		CreatedBy: gptr.Of(p.CreatedBy),
		CreatedAt: gptr.Of(p.CreatedAt.UnixMilli()),
		UpdatedBy: gptr.Of(p.UpdatedBy),
		UpdatedAt: gptr.Of(p.UpdatedAt.UnixMilli()),
	}

	m.Progress = &entity.DatasetIOJobProgress{
		Total:     gptr.Of(p.ProgressTotal),
		Processed: gptr.Of(p.ProgressProcessed),
		Added:     gptr.Of(p.ProgressAdded),
	}
	if t, err := entity.JobTypeFromString(p.JobType); err != nil {
		return nil, errors.WithMessagef(err, "unknown job_type '%s'", p.JobType)
	} else {
		m.JobType = t
	}

	if s, err := entity.JobStatusFromString(p.Status); err != nil {
		return nil, errors.WithMessagef(err, "unknown job_status '%s'", p.Status)
	} else {
		m.Status = &s
	}

	var err error
	m.Source.Dataset, err = unmarshalFiled[*entity.DatasetIODataset]("source_dataset", p.SourceDataset)
	if err != nil {
		return nil, err
	}
	m.Source.File, err = unmarshalFiled[*entity.DatasetIOFile]("source_file", p.SourceFile)
	if err != nil {
		return nil, err
	}
	m.Target.Dataset, err = unmarshalFiled[*entity.DatasetIODataset]("target_dataset", p.TargetDataset)
	if err != nil {
		return nil, err
	}
	m.Target.File, err = unmarshalFiled[*entity.DatasetIOFile]("target_file", p.TargetFile)
	if err != nil {
		return nil, err
	}
	m.FieldMappings, err = unmarshalFiled[[]*entity.FieldMapping]("field_mappings", p.FieldMappings)
	if err != nil {
		return nil, err
	}
	m.Option, err = unmarshalFiled[*entity.DatasetIOJobOption]("option", p.Option)
	if err != nil {
		return nil, err
	}
	m.Progress.SubProgresses, err = unmarshalFiled[[]*entity.DatasetIOJobProgress]("sub_progresses", p.SubProgresses)
	if err != nil {
		return nil, err
	}
	m.Errors, err = unmarshalFiled[[]*entity.ItemErrorGroup]("errors", p.Errors)
	if err != nil {
		return nil, err
	}
	if p.StartedAt != nil {
		m.StartedAt = gptr.Of(p.StartedAt.UnixMilli())
	}
	if p.EndedAt != nil {
		m.EndedAt = gptr.Of(p.EndedAt.UnixMicro())
	}
	return m, nil
}

func ConvertIoJobDOToPO(m *entity.IOJob) (p *model.DatasetIOJob, err error) {
	p = &model.DatasetIOJob{
		ID:        m.ID,
		AppID:     gptr.Indirect(m.AppID),
		SpaceID:   m.SpaceID,
		DatasetID: m.DatasetID,
		JobType:   m.JobType.String(),
		Status:    m.Status.String(),
		CreatedBy: gptr.Indirect(m.CreatedBy),
		CreatedAt: unixMilliToTime(gptr.Indirect(m.CreatedAt)),
		UpdatedBy: gptr.Indirect(m.UpdatedBy),
		UpdatedAt: unixMilliToTime(gptr.Indirect(m.UpdatedAt)),
		StartedAt: gcond.If(m.IsSetStartedAt(), gptr.Of(unixMilliToTime(gptr.Indirect(m.StartedAt))), nil),
		EndedAt:   gcond.If(m.IsSetEndedAt(), gptr.Of(unixMilliToTime(gptr.Indirect(m.EndedAt))), nil),
	}
	if m.Source != nil {
		if s := m.Source.Dataset; s != nil {
			p.SourceDataset, err = sonic.Marshal(s)
			if err != nil {
				return nil, errors.WithMessage(err, "marshal source dataset")
			}
		}
		if s := m.Source.File; s != nil {
			p.SourceFile, err = sonic.Marshal(s)
			if err != nil {
				return nil, errors.WithMessage(err, "marshal source file")
			}
		}
	}
	if m.Target != nil {
		if t := m.Target.Dataset; t != nil {
			p.TargetDataset, err = sonic.Marshal(t)
			if err != nil {
				return nil, errors.WithMessage(err, "marshal target dataset")
			}
		}
		if t := m.Target.File; t != nil {
			p.TargetFile, err = sonic.Marshal(t)
			if err != nil {
				return nil, errors.WithMessage(err, "marshal target file")
			}
		}
	}
	if f := m.FieldMappings; len(f) > 0 {
		p.FieldMappings, err = sonic.Marshal(f)
		if err != nil {
			return nil, errors.WithMessage(err, "marshal field mappings")
		}
	}
	if o := m.Option; o != nil {
		p.Option, err = sonic.Marshal(o)
		if err != nil {
			return nil, errors.WithMessage(err, "marshal option")
		}
	}

	if prog := m.Progress; prog != nil {
		p.ProgressTotal = gptr.Indirect(prog.Total)
		p.ProgressProcessed = gptr.Indirect(prog.Processed)
		p.ProgressAdded = gptr.Indirect(prog.Added)
		if len(prog.SubProgresses) != 0 {
			p.SubProgresses, err = sonic.Marshal(prog.SubProgresses)
			if err != nil {
				return nil, errors.WithMessage(err, "marshal sub progress")
			}
		}
	}
	if l := m.Errors; len(l) > 0 {
		p.Errors, err = sonic.Marshal(l)
		if err != nil {
			return nil, errors.WithMessage(err, "marshal errors")
		}
	}
	return p, nil
}

func unixMilliToTime(ts int64) time.Time {
	if ts == 0 {
		return time.Time{}
	}
	return time.UnixMilli(ts)
}

func unmarshalFiled[T any](name string, data []byte) (T, error) {
	var t T
	if len(data) == 0 {
		return t, nil
	}
	if err := sonic.Unmarshal(data, &t); err != nil {
		return t, errno.JSONErr(err, "unmarshal %s", name)
	}
	return t, nil
}
