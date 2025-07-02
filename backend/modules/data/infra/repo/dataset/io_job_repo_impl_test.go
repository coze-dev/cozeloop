// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bytedance/gg/gptr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/datatypes"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/infra/idgen/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	mysqlmocks "github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/mocks"
)

func TestDatasetRepo_CreateIOJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIOJobDAO := mysqlmocks.NewMockIIOJobDAO(ctrl)
	mockIDGen := mocks.NewMockIIDGenerator(ctrl)
	repo := &DatasetRepo{ioJobDAO: mockIOJobDAO, idGen: mockIDGen}

	tests := []struct {
		name    string
		job     *entity.IOJob
		mockFn  func()
		wantErr bool
	}{
		{
			name: "success",
			job: &entity.IOJob{
				SpaceID:   1,
				DatasetID: 100,
				JobType:   1,
				Status:    gptr.Of(entity.JobStatus_Pending),
			},
			mockFn: func() {
				mockIDGen.EXPECT().
					GenMultiIDs(gomock.Any(), gomock.Any()).
					Return([]int64{1}, nil)
				mockIOJobDAO.EXPECT().
					CreateIOJob(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "dao error",
			job: &entity.IOJob{
				SpaceID:   1,
				DatasetID: 100,
				JobType:   1,
				Status:    gptr.Of(entity.JobStatus_Pending),
			},
			mockFn: func() {
				mockIDGen.EXPECT().
					GenMultiIDs(gomock.Any(), gomock.Any()).
					Return([]int64{1}, nil)
				mockIOJobDAO.EXPECT().
					CreateIOJob(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("dao error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := repo.CreateIOJob(context.Background(), tt.job)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, int64(1), tt.job.ID)
				assert.NotNil(t, tt.job.CreatedAt)
				assert.NotNil(t, tt.job.UpdatedAt)
			}
		})
	}
}

func TestDatasetRepo_GetIOJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIOJobDAO := mysqlmocks.NewMockIIOJobDAO(ctrl)
	// idGen 和其他 DAO 不是 GetIOJob 直接依赖的，可以不 mock 或传入 nil
	// 但为了 DatasetRepo 结构完整性，可以传入 mock
	repoInstance := &DatasetRepo{
		ioJobDAO: mockIOJobDAO,
		// 其他字段可以根据需要初始化或留空/nil，如果 GetIOJob 不使用它们
	}

	ctx := context.Background()
	now := time.Now()
	validJobTypeStr := entity.JobType_ImportFromFile.String()
	validJobStatusStr := entity.JobStatus_Pending.String()

	// 预期的成功转换后的 entity.IOJob
	expectedEntityJob := &entity.IOJob{
		ID:        1,
		AppID:     gptr.Of(int32(1)),
		SpaceID:   10,
		DatasetID: 100,
		JobType:   entity.JobType_ImportFromFile,
		Status:    gptr.Of(entity.JobStatus_Pending),
		Source:    &entity.DatasetIOEndpoint{},
		Target:    &entity.DatasetIOEndpoint{},
		Progress: &entity.DatasetIOJobProgress{
			Total:     gptr.Of(int64(0)),
			Processed: gptr.Of(int64(0)),
			Added:     gptr.Of(int64(0)),
		},
		CreatedBy: gptr.Of("test_user"),
		CreatedAt: gptr.Of(now.UnixMilli()),
		UpdatedBy: gptr.Of("test_user"),
		UpdatedAt: gptr.Of(now.UnixMilli()),
		StartedAt: nil, // 假设 StartedAt 为 nil
		EndedAt:   nil, // 假设 EndedAt 为 nil
	}

	tests := []struct {
		name      string
		jobID     int64
		repoOpts  []repo.Option // repo.Option for the main function call
		dbOpts    []db.Option   // Expected db.Option for the DAO call
		setupMock func()
		wantJob   *entity.IOJob
		wantErr   bool
		errText   string // 可选，用于检查特定错误信息
	}{
		{
			name:     "success - get io job",
			jobID:    1,
			repoOpts: nil, // No specific repo options
			dbOpts:   []db.Option{},
			setupMock: func() {
				mockIOJobDAO.EXPECT().GetIOJob(ctx, int64(1), gomock.Any()). // gomock.Any() for db.Option for simplicity
												Return(&model.DatasetIOJob{
						ID:                1,
						AppID:             1,
						SpaceID:           10,
						DatasetID:         100,
						JobType:           validJobTypeStr,
						Status:            validJobStatusStr,
						CreatedBy:         "test_user",
						CreatedAt:         now,
						UpdatedBy:         "test_user",
						UpdatedAt:         now,
						SourceDataset:     datatypes.JSON("null"),
						SourceFile:        datatypes.JSON("null"),
						TargetDataset:     datatypes.JSON("null"),
						TargetFile:        datatypes.JSON("null"),
						FieldMappings:     datatypes.JSON("null"),
						Option:            datatypes.JSON("null"),
						SubProgresses:     datatypes.JSON("null"),
						Errors:            datatypes.JSON("null"),
						ProgressTotal:     0,
						ProgressProcessed: 0,
						ProgressAdded:     0,
					}, nil)
			},
			wantJob: expectedEntityJob,
			wantErr: false,
		},
		{
			name:     "success - with repo options",
			jobID:    2,
			repoOpts: []repo.Option{func(opt *repo.Opt) { opt.WithMaster = true }},
			dbOpts:   []db.Option{db.WithMaster()}, // 期望 Opt2DBOpt 转换后的结果
			setupMock: func() {
				// 注意：这里我们期望 Opt2DBOpt(tt.repoOpts...) 的结果作为参数传递
				// 为了简化，我们可以在这里直接使用转换后的 dbOpts，或者使用 gomock.Any()
				// 如果要精确匹配，需要确保 Opt2DBOpt 的逻辑是稳定的
				// 这里使用 gomock.Any() 来避免对 Opt2DBOpt 内部逻辑的强依赖测试
				mockIOJobDAO.EXPECT().GetIOJob(ctx, int64(2), gomock.AssignableToTypeOf([]db.Option{})).
					DoAndReturn(func(_ context.Context, _ int64, opts ...db.Option) (*model.DatasetIOJob, error) {
						// 可以在这里断言 opts 的内容是否符合预期，如果需要的话
						// 例如，检查是否包含 db.WithMaster()
						// 对于本测试，主要关注 GetIOJob 的行为
						return &model.DatasetIOJob{
							ID:            2,
							AppID:         1,
							SpaceID:       10,
							DatasetID:     100,
							JobType:       validJobTypeStr,
							Status:        validJobStatusStr,
							CreatedBy:     "test_user",
							CreatedAt:     now,
							UpdatedBy:     "test_user",
							UpdatedAt:     now,
							SourceDataset: datatypes.JSON("null"), // 确保所有 JSON 字段有效或为 null
							SourceFile:    datatypes.JSON("null"),
							TargetDataset: datatypes.JSON("null"),
							TargetFile:    datatypes.JSON("null"),
							FieldMappings: datatypes.JSON("null"),
							Option:        datatypes.JSON("null"),
							SubProgresses: datatypes.JSON("null"),
							Errors:        datatypes.JSON("null"),
						}, nil
					})
			},
			wantJob: func() *entity.IOJob { // 动态生成期望值，因为 ID 不同
				job := *expectedEntityJob // 浅拷贝
				job.ID = 2
				return &job
			}(),
			wantErr: false,
		},
		{
			name:   "error - dao returns error",
			jobID:  3,
			dbOpts: []db.Option{},
			setupMock: func() {
				mockIOJobDAO.EXPECT().GetIOJob(ctx, int64(3), gomock.Any()).
					Return(nil, errors.New("dao GetIOJob error"))
			},
			wantJob: nil,
			wantErr: true,
			errText: "dao GetIOJob error",
		},
		{
			name:   "error - convertor.IoJobPO2DO fails due to invalid job type",
			jobID:  4,
			dbOpts: []db.Option{},
			setupMock: func() {
				mockIOJobDAO.EXPECT().GetIOJob(ctx, int64(4), gomock.Any()).
					Return(&model.DatasetIOJob{
						ID:        4,
						AppID:     1,
						SpaceID:   10,
						DatasetID: 100,
						JobType:   "invalid_job_type_string", // 这将导致 IoJobPO2DO 失败
						Status:    validJobStatusStr,
						CreatedBy: "test_user",
						CreatedAt: now,
						UpdatedBy: "test_user",
						UpdatedAt: now,
					}, nil)
			},
			wantJob: nil,
			wantErr: true,
			errText: "unknown job_type 'invalid_job_type_string'", // 期望 convertor 抛出的错误
		},
		{
			name:   "error - convertor.IoJobPO2DO fails due to invalid status",
			jobID:  5,
			dbOpts: []db.Option{},
			setupMock: func() {
				mockIOJobDAO.EXPECT().GetIOJob(ctx, int64(5), gomock.Any()).
					Return(&model.DatasetIOJob{
						ID:        5,
						AppID:     1,
						SpaceID:   10,
						DatasetID: 100,
						JobType:   validJobTypeStr,
						Status:    "invalid_status_string", // 这将导致 IoJobPO2DO 失败
						CreatedBy: "test_user",
						CreatedAt: now,
						UpdatedBy: "test_user",
						UpdatedAt: now,
					}, nil)
			},
			wantJob: nil,
			wantErr: true,
			errText: "unknown job_status 'invalid_status_string'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			gotJob, err := repoInstance.GetIOJob(ctx, tt.jobID, tt.repoOpts...)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errText != "" {
					assert.Contains(t, err.Error(), tt.errText)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantJob, gotJob)
		})
	}
}

