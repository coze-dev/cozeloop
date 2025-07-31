// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/coze-loop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	mysqlmocks "github.com/coze-dev/coze-loop/backend/modules/data/infra/repo/dataset/mysql/mocks"
	redismocks "github.com/coze-dev/coze-loop/backend/modules/data/infra/repo/dataset/redis/mocks"
)

func TestDatasetRepo_SetItemCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedisDAO := redismocks.NewMockDatasetDAO(ctrl)
	repo := &DatasetRepo{datasetRedisDAO: mockRedisDAO}

	tests := []struct {
		name      string
		datasetID int64
		count     int64
		mockFn    func()
		wantErr   bool
	}{
		{
			name:      "success",
			datasetID: 1,
			count:     100,
			mockFn: func() {
				mockRedisDAO.EXPECT().
					SetItemCount(gomock.Any(), int64(1), int64(100)).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "redis error",
			datasetID: 2,
			count:     50,
			mockFn: func() {
				mockRedisDAO.EXPECT().
					SetItemCount(gomock.Any(), int64(2), int64(50)).
					Return(errors.New("redis error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := repo.SetItemCount(context.Background(), tt.datasetID, tt.count)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestDatasetRepo_IncrItemCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedisDAO := redismocks.NewMockDatasetDAO(ctrl)
	repo := &DatasetRepo{datasetRedisDAO: mockRedisDAO}

	tests := []struct {
		name      string
		datasetID int64
		incr      int64
		mockFn    func()
		wantCount int64
		wantErr   bool
	}{
		{
			name:      "success",
			datasetID: 1,
			incr:      10,
			mockFn: func() {
				mockRedisDAO.EXPECT().
					IncrItemCount(gomock.Any(), int64(1), int64(10)).
					Return(int64(110), nil)
			},
			wantCount: 110,
			wantErr:   false,
		},
		{
			name:      "redis error",
			datasetID: 2,
			incr:      -5,
			mockFn: func() {
				mockRedisDAO.EXPECT().
					IncrItemCount(gomock.Any(), int64(2), int64(-5)).
					Return(int64(0), errors.New("redis error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			count, err := repo.IncrItemCount(context.Background(), tt.datasetID, tt.incr)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantCount, count)
		})
	}
}

func TestDatasetRepo_GetItemCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedisDAO := redismocks.NewMockDatasetDAO(ctrl)
	repo := &DatasetRepo{datasetRedisDAO: mockRedisDAO}

	tests := []struct {
		name      string
		datasetID int64
		mockFn    func()
		wantCount int64
		wantErr   bool
	}{
		{
			name:      "success",
			datasetID: 1,
			mockFn: func() {
				mockRedisDAO.EXPECT().
					GetItemCount(gomock.Any(), int64(1)).
					Return(int64(100), nil)
			},
			wantCount: 100,
			wantErr:   false,
		},
		{
			name:      "not found",
			datasetID: 2,
			mockFn: func() {
				mockRedisDAO.EXPECT().
					GetItemCount(gomock.Any(), int64(2)).
					Return(int64(0), errors.New("not found"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			count, err := repo.GetItemCount(context.Background(), tt.datasetID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantCount, count)
		})
	}
}

func TestDatasetRepo_MGetItemCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedisDAO := redismocks.NewMockDatasetDAO(ctrl)
	repo := &DatasetRepo{datasetRedisDAO: mockRedisDAO}

	tests := []struct {
		name       string
		datasetIDs []int64
		mockFn     func()
		wantCounts map[int64]int64
		wantErr    bool
	}{
		{
			name:       "success",
			datasetIDs: []int64{1, 2},
			mockFn: func() {
				mockRedisDAO.EXPECT().
					MGetItemCount(gomock.Any(), []int64{1, 2}).
					Return(map[int64]int64{1: 100, 2: 200}, nil)
			},
			wantCounts: map[int64]int64{1: 100, 2: 200},
			wantErr:    false,
		},
		{
			name:       "empty ids",
			datasetIDs: []int64{},
			mockFn: func() {
				mockRedisDAO.EXPECT().
					MGetItemCount(gomock.Any(), []int64{}).
					Return(map[int64]int64{}, nil)
			},
			wantCounts: map[int64]int64{},
			wantErr:    false,
		},
		{
			name:       "redis error",
			datasetIDs: []int64{3, 4},
			mockFn: func() {
				mockRedisDAO.EXPECT().
					MGetItemCount(gomock.Any(), []int64{3, 4}).
					Return(nil, errors.New("redis error"))
			},
			wantCounts: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			counts, err := repo.MGetItemCount(context.Background(), tt.datasetIDs...)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantCounts, counts)
		})
	}
}

func TestDatasetRepo_GetDataset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMySQLDAO := mysqlmocks.NewMockIDatasetDAO(ctrl)
	repo := &DatasetRepo{datasetDAO: mockMySQLDAO}

	tests := []struct {
		name       string
		datasetID  int64
		spaceID    int64
		mockFn     func()
		wantResult *entity.Dataset
		wantErr    bool
	}{
		{
			name:      "success",
			datasetID: 1,
			spaceID:   1,
			mockFn: func() {
				mockMySQLDAO.EXPECT().
					GetDataset(gomock.Any(), int64(1), gomock.Any()).
					Return(&model.Dataset{ID: 1, Spec: []byte("{}"), Features: []byte("{}")}, nil)
			},
			wantResult: &entity.Dataset{ID: 1},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			result, err := repo.GetDataset(context.Background(), tt.spaceID, tt.datasetID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantResult.ID, result.ID)
		})
	}
}

func TestDatasetRepo_MGetDatasets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMySQLDAO := mysqlmocks.NewMockIDatasetDAO(ctrl)
	repo := &DatasetRepo{datasetDAO: mockMySQLDAO}

	tests := []struct {
		name        string
		datasetIDs  []int64
		spaceID     int64
		mockFn      func()
		wantResults []*entity.Dataset
		wantErr     bool
	}{
		{
			name:       "success",
			datasetIDs: []int64{1, 2},
			spaceID:    1,
			mockFn: func() {
				mockMySQLDAO.EXPECT().
					MGetDatasets(gomock.Any(), gomock.Any(), []int64{1, 2}).
					Return([]*model.Dataset{{ID: 1, Spec: []byte("{}"), Features: []byte("{}")}, {ID: 2, Spec: []byte("{}"), Features: []byte("{}")}}, nil)
			},
			wantResults: []*entity.Dataset{{ID: 1}, {ID: 2}},
			wantErr:     false,
		},
		{
			name:       "empty ids",
			datasetIDs: []int64{},
			spaceID:    1,
			mockFn: func() {
				mockMySQLDAO.EXPECT().
					MGetDatasets(gomock.Any(), gomock.Any(), []int64{}).
					Return([]*model.Dataset{}, nil)
			},
			wantResults: []*entity.Dataset{},
			wantErr:     false,
		},
		{
			name:       "redis error",
			datasetIDs: []int64{3, 4},
			spaceID:    1,
			mockFn: func() {
				mockMySQLDAO.EXPECT().
					MGetDatasets(gomock.Any(), gomock.Any(), []int64{3, 4}).
					Return(nil, errors.New("redis error"))
			},
			wantResults: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			results, err := repo.MGetDatasets(context.Background(), tt.spaceID, tt.datasetIDs)
			assert.Equal(t, tt.wantErr, err != nil)
			if len(tt.wantResults) > 0 {
				assert.Equal(t, len(tt.wantResults), len(results))
			}
		})
	}
}

func TestDatasetRepo_PatchDataset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMySQLDAO := mysqlmocks.NewMockIDatasetDAO(ctrl)
	r := &DatasetRepo{datasetDAO: mockMySQLDAO}

	tests := []struct {
		name    string
		patch   *entity.Dataset
		where   *entity.Dataset
		opts    []repo.Option
		mockFn  func()
		wantErr bool
	}{
		{
			name:  "success",
			patch: &entity.Dataset{ID: 1, Name: "new_name", Features: &entity.DatasetFeatures{}, Spec: &entity.DatasetSpec{}},
			where: &entity.Dataset{ID: 1},
			opts:  []repo.Option{},
			mockFn: func() {
				mockMySQLDAO.EXPECT().
					PatchDataset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "patch_nil",
			patch: nil,
			where: &entity.Dataset{ID: 1},
			opts:  []repo.Option{},
			mockFn: func() {
				mockMySQLDAO.EXPECT().
					PatchDataset(gomock.Any(), nil, gomock.Any(), gomock.Any()).
					Return(errors.New("patch cannot be nil"))
			},
			wantErr: true,
		},
		{
			name:  "where_nil",
			patch: &entity.Dataset{ID: 1, Name: "new_name"},
			where: nil,
			opts:  []repo.Option{},
			mockFn: func() {
				mockMySQLDAO.EXPECT().
					PatchDataset(gomock.Any(), gomock.Any(), nil, gomock.Any()).
					Return(errors.New("where cannot be nil"))
			},
			wantErr: true,
		},
		{
			name:  "mysql_error",
			patch: &entity.Dataset{ID: 1, Name: "new_name"},
			where: &entity.Dataset{ID: 1},
			opts:  []repo.Option{},
			mockFn: func() {
				mockMySQLDAO.EXPECT().
					PatchDataset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("mysql error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := r.PatchDataset(context.Background(), tt.patch, tt.where, tt.opts...)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestDatasetRepo_DeleteDataset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := mysqlmocks.NewMockIDatasetDAO(ctrl)
	repo := &DatasetRepo{datasetDAO: mockDAO}

	tests := []struct {
		name      string
		datasetID int64
		mockFn    func()
		wantErr   bool
	}{
		{
			name:      "success",
			datasetID: 1,
			mockFn: func() {
				mockDAO.EXPECT().DeleteDataset(gomock.Any(), gomock.Any(), int64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "mysql error",
			datasetID: 2,
			mockFn: func() {
				mockDAO.EXPECT().DeleteDataset(gomock.Any(), gomock.Any(), int64(2)).Return(errors.New("mysql error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := repo.DeleteDataset(context.Background(), 1, tt.datasetID)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestDatasetRepo_ListDatasets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := mysqlmocks.NewMockIDatasetDAO(ctrl)
	r := &DatasetRepo{datasetDAO: mockDAO}

	params := &repo.ListDatasetsParams{}

	tests := []struct {
		name         string
		params       *repo.ListDatasetsParams
		opt          []repo.Option
		mockFn       func()
		wantDatasets []*entity.Dataset
		wantErr      bool
	}{
		{
			name:   "success",
			params: params,
			opt:    nil,
			mockFn: func() {
				mockDAO.EXPECT().ListDatasets(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*model.Dataset{}, nil, nil)
			},
			wantDatasets: []*entity.Dataset{},
			wantErr:      false,
		},
		{
			name:   "mysql error",
			params: params,
			opt:    nil,
			mockFn: func() {
				mockDAO.EXPECT().ListDatasets(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, errors.New("mysql error"))
			},
			wantDatasets: nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			datasets, _, err := r.ListDatasets(context.Background(), tt.params, tt.opt...)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.wantDatasets, datasets)
			}
		})
	}
}

func TestDatasetRepo_CountDatasets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := mysqlmocks.NewMockIDatasetDAO(ctrl)
	r := &DatasetRepo{datasetDAO: mockDAO}

	params := &repo.ListDatasetsParams{}

	tests := []struct {
		name      string
		params    *repo.ListDatasetsParams
		opt       []repo.Option
		mockFn    func()
		wantCount int64
		wantErr   bool
	}{
		{
			name:   "success",
			params: params,
			opt:    nil,
			mockFn: func() {
				mockDAO.EXPECT().CountDatasets(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(10), nil)
			},
			wantCount: 10,
			wantErr:   false,
		},
		{
			name:   "mysql error",
			params: params,
			opt:    nil,
			mockFn: func() {
				mockDAO.EXPECT().CountDatasets(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), errors.New("mysql error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			count, err := r.CountDatasets(context.Background(), tt.params, tt.opt...)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.wantCount, count)
			}
		})
	}
}

// ... existing code ...
