// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"reflect"
	"testing"
	"time"

	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/sonic"
	"gorm.io/datatypes"

	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	domainEntity "github.com/coze-dev/coze-loop/backend/modules/data/domain/entity"
	"github.com/coze-dev/coze-loop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
)

// Helper to create datatypes.JSON for expected values
func mustMarshalJSON(v interface{}) datatypes.JSON {
	b, err := sonic.Marshal(v)
	if err != nil {
		panic("test setup: failed to marshal to JSON: " + err.Error())
	}
	return datatypes.JSON(b)
}

// Helper to create *time.Time from int64 timestamp
func timePtr(ts int64) *time.Time {
	if ts == 0 {
		return nil
	}
	t := time.UnixMilli(ts)
	return &t
}

// Helper to create time.Time from int64 timestamp
func timeVal(ts int64) time.Time {
	if ts == 0 {
		return time.Time{}
	}
	return time.UnixMilli(ts)
}

func TestConvertIoJobDOToPO(t *testing.T) {
	// Define some common timestamps
	now := time.Now().UnixMilli()
	later := time.Now().Add(time.Hour).UnixMilli()

	// Define common entity parts
	trueVal := true
	var appID int32 = 123
	var spaceID int64 = 456
	var datasetID int64 = 789
	createdBy := "user_a"
	updatedBy := "user_b"

	sourceFile := &entity.DatasetIOFile{
		Provider: domainEntity.Provider("test_provider"),
		Path:     "/source/path",
		Files:    []string{"file1.txt", "file2.txt"},
	}
	targetDataset := &entity.DatasetIODataset{
		SpaceID:   gptr.Of[int64](spaceID),
		DatasetID: datasetID + 1,
		VersionID: gptr.Of[int64](2),
	}
	fieldMappings := []*entity.FieldMapping{
		{Source: "s_col1", Target: "t_col1"},
		{Source: "s_col2", Target: "t_col2"},
	}
	jobOption := &entity.DatasetIOJobOption{
		OverwriteDataset: &trueVal,
	}
	jobProgress := &entity.DatasetIOJobProgress{
		Total:     gptr.Of[int64](100),
		Processed: gptr.Of[int64](50),
		Added:     gptr.Of[int64](40),
		SubProgresses: []*entity.DatasetIOJobProgress{
			{Name: gptr.Of("sub1"), Total: gptr.Of[int64](10)},
		},
	}
	itemErrors := []*entity.ItemErrorGroup{
		{
			Type:       gptr.Of(entity.ItemErrorType(1)),
			Summary:    gptr.Of("Error summary"),
			ErrorCount: gptr.Of[int32](5),
		},
	}

	tests := []struct {
		name    string
		args    *entity.IOJob
		wantP   *model.DatasetIOJob
		wantErr bool
		// If we were mocking sonic.Marshal, we'd add mock setup here.
		// For now, we assume sonic.Marshal works or test its error propagation if applicable.
	}{
		{
			name: "Full data conversion",
			args: &entity.IOJob{
				ID:            1,
				AppID:         &appID,
				SpaceID:       spaceID,
				DatasetID:     datasetID,
				JobType:       entity.JobType_ImportFromFile,
				Source:        &entity.DatasetIOEndpoint{File: sourceFile},
				Target:        &entity.DatasetIOEndpoint{Dataset: targetDataset},
				FieldMappings: fieldMappings,
				Option:        jobOption,
				Status:        gptr.Of(entity.JobStatus_Running),
				Progress:      jobProgress,
				Errors:        itemErrors,
				CreatedBy:     &createdBy,
				CreatedAt:     &now,
				UpdatedBy:     &updatedBy,
				UpdatedAt:     &later,
				StartedAt:     &now,
				EndedAt:       &later,
			},
			wantP: &model.DatasetIOJob{
				ID:                1,
				AppID:             appID,
				SpaceID:           spaceID,
				DatasetID:         datasetID,
				JobType:           entity.JobType_ImportFromFile.String(),
				SourceFile:        mustMarshalJSON(sourceFile),
				SourceDataset:     nil, // Source is File
				TargetFile:        nil, // Target is Dataset
				TargetDataset:     mustMarshalJSON(targetDataset),
				FieldMappings:     mustMarshalJSON(fieldMappings),
				Option:            mustMarshalJSON(jobOption),
				Status:            entity.JobStatus_Running.String(),
				ProgressTotal:     gptr.Indirect(jobProgress.Total),
				ProgressProcessed: gptr.Indirect(jobProgress.Processed),
				ProgressAdded:     gptr.Indirect(jobProgress.Added),
				SubProgresses:     mustMarshalJSON(jobProgress.SubProgresses),
				Errors:            mustMarshalJSON(itemErrors),
				CreatedBy:         createdBy,
				CreatedAt:         timeVal(now),
				UpdatedBy:         updatedBy,
				UpdatedAt:         timeVal(later),
				StartedAt:         timePtr(now),
				EndedAt:           timePtr(later),
			},
			wantErr: false,
		},
		{
			name: "Minimal data conversion (many nils and zeros)",
			args: &entity.IOJob{
				ID: 2,
				// AppID is nil by default if not set in struct literal
				SpaceID:   spaceID,
				DatasetID: datasetID,
				JobType:   entity.JobType_ExportToFile,
				// Source is nil
				// Target is nil
				// FieldMappings is nil
				// Option is nil
				Status: gptr.Of(entity.JobStatus_Pending),
				// Progress is nil
				// Errors is nil
				// CreatedBy is nil
				// CreatedAt is 0 (nil pointer)
				// UpdatedBy is nil
				// UpdatedAt is 0 (nil pointer)
				// StartedAt is nil
				// EndedAt is nil
			},
			wantP: &model.DatasetIOJob{
				ID:                2,
				AppID:             0, // gptr.Indirect(nil int32) is 0
				SpaceID:           spaceID,
				DatasetID:         datasetID,
				JobType:           entity.JobType_ExportToFile.String(),
				SourceFile:        nil,
				SourceDataset:     nil,
				TargetFile:        nil,
				TargetDataset:     nil,
				FieldMappings:     nil,
				Option:            nil,
				Status:            entity.JobStatus_Pending.String(),
				ProgressTotal:     0,
				ProgressProcessed: 0,
				ProgressAdded:     0,
				SubProgresses:     nil,
				Errors:            nil,
				CreatedBy:         "",          // gptr.Indirect(nil string) is ""
				CreatedAt:         time.Time{}, // unixMilliToTime(0)
				UpdatedBy:         "",
				UpdatedAt:         time.Time{},
				StartedAt:         nil,
				EndedAt:           nil,
			},
			wantErr: false,
		},
		// Add a case for sonic.Marshal failure if we could reliably trigger it
		// or if we were to mock it.
		// For example, if sonic.Marshal for SourceDataset fails:
		// {
		// 	name: "Error on marshalling SourceDataset",
		// 	args: &entity.IOJob{
		// 		ID:      4,
		// 		JobType: entity.JobType_ImportFromFile,
		// 		Status:  gptr.Of(entity.JobStatus_Pending),
		// 		Source:  &entity.DatasetIOEndpoint{Dataset: &entity.DatasetIODataset{ /* construct problematic data if possible */ }},
		// 	},
		// 	wantP:   nil,
		// 	wantErr: true,
		//  // If mocking, setup mock for sonic.Marshal to return an error for this specific call.
		// },
		// Since we are not mocking sonic.Marshal to fail, we'll skip the above.
		// The function is designed to propagate errors from sonic.Marshal.
		// If sonic.Marshal fails (e.g. out of memory, or truly unmarshalable type not present in entity),
		// the function will return (nil, wrappedError), which is the correct behavior.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// If using gomock and mocking interfaces, controller would be initialized here.
			// ctrl := gomock.NewController(t)
			// defer ctrl.Finish()
			// Mocks would be created using ctrl.

			gotP, err := ConvertIoJobDOToPO(tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertIoJobDOToPO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// For error cases, we might want to check the error message/type.
			// if tt.wantErr && err != nil {
			//    // Example: if errors.Is(err, expectedSpecificError)
			// }

			// Using reflect.DeepEqual for comparison.
			// Note: For JSON fields, DeepEqual compares the []byte slices.
			// This is fine if mustMarshalJSON produces consistent output (which sonic usually does for same input).
			if !reflect.DeepEqual(gotP, tt.wantP) {
				// For better diff, can marshal both to JSON string and compare, or print field by field.
				t.Errorf("ConvertIoJobDOToPO() gotP = %v, want %v", gotP, tt.wantP)
				// Detailed diff (optional, can be verbose)
				// if gotP != nil && tt.wantP != nil {
				//  	gotJSON, _ := json.MarshalIndent(gotP, "", "  ")
				//  	wantJSON, _ := json.MarshalIndent(tt.wantP, "", "  ")
				//  	if string(gotJSON) != string(wantJSON) {
				//  		t.Logf("GOT JSON:\n%s\n", string(gotJSON))
				//  		t.Logf("WANT JSON:\n%s\n", string(wantJSON))
				//  	}
				// }
			}
		})
	}
}