// TestDatasetRepo_UpdateIOJob tests the UpdateIOJob method of DatasetRepo.
func TestDatasetRepo_UpdateIOJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIOJobDAO := mysqlmocks.NewMockIIOJobDAO(ctrl)
	// idGen is part of DatasetRepo, mock it for completeness, though not directly used by UpdateIOJob.
	mockIDGen := mocks.NewMockIIDGenerator(ctrl)

	repoInstance := &DatasetRepo{
		ioJobDAO: mockIOJobDAO,
		idGen:    mockIDGen,
		// Other DAO fields (datasetDAO, schemaDAO, etc.) are not directly used by UpdateIOJob
		// and can be nil or their respective mocks if a more complete repoInstance is desired.
		// For this specific test, only ioJobDAO is essential.
	}

	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name           string
		jobID          int64
		delta          *repo.DeltaDatasetIOJob
		repoOpts       []repo.Option // Options for the UpdateIOJob method itself
		setupMock      func(mockCtrl *gomock.Controller, jobID int64, delta *repo.DeltaDatasetIOJob, repoOpts []repo.Option)
		wantErr        bool
		expectedErrMsg string // Optional: for checking specific error messages if wantErr is true
	}{
		{
			name:  "成功 - 所有字段均提供",
			jobID: 1,
			delta: &repo.DeltaDatasetIOJob{
				Total:          gptr.Of(int64(100)),
				Status:         gptr.Of(entity.JobStatus_Running),
				PreProcessed:   gptr.Of(int64(10)),
				DeltaProcessed: 5,
				DeltaAdded:     3,
				SubProgresses:  []*entity.DatasetIOJobProgress{{Name: gptr.Of("step1_progress")}},
				Errors:         []*entity.ItemErrorGroup{{Summary: gptr.Of("error_summary_1")}},
				StartedAt:      &now,
				EndedAt:        nil, // Job still running or endedAt not yet set
			},
			repoOpts: nil, // No specific repo options
			setupMock: func(_ *gomock.Controller, jobID int64, delta *repo.DeltaDatasetIOJob, _ []repo.Option) {
				// Opt2DBOpt will convert nil repoOpts to an empty slice of db.Option.
				// Using AssignableToTypeOf to match any slice of db.Option.
				dbOptsMatcher := gomock.AssignableToTypeOf([]db.Option{})
				mockIOJobDAO.EXPECT().UpdateIOJob(ctx, jobID, gomock.Any(), dbOptsMatcher).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "成功 - SubProgresses 为空",
			jobID: 2,
			delta: &repo.DeltaDatasetIOJob{
				Status:        gptr.Of(entity.JobStatus_Completed),
				EndedAt:       &now,
				SubProgresses: []*entity.DatasetIOJobProgress{}, // Explicitly empty
			},
			repoOpts: nil,
			setupMock: func(_ *gomock.Controller, jobID int64, delta *repo.DeltaDatasetIOJob, _ []repo.Option) {
				dbOptsMatcher := gomock.AssignableToTypeOf([]db.Option{})
				mockIOJobDAO.EXPECT().UpdateIOJob(ctx, jobID, gomock.Any(), dbOptsMatcher).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "成功 - Status 为 nil",
			jobID: 3,
			delta: &repo.DeltaDatasetIOJob{
				Total:         gptr.Of(int64(50)),
				Status:        nil, // Status field in delta is nil
				SubProgresses: []*entity.DatasetIOJobProgress{{Name: gptr.Of("step2_progress")}},
			},
			repoOpts: nil,
			setupMock: func(_ *gomock.Controller, jobID int64, delta *repo.DeltaDatasetIOJob, _ []repo.Option) {
				dbOptsMatcher := gomock.AssignableToTypeOf([]db.Option{})
				mockIOJobDAO.EXPECT().UpdateIOJob(ctx, jobID, gomock.Any(), dbOptsMatcher).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "成功 - 带有 repo.Option (例如 WithMaster)",
			jobID: 4,
			delta: &repo.DeltaDatasetIOJob{Status: gptr.Of(entity.JobStatus_Pending)},
			repoOpts: []repo.Option{ // Provide a repo.Option
				func(opt *repo.Opt) {
					opt.WithMaster = true // This should be converted by Opt2DBOpt
				},
			},
			setupMock: func(_ *gomock.Controller, jobID int64, delta *repo.DeltaDatasetIOJob, _ []repo.Option) {
				// Opt2DBOpt(repoOpts...) will be called internally.
				// We expect the db.Option slice passed to the DAO to reflect this.
				// AssignableToTypeOf checks the type. Verifying specific options
				// would require more complex matchers or DoAndReturn with inspection.
				dbOptsMatcher := gomock.AssignableToTypeOf([]db.Option{})
				mockIOJobDAO.EXPECT().UpdateIOJob(ctx, jobID, gomock.Any(), dbOptsMatcher).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "失败 - ioJobDAO.UpdateIOJob 返回错误",
			jobID:    5,
			delta:    &repo.DeltaDatasetIOJob{Status: gptr.Of(entity.JobStatus_Failed)},
			repoOpts: nil,
			setupMock: func(_ *gomock.Controller, jobID int64, delta *repo.DeltaDatasetIOJob, _ []repo.Option) {
				dbOptsMatcher := gomock.AssignableToTypeOf([]db.Option{})
				mockIOJobDAO.EXPECT().UpdateIOJob(ctx, jobID, gomock.Any(), dbOptsMatcher).Return(errors.New("dao update error"))
			},
			wantErr:        true,
			expectedErrMsg: "dao update error",
		},
		// Note: A test case for sonic.MarshalString returning an error is omitted.
		// This is because sonic.MarshalString is a global function from an external package,
		// and gomock is primarily designed for mocking interface methods.
		// Without the ability to mock sonic.MarshalString directly (e.g., via mockey or if it were an interface method),
		// reliably causing it to fail for valid []*domainEntity.DatasetIOJobProgress input in a unit test is difficult.
		// We assume sonic.MarshalString behaves correctly for valid inputs.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks specific to this test case.
			// A new controller for each sub-test is an option for stricter isolation,
			// but gomock is designed to work with a controller per test function (TestXxx).
			// Here, we reuse the main ctrl and set expectations for each sub-test.
			// Ensure mocks are reset or expectations are specific enough if tests interfere.
			// For table-driven tests with gomock, ensure EXPECT calls are made right before the action.
			if tt.setupMock != nil {
				tt.setupMock(ctrl, tt.jobID, tt.delta, tt.repoOpts)
			}

			err := repoInstance.UpdateIOJob(ctx, tt.jobID, tt.delta, tt.repoOpts...)

			if tt.wantErr {
				assert.Error(t, err, "Expected an error but got none")
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg, "Error message mismatch")
				}
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}

func TestDatasetRepo_ListIOJobs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIOJobDAO := mysqlmocks.NewMockIIOJobDAO(ctrl)
	r := &DatasetRepo{ioJobDAO: mockIOJobDAO}

	params := &repo.ListIOJobsParams{}

	now := time.Now()
	validJobTypeStr := entity.JobType_ImportFromFile.String()
	validJobStatusStr := entity.JobStatus_Pending.String()

	// 模拟返回的 model.DatasetIOJob 列表
	mockModels := []*model.DatasetIOJob{
		{
			ID:                1,
			AppID:             1,
			SpaceID:           10,
			DatasetID:         100,
			JobType:           validJobTypeStr,
			Status:            validJobStatusStr,
			CreatedBy:         "test_user",
			CreatedAt:         now,
			UpdatedBy:         "test_user",
			UpdatedAt:         now,
			SourceDataset:     datatypes.JSON("null"),
			SourceFile:        datatypes.JSON("null"),
			TargetDataset:     datatypes.JSON("null"),
			TargetFile:        datatypes.JSON("null"),
			FieldMappings:     datatypes.JSON("null"),
			Option:            datatypes.JSON("null"),
			SubProgresses:     datatypes.JSON("null"),
			Errors:            datatypes.JSON("null"),
			ProgressTotal:     0,
			ProgressProcessed: 0,
			ProgressAdded:     0,
		},
	}

	// 转换后的 entity.IOJob 列表
	expectedEntities := []*entity.IOJob{
		{
			ID:        1,
			AppID:     gptr.Of(int32(1)),
			SpaceID:   10,
			DatasetID: 100,
			JobType:   entity.JobType_ImportFromFile,
			Status:    gptr.Of(entity.JobStatus_Pending),
			Source:    &entity.DatasetIOEndpoint{},
			Target:    &entity.DatasetIOEndpoint{},
			Progress: &entity.DatasetIOJobProgress{
				Total:     gptr.Of(int64(0)),
				Processed: gptr.Of(int64(0)),
				Added:     gptr.Of(int64(0)),
			},
			CreatedBy: gptr.Of("test_user"),
			CreatedAt: gptr.Of(now.UnixMilli()),
			UpdatedBy: gptr.Of("test_user"),
			UpdatedAt: gptr.Of(now.UnixMilli()),
			StartedAt: nil,
			EndedAt:   nil,
		},
	}

	tests := []struct {
		name     string
		params   *repo.ListIOJobsParams
		opt      []repo.Option
		mockFn   func()
		wantJobs []*entity.IOJob
		wantErr  bool
		errText  string
	}{
		{
			name:   "success",
			params: params,
			opt:    nil,
			mockFn: func() {
				mockIOJobDAO.EXPECT().ListIOJobs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockModels, nil)
			},
			wantJobs: expectedEntities,
			wantErr:  false,
		},
		{
			name:   "dao error",
			params: params,
			opt:    nil,
			mockFn: func() {
				mockIOJobDAO.EXPECT().ListIOJobs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("dao error"))
			},
			wantJobs: nil,
			wantErr:  true,
			errText:  "dao error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			jobs, err := r.ListIOJobs(context.Background(), tt.params, tt.opt...)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errText != "" {
					assert.Contains(t, err.Error(), tt.errText)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantJobs, jobs)
		})
	}
}
