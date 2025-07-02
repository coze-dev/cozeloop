// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/bytedance/gg/gptr"
	"go.uber.org/mock/gomock"

	idgenMocks "github.com/coze-dev/cozeloop/backend/infra/idgen/mocks"
	lwtMocks "github.com/coze-dev/cozeloop/backend/infra/platestwrite/mocks"
	metricsMocks "github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/component/metrics/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
	eventsMocks "github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/events/mocks"
	repoMocks "github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/repo/mocks"
	svcMocks "github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/service/mocks"
)

func TestExptResultServiceImpl_MGetStats(t *testing.T) {
	tests := []struct {
		name    string
		exptIDs []int64
		spaceID int64
		session *entity.Session
		setup   func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo)
		want    []*entity.ExptStats
		wantErr bool
	}{
		{
			name:    "正常获取多个实验统计",
			exptIDs: []int64{1, 2},
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo) {
				mockExptStatsRepo.EXPECT().
					MGet(gomock.Any(), []int64{1, 2}, int64(100)).
					Return([]*entity.ExptStats{
						{
							ID:      1,
							ExptID:  1,
							SpaceID: 100,
						},
						{
							ID:      2,
							ExptID:  2,
							SpaceID: 100,
						},
					}, nil).
					Times(1)
			},
			want: []*entity.ExptStats{
				{
					ID:      1,
					ExptID:  1,
					SpaceID: 100,
				},
				{
					ID:      2,
					ExptID:  2,
					SpaceID: 100,
				},
			},
			wantErr: false,
		},
		{
			name:    "获取空列表",
			exptIDs: []int64{},
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo) {
				mockExptStatsRepo.EXPECT().
					MGet(gomock.Any(), []int64{}, int64(100)).
					Return([]*entity.ExptStats{}, nil).
					Times(1)
			},
			want:    []*entity.ExptStats{},
			wantErr: false,
		},
		{
			name:    "数据库错误",
			exptIDs: []int64{1},
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo) {
				mockExptStatsRepo.EXPECT().
					MGet(gomock.Any(), []int64{1}, int64(100)).
					Return(nil, fmt.Errorf("db error")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
			svc := ExptResultServiceImpl{
				ExptStatsRepo: mockExptStatsRepo,
			}

			tt.setup(mockExptStatsRepo)

			got, err := svc.MGetStats(context.Background(), tt.exptIDs, tt.spaceID, tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("MGetStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("MGetStats() got length = %v, want %v", len(got), len(tt.want))
					return
				}
				for i := range got {
					if got[i].ID != tt.want[i].ID {
						t.Errorf("MGetStats() got[%d].ID = %v, want %v", i, got[i].ID, tt.want[i].ID)
					}
					if got[i].ExptID != tt.want[i].ExptID {
						t.Errorf("MGetStats() got[%d].ExptID = %v, want %v", i, got[i].ExptID, tt.want[i].ExptID)
					}
					if got[i].SpaceID != tt.want[i].SpaceID {
						t.Errorf("MGetStats() got[%d].SpaceID = %v, want %v", i, got[i].SpaceID, tt.want[i].SpaceID)
					}
				}
			}
		})
	}
}

func TestExptResultServiceImpl_GetStats(t *testing.T) {
	tests := []struct {
		name    string
		exptID  int64
		spaceID int64
		session *entity.Session
		setup   func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo)
		want    *entity.ExptStats
		wantErr bool
	}{
		{
			name:    "正常获取单个实验统计",
			exptID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo) {
				mockExptStatsRepo.EXPECT().
					MGet(gomock.Any(), []int64{1}, int64(100)).
					Return([]*entity.ExptStats{
						{
							ID:      1,
							ExptID:  1,
							SpaceID: 100,
						},
					}, nil).
					Times(1)
			},
			want: &entity.ExptStats{
				ID:      1,
				ExptID:  1,
				SpaceID: 100,
			},
			wantErr: false,
		},
		{
			name:    "数据库错误",
			exptID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo) {
				mockExptStatsRepo.EXPECT().
					MGet(gomock.Any(), []int64{1}, int64(100)).
					Return(nil, fmt.Errorf("db error")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
			svc := ExptResultServiceImpl{
				ExptStatsRepo: mockExptStatsRepo,
			}

			tt.setup(mockExptStatsRepo)

			got, err := svc.GetStats(context.Background(), tt.exptID, tt.spaceID, tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.ID != tt.want.ID {
					t.Errorf("GetStats() got.ID = %v, want %v", got.ID, tt.want.ID)
				}
				if got.ExptID != tt.want.ExptID {
					t.Errorf("GetStats() got.ExptID = %v, want %v", got.ExptID, tt.want.ExptID)
				}
				if got.SpaceID != tt.want.SpaceID {
					t.Errorf("GetStats() got.SpaceID = %v, want %v", got.SpaceID, tt.want.SpaceID)
				}
			}
		})
	}
}

