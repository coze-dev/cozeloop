// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package target

import (
	"context"
	"testing"
	"time"

	"github.com/bytedance/gg/gptr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	dbmock "github.com/coze-dev/cozeloop/backend/infra/db/mocks"
	idgen "github.com/coze-dev/cozeloop/backend/infra/idgen/mocks"
	"github.com/coze-dev/cozeloop/backend/infra/platestwrite"
	platestwrite_mocks "github.com/coze-dev/cozeloop/backend/infra/platestwrite/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/infra/repo/target/mysql/gorm_gen/model"
	mysqlmocks "github.com/coze-dev/cozeloop/backend/modules/evaluation/infra/repo/target/mysql/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
)

func TestEvalTargetRepoImpl_CreateEvalTarget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mocks
	mockEvalTargetDao := mysqlmocks.NewMockEvalTargetDAO(ctrl)
	mockEvalTargetVersionDao := mysqlmocks.NewMockEvalTargetVersionDAO(ctrl)
	mockEvalTargetRecordDao := mysqlmocks.NewMockEvalTargetRecordDAO(ctrl)
	mockIDGen := idgen.NewMockIIDGenerator(ctrl)
	mockDBProvider := dbmock.NewMockProvider(ctrl)
	mockLWT := platestwrite_mocks.NewMockILatestWriteTracker(ctrl)

	repo := &EvalTargetRepoImpl{
		evalTargetDao:        mockEvalTargetDao,
		evalTargetVersionDao: mockEvalTargetVersionDao,
		evalTargetRecordDao:  mockEvalTargetRecordDao,
		idgen:                mockIDGen,
		dbProvider:           mockDBProvider,
		lwt:                  mockLWT,
	}

	// Test data
	validSpaceID := int64(123)
	validSourceTargetID := "source-123"
	validSourceTargetVersion := "v1.0"
	validEvalTargetType := int32(1)
	validTargetID := int64(456)
	validVersionID := int64(789)

	tests := []struct {
		name        string
		do          *entity.EvalTarget
		mockSetup   func()
		wantID      int64
		wantVersion int64
		wantErr     bool
		wantErrCode int32
	}{
		{
			name: "success - create new target and version",
			do: &entity.EvalTarget{
				SpaceID:        validSpaceID,
				SourceTargetID: validSourceTargetID,
				EvalTargetType: entity.EvalTargetType(validEvalTargetType),
				EvalTargetVersion: &entity.EvalTargetVersion{
					SpaceID:             validSpaceID,
					SourceTargetVersion: validSourceTargetVersion,
					BaseInfo: &entity.BaseInfo{
						CreatedBy: &entity.UserInfo{},
						UpdatedBy: &entity.UserInfo{},
						CreatedAt: gptr.Of(time.Now().UnixMilli()),
						UpdatedAt: gptr.Of(time.Now().UnixMilli()),
					},
				},
				BaseInfo: &entity.BaseInfo{
					CreatedBy: &entity.UserInfo{},
					UpdatedBy: &entity.UserInfo{},
					CreatedAt: gptr.Of(time.Now().UnixMilli()),
					UpdatedAt: gptr.Of(time.Now().UnixMilli()),
				},
			},
			mockSetup: func() {
				// Mock ID generation
				mockIDGen.EXPECT().
					GenMultiIDs(gomock.Any(), 2).
					Return([]int64{validTargetID, validVersionID}, nil)

				// Mock transaction

				mockDBProvider.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
					return fc(nil)
				})

				// Mock get target by source
				mockEvalTargetDao.EXPECT().
					GetEvalTargetBySourceID(gomock.Any(), validSpaceID, validSourceTargetID, validEvalTargetType, gomock.Any(), gomock.Any()).
					Return(nil, nil)

				// Mock create target
				mockEvalTargetDao.EXPECT().
					CreateEvalTarget(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				// Mock get version by target
				mockEvalTargetVersionDao.EXPECT().
					GetEvalTargetVersionByTarget(gomock.Any(), validSpaceID, validTargetID, validSourceTargetVersion, gomock.Any(), gomock.Any()).
					Return(nil, nil)

				// Mock create version
				mockEvalTargetVersionDao.EXPECT().
					CreateEvalTargetVersion(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				// Mock latest write tracker
				mockLWT.EXPECT().
					SetWriteFlag(gomock.Any(), platestwrite.ResourceTypeTarget, gomock.Any())
				mockLWT.EXPECT().
					SetWriteFlag(gomock.Any(), platestwrite.ResourceTypeTargetVersion, gomock.Any())
			},
			wantID:      validTargetID,
			wantVersion: validVersionID,
			wantErr:     false,
		},
		{
			name:        "error - nil target",
			do:          nil,
			mockSetup:   func() {},
			wantID:      0,
			wantVersion: 0,
			wantErr:     true,
			wantErrCode: errno.CommonInvalidParamCode,
		},
		{
			name: "error - nil version",
			do: &entity.EvalTarget{
				SpaceID:        validSpaceID,
				SourceTargetID: validSourceTargetID,
				EvalTargetType: entity.EvalTargetType(validEvalTargetType),
			},
			mockSetup:   func() {},
			wantID:      0,
			wantVersion: 0,
			wantErr:     true,
			wantErrCode: errno.CommonInvalidParamCode,
		},
		{
			name: "error - ID generation failed",
			do: &entity.EvalTarget{
				SpaceID:        validSpaceID,
				SourceTargetID: validSourceTargetID,
				EvalTargetType: entity.EvalTargetType(validEvalTargetType),
				EvalTargetVersion: &entity.EvalTargetVersion{
					SpaceID:             validSpaceID,
					SourceTargetVersion: validSourceTargetVersion,
				},
			},
			mockSetup: func() {
				mockIDGen.EXPECT().
					GenMultiIDs(gomock.Any(), 2).
					Return(nil, errorx.NewByCode(errno.CommonInternalErrorCode))
			},
			wantID:      0,
			wantVersion: 0,
			wantErr:     true,
			wantErrCode: errno.CommonInternalErrorCode,
		},
		{
			name: "error - target already exists",
			do: &entity.EvalTarget{
				SpaceID:        validSpaceID,
				SourceTargetID: validSourceTargetID,
				EvalTargetType: entity.EvalTargetType(validEvalTargetType),
				EvalTargetVersion: &entity.EvalTargetVersion{
					SpaceID:             validSpaceID,
					SourceTargetVersion: validSourceTargetVersion,
				},
			},
			mockSetup: func() {
				mockIDGen.EXPECT().
					GenMultiIDs(gomock.Any(), 2).
					Return([]int64{validTargetID, validVersionID}, nil)

				mockDBProvider.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
					return fc(nil)
				})

				mockEvalTargetDao.EXPECT().
					GetEvalTargetBySourceID(gomock.Any(), validSpaceID, validSourceTargetID, validEvalTargetType, gomock.Any(), gomock.Any()).
					Return(&model.Target{ID: validTargetID}, nil)

				mockEvalTargetVersionDao.EXPECT().
					GetEvalTargetVersionByTarget(gomock.Any(), validSpaceID, validTargetID, validSourceTargetVersion, gomock.Any(), gomock.Any()).
					Return(nil, nil)

				mockEvalTargetVersionDao.EXPECT().
					CreateEvalTargetVersion(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				mockLWT.EXPECT().
					SetWriteFlag(gomock.Any(), platestwrite.ResourceTypeTarget, gomock.Any())
				mockLWT.EXPECT().
					SetWriteFlag(gomock.Any(), platestwrite.ResourceTypeTargetVersion, gomock.Any())
			},
			wantID:      validTargetID,
			wantVersion: validVersionID,
			wantErr:     false,
		},
		{
			name: "error - version already exists",
			do: &entity.EvalTarget{
				SpaceID:        validSpaceID,
				SourceTargetID: validSourceTargetID,
				EvalTargetType: entity.EvalTargetType(validEvalTargetType),
				EvalTargetVersion: &entity.EvalTargetVersion{
					SpaceID:             validSpaceID,
					SourceTargetVersion: validSourceTargetVersion,
				},
			},
			mockSetup: func() {
				mockIDGen.EXPECT().
					GenMultiIDs(gomock.Any(), 2).
					Return([]int64{validTargetID, validVersionID}, nil)

				mockDBProvider.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
					return fc(nil)
				})

				mockEvalTargetDao.EXPECT().
					GetEvalTargetBySourceID(gomock.Any(), validSpaceID, validSourceTargetID, validEvalTargetType, gomock.Any(), gomock.Any()).
					Return(&model.Target{ID: validTargetID}, nil)

				mockEvalTargetVersionDao.EXPECT().
					GetEvalTargetVersionByTarget(gomock.Any(), validSpaceID, validTargetID, validSourceTargetVersion, gomock.Any(), gomock.Any()).
					Return(&model.TargetVersion{ID: validVersionID}, nil)

				mockLWT.EXPECT().
					SetWriteFlag(gomock.Any(), platestwrite.ResourceTypeTarget, gomock.Any())
				mockLWT.EXPECT().
					SetWriteFlag(gomock.Any(), platestwrite.ResourceTypeTargetVersion, gomock.Any())
			},
			wantID:      validTargetID,
			wantVersion: validVersionID,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			id, versionID, err := repo.CreateEvalTarget(context.Background(), tt.do)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrCode != 0 {
					statusErr, ok := errorx.FromStatusError(err)
					assert.True(t, ok)
					assert.Equal(t, tt.wantErrCode, statusErr.Code())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, id)
				assert.Equal(t, tt.wantVersion, versionID)
			}
		})
	}
}

