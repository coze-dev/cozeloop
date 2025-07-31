// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gslice"

	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/data/domain/dataset"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/data/domain/dataset_job"
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	common_entity "github.com/coze-dev/coze-loop/backend/modules/data/domain/entity"
)

func IOJobDO2DTO(job *entity.IOJob) *dataset_job.DatasetIOJob {
	if job == nil {
		return nil
	}
	return &dataset_job.DatasetIOJob{
		ID:            job.GetID(),
		AppID:         job.AppID,
		SpaceID:       job.SpaceID,
		DatasetID:     job.DatasetID,
		JobType:       dataset_job.JobType(job.JobType),
		Source:        DatasetIOEndpointDO2DTO(job.Source),
		Target:        DatasetIOEndpointDO2DTO(job.Target),
		FieldMappings: gslice.Map(job.FieldMappings, FieldMappingDO2DTO),
		Option:        DatasetIOJobOptionDO2DTO(job.Option),
		Status:        gptr.Of(dataset_job.JobStatus(gptr.Indirect(job.Status))),
		Progress:      DatasetIOJobProgressDO2DTO(job.Progress),
		Errors:        gslice.Map(job.Errors, ItemErrorGroupDO2DTO),
		CreatedBy:     job.CreatedBy,
		CreatedAt:     job.CreatedAt,
		UpdatedBy:     job.UpdatedBy,
		UpdatedAt:     job.UpdatedAt,
		StartedAt:     job.StartedAt,
		EndedAt:       job.EndedAt,
	}
}

func DatasetIOJobProgressDO2DTO(dto *entity.DatasetIOJobProgress) *dataset_job.DatasetIOJobProgress {
	if dto == nil {
		return nil
	}
	return &dataset_job.DatasetIOJobProgress{
		Total:         dto.Total,
		Processed:     dto.Processed,
		Added:         dto.Added,
		Name:          dto.Name,
		SubProgresses: gslice.Map(dto.SubProgresses, DatasetIOJobProgressDO2DTO),
	}
}

func DatasetIOJobOptionDO2DTO(dto *entity.DatasetIOJobOption) *dataset_job.DatasetIOJobOption {
	if dto == nil {
		return nil
	}
	return &dataset_job.DatasetIOJobOption{
		OverwriteDataset: dto.OverwriteDataset,
	}
}

func FieldMappingDO2DTO(dto *entity.FieldMapping) *dataset_job.FieldMapping {
	if dto == nil {
		return nil
	}
	return &dataset_job.FieldMapping{
		Source: dto.Source,
		Target: dto.Target,
	}
}

func DatasetIOEndpointDO2DTO(dto *entity.DatasetIOEndpoint) *dataset_job.DatasetIOEndpoint {
	if dto == nil {
		return nil
	}
	return &dataset_job.DatasetIOEndpoint{
		File:    DatasetIOFileDO2DTO(dto.File),
		Dataset: DatasetIODatasetDO2DTO(dto.Dataset),
	}
}

func DatasetIOFileDO2DTO(dto *entity.DatasetIOFile) *dataset_job.DatasetIOFile {
	if dto == nil {
		return nil
	}
	fromString, _ := dataset.StorageProviderFromString(string(dto.Provider))
	return &dataset_job.DatasetIOFile{
		Provider:       fromString,
		Path:           dto.Path,
		Format:         gptr.Of(dataset_job.FileFormat(gptr.Indirect(dto.Format))),
		CompressFormat: gptr.Of(dataset_job.FileFormat(gptr.Indirect(dto.CompressFormat))),
		Files:          dto.Files,
	}
}

func DatasetIODatasetDO2DTO(dto *entity.DatasetIODataset) *dataset_job.DatasetIODataset {
	if dto == nil {
		return nil
	}
	return &dataset_job.DatasetIODataset{
		SpaceID:   dto.SpaceID,
		DatasetID: dto.DatasetID,
		VersionID: dto.VersionID,
	}
}

func IOJobDTO2DO(dto *dataset_job.DatasetIOJob) *entity.IOJob {
	if dto == nil {
		return nil
	}
	return &entity.IOJob{
		ID:            dto.GetID(),
		AppID:         dto.AppID,
		SpaceID:       dto.GetSpaceID(),
		DatasetID:     dto.GetDatasetID(),
		JobType:       entity.JobType(dto.GetJobType()),
		Source:        DatasetIOEndpointDTO2DO(dto.GetSource()),
		Target:        DatasetIOEndpointDTO2DO(dto.GetTarget()),
		FieldMappings: gslice.Map(dto.GetFieldMappings(), FieldMappingDTO2DO),
		Option:        DatasetIOJobOptionDTO2DO(dto.GetOption()),
		Status:        gptr.Of(entity.JobStatus(dto.GetStatus())),
		Progress:      DatasetIOJobProgressDTO2DO(dto.GetProgress()),
		Errors:        gslice.Map(dto.GetErrors(), ItemErrorGroupDTO2DO),
		CreatedBy:     dto.CreatedBy,
		CreatedAt:     dto.CreatedAt,
		UpdatedBy:     dto.UpdatedBy,
		UpdatedAt:     dto.UpdatedAt,
		StartedAt:     dto.StartedAt,
		EndedAt:       dto.EndedAt,
	}
}

func DatasetIOJobProgressDTO2DO(dto *dataset_job.DatasetIOJobProgress) *entity.DatasetIOJobProgress {
	if dto == nil {
		return nil
	}
	return &entity.DatasetIOJobProgress{
		Total:         dto.Total,
		Processed:     dto.Processed,
		Added:         dto.Added,
		Name:          dto.Name,
		SubProgresses: gslice.Map(dto.GetSubProgresses(), DatasetIOJobProgressDTO2DO),
	}
}

func DatasetIOJobOptionDTO2DO(dto *dataset_job.DatasetIOJobOption) *entity.DatasetIOJobOption {
	if dto == nil {
		return nil
	}
	return &entity.DatasetIOJobOption{
		OverwriteDataset: dto.OverwriteDataset,
	}
}

func FieldMappingDTO2DO(dto *dataset_job.FieldMapping) *entity.FieldMapping {
	if dto == nil {
		return nil
	}
	return &entity.FieldMapping{
		Source: dto.GetSource(),
		Target: dto.GetTarget(),
	}
}

func DatasetIOEndpointDTO2DO(dto *dataset_job.DatasetIOEndpoint) *entity.DatasetIOEndpoint {
	if dto == nil {
		return nil
	}
	return &entity.DatasetIOEndpoint{
		File:    DatasetIOFileDTO2DO(dto.GetFile()),
		Dataset: DatasetIODatasetDTO2DO(dto.GetDataset()),
	}
}

func DatasetIOFileDTO2DO(dto *dataset_job.DatasetIOFile) *entity.DatasetIOFile {
	if dto == nil {
		return nil
	}
	return &entity.DatasetIOFile{
		Provider:       common_entity.Provider(dto.Provider.String()),
		Path:           dto.GetPath(),
		Format:         gptr.Of(entity.FileFormat(dto.GetFormat())),
		CompressFormat: gptr.Of(entity.FileFormat(dto.GetCompressFormat())),
		Files:          dto.GetFiles(),
	}
}

func DatasetIODatasetDTO2DO(dto *dataset_job.DatasetIODataset) *entity.DatasetIODataset {
	if dto == nil {
		return nil
	}
	return &entity.DatasetIODataset{
		SpaceID:   dto.SpaceID,
		DatasetID: dto.GetDatasetID(),
		VersionID: dto.VersionID,
	}
}