func TestExptResultServiceImpl_CreateStats(t *testing.T) {
	tests := []struct {
		name    string
		stats   *entity.ExptStats
		session *entity.Session
		setup   func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo)
		wantErr bool
	}{
		{
			name: "正常创建统计",
			stats: &entity.ExptStats{
				ID:      1,
				ExptID:  1,
				SpaceID: 100,
			},
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo) {
				mockExptStatsRepo.EXPECT().
					Create(gomock.Any(), &entity.ExptStats{
						ID:      1,
						ExptID:  1,
						SpaceID: 100,
					}).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "数据库错误",
			stats: &entity.ExptStats{
				ID:      1,
				ExptID:  1,
				SpaceID: 100,
			},
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo) {
				mockExptStatsRepo.EXPECT().
					Create(gomock.Any(), &entity.ExptStats{
						ID:      1,
						ExptID:  1,
						SpaceID: 100,
					}).
					Return(fmt.Errorf("db error")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
			svc := ExptResultServiceImpl{
				ExptStatsRepo: mockExptStatsRepo,
			}

			tt.setup(mockExptStatsRepo)

			err := svc.CreateStats(context.Background(), tt.stats, tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateStats() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExptResultServiceImpl_GetExptItemTurnResults(t *testing.T) {
	tests := []struct {
		name    string
		exptID  int64
		itemID  int64
		spaceID int64
		session *entity.Session
		setup   func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo)
		want    []*entity.ExptTurnResult
		wantErr bool
	}{
		{
			name:    "正常获取实验结果",
			exptID:  1,
			itemID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo) {
				// 设置 GetItemTurnResults 的 mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResult{
						{
							ID:     1,
							ExptID: 1,
							ItemID: 1,
						},
					}, nil).
					Times(1)

				// 设置 BatchGetTurnEvaluatorResultRef 的 mock
				mockExptTurnResultRepo.EXPECT().
					BatchGetTurnEvaluatorResultRef(gomock.Any(), int64(100), []int64{1}).
					Return([]*entity.ExptTurnEvaluatorResultRef{
						{
							ExptTurnResultID:   1,
							EvaluatorVersionID: 1,
						},
					}, nil).
					Times(1)
			},
			want: []*entity.ExptTurnResult{
				{
					ID:     1,
					ExptID: 1,
					ItemID: 1,
					EvaluatorResults: &entity.EvaluatorResults{
						EvalVerIDToResID: map[int64]int64{
							1: 1,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "获取空结果",
			exptID:  1,
			itemID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo) {
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResult{}, nil).
					Times(1)

				mockExptTurnResultRepo.EXPECT().
					BatchGetTurnEvaluatorResultRef(gomock.Any(), int64(100), []int64{}).
					Return([]*entity.ExptTurnEvaluatorResultRef{}, nil).
					Times(1)
			},
			want:    []*entity.ExptTurnResult{},
			wantErr: false,
		},
		{
			name:    "数据库错误",
			exptID:  1,
			itemID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo) {
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(1), int64(1), int64(100)).
					Return(nil, fmt.Errorf("db error")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
			svc := ExptResultServiceImpl{
				ExptTurnResultRepo: mockExptTurnResultRepo,
			}

			tt.setup(mockExptTurnResultRepo)

			got, err := svc.GetExptItemTurnResults(context.Background(), tt.exptID, tt.itemID, tt.spaceID, tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExptItemTurnResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("GetExptItemTurnResults() got length = %v, want %v", len(got), len(tt.want))
					return
				}
				for i := range got {
					if got[i].ID != tt.want[i].ID {
						t.Errorf("GetExptItemTurnResults() got[%d].ID = %v, want %v", i, got[i].ID, tt.want[i].ID)
					}
					if got[i].ExptID != tt.want[i].ExptID {
						t.Errorf("GetExptItemTurnResults() got[%d].ExptID = %v, want %v", i, got[i].ExptID, tt.want[i].ExptID)
					}
					if got[i].ItemID != tt.want[i].ItemID {
						t.Errorf("GetExptItemTurnResults() got[%d].ItemID = %v, want %v", i, got[i].ItemID, tt.want[i].ItemID)
					}
				}
			}
		})
	}
}