func TestEvalTargetRepoImpl_GetEvalTarget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mocks
	mockEvalTargetDao := mysqlmocks.NewMockEvalTargetDAO(ctrl)
	mockEvalTargetVersionDao := mysqlmocks.NewMockEvalTargetVersionDAO(ctrl)
	mockEvalTargetRecordDao := mysqlmocks.NewMockEvalTargetRecordDAO(ctrl)
	mockIDGen := idgen.NewMockIIDGenerator(ctrl)
	mockDBProvider := dbmock.NewMockProvider(ctrl)
	mockLWT := platestwrite_mocks.NewMockILatestWriteTracker(ctrl)

	repo := &EvalTargetRepoImpl{
		evalTargetDao:        mockEvalTargetDao,
		evalTargetVersionDao: mockEvalTargetVersionDao,
		evalTargetRecordDao:  mockEvalTargetRecordDao,
		idgen:                mockIDGen,
		dbProvider:           mockDBProvider,
		lwt:                  mockLWT,
	}

	// Test data
	validTargetID := int64(123)
	validSpaceID := int64(456)
	validSourceTargetID := "source-123"
	validEvalTargetType := int32(1)

	tests := []struct {
		name      string
		targetID  int64
		mockSetup func()
		want      *entity.EvalTarget
		wantErr   bool
	}{
		{
			name:     "success - target exists",
			targetID: validTargetID,
			mockSetup: func() {
				// Mock latest write tracker check
				mockLWT.EXPECT().
					CheckWriteFlagByID(gomock.Any(), platestwrite.ResourceTypeTarget, validTargetID).
					Return(false)

				// Mock get target
				mockEvalTargetDao.EXPECT().
					GetEvalTarget(gomock.Any(), validTargetID, gomock.Any()).
					Return(&model.Target{
						ID:             validTargetID,
						SpaceID:        validSpaceID,
						SourceTargetID: validSourceTargetID,
						TargetType:     validEvalTargetType,
					}, nil)
			},
			want: &entity.EvalTarget{
				ID:             validTargetID,
				SpaceID:        validSpaceID,
				SourceTargetID: validSourceTargetID,
				EvalTargetType: entity.EvalTargetType(validEvalTargetType),
			},
			wantErr: false,
		},
		{
			name:     "success - target exists with latest write",
			targetID: validTargetID,
			mockSetup: func() {
				// Mock latest write tracker check
				mockLWT.EXPECT().
					CheckWriteFlagByID(gomock.Any(), platestwrite.ResourceTypeTarget, validTargetID).
					Return(true)

				// Mock get target with master option
				mockEvalTargetDao.EXPECT().
					GetEvalTarget(gomock.Any(), validTargetID, gomock.Any()).
					Return(&model.Target{
						ID:             validTargetID,
						SpaceID:        validSpaceID,
						SourceTargetID: validSourceTargetID,
						TargetType:     validEvalTargetType,
					}, nil)
			},
			want: &entity.EvalTarget{
				ID:             validTargetID,
				SpaceID:        validSpaceID,
				SourceTargetID: validSourceTargetID,
				EvalTargetType: entity.EvalTargetType(validEvalTargetType),
			},
			wantErr: false,
		},
		{
			name:     "error - dao error",
			targetID: validTargetID,
			mockSetup: func() {
				// Mock latest write tracker check
				mockLWT.EXPECT().
					CheckWriteFlagByID(gomock.Any(), platestwrite.ResourceTypeTarget, validTargetID).
					Return(false)

				// Mock get target returns error
				mockEvalTargetDao.EXPECT().
					GetEvalTarget(gomock.Any(), validTargetID).
					Return(nil, errorx.NewByCode(errno.CommonInternalErrorCode))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repo.GetEvalTarget(context.Background(), tt.targetID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, got.ID)
			}
		})
	}
}

