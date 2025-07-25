// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/bytedance/gg/gptr"

	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	"github.com/coze-dev/cozeloop/backend/infra/middleware/session"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/component/userinfo"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/events"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/repo"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

var (
	evaluatorRecordServiceOnce      = sync.Once{}
	singletonEvaluatorRecordService EvaluatorRecordService
)

// NewEvaluatorServiceImpl 创建 EvaluatorService 实例
func NewEvaluatorRecordServiceImpl(idgen idgen.IIDGenerator,
	evaluatorRecordRepo repo.IEvaluatorRecordRepo,
	exptPublisher events.ExptEventPublisher,
	evaluatorPublisher events.EvaluatorEventPublisher,
	userInfoService userinfo.UserInfoService,
) EvaluatorRecordService {
	evaluatorRecordServiceOnce.Do(func() {
		singletonEvaluatorRecordService = &EvaluatorRecordServiceImpl{
			evaluatorRecordRepo: evaluatorRecordRepo,
			idgen:               idgen,
			exptPublisher:       exptPublisher,
			evaluatorPublisher:  evaluatorPublisher,
			userInfoService:     userInfoService,
		}
	})
	return singletonEvaluatorRecordService
}

// EvaluatorRecordServiceImpl 实现 EvaluatorService 接口
type EvaluatorRecordServiceImpl struct {
	idgen               idgen.IIDGenerator
	evaluatorRecordRepo repo.IEvaluatorRecordRepo
	exptPublisher       events.ExptEventPublisher
	evaluatorPublisher  events.EvaluatorEventPublisher
	userInfoService     userinfo.UserInfoService
}

// CorrectEvaluatorRecord 创建 evaluator_version 运行结果
func (s *EvaluatorRecordServiceImpl) CorrectEvaluatorRecord(ctx context.Context, evaluatorRecordDO *entity.EvaluatorRecord, correctionDO *entity.Correction) error {
	userIDInContext := session.UserIDInCtxOrEmpty(ctx)
	correctionDO.UpdatedBy = userIDInContext
	if evaluatorRecordDO.EvaluatorOutputData == nil {
		evaluatorRecordDO.EvaluatorOutputData = &entity.EvaluatorOutputData{}
	}
	if evaluatorRecordDO.EvaluatorOutputData.EvaluatorResult == nil {
		evaluatorRecordDO.EvaluatorOutputData.EvaluatorResult = &entity.EvaluatorResult{}
	}
	evaluatorRecordDO.EvaluatorOutputData.EvaluatorResult.Correction = correctionDO
	if evaluatorRecordDO.BaseInfo == nil {
		evaluatorRecordDO.BaseInfo = &entity.BaseInfo{}
	}
	evaluatorRecordDO.BaseInfo.UpdatedBy = &entity.UserInfo{
		UserID: gptr.Of(userIDInContext),
	}
	evaluatorRecordDO.BaseInfo.UpdatedAt = gptr.Of(time.Now().UnixMilli())
	err := s.evaluatorRecordRepo.CorrectEvaluatorRecord(ctx, evaluatorRecordDO)
	if err != nil {
		return err
	}
	// 发送聚合报告计算消息
	evaluatorVersionIDStr := strconv.FormatInt(evaluatorRecordDO.EvaluatorVersionID, 10)
	if err = s.exptPublisher.PublishExptAggrCalculateEvent(ctx, []*entity.AggrCalculateEvent{{
		ExperimentID:  evaluatorRecordDO.ExperimentID,
		SpaceID:       evaluatorRecordDO.SpaceID,
		CalculateMode: entity.UpdateSpecificField,
		SpecificFieldInfo: &entity.SpecificFieldInfo{
			FieldKey:  evaluatorVersionIDStr,
			FieldType: entity.FieldType_EvaluatorScore,
		},
	}}, gptr.Of(time.Second*3)); err != nil {
		logs.CtxError(ctx, "Failed to send AggrCalculateEvent, evaluatorVersionIDStr: %s, experimentID: %s, err: %v", evaluatorVersionIDStr, evaluatorRecordDO.ExperimentID, err)
	}
	if err = s.evaluatorPublisher.PublishEvaluatorRecordCorrection(ctx, &entity.EvaluatorRecordCorrectionEvent{
		EvaluatorResult:    evaluatorRecordDO.EvaluatorOutputData.EvaluatorResult,
		EvaluatorRecordID:  evaluatorRecordDO.ID,
		EvaluatorVersionID: evaluatorRecordDO.EvaluatorVersionID,
		Ext:                evaluatorRecordDO.Ext,
		CreatedAt:          gptr.Indirect(evaluatorRecordDO.BaseInfo.CreatedAt),
		UpdatedAt:          gptr.Indirect(evaluatorRecordDO.BaseInfo.UpdatedAt),
	}, gptr.Of(time.Second*3)); err != nil {
		return err
	}
	return nil
}

func (s *EvaluatorRecordServiceImpl) GetEvaluatorRecord(ctx context.Context, evaluatorRecordID int64, includeDeleted bool) (*entity.EvaluatorRecord, error) {
	return s.evaluatorRecordRepo.GetEvaluatorRecord(ctx, evaluatorRecordID, includeDeleted)
}

func (s *EvaluatorRecordServiceImpl) BatchGetEvaluatorRecord(ctx context.Context, evaluatorRecordIDs []int64, includeDeleted bool) ([]*entity.EvaluatorRecord, error) {
	records, err := s.evaluatorRecordRepo.BatchGetEvaluatorRecord(ctx, evaluatorRecordIDs, includeDeleted)
	if err != nil {
		return nil, err
	}
	s.userInfoService.PackUserInfo(ctx, userinfo.BatchConvertDO2UserInfoDomainCarrier(records))
	return records, nil
}