func TestExptResultServiceImpl_CalculateStats(t *testing.T) {
	tests := []struct {
		name    string
		exptID  int64
		spaceID int64
		session *entity.Session
		setup   func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo)
		want    *entity.ExptCalculateStats
		wantErr bool
	}{
		{
			name:    "正常计算统计",
			exptID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo) {
				mockExptTurnResultRepo.EXPECT().
					ListTurnResult(
						gomock.Any(),
						int64(100),
						int64(1),
						gomock.Any(),
						gomock.Any(),
						false,
					).
					Return([]*entity.ExptTurnResult{
						{
							ID:     1,
							Status: int32(entity.TurnRunState_Success),
						},
						{
							ID:     2,
							Status: int32(entity.TurnRunState_Fail),
						},
					}, int64(2), nil).
					Times(1)
			},
			want: &entity.ExptCalculateStats{
				SuccessTurnCnt: 1,
				FailTurnCnt:    1,
			},
			wantErr: false,
		},
		{
			name:    "数据库错误",
			exptID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo) {
				mockExptTurnResultRepo.EXPECT().
					ListTurnResult(
						gomock.Any(),
						int64(100),
						int64(1),
						gomock.Any(),
						gomock.Any(),
						false,
					).
					Return(nil, int64(0), fmt.Errorf("db error")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "处理中状态",
			exptID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo) {
				mockExptTurnResultRepo.EXPECT().
					ListTurnResult(
						gomock.Any(),
						int64(100),
						int64(1),
						gomock.Any(),
						gomock.Any(),
						false,
					).
					Return([]*entity.ExptTurnResult{
						{
							ID:     1,
							Status: int32(entity.TurnRunState_Processing),
						},
						{
							ID:     2,
							Status: int32(entity.TurnRunState_Queueing),
						},
					}, int64(2), nil).
					Times(1)
			},
			want: &entity.ExptCalculateStats{
				ProcessingTurnCnt: 1,
				PendingTurnCnt:    1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
			svc := ExptResultServiceImpl{
				ExptTurnResultRepo: mockExptTurnResultRepo,
			}

			tt.setup(mockExptTurnResultRepo)

			got, err := svc.CalculateStats(context.Background(), tt.exptID, tt.spaceID, tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.SuccessTurnCnt != tt.want.SuccessTurnCnt {
					t.Errorf("CalculateStats() got.SuccessTurnCnt = %v, want %v", got.SuccessTurnCnt, tt.want.SuccessTurnCnt)
				}
				if got.FailTurnCnt != tt.want.FailTurnCnt {
					t.Errorf("CalculateStats() got.FailTurnCnt = %v, want %v", got.FailTurnCnt, tt.want.FailTurnCnt)
				}
				if got.ProcessingTurnCnt != tt.want.ProcessingTurnCnt {
					t.Errorf("CalculateStats() got.ProcessingTurnCnt = %v, want %v", got.ProcessingTurnCnt, tt.want.ProcessingTurnCnt)
				}
				if got.PendingTurnCnt != tt.want.PendingTurnCnt {
					t.Errorf("CalculateStats() got.PendingTurnCnt = %v, want %v", got.PendingTurnCnt, tt.want.PendingTurnCnt)
				}
			}
		})
	}
}