func TestEvalTargetRepoImpl_GetEvalTargetVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mocks
	mockEvalTargetDao := mysqlmocks.NewMockEvalTargetDAO(ctrl)
	mockEvalTargetVersionDao := mysqlmocks.NewMockEvalTargetVersionDAO(ctrl)
	mockEvalTargetRecordDao := mysqlmocks.NewMockEvalTargetRecordDAO(ctrl)
	mockIDGen := idgen.NewMockIIDGenerator(ctrl)
	mockDBProvider := dbmock.NewMockProvider(ctrl)
	mockLWT := platestwrite_mocks.NewMockILatestWriteTracker(ctrl)

	repo := &EvalTargetRepoImpl{
		evalTargetDao:        mockEvalTargetDao,
		evalTargetVersionDao: mockEvalTargetVersionDao,
		evalTargetRecordDao:  mockEvalTargetRecordDao,
		idgen:                mockIDGen,
		dbProvider:           mockDBProvider,
		lwt:                  mockLWT,
	}

	// Test data
	validSpaceID := int64(123)
	validVersionID := int64(456)
	validTargetID := int64(789)
	validSourceTargetID := "source-123"
	validEvalTargetType := int32(1)
	validSourceTargetVersion := "v1.0"

	tests := []struct {
		name        string
		spaceID     int64
		versionID   int64
		mockSetup   func()
		want        *entity.EvalTarget
		wantErr     bool
		wantErrCode int32
	}{
		{
			name:      "success - version exists",
			spaceID:   validSpaceID,
			versionID: validVersionID,
			mockSetup: func() {
				// Mock latest write tracker check for version
				mockLWT.EXPECT().
					CheckWriteFlagByID(gomock.Any(), platestwrite.ResourceTypeTargetVersion, validVersionID).
					Return(false)

				// Mock get version
				mockEvalTargetVersionDao.EXPECT().
					GetEvalTargetVersion(gomock.Any(), validSpaceID, validVersionID, gomock.Any()).
					Return(&model.TargetVersion{
						ID:                  validVersionID,
						SpaceID:             validSpaceID,
						TargetID:            validTargetID,
						SourceTargetVersion: validSourceTargetVersion,
						InputSchema:         gptr.Of([]byte("[]")),
						OutputSchema:        gptr.Of([]byte("[]")),
						TargetMeta:          gptr.Of([]byte("{}")),
					}, nil)

				// Mock latest write tracker check for target
				mockLWT.EXPECT().
					CheckWriteFlagByID(gomock.Any(), platestwrite.ResourceTypeTarget, validTargetID).
					Return(false)

				// Mock get target
				mockEvalTargetDao.EXPECT().
					GetEvalTarget(gomock.Any(), validTargetID, gomock.Any()).
					Return(&model.Target{
						ID:             validTargetID,
						SpaceID:        validSpaceID,
						SourceTargetID: validSourceTargetID,
						TargetType:     validEvalTargetType,
					}, nil)
			},
			want: &entity.EvalTarget{
				ID:             validTargetID,
				SpaceID:        validSpaceID,
				SourceTargetID: validSourceTargetID,
				EvalTargetType: entity.EvalTargetType(validEvalTargetType),
				EvalTargetVersion: &entity.EvalTargetVersion{
					ID:                  validVersionID,
					SpaceID:             validSpaceID,
					TargetID:            validTargetID,
					SourceTargetVersion: validSourceTargetVersion,
				},
			},
			wantErr: false,
		},
		{
			name:      "error - version not found",
			spaceID:   validSpaceID,
			versionID: validVersionID,
			mockSetup: func() {
				// Mock latest write tracker check for version
				mockLWT.EXPECT().
					CheckWriteFlagByID(gomock.Any(), platestwrite.ResourceTypeTargetVersion, validVersionID).
					Return(false)

				// Mock get version returns nil
				mockEvalTargetVersionDao.EXPECT().
					GetEvalTargetVersion(gomock.Any(), validSpaceID, validVersionID, gomock.Any()).
					Return(nil, nil)
			},
			want:        nil,
			wantErr:     true,
			wantErrCode: errno.ResourceNotFoundCode,
		},
		{
			name:      "error - target not found",
			spaceID:   validSpaceID,
			versionID: validVersionID,
			mockSetup: func() {
				// Mock latest write tracker check for version
				mockLWT.EXPECT().
					CheckWriteFlagByID(gomock.Any(), platestwrite.ResourceTypeTargetVersion, validVersionID).
					Return(false)

				// Mock get version
				mockEvalTargetVersionDao.EXPECT().
					GetEvalTargetVersion(gomock.Any(), validSpaceID, validVersionID, gomock.Any()).
					Return(&model.TargetVersion{
						ID:                  validVersionID,
						SpaceID:             validSpaceID,
						TargetID:            validTargetID,
						SourceTargetVersion: validSourceTargetVersion,
					}, nil)

				// Mock latest write tracker check for target
				mockLWT.EXPECT().
					CheckWriteFlagByID(gomock.Any(), platestwrite.ResourceTypeTarget, validTargetID).
					Return(false)

				// Mock get target returns nil
				mockEvalTargetDao.EXPECT().
					GetEvalTarget(gomock.Any(), validTargetID, gomock.Any()).
					Return(nil, nil)
			},
			want:        nil,
			wantErr:     true,
			wantErrCode: errno.ResourceNotFoundCode,
		},
		{
			name:      "error - version dao error",
			spaceID:   validSpaceID,
			versionID: validVersionID,
			mockSetup: func() {
				// Mock latest write tracker check for version
				mockLWT.EXPECT().
					CheckWriteFlagByID(gomock.Any(), platestwrite.ResourceTypeTargetVersion, validVersionID).
					Return(false)

				// Mock get version returns error
				mockEvalTargetVersionDao.EXPECT().
					GetEvalTargetVersion(gomock.Any(), validSpaceID, validVersionID, gomock.Any()).
					Return(nil, errorx.NewByCode(errno.CommonInternalErrorCode))
			},
			want:        nil,
			wantErr:     true,
			wantErrCode: errno.CommonInternalErrorCode,
		},
		{
			name:      "error - target dao error",
			spaceID:   validSpaceID,
			versionID: validVersionID,
			mockSetup: func() {
				// Mock latest write tracker check for version
				mockLWT.EXPECT().
					CheckWriteFlagByID(gomock.Any(), platestwrite.ResourceTypeTargetVersion, validVersionID).
					Return(false)

				// Mock get version
				mockEvalTargetVersionDao.EXPECT().
					GetEvalTargetVersion(gomock.Any(), validSpaceID, validVersionID, gomock.Any()).
					Return(&model.TargetVersion{
						ID:                  validVersionID,
						SpaceID:             validSpaceID,
						TargetID:            validTargetID,
						SourceTargetVersion: validSourceTargetVersion,
					}, nil)

				// Mock latest write tracker check for target
				mockLWT.EXPECT().
					CheckWriteFlagByID(gomock.Any(), platestwrite.ResourceTypeTarget, validTargetID).
					Return(false)

				// Mock get target returns error
				mockEvalTargetDao.EXPECT().
					GetEvalTarget(gomock.Any(), validTargetID, gomock.Any()).
					Return(nil, errorx.NewByCode(errno.CommonInternalErrorCode))
			},
			want:        nil,
			wantErr:     true,
			wantErrCode: errno.CommonInternalErrorCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repo.GetEvalTargetVersion(context.Background(), tt.spaceID, tt.versionID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrCode != 0 {
					statusErr, ok := errorx.FromStatusError(err)
					assert.True(t, ok)
					assert.Equal(t, tt.wantErrCode, statusErr.Code())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.EvalTargetVersion.ID, got.EvalTargetVersion.ID)
			}
		})
	}
}

func TestEvalTargetRepoImpl_CreateEvalTargetRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mocks
	mockEvalTargetDao := mysqlmocks.NewMockEvalTargetDAO(ctrl)
	mockEvalTargetVersionDao := mysqlmocks.NewMockEvalTargetVersionDAO(ctrl)
	mockEvalTargetRecordDao := mysqlmocks.NewMockEvalTargetRecordDAO(ctrl)
	mockIDGen := idgen.NewMockIIDGenerator(ctrl)
	mockDBProvider := dbmock.NewMockProvider(ctrl)
	mockLWT := platestwrite_mocks.NewMockILatestWriteTracker(ctrl)

	repo := &EvalTargetRepoImpl{
		evalTargetDao:        mockEvalTargetDao,
		evalTargetVersionDao: mockEvalTargetVersionDao,
		evalTargetRecordDao:  mockEvalTargetRecordDao,
		idgen:                mockIDGen,
		dbProvider:           mockDBProvider,
		lwt:                  mockLWT,
	}

	// Test data
	validSpaceID := int64(123)
	validTargetID := int64(456)
	validVersionID := int64(789)
	validRecordID := int64(101)

	tests := []struct {
		name        string
		record      *entity.EvalTargetRecord
		mockSetup   func()
		wantID      int64
		wantErr     bool
		wantErrCode int32
	}{
		{
			name: "success - create record",
			record: &entity.EvalTargetRecord{
				SpaceID:              validSpaceID,
				TargetID:             validTargetID,
				TargetVersionID:      validVersionID,
				EvalTargetInputData:  &entity.EvalTargetInputData{},
				EvalTargetOutputData: &entity.EvalTargetOutputData{},
				BaseInfo: &entity.BaseInfo{
					CreatedBy: &entity.UserInfo{},
					UpdatedBy: &entity.UserInfo{},
					CreatedAt: gptr.Of(time.Now().UnixMilli()),
					UpdatedAt: gptr.Of(time.Now().UnixMilli()),
				},
			},
			mockSetup: func() {
				// Mock create record
				mockEvalTargetRecordDao.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(validRecordID, nil)
			},
			wantID:  validRecordID,
			wantErr: false,
		},
		{
			name: "error - create record failed",
			record: &entity.EvalTargetRecord{
				SpaceID:         validSpaceID,
				TargetID:        validTargetID,
				TargetVersionID: validVersionID,
				BaseInfo: &entity.BaseInfo{
					CreatedBy: &entity.UserInfo{},
					UpdatedBy: &entity.UserInfo{},
					CreatedAt: gptr.Of(time.Now().UnixMilli()),
					UpdatedAt: gptr.Of(time.Now().UnixMilli()),
				},
			},
			mockSetup: func() {
				// Mock create record returns error
				mockEvalTargetRecordDao.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(int64(0), errorx.NewByCode(errno.CommonInternalErrorCode))
			},
			wantID:      0,
			wantErr:     true,
			wantErrCode: errno.CommonInternalErrorCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repo.CreateEvalTargetRecord(context.Background(), tt.record)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrCode != 0 {
					statusErr, ok := errorx.FromStatusError(err)
					assert.True(t, ok)
					assert.Equal(t, tt.wantErrCode, statusErr.Code())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, got)
			}
		})
	}
}