func TestExptResultServiceImpl_MGetExperimentResult(t *testing.T) {
	tests := []struct {
		name    string
		param   *entity.MGetExperimentResultParam
		setup   func(ctrl *gomock.Controller) ExptResultServiceImpl
		want    []*entity.ColumnEvaluator
		wantErr bool
	}{
		{
			name: "正常获取实验结果",
			param: &entity.MGetExperimentResultParam{
				SpaceID: 100,
				ExptIDs: []int64{1},
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
				mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
				mockMetric := metricsMocks.NewMockExptMetric(ctrl)
				mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockEvaluatorService := svcMocks.NewMockEvaluatorService(ctrl)
				mockEvaluationSetItemService := svcMocks.NewMockEvaluationSetItemService(ctrl)
				mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)
				mockEvalTargetService := svcMocks.NewMockIEvalTargetService(ctrl)
				mockEvaluationSetService := svcMocks.NewMockIEvaluationSetService(ctrl)
				mockEvaluationSetVersionService := svcMocks.NewMockEvaluationSetVersionService(ctrl)

				mockExperimentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Experiment{EvalSetVersionID: 1}, nil).AnyTimes()
				mockExptTurnResultRepo.EXPECT().ListTurnResult(gomock.Any(), int64(100), int64(1), gomock.Any(), gomock.Any(), false).Return([]*entity.ExptTurnResult{{ID: 1, ItemID: 1}}, int64(1), nil)
				mockMetric.EXPECT().EmitGetExptResult(gomock.Any(), gomock.Any()).AnyTimes()
				mockLWT.EXPECT().CheckWriteFlagByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).AnyTimes()
				mockExptStatsRepo.EXPECT().MGet(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptStats{}, nil).AnyTimes()
				mockExperimentRepo.EXPECT().GetEvaluatorRefByExptIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptEvaluatorRef{}, nil).AnyTimes()
				mockEvaluatorService.EXPECT().BatchGetEvaluatorVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.Evaluator{
					{
						ID:            1,
						Name:          "test_evaluator",
						Description:   "test description",
						EvaluatorType: entity.EvaluatorTypePrompt,
						PromptEvaluatorVersion: &entity.PromptEvaluatorVersion{
							ID:      1,
							Version: "v1",
						},
					},
				}, nil).AnyTimes()
				mockEvaluationSetItemService.EXPECT().BatchGetEvaluationSetItems(gomock.Any(), gomock.Any()).Return([]*entity.EvaluationSetItem{}, nil).AnyTimes()
				mockEvaluatorRecordService.EXPECT().BatchGetEvaluatorRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.EvaluatorRecord{}, nil).AnyTimes()
				mockEvalTargetService.EXPECT().BatchGetRecordByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.EvalTargetRecord{}, nil).AnyTimes()
				mockEvaluationSetService.EXPECT().GetEvaluationSet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSet{}, nil).AnyTimes()
				mockEvaluationSetVersionService.EXPECT().GetEvaluationSetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSetVersion{}, nil, nil).AnyTimes()
				mockExptItemResultRepo.EXPECT().BatchGet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptItemResult{}, nil).AnyTimes()
				mockExptTurnResultRepo.EXPECT().BatchGetTurnEvaluatorResultRef(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnEvaluatorResultRef{}, nil).AnyTimes()
				mockExptItemResultRepo.EXPECT().GetItemTurnResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResult{}, nil).AnyTimes()

				return ExptResultServiceImpl{
					ExptTurnResultRepo:          mockExptTurnResultRepo,
					ExperimentRepo:              mockExperimentRepo,
					ExptStatsRepo:               mockExptStatsRepo,
					Metric:                      mockMetric,
					lwt:                         mockLWT,
					ExptItemResultRepo:          mockExptItemResultRepo,
					evaluatorService:            mockEvaluatorService,
					evaluationSetItemService:    mockEvaluationSetItemService,
					evaluatorRecordService:      mockEvaluatorRecordService,
					evalTargetService:           mockEvalTargetService,
					evaluationSetService:        mockEvaluationSetService,
					evaluationSetVersionService: mockEvaluationSetVersionService,
				}
			},
			want:    []*entity.ColumnEvaluator{},
			wantErr: false,
		},
		{
			name: "获取实验失败",
			param: &entity.MGetExperimentResultParam{
				SpaceID: 100,
				ExptIDs: []int64{1},
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
				mockMetric := metricsMocks.NewMockExptMetric(ctrl)
				mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
				mockEvaluationSetService := svcMocks.NewMockIEvaluationSetService(ctrl)

				mockExperimentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("get experiment error"))
				mockMetric.EXPECT().EmitGetExptResult(gomock.Any(), gomock.Any()).AnyTimes()
				mockLWT.EXPECT().CheckWriteFlagByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).AnyTimes()

				return ExptResultServiceImpl{
					ExperimentRepo:       mockExperimentRepo,
					Metric:               mockMetric,
					lwt:                  mockLWT,
					evaluationSetService: mockEvaluationSetService,
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "获取轮次结果失败",
			param: &entity.MGetExperimentResultParam{
				SpaceID: 100,
				ExptIDs: []int64{1},
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
				mockMetric := metricsMocks.NewMockExptMetric(ctrl)
				mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
				mockEvaluatorService := svcMocks.NewMockEvaluatorService(ctrl)
				mockEvaluationSetService := svcMocks.NewMockIEvaluationSetService(ctrl)
				mockEvaluationSetVersionService := svcMocks.NewMockEvaluationSetVersionService(ctrl)

				mockExperimentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Experiment{
					EvalSetVersionID: 1,
					EvalSetID:        1,
				}, nil)
				mockExptTurnResultRepo.EXPECT().ListTurnResult(gomock.Any(), int64(100), int64(1), gomock.Any(), gomock.Any(), false).Return(nil, int64(0), fmt.Errorf("list turn result error"))
				mockMetric.EXPECT().EmitGetExptResult(gomock.Any(), gomock.Any()).AnyTimes()
				mockLWT.EXPECT().CheckWriteFlagByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).AnyTimes()
				mockExperimentRepo.EXPECT().GetEvaluatorRefByExptIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptEvaluatorRef{}, nil)
				mockEvaluatorService.EXPECT().BatchGetEvaluatorVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.Evaluator{}, nil)
				mockEvaluationSetService.EXPECT().GetEvaluationSet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSet{
					EvaluationSetVersion: &entity.EvaluationSetVersion{
						EvaluationSetSchema: &entity.EvaluationSetSchema{
							FieldSchemas: []*entity.FieldSchema{},
						},
					},
				}, nil)
				mockEvaluationSetVersionService.EXPECT().GetEvaluationSetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSetVersion{
					EvaluationSetSchema: &entity.EvaluationSetSchema{
						FieldSchemas: []*entity.FieldSchema{},
					},
				}, nil, nil).AnyTimes()

				return ExptResultServiceImpl{
					ExptTurnResultRepo:          mockExptTurnResultRepo,
					ExperimentRepo:              mockExperimentRepo,
					Metric:                      mockMetric,
					lwt:                         mockLWT,
					evaluatorService:            mockEvaluatorService,
					evaluationSetService:        mockEvaluationSetService,
					evaluationSetVersionService: mockEvaluationSetVersionService,
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "在线实验对比场景",
			param: &entity.MGetExperimentResultParam{
				SpaceID: 100,
				ExptIDs: []int64{1, 2},
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
				mockMetric := metricsMocks.NewMockExptMetric(ctrl)
				mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
				mockEvaluatorService := svcMocks.NewMockEvaluatorService(ctrl)
				mockEvaluationSetService := svcMocks.NewMockIEvaluationSetService(ctrl)
				mockEvaluationSetVersionService := svcMocks.NewMockEvaluationSetVersionService(ctrl)

				mockExperimentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Experiment{
					ExptType:         entity.ExptType_Online,
					EvalSetVersionID: 1,
					EvalSetID:        1,
				}, nil)
				mockMetric.EXPECT().EmitGetExptResult(gomock.Any(), gomock.Any()).AnyTimes()
				mockLWT.EXPECT().CheckWriteFlagByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).AnyTimes()
				mockExperimentRepo.EXPECT().GetEvaluatorRefByExptIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptEvaluatorRef{}, nil)
				mockEvaluatorService.EXPECT().BatchGetEvaluatorVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.Evaluator{}, nil)
				mockEvaluationSetService.EXPECT().GetEvaluationSet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSet{
					EvaluationSetVersion: &entity.EvaluationSetVersion{
						EvaluationSetSchema: &entity.EvaluationSetSchema{
							FieldSchemas: []*entity.FieldSchema{},
						},
					},
				}, nil)
				mockEvaluationSetVersionService.EXPECT().GetEvaluationSetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSetVersion{
					EvaluationSetSchema: &entity.EvaluationSetSchema{
						FieldSchemas: []*entity.FieldSchema{},
					},
				}, nil, nil).AnyTimes()

				return ExptResultServiceImpl{
					ExperimentRepo:              mockExperimentRepo,
					Metric:                      mockMetric,
					lwt:                         mockLWT,
					evaluatorService:            mockEvaluatorService,
					evaluationSetService:        mockEvaluationSetService,
					evaluationSetVersionService: mockEvaluationSetVersionService,
				}
			},
			want:    []*entity.ColumnEvaluator{},
			wantErr: false,
		},
		{
			name: "正常获取离线实验结果",
			param: &entity.MGetExperimentResultParam{
				SpaceID: 100,
				ExptIDs: []int64{1},
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
				mockMetric := metricsMocks.NewMockExptMetric(ctrl)
				mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockEvaluatorService := svcMocks.NewMockEvaluatorService(ctrl)
				mockEvaluationSetItemService := svcMocks.NewMockEvaluationSetItemService(ctrl)
				mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)
				mockEvalTargetService := svcMocks.NewMockIEvalTargetService(ctrl)
				mockEvaluationSetService := svcMocks.NewMockIEvaluationSetService(ctrl)
				mockEvaluationSetVersionService := svcMocks.NewMockEvaluationSetVersionService(ctrl)
				mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)

				mockExperimentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Experiment{
					EvalSetVersionID: 1,
					EvalSetID:        1,
					ExptType:         entity.ExptType_Offline,
				}, nil).AnyTimes()
				mockExptTurnResultRepo.EXPECT().ListTurnResult(gomock.Any(), int64(100), int64(1), gomock.Any(), gomock.Any(), false).Return([]*entity.ExptTurnResult{{ID: 1, ItemID: 1}}, int64(1), nil)
				mockMetric.EXPECT().EmitGetExptResult(gomock.Any(), gomock.Any()).AnyTimes()
				mockLWT.EXPECT().CheckWriteFlagByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).AnyTimes()
				mockExptStatsRepo.EXPECT().MGet(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptStats{}, nil).AnyTimes()
				mockExperimentRepo.EXPECT().GetEvaluatorRefByExptIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptEvaluatorRef{
					{
						EvaluatorVersionID: 1,
						EvaluatorID:        1,
					},
				}, nil).AnyTimes()
				mockEvaluatorService.EXPECT().BatchGetEvaluatorVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.Evaluator{
					{
						ID:            1,
						Name:          "test_evaluator",
						Description:   "test description",
						EvaluatorType: entity.EvaluatorTypePrompt,
						PromptEvaluatorVersion: &entity.PromptEvaluatorVersion{
							ID:      1,
							Version: "v1",
						},
					},
				}, nil).AnyTimes()
				mockEvaluationSetItemService.EXPECT().BatchGetEvaluationSetItems(gomock.Any(), gomock.Any()).Return([]*entity.EvaluationSetItem{}, nil).AnyTimes()
				mockEvaluatorRecordService.EXPECT().BatchGetEvaluatorRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.EvaluatorRecord{}, nil).AnyTimes()
				mockEvalTargetService.EXPECT().BatchGetRecordByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.EvalTargetRecord{}, nil).AnyTimes()
				mockEvaluationSetService.EXPECT().GetEvaluationSet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSet{
					EvaluationSetVersion: &entity.EvaluationSetVersion{
						EvaluationSetSchema: &entity.EvaluationSetSchema{
							FieldSchemas: []*entity.FieldSchema{},
						},
					},
				}, nil).AnyTimes()
				mockEvaluationSetVersionService.EXPECT().GetEvaluationSetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSetVersion{
					EvaluationSetSchema: &entity.EvaluationSetSchema{
						FieldSchemas: []*entity.FieldSchema{},
					},
				}, nil, nil).AnyTimes()
				mockExptItemResultRepo.EXPECT().BatchGet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptItemResult{
					{
						ItemID: 1,
						Status: 1,
					},
				}, nil).AnyTimes()
				mockExptTurnResultRepo.EXPECT().BatchGetTurnEvaluatorResultRef(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnEvaluatorResultRef{}, nil).AnyTimes()
				mockExptItemResultRepo.EXPECT().GetItemTurnResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResult{}, nil).AnyTimes()

				return ExptResultServiceImpl{
					ExptTurnResultRepo:          mockExptTurnResultRepo,
					ExperimentRepo:              mockExperimentRepo,
					ExptStatsRepo:               mockExptStatsRepo,
					Metric:                      mockMetric,
					lwt:                         mockLWT,
					ExptItemResultRepo:          mockExptItemResultRepo,
					evaluatorService:            mockEvaluatorService,
					evaluationSetItemService:    mockEvaluationSetItemService,
					evaluatorRecordService:      mockEvaluatorRecordService,
					evalTargetService:           mockEvalTargetService,
					evaluationSetService:        mockEvaluationSetService,
					evaluationSetVersionService: mockEvaluationSetVersionService,
				}
			},
			want: []*entity.ColumnEvaluator{
				{
					EvaluatorVersionID: 1,
					EvaluatorID:        1,
					EvaluatorType:      entity.EvaluatorTypePrompt,
					Name:               gptr.Of("test_evaluator"),
					Version:            gptr.Of("v1"),
					Description:        gptr.Of("test description"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := tt.setup(ctrl)
			got, _, _, _, err := svc.MGetExperimentResult(context.Background(), tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("MGetExperimentResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MGetExperimentResult() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExptResultServiceImpl_RecordItemRunLogs(t *testing.T) {
	tests := []struct {
		name      string
		exptID    int64
		exptRunID int64
		itemID    int64
		spaceID   int64
		session   *entity.Session
		setup     func(ctrl *gomock.Controller) ExptResultServiceImpl
		wantErr   bool
	}{
		{
			name:      "正常记录运行日志",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
				mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)
				mockPublisher := eventsMocks.NewMockExptEventPublisher(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// GetItemTurnRunLogs mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{{TurnID: 1, Status: entity.TurnRunState_Success}}, nil)

				// GetItemTurnResults mock
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(100), int64(1), int64(1)).
					Return([]*entity.ExptTurnResult{{
						ID:     1,
						TurnID: 1,
						Status: int32(entity.TurnRunState_Success),
					}}, nil)

				// SaveTurnResults mock
				mockExptTurnResultRepo.EXPECT().
					SaveTurnResults(gomock.Any(), gomock.Any()).
					Return(nil)

				// UpdateItemsResult mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemsResult(gomock.Any(), int64(100), int64(1), []int64{1}, gomock.Any()).
					Return(nil)

				// UpdateItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemRunLog(gomock.Any(), int64(1), int64(1), []int64{1}, gomock.Any(), int64(100)).
					Return(nil)

				// ArithOperateCount mock
				mockExptStatsRepo.EXPECT().
					ArithOperateCount(gomock.Any(), int64(1), int64(100), gomock.Any()).
					Return(nil)

				// BatchGetEvaluatorRecord mock
				mockEvaluatorRecordService.EXPECT().
					BatchGetEvaluatorRecord(gomock.Any(), gomock.Any(), true).
					Return([]*entity.EvaluatorRecord{}, nil)

				// PublishExptOnlineEvalResult mock
				mockPublisher.EXPECT().
					PublishExptOnlineEvalResult(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				return ExptResultServiceImpl{
					ExptItemResultRepo:     mockExptItemResultRepo,
					ExptTurnResultRepo:     mockExptTurnResultRepo,
					ExptStatsRepo:          mockExptStatsRepo,
					evaluatorRecordService: mockEvaluatorRecordService,
					publisher:              mockPublisher,
				}
			},
			wantErr: false,
		},
		{
			name:      "获取运行日志失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)

				// GetItemRunLog mock 返回错误
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(nil, fmt.Errorf("get item run log error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo: mockExptItemResultRepo,
				}
			},
			wantErr: true,
		},
		{
			name:      "获取轮次运行日志失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// GetItemTurnRunLogs mock 返回错误
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(nil, fmt.Errorf("get turn run logs error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo: mockExptItemResultRepo,
					ExptTurnResultRepo: mockExptTurnResultRepo,
				}
			},
			wantErr: true,
		},
		{
			name:      "获取轮次结果失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// GetItemTurnRunLogs mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{{TurnID: 1, Status: entity.TurnRunState_Success}}, nil)

				// GetItemTurnResults mock 返回错误
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(100), int64(1), int64(1)).
					Return(nil, fmt.Errorf("get turn results error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo: mockExptItemResultRepo,
					ExptTurnResultRepo: mockExptTurnResultRepo,
				}
			},
			wantErr: true,
		},
		{
			name:      "保存轮次结果失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// GetItemTurnRunLogs mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{{TurnID: 1, Status: entity.TurnRunState_Success}}, nil)

				// GetItemTurnResults mock
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(100), int64(1), int64(1)).
					Return([]*entity.ExptTurnResult{{
						ID:     1,
						TurnID: 1,
						Status: int32(entity.TurnRunState_Success),
					}}, nil)

				// SaveTurnResults mock 返回错误
				mockExptTurnResultRepo.EXPECT().
					SaveTurnResults(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("save turn results error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo: mockExptItemResultRepo,
					ExptTurnResultRepo: mockExptTurnResultRepo,
				}
			},
			wantErr: true,
		},
		{
			name:      "更新项目结果失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// GetItemTurnRunLogs mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{{TurnID: 1, Status: entity.TurnRunState_Success}}, nil)

				// GetItemTurnResults mock
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(100), int64(1), int64(1)).
					Return([]*entity.ExptTurnResult{{
						ID:     1,
						TurnID: 1,
						Status: int32(entity.TurnRunState_Success),
					}}, nil)

				// SaveTurnResults mock
				mockExptTurnResultRepo.EXPECT().
					SaveTurnResults(gomock.Any(), gomock.Any()).
					Return(nil)

				// UpdateItemsResult mock 返回错误
				mockExptItemResultRepo.EXPECT().
					UpdateItemsResult(gomock.Any(), int64(100), int64(1), []int64{1}, gomock.Any()).
					Return(fmt.Errorf("update items result error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo: mockExptItemResultRepo,
					ExptTurnResultRepo: mockExptTurnResultRepo,
				}
			},
			wantErr: true,
		},
		{
			name:      "更新运行日志失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// GetItemTurnRunLogs mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{{TurnID: 1, Status: entity.TurnRunState_Success}}, nil)

				// GetItemTurnResults mock
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(100), int64(1), int64(1)).
					Return([]*entity.ExptTurnResult{{
						ID:     1,
						TurnID: 1,
						Status: int32(entity.TurnRunState_Success),
					}}, nil)

				// SaveTurnResults mock
				mockExptTurnResultRepo.EXPECT().
					SaveTurnResults(gomock.Any(), gomock.Any()).
					Return(nil)

				// UpdateItemsResult mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemsResult(gomock.Any(), int64(100), int64(1), []int64{1}, gomock.Any()).
					Return(nil)

				// UpdateItemRunLog mock 返回错误
				mockExptItemResultRepo.EXPECT().
					UpdateItemRunLog(gomock.Any(), int64(1), int64(1), []int64{1}, gomock.Any(), int64(100)).
					Return(fmt.Errorf("update run log error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo: mockExptItemResultRepo,
					ExptTurnResultRepo: mockExptTurnResultRepo,
				}
			},
			wantErr: true,
		},
		{
			name:      "统计操作失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// GetItemTurnRunLogs mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{{TurnID: 1, Status: entity.TurnRunState_Success}}, nil)

				// GetItemTurnResults mock
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(100), int64(1), int64(1)).
					Return([]*entity.ExptTurnResult{{
						ID:     1,
						TurnID: 1,
						Status: int32(entity.TurnRunState_Success),
					}}, nil)

				// SaveTurnResults mock
				mockExptTurnResultRepo.EXPECT().
					SaveTurnResults(gomock.Any(), gomock.Any()).
					Return(nil)

				// UpdateItemsResult mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemsResult(gomock.Any(), int64(100), int64(1), []int64{1}, gomock.Any()).
					Return(nil)

				// UpdateItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemRunLog(gomock.Any(), int64(1), int64(1), []int64{1}, gomock.Any(), int64(100)).
					Return(nil)

				// ArithOperateCount mock 返回错误
				mockExptStatsRepo.EXPECT().
					ArithOperateCount(gomock.Any(), int64(1), int64(100), gomock.Any()).
					Return(fmt.Errorf("stats operation error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo: mockExptItemResultRepo,
					ExptTurnResultRepo: mockExptTurnResultRepo,
					ExptStatsRepo:      mockExptStatsRepo,
				}
			},
			wantErr: true,
		},
		{
			name:      "获取评估记录失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
				mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// GetItemTurnRunLogs mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{{TurnID: 1, Status: entity.TurnRunState_Success}}, nil)

				// GetItemTurnResults mock
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(100), int64(1), int64(1)).
					Return([]*entity.ExptTurnResult{{
						ID:     1,
						TurnID: 1,
						Status: int32(entity.TurnRunState_Success),
					}}, nil)

				// SaveTurnResults mock
				mockExptTurnResultRepo.EXPECT().
					SaveTurnResults(gomock.Any(), gomock.Any()).
					Return(nil)

				// UpdateItemsResult mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemsResult(gomock.Any(), int64(100), int64(1), []int64{1}, gomock.Any()).
					Return(nil)

				// UpdateItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemRunLog(gomock.Any(), int64(1), int64(1), []int64{1}, gomock.Any(), int64(100)).
					Return(nil)

				// ArithOperateCount mock
				mockExptStatsRepo.EXPECT().
					ArithOperateCount(gomock.Any(), int64(1), int64(100), gomock.Any()).
					Return(nil)

				// BatchGetEvaluatorRecord mock 返回错误
				mockEvaluatorRecordService.EXPECT().
					BatchGetEvaluatorRecord(gomock.Any(), gomock.Any(), true).
					Return(nil, fmt.Errorf("get evaluator record error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo:     mockExptItemResultRepo,
					ExptTurnResultRepo:     mockExptTurnResultRepo,
					ExptStatsRepo:          mockExptStatsRepo,
					evaluatorRecordService: mockEvaluatorRecordService,
				}
			},
			wantErr: true,
		},
		{
			name:      "发布评估结果失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
				mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)
				mockPublisher := eventsMocks.NewMockExptEventPublisher(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// GetItemTurnRunLogs mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{
						{Status: entity.TurnRunState_Success},
					}, nil)

				// GetItemTurnResults mock
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.ExptTurnResult{
						{Status: int32(entity.TurnRunState_Success)},
					}, nil)

				// SaveTurnResults mock
				mockExptTurnResultRepo.EXPECT().
					SaveTurnResults(gomock.Any(), gomock.Any()).
					Return(nil)

				// UpdateItemsResult mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemsResult(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				// UpdateItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemRunLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				// ArithOperateCount mock
				mockExptStatsRepo.EXPECT().
					ArithOperateCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				// BatchGetEvaluatorRecord mock
				mockEvaluatorRecordService.EXPECT().
					BatchGetEvaluatorRecord(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.EvaluatorRecord{}, nil)

				// PublishExptOnlineEvalResult mock 返回错误
				mockPublisher.EXPECT().
					PublishExptOnlineEvalResult(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("publish result error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo:     mockExptItemResultRepo,
					ExptTurnResultRepo:     mockExptTurnResultRepo,
					ExptStatsRepo:          mockExptStatsRepo,
					evaluatorRecordService: mockEvaluatorRecordService,
					publisher:              mockPublisher,
				}
			},
			wantErr: true,
		},
		{
			name:      "轮次日志为空",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
				mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)
				mockPublisher := eventsMocks.NewMockExptEventPublisher(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// GetItemTurnRunLogs mock 返回空结果
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{}, nil)

				// GetItemTurnResults mock
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(100), int64(1), int64(1)).
					Return([]*entity.ExptTurnResult{}, nil)

				// SaveTurnResults mock
				mockExptTurnResultRepo.EXPECT().
					SaveTurnResults(gomock.Any(), gomock.Any()).
					Return(nil)

				// UpdateItemsResult mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemsResult(gomock.Any(), int64(100), int64(1), []int64{1}, gomock.Any()).
					Return(nil)

				// UpdateItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemRunLog(gomock.Any(), int64(1), int64(1), []int64{1}, gomock.Any(), int64(100)).
					Return(nil)

				// ArithOperateCount mock
				mockExptStatsRepo.EXPECT().
					ArithOperateCount(gomock.Any(), int64(1), int64(100), gomock.Any()).
					Return(nil)

				// BatchGetEvaluatorRecord mock
				mockEvaluatorRecordService.EXPECT().
					BatchGetEvaluatorRecord(gomock.Any(), []int64{}, true).
					Return([]*entity.EvaluatorRecord{}, nil)

				// PublishExptOnlineEvalResult mock
				mockPublisher.EXPECT().
					PublishExptOnlineEvalResult(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				return ExptResultServiceImpl{
					ExptItemResultRepo:     mockExptItemResultRepo,
					ExptTurnResultRepo:     mockExptTurnResultRepo,
					ExptStatsRepo:          mockExptStatsRepo,
					evaluatorRecordService: mockEvaluatorRecordService,
					publisher:              mockPublisher,
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := tt.setup(ctrl)
			err := svc.RecordItemRunLogs(context.Background(), tt.exptID, tt.exptRunID, tt.itemID, tt.spaceID, tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordItemRunLogs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewExptResultService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建所有依赖的 mock
	mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
	mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
	mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
	mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
	mockMetric := metricsMocks.NewMockExptMetric(ctrl)
	mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
	mockIDGen := idgenMocks.NewMockIIDGenerator(ctrl)
	mockEvaluatorService := svcMocks.NewMockEvaluatorService(ctrl)
	mockEvalTargetService := svcMocks.NewMockIEvalTargetService(ctrl)
	mockEvaluationSetVersionService := svcMocks.NewMockEvaluationSetVersionService(ctrl)
	mockEvaluationSetService := svcMocks.NewMockIEvaluationSetService(ctrl)
	mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)
	mockEvaluationSetItemService := svcMocks.NewMockEvaluationSetItemService(ctrl)
	mockPublisher := eventsMocks.NewMockExptEventPublisher(ctrl)

	svc := NewExptResultService(
		mockExptItemResultRepo,
		mockExptTurnResultRepo,
		mockExptStatsRepo,
		mockExperimentRepo,
		mockMetric,
		mockLWT,
		mockIDGen,
		mockEvaluatorService,
		mockEvalTargetService,
		mockEvaluationSetVersionService,
		mockEvaluationSetService,
		mockEvaluatorRecordService,
		mockEvaluationSetItemService,
		mockPublisher,
	)

	impl, ok := svc.(ExptResultServiceImpl)
	if !ok {
		t.Fatalf("NewExptResultService should return ExptResultServiceImpl")
	}

	// 断言每个依赖都被正确赋值
	if impl.ExptItemResultRepo != mockExptItemResultRepo {
		t.Errorf("ExptItemResultRepo not set correctly")
	}
	if impl.ExptTurnResultRepo != mockExptTurnResultRepo {
		t.Errorf("ExptTurnResultRepo not set correctly")
	}
	if impl.ExptStatsRepo != mockExptStatsRepo {
		t.Errorf("ExptStatsRepo not set correctly")
	}
	if impl.ExperimentRepo != mockExperimentRepo {
		t.Errorf("ExperimentRepo not set correctly")
	}
	if impl.Metric != mockMetric {
		t.Errorf("Metric not set correctly")
	}
	if impl.lwt != mockLWT {
		t.Errorf("lwt not set correctly")
	}
	if impl.idgen != mockIDGen {
		t.Errorf("idgen not set correctly")
	}
	if impl.evaluatorService != mockEvaluatorService {
		t.Errorf("evaluatorService not set correctly")
	}
	if impl.evalTargetService != mockEvalTargetService {
		t.Errorf("evalTargetService not set correctly")
	}
	if impl.evaluationSetVersionService != mockEvaluationSetVersionService {
		t.Errorf("evaluationSetVersionService not set correctly")
	}
	if impl.evaluationSetService != mockEvaluationSetService {
		t.Errorf("evaluationSetService not set correctly")
	}
	if impl.evaluatorRecordService != mockEvaluatorRecordService {
		t.Errorf("evaluatorRecordService not set correctly")
	}
	if impl.evaluationSetItemService != mockEvaluationSetItemService {
		t.Errorf("evaluationSetItemService not set correctly")
	}
	if impl.publisher != mockPublisher {
		t.Errorf("publisher not set correctly")
	}
}