func TestEvalTargetRepoImpl_GetEvalTargetRecordByIDAndSpaceID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mocks
	mockEvalTargetDao := mysqlmocks.NewMockEvalTargetDAO(ctrl)
	mockEvalTargetVersionDao := mysqlmocks.NewMockEvalTargetVersionDAO(ctrl)
	mockEvalTargetRecordDao := mysqlmocks.NewMockEvalTargetRecordDAO(ctrl)
	mockIDGen := idgen.NewMockIIDGenerator(ctrl)
	mockDBProvider := dbmock.NewMockProvider(ctrl)
	mockLWT := platestwrite_mocks.NewMockILatestWriteTracker(ctrl)

	repo := &EvalTargetRepoImpl{
		evalTargetDao:        mockEvalTargetDao,
		evalTargetVersionDao: mockEvalTargetVersionDao,
		evalTargetRecordDao:  mockEvalTargetRecordDao,
		idgen:                mockIDGen,
		dbProvider:           mockDBProvider,
		lwt:                  mockLWT,
	}

	// Test data
	validSpaceID := int64(123)
	validRecordID := int64(456)
	validTargetID := int64(789)
	validVersionID := int64(101)

	tests := []struct {
		name        string
		spaceID     int64
		recordID    int64
		mockSetup   func()
		want        *entity.EvalTargetRecord
		wantErr     bool
		wantErrCode int32
	}{
		{
			name:     "success - record exists",
			spaceID:  validSpaceID,
			recordID: validRecordID,
			mockSetup: func() {
				// Mock get record
				mockEvalTargetRecordDao.EXPECT().
					GetByIDAndSpaceID(gomock.Any(), validRecordID, validSpaceID).
					Return(&model.TargetRecord{
						ID:              validRecordID,
						SpaceID:         validSpaceID,
						TargetID:        validTargetID,
						TargetVersionID: validVersionID,
						InputData:       gptr.Of([]byte("{}")),
						OutputData:      gptr.Of([]byte("{}")),
					}, nil)
			},
			want: &entity.EvalTargetRecord{
				ID:                   validRecordID,
				SpaceID:              validSpaceID,
				TargetID:             validTargetID,
				TargetVersionID:      validVersionID,
				EvalTargetInputData:  &entity.EvalTargetInputData{},
				EvalTargetOutputData: &entity.EvalTargetOutputData{},
			},
			wantErr: false,
		},
		{
			name:     "success - record not found",
			spaceID:  validSpaceID,
			recordID: validRecordID,
			mockSetup: func() {
				// Mock get record returns nil
				mockEvalTargetRecordDao.EXPECT().
					GetByIDAndSpaceID(gomock.Any(), validRecordID, validSpaceID).
					Return(nil, nil)
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:     "error - dao error",
			spaceID:  validSpaceID,
			recordID: validRecordID,
			mockSetup: func() {
				// Mock get record returns error
				mockEvalTargetRecordDao.EXPECT().
					GetByIDAndSpaceID(gomock.Any(), validRecordID, validSpaceID).
					Return(nil, errorx.NewByCode(errno.CommonInternalErrorCode))
			},
			want:        nil,
			wantErr:     true,
			wantErrCode: errno.CommonInternalErrorCode,
		},
		{
			name:     "error - convert to DO failed",
			spaceID:  validSpaceID,
			recordID: validRecordID,
			mockSetup: func() {
				// Mock get record returns invalid data
				mockEvalTargetRecordDao.EXPECT().
					GetByIDAndSpaceID(gomock.Any(), validRecordID, validSpaceID).
					Return(&model.TargetRecord{
						ID:              validRecordID,
						SpaceID:         validSpaceID,
						TargetID:        validTargetID,
						TargetVersionID: validVersionID,
						// Invalid data to trigger conversion error
						InputData:  gptr.Of([]byte("1")),
						OutputData: gptr.Of([]byte("1")),
					}, nil)
			},
			want:        nil,
			wantErr:     true,
			wantErrCode: errno.CommonInternalErrorCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repo.GetEvalTargetRecordByIDAndSpaceID(context.Background(), tt.spaceID, tt.recordID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrCode != 0 {
					statusErr, ok := errorx.FromStatusError(err)
					assert.True(t, ok)
					assert.Equal(t, tt.wantErrCode, statusErr.Code())
				}
			} else {
				assert.NoError(t, err)
				if tt.want == nil {
					assert.Nil(t, got)
				} else {
					assert.Equal(t, tt.want.ID, got.ID)
					assert.Equal(t, tt.want.SpaceID, got.SpaceID)
					assert.Equal(t, tt.want.TargetID, got.TargetID)
					assert.Equal(t, tt.want.TargetVersionID, got.TargetVersionID)
				}
			}
		})
	}
}
