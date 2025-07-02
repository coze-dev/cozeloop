// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bytedance/gg/gcond"
	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gslice"

	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	"github.com/coze-dev/cozeloop/backend/infra/platestwrite"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/component/metrics"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/events"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/repo"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/pkg/contexts"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/json"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/goroutine"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/maps"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

func NewExptResultService(
	exptItemResultRepo repo.IExptItemResultRepo,
	exptTurnResultRepo repo.IExptTurnResultRepo,
	exptStatsRepo repo.IExptStatsRepo,
	experimentRepo repo.IExperimentRepo,
	metric metrics.ExptMetric,
	lwt platestwrite.ILatestWriteTracker,
	idgen idgen.IIDGenerator,
	evaluatorService EvaluatorService,
	evalTargetService IEvalTargetService,
	evaluationSetVersionService EvaluationSetVersionService,
	evaluationSetService IEvaluationSetService,
	evaluatorRecordService EvaluatorRecordService,
	evaluationSetItemService EvaluationSetItemService,
	publisher events.ExptEventPublisher,
) ExptResultService {
	return ExptResultServiceImpl{
		ExptItemResultRepo:          exptItemResultRepo,
		ExptTurnResultRepo:          exptTurnResultRepo,
		ExptStatsRepo:               exptStatsRepo,
		ExperimentRepo:              experimentRepo,
		Metric:                      metric,
		lwt:                         lwt,
		idgen:                       idgen,
		evalTargetService:           evalTargetService,
		evaluationSetVersionService: evaluationSetVersionService,
		evaluationSetService:        evaluationSetService,
		evaluatorService:            evaluatorService,
		evaluatorRecordService:      evaluatorRecordService,
		evaluationSetItemService:    evaluationSetItemService,
		publisher:                   publisher,
	}
}

type ExptResultServiceImpl struct {
	ExptItemResultRepo repo.IExptItemResultRepo
	ExptTurnResultRepo repo.IExptTurnResultRepo
	ExptStatsRepo      repo.IExptStatsRepo
	ExperimentRepo     repo.IExperimentRepo
	Metric             metrics.ExptMetric
	lwt                platestwrite.ILatestWriteTracker
	idgen              idgen.IIDGenerator

	evalTargetService           IEvalTargetService
	evaluationSetVersionService EvaluationSetVersionService
	evaluationSetService        IEvaluationSetService
	evaluatorService            EvaluatorService
	evaluatorRecordService      EvaluatorRecordService
	evaluationSetItemService    EvaluationSetItemService

	publisher events.ExptEventPublisher
}

func (e ExptResultServiceImpl) GetExptItemTurnResults(ctx context.Context, exptID, itemID int64, spaceID int64, session *entity.Session) ([]*entity.ExptTurnResult, error) {
	turnResults, err := e.ExptTurnResultRepo.GetItemTurnResults(ctx, exptID, itemID, spaceID)
	if err != nil {
		return nil, err
	}

	turnResultIDs := make([]int64, 0, len(turnResults))
	for _, tr := range turnResults {
		turnResultIDs = append(turnResultIDs, tr.ID)
	}
	refs, err := e.ExptTurnResultRepo.BatchGetTurnEvaluatorResultRef(ctx, spaceID, turnResultIDs)
	if err != nil {
		return nil, err
	}

	turnEvaluatorVerIDToResultID := make(map[int64]map[int64]int64, len(turnResults))
	for _, ref := range refs {
		if turnEvaluatorVerIDToResultID[ref.ExptTurnResultID] == nil {
			turnEvaluatorVerIDToResultID[ref.ExptTurnResultID] = make(map[int64]int64)
		}
		turnEvaluatorVerIDToResultID[ref.ExptTurnResultID][ref.EvaluatorVersionID] = ref.EvaluatorVersionID
	}

	res := make([]*entity.ExptTurnResult, 0, len(turnResults))
	for _, tr := range turnResults {
		evalVerID2ResultID := turnEvaluatorVerIDToResultID[tr.ID]
		tr.EvaluatorResults = &entity.EvaluatorResults{EvalVerIDToResID: evalVerID2ResultID}
		res = append(res, tr)
	}

	return res, nil
}

func (e ExptResultServiceImpl) RecordItemRunLogs(ctx context.Context, exptID, exptRunID int64, itemID int64, spaceID int64, session *entity.Session) error {
	itemRunLog, err := e.ExptItemResultRepo.GetItemRunLog(ctx, exptID, exptRunID, itemID, spaceID)
	if err != nil {
		return err
	}

	turnRunLogs, err := e.ExptTurnResultRepo.GetItemTurnRunLogs(ctx, exptID, exptRunID, itemID, spaceID)
	if err != nil {
		return err
	}

	turnResults, err := e.ExptItemResultRepo.GetItemTurnResults(ctx, spaceID, exptID, itemID)
	if err != nil {
		return err
	}

	statsCntOp := &entity.StatsCntArithOp{OpStatusCnt: make(map[entity.TurnRunState]int)}
	for _, tr := range turnResults {
		statsCntOp.OpStatusCnt[entity.TurnRunState(tr.Status)] = statsCntOp.OpStatusCnt[entity.TurnRunState(tr.Status)] - 1
	}

	turn2RunLog := make(map[int64]*entity.ExptTurnResultRunLog, len(turnRunLogs))
	for _, trl := range turnRunLogs {
		turn2RunLog[trl.TurnID] = trl
		statsCntOp.OpStatusCnt[trl.Status] = statsCntOp.OpStatusCnt[trl.Status] + 1
	}

	logs.CtxInfo(ctx, "[ExptEval] expt item result with recording run_log, expt_id=%v, expt_run_id=%v, item_id=%v, cnt_op: %v", exptID, exptRunID, itemID, json.Jsonify(statsCntOp))

	var (
		turnEvaluatorRefs []*entity.ExptTurnEvaluatorResultRef
		turn2Result       = gslice.ToMap(turnResults, func(t *entity.ExptTurnResult) (int64, *entity.ExptTurnResult) { return t.TurnID, t })
	)

	for tid, result := range turn2Result {
		rl := turn2RunLog[tid]
		if rl == nil {
			return fmt.Errorf("RecordItemRunLogs found null turn log result, expt_id: %v, expt_run_id: %v, item: %v, tid: %v", exptID, exptRunID, itemID, tid)
		}

		result.Status = int32(rl.Status)
		result.TargetResultID = rl.TargetResultID
		result.ErrMsg = rl.ErrMsg
		result.LogID = rl.LogID
		result.ExptRunID = rl.ExptRunID

		turnEvaluatorRefs = append(turnEvaluatorRefs, NewTurnEvaluatorResultRefs(0, result.ExptID, result.ID, spaceID, rl.EvaluatorResultIds)...)
	}

	if len(turnEvaluatorRefs) > 0 {
		ids, err := e.idgen.GenMultiIDs(ctx, len(turnEvaluatorRefs))
		if err != nil {
			return err
		}

		for idx, ref := range turnEvaluatorRefs {
			ref.ID = ids[idx]
		}

		if err := e.ExptTurnResultRepo.CreateTurnEvaluatorRefs(ctx, turnEvaluatorRefs); err != nil {
			return err
		}
	}

	if err := e.ExptTurnResultRepo.SaveTurnResults(ctx, turnResults); err != nil {
		return err
	}

	if err := e.ExptItemResultRepo.UpdateItemsResult(ctx, spaceID, exptID, []int64{itemID}, map[string]any{
		"status":  itemRunLog.Status,
		"log_id":  itemRunLog.LogID,
		"err_msg": itemRunLog.ErrMsg,
	}); err != nil {
		return err
	}

	if err := e.ExptItemResultRepo.UpdateItemRunLog(ctx, exptID, exptRunID, []int64{itemID}, map[string]any{
		"result_state": int32(entity.ExptItemResultStateResulted),
	}, spaceID); err != nil {
		return err
	}

	if err := e.ExptStatsRepo.ArithOperateCount(ctx, exptID, spaceID, statsCntOp); err != nil {
		return err
	}

	evaluatorResultIDs := make([]int64, 0, len(turnEvaluatorRefs))
	for _, ref := range turnEvaluatorRefs {
		evaluatorResultIDs = append(evaluatorResultIDs, ref.EvaluatorResultID)
	}
	evaluatorRecords, err := e.evaluatorRecordService.BatchGetEvaluatorRecord(ctx, evaluatorResultIDs, true)
	if err != nil {
		return err
	}
	onlineExptTurnEvalResults := make([]*entity.OnlineExptTurnEvalResult, 0, len(evaluatorRecords))
	for _, record := range evaluatorRecords {
		onlineExptTurnEvalResult := &entity.OnlineExptTurnEvalResult{
			EvaluatorVersionId: record.EvaluatorVersionID,
			EvaluatorRecordId:  record.ID,
			Status:             int32(record.Status),
			Ext:                record.Ext,
			BaseInfo:           record.BaseInfo,
		}
		if record.EvaluatorOutputData != nil {
			if record.Status == entity.EvaluatorRunStatusFail && record.EvaluatorOutputData.EvaluatorRunError != nil {
				onlineExptTurnEvalResult.EvaluatorRunError = &entity.EvaluatorRunError{
					Code:    record.EvaluatorOutputData.EvaluatorRunError.Code,
					Message: record.EvaluatorOutputData.EvaluatorRunError.Message,
				}
			} else if record.Status == entity.EvaluatorRunStatusSuccess && record.EvaluatorOutputData.EvaluatorResult != nil {
				onlineExptTurnEvalResult.Score = gptr.Indirect(record.EvaluatorOutputData.EvaluatorResult.Score)
				onlineExptTurnEvalResult.Reasoning = record.EvaluatorOutputData.EvaluatorResult.Reasoning
			}
		}

		onlineExptTurnEvalResults = append(onlineExptTurnEvalResults, onlineExptTurnEvalResult)
	}

	// 发送评估结果Event
	err = e.publisher.PublishExptOnlineEvalResult(ctx, &entity.OnlineExptEvalResultEvent{
		ExptId:          exptID,
		TurnEvalResults: onlineExptTurnEvalResults,
	}, gptr.Of(time.Second*3))
	if err != nil {
		return err
	}

	return nil
}

func NewTurnEvaluatorResultRefs(id, exptID, turnResultID, spaceID int64, evaluatorResults *entity.EvaluatorResults) []*entity.ExptTurnEvaluatorResultRef {
	if evaluatorResults == nil {
		return nil
	}

	refs := make([]*entity.ExptTurnEvaluatorResultRef, 0, len(evaluatorResults.EvalVerIDToResID))
	for evalVerID, evalResID := range evaluatorResults.EvalVerIDToResID {
		refs = append(refs, &entity.ExptTurnEvaluatorResultRef{
			ID:                 id,
			ExptID:             exptID,
			SpaceID:            spaceID,
			ExptTurnResultID:   turnResultID,
			EvaluatorVersionID: evalVerID,
			EvaluatorResultID:  evalResID,
		})
	}
	return refs
}

func (e ExptResultServiceImpl) MGetExperimentResult(ctx context.Context, param *entity.MGetExperimentResultParam) (
	columnEvaluators []*entity.ColumnEvaluator, columnEvalSetFields []*entity.ColumnEvalSetField, itemResults []*entity.ItemResult, total int64, err error,
) {
	var (
		spaceID        = param.SpaceID
		exptIDs        = param.ExptIDs
		baselineExptID = param.BaseExptID
		filters        = param.Filters
		page           = param.Page
	)

	defer e.Metric.EmitGetExptResult(spaceID, err != nil)

	if len(exptIDs) == 1 && e.lwt.CheckWriteFlagByID(ctx, platestwrite.ResourceTypeExperiment, exptIDs[0]) {
		ctx = contexts.WithCtxWriteDB(ctx)
	}

	var baseExptID int64
	if baselineExptID != nil {
		baseExptID = *baselineExptID
	}
	// 只有一个实验，且没有指定baseline
	if len(exptIDs) == 1 && baselineExptID == nil {
		baseExptID = exptIDs[0]
	}

	baseExpt, err := e.ExperimentRepo.GetByID(ctx, baseExptID, spaceID)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	baseExptEvalSetVersionID := baseExpt.EvalSetVersionID

	var filter *entity.ExptTurnResultFilter
	if len(filters) != 0 && filters[baseExptID] != nil {
		filter = filters[baseExptID]
	}

	columnEvaluators, err = e.getColumnEvaluators(ctx, spaceID, exptIDs)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	columnEvalSetFields, err = e.getColumnEvalSetFields(ctx, spaceID, baseExpt.EvalSetID, baseExptEvalSetVersionID)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	if baseExpt.ExptType == entity.ExptType_Online && len(exptIDs) > 1 {
		// 在线实验对比场景，不返回行级结果
		return columnEvaluators, columnEvalSetFields, nil, 0, nil
	}

	// 获取baseline 该分页的turn_result
	turnResultDAOs, total, err := e.ExptTurnResultRepo.ListTurnResult(ctx, spaceID, baseExptID, filter, page, gcond.If(baseExpt.ExptType == entity.ExptType_Online, true, false))
	if err != nil {
		return nil, nil, nil, 0, err
	}

	if len(turnResultDAOs) == 0 {
		return columnEvaluators, columnEvalSetFields, nil, 0, nil
	}

	itemIDMap := make(map[int64]bool)
	for _, turnResult := range turnResultDAOs {
		itemIDMap[turnResult.ItemID] = true
	}
	itemIDs := maps.ToSlice(itemIDMap, func(k int64, v bool) int64 {
		return k
	})
	itemResultDAOs, err := e.ExptItemResultRepo.BatchGet(ctx, spaceID, baseExptID, itemIDs)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	payloadBuilder := NewPayloadBuilder(ctx, param, baseExptID, turnResultDAOs, itemResultDAOs, e.ExperimentRepo, e.ExptTurnResultRepo, e.evalTargetService, e.evaluatorRecordService, e.evaluationSetItemService)

	itemResults, err = payloadBuilder.Build(ctx)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	return columnEvaluators, columnEvalSetFields, itemResults, total, nil
}

// getColumnEvaluators 试验对比无需返回多试验的评估器合集,没有评估器的column,前端从实验接口获取评估器数据
func (e ExptResultServiceImpl) getColumnEvaluators(ctx context.Context, spaceID int64, exptIDs []int64) ([]*entity.ColumnEvaluator, error) {
	evaluatorRef, err := e.ExperimentRepo.GetEvaluatorRefByExptIDs(ctx, exptIDs, spaceID)
	if err != nil {
		return nil, err
	}
	// 去重
	evaluatorVersionIDMap := make(map[int64]bool)
	evaluatorIDMap := make(map[int64]bool)
	versionID2evaluatorID := make(map[int64]int64)
	for _, ref := range evaluatorRef {
		evaluatorVersionIDMap[ref.EvaluatorVersionID] = true
		evaluatorIDMap[ref.EvaluatorID] = true
		versionID2evaluatorID[ref.EvaluatorVersionID] = ref.EvaluatorID
	}

	evaluatorVersionIDs := maps.ToSlice(evaluatorVersionIDMap, func(k int64, v bool) int64 {
		return k
	})

	evaluatorVersions, err := e.evaluatorService.BatchGetEvaluatorVersion(ctx, evaluatorVersionIDs, true)
	if err != nil {
		return nil, err
	}

	columnEvaluators := make([]*entity.ColumnEvaluator, 0)
	for _, e := range evaluatorVersions {
		evaluatorVersion := e.GetEvaluatorVersion()
		if evaluatorVersion == nil || !gslice.Contains(evaluatorVersionIDs, evaluatorVersion.GetID()) {
			continue
		}

		columnEvaluator := &entity.ColumnEvaluator{
			EvaluatorVersionID: evaluatorVersion.GetID(),
			EvaluatorID:        e.ID,
			EvaluatorType:      e.EvaluatorType,
			Name:               gptr.Of(e.Name),
			Version:            gptr.Of(e.GetEvaluatorVersion().GetVersion()),
			Description:        gptr.Of(e.Description),
		}
		columnEvaluators = append(columnEvaluators, columnEvaluator)
	}

	return columnEvaluators, nil
}

func (e ExptResultServiceImpl) getColumnEvalSetFields(ctx context.Context, spaceID int64, evalSetID, evalSetVersionID int64) ([]*entity.ColumnEvalSetField, error) {
	var version *entity.EvaluationSetVersion
	if evalSetID == evalSetVersionID {
		evalSet, err := e.evaluationSetService.GetEvaluationSet(ctx, gptr.Of(spaceID), evalSetID, gptr.Of(true))
		if err != nil {
			return nil, err
		}
		version = evalSet.EvaluationSetVersion
	} else {
		var err error
		version, _, err = e.evaluationSetVersionService.GetEvaluationSetVersion(ctx, spaceID, evalSetVersionID, gptr.Of(true))
		if err != nil {
			return nil, err
		}
	}

	var fieldSchema []*entity.FieldSchema
	if version != nil && version.EvaluationSetSchema != nil {
		fieldSchema = version.EvaluationSetSchema.FieldSchemas
	}

	columnEvalSetFields := make([]*entity.ColumnEvalSetField, 0)
	for _, field := range fieldSchema {
		columnEvalSetFields = append(columnEvalSetFields, &entity.ColumnEvalSetField{
			Key:         gptr.Of(field.Key),
			Name:        gptr.Of(field.Name),
			Description: gptr.Of(field.Description),
			ContentType: field.ContentType,
		})
	}

	return columnEvalSetFields, nil
}

type PayloadBuilder struct {
	BaselineExptID       int64
	SpaceID              int64
	ExptIDs              []int64
	BaseExptTurnResultDO []*entity.ExptTurnResult

	ItemIDs   []int64
	TurnIDMap map[int64]bool

	ItemResults        []*entity.ItemResult // 最终结果
	ExptResultBuilders []*ExptResultBuilder // 每个实验的结果builder以及build result

	ExperimentRepo     repo.IExperimentRepo
	ExptTurnResultRepo repo.IExptTurnResultRepo

	EvaluationSetItemService EvaluationSetItemService
	EvalTargetService        IEvalTargetService
	EvaluatorRecordService   EvaluatorRecordService
}

func NewPayloadBuilder(ctx context.Context, param *entity.MGetExperimentResultParam, baselineExptID int64, baselineTurnResults []*entity.ExptTurnResult,
	baselineItemResults []*entity.ExptItemResult, experimentRepo repo.IExperimentRepo,
	exptTurnResultRepo repo.IExptTurnResultRepo,
	evalTargetService IEvalTargetService,
	evaluatorRecordService EvaluatorRecordService,
	evaluationSetItemService EvaluationSetItemService,
) *PayloadBuilder {
	builder := &PayloadBuilder{
		BaselineExptID:           baselineExptID,
		SpaceID:                  param.SpaceID,
		ExptIDs:                  param.ExptIDs,
		BaseExptTurnResultDO:     baselineTurnResults,
		ExperimentRepo:           experimentRepo,
		ExptTurnResultRepo:       exptTurnResultRepo,
		EvaluationSetItemService: evaluationSetItemService,
		EvalTargetService:        evalTargetService,
		EvaluatorRecordService:   evaluatorRecordService,
	}

	builder.ItemResults = make([]*entity.ItemResult, 0)

	// 需要分实验获取的数据范围
	itemIDs := make([]int64, 0)                              // itemID列表 有序
	itemID2TurnIDs := make(map[int64][]int64)                // itemID -> turnIDs列表 turnIDs有序
	itemIDMap := make(map[int64]bool)                        // 去重
	itemIDTurnIDTurnIndex := make(map[int64]map[int64]int64) // itemID -> turnID -> turnIndex
	itemIDItemResultPO := make(map[int64]*entity.ExptItemResult)

	turnIDMap := make(map[int64]bool)
	turnID2ItemID := make(map[int64]int64)

	for _, itemResult := range baselineItemResults {
		itemIDItemResultPO[itemResult.ItemID] = itemResult
	}

	for _, turnResultDO := range builder.BaseExptTurnResultDO {
		if _, ok := itemIDMap[turnResultDO.ItemID]; !ok {
			itemIDs = append(itemIDs, turnResultDO.ItemID) // 使用turnResultDO中的itemID append确保item有序
		}
		itemIDMap[turnResultDO.ItemID] = true

		if itemIDTurnIDTurnIndex[turnResultDO.ItemID] == nil {
			itemIDTurnIDTurnIndex[turnResultDO.ItemID] = make(map[int64]int64)
		}
		itemIDTurnIDTurnIndex[turnResultDO.ItemID][turnResultDO.TurnID] = int64(turnResultDO.TurnIdx)

		if turnResultDO.TurnID != 0 {
			turnIDMap[turnResultDO.TurnID] = true
			turnID2ItemID[turnResultDO.TurnID] = turnResultDO.ItemID
		}

		if _, ok := itemID2TurnIDs[turnResultDO.ItemID]; !ok {
			itemID2TurnIDs[turnResultDO.ItemID] = make([]int64, 0)
		}
		itemID2TurnIDs[turnResultDO.ItemID] = append(itemID2TurnIDs[turnResultDO.ItemID], turnResultDO.TurnID)
	}

	builder.ItemIDs = itemIDs
	builder.TurnIDMap = turnIDMap

	// 初始化payload结构
	for _, itemID := range itemIDs {
		if itemIDItemResultPO[itemID] == nil {
			continue
		}
		itemResultPO := itemIDItemResultPO[itemID]

		itemResult := &entity.ItemResult{
			ItemID:      itemID,
			TurnResults: make([]*entity.TurnResult, 0),
			ItemIndex:   gptr.Of(int64(itemResultPO.ItemIdx)),
			SystemInfo: &entity.ItemSystemInfo{ // 只有基准实验有ItemSystemInfo
				RunState: itemResultPO.Status,
				Error:    nil,
			},
		}

		for _, turnID := range itemID2TurnIDs[itemID] {
			turnIndex := int64(0)
			if itemIDTurnIDTurnIndex[itemID] != nil {
				turnIndex = itemIDTurnIDTurnIndex[itemID][turnID]
			}
			itemResult.TurnResults = append(itemResult.TurnResults, &entity.TurnResult{
				TurnID:            turnID,
				ExperimentResults: make([]*entity.ExperimentResult, 0),
				TurnIndex:         gptr.Of(turnIndex),
			})

		}

		builder.ItemResults = append(builder.ItemResults, itemResult)
	}

	return builder
}

// ExptResultBuilder 构建单实验结果
type ExptResultBuilder struct {
	ExptID                    int64
	BaselineExptID            int64
	SpaceID                   int64
	ItemIDs                   []int64        // 基准实验的itemID, 未匹配的不展示
	TurnIDMap                 map[int64]bool // 由于是itemID查询，对于多轮需要用turnID过滤. 对于单轮长度为0
	ItemIDTurnID2TurnResultID map[int64]map[int64]int64

	exptDO       *entity.Experiment
	turnResultDO []*entity.ExptTurnResult

	// 获取的结果
	turnResultID2EvaluatorVersionID2Result map[int64]map[int64]*entity.EvaluatorRecord // turn_result_id -> evaluator_version_id -> result
	turnResultID2TargetOutput              map[int64]*entity.TurnTargetOutput
	itemIDTurnID2Turn                      map[int64]map[int64]*entity.TurnEvalSet

	// 错误信息
	Err error

	ExperimentRepo     repo.IExperimentRepo
	ExptTurnResultRepo repo.IExptTurnResultRepo

	evaluationSetItemService EvaluationSetItemService
	evalTargetService        IEvalTargetService
	evaluatorRecordService   EvaluatorRecordService
}

// 1.确定当前分页下数据范围
// 2.分实验batch get 所需数据
// 3.组装数据
func (b *PayloadBuilder) Build(ctx context.Context) ([]*entity.ItemResult, error) {
	// 分实验获取数据
	exptResultBuilders := make([]*ExptResultBuilder, 0)
	for _, exptID := range b.ExptIDs {
		exptResultBuilder := &ExptResultBuilder{
			ExptID:                   exptID,
			BaselineExptID:           b.BaselineExptID,
			SpaceID:                  b.SpaceID,
			ItemIDs:                  b.ItemIDs,
			TurnIDMap:                b.TurnIDMap,
			ExperimentRepo:           b.ExperimentRepo,
			ExptTurnResultRepo:       b.ExptTurnResultRepo,
			evalTargetService:        b.EvalTargetService,
			evaluatorRecordService:   b.EvaluatorRecordService,
			evaluationSetItemService: b.EvaluationSetItemService,
		}

		if exptID == b.BaselineExptID {
			// 不用重复获取基准实验的数据
			exptResultBuilder.turnResultDO = b.BaseExptTurnResultDO
		}

		exptResultBuilders = append(exptResultBuilders, exptResultBuilder)
	}

	var wg sync.WaitGroup
	resultCh := make(chan *ExptResultBuilder, len(exptResultBuilders)) // 缓冲通道，收集结果和错误

	for _, exptResultBuilder := range exptResultBuilders {
		wg.Add(1)
		go func(builder *ExptResultBuilder) {
			defer wg.Done()
			defer goroutine.Recovery(ctx)
			err := builder.build(ctx)
			builder.Err = err
			resultCh <- builder
		}(exptResultBuilder)
	}

	wg.Wait()
	close(resultCh)

	var exptResultBuildersWithResult []*ExptResultBuilder // 任务结果
	var errors []error
	for exptResultBuilder := range resultCh {
		if exptResultBuilder.Err != nil {
			errors = append(errors, fmt.Errorf("ExptID %d: %v", exptResultBuilder.ExptID, exptResultBuilder.Err))
		} else {
			exptResultBuildersWithResult = append(exptResultBuildersWithResult, exptResultBuilder)
		}
	}
	if len(errors) > 0 {
		logs.CtxError(ctx, "build expt result fail, errors:%v", errors)
		return nil, fmt.Errorf("build expt result fail, errors:%v", errors)
	}

	b.ExptResultBuilders = exptResultBuildersWithResult

	// 填充数据
	err := b.fillData(ctx)
	if err != nil {
		return nil, err
	}

	return b.ItemResults, nil
}

func (b *PayloadBuilder) fillData(ctx context.Context) error {
	for i := range b.ItemResults {
		itemResult := b.ItemResults[i]
		itemID := itemResult.ItemID
		for j := range itemResult.TurnResults {
			turnResult := itemResult.TurnResults[j]
			if turnResult.ExperimentResults == nil {
				turnResult.ExperimentResults = make([]*entity.ExperimentResult, 0)
			}

			turnID := turnResult.TurnID

			for _, exptResultBuilder := range b.ExptResultBuilders {
				exptID := exptResultBuilder.ExptID
				exptResult := &entity.ExperimentResult{
					ExperimentID: exptID,
					Payload:      &entity.ExperimentTurnPayload{},
				}
				exptResult.Payload.TurnID = turnID
				exptResult.Payload.EvaluatorOutput = exptResultBuilder.getTurnEvaluatorResult(ctx, itemID, turnID)
				exptResult.Payload.EvalSet = exptResultBuilder.getTurnEvalSet(ctx, itemID, turnID)
				exptResult.Payload.TargetOutput = exptResultBuilder.getTurnTargetOutput(ctx, itemID, turnID)
				exptResult.Payload.SystemInfo = exptResultBuilder.getTurnSystemInfo(ctx, itemID, turnID)

				itemResult.TurnResults[j].ExperimentResults = append(itemResult.TurnResults[j].ExperimentResults, exptResult)
			}
		}
	}

	return nil
}

func (e *ExptResultBuilder) build(ctx context.Context) error {
	exptDO, err := e.ExperimentRepo.GetByID(ctx, e.ExptID, e.SpaceID)
	if err != nil {
		return err
	}
	e.exptDO = exptDO

	// 查询非基准实验的, turn_result. 基准实验跳过查询
	if e.ExptID != e.BaselineExptID {
		// 单轮的turnID始终是0
		// 索引（space_id, expt_id, item_id, turn_id）用item_id查询后过滤
		itemTurnResults, err := e.ExptTurnResultRepo.BatchGet(ctx, e.SpaceID, e.ExptID, e.ItemIDs)
		if err != nil {
			return err
		}

		// 由于是itemID查询，对于多轮需要用turnID过滤
		turnResults := make([]*entity.ExptTurnResult, 0)
		for _, itemTurnResult := range itemTurnResults {
			if itemTurnResult.TurnID == 0 {
				turnResults = append(turnResults, itemTurnResult)
				continue
			}
			if len(e.TurnIDMap) > 0 && e.TurnIDMap[itemTurnResult.ItemID] {
				turnResults = append(turnResults, itemTurnResult)
			}
		}
		e.turnResultDO = turnResults
	}

	if len(e.turnResultDO) == 0 {
		return nil
	}

	// 由于turnID可能为0，以turn_result_id为行的唯一标识聚合数据，组装payload数据时再通过turn_result_id与item_id(单轮)或turn_id(多轮)映射进行组装
	e.ItemIDTurnID2TurnResultID = make(map[int64]map[int64]int64) // itemID -> turnID -> turn_result_id
	for _, turnResult := range e.turnResultDO {
		if e.ItemIDTurnID2TurnResultID[turnResult.ItemID] == nil {
			e.ItemIDTurnID2TurnResultID[turnResult.ItemID] = make(map[int64]int64)
		}
		e.ItemIDTurnID2TurnResultID[turnResult.ItemID][turnResult.TurnID] = turnResult.ID
	}

	err = e.buildEvaluatorResult(ctx)
	if err != nil {
		return err
	}
	err = e.buildEvalSet(ctx)
	if err != nil {
		return err
	}
	err = e.buildTargetOutput(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (e *ExptResultBuilder) buildEvaluatorResult(ctx context.Context) error {
	turnResultIDs := make([]int64, 0)
	for _, turnResult := range e.turnResultDO {
		turnResultIDs = append(turnResultIDs, turnResult.ID)
	}
	turnEvaluatorResultRefs, err := e.ExptTurnResultRepo.BatchGetTurnEvaluatorResultRef(ctx, e.SpaceID, turnResultIDs)
	if err != nil {
		return err
	}

	evaluatorResultIDs := make([]int64, 0)
	evaluatorResultID2TurnResultID := make(map[int64]int64)
	for _, turnEvaluatorResultRef := range turnEvaluatorResultRefs {
		evaluatorResultIDs = append(evaluatorResultIDs, turnEvaluatorResultRef.EvaluatorResultID)

		evaluatorResultID2TurnResultID[turnEvaluatorResultRef.EvaluatorResultID] = turnEvaluatorResultRef.ExptTurnResultID
	}

	evaluatorRecords, err := e.evaluatorRecordService.BatchGetEvaluatorRecord(ctx, evaluatorResultIDs, false)
	if err != nil {
		return err
	}

	turnResultID2VersionID2Result := make(map[int64]map[int64]*entity.EvaluatorRecord) // turn_result_id -> version_id -> result
	for _, evaluatorRecord := range evaluatorRecords {
		turnResultID, ok := evaluatorResultID2TurnResultID[evaluatorRecord.ID]
		if !ok {
			continue
		}
		if _, ok := turnResultID2VersionID2Result[turnResultID]; !ok {
			turnResultID2VersionID2Result[turnResultID] = make(map[int64]*entity.EvaluatorRecord)
		}
		turnResultID2VersionID2Result[turnResultID][evaluatorRecord.EvaluatorVersionID] = evaluatorRecord
	}

	e.turnResultID2EvaluatorVersionID2Result = turnResultID2VersionID2Result

	return nil
}

func (e *ExptResultBuilder) getTurnEvaluatorResult(ctx context.Context, itemID, turnID int64) *entity.TurnEvaluatorOutput {
	turnID2TurnResultID, ok := e.ItemIDTurnID2TurnResultID[itemID]
	if !ok {
		return &entity.TurnEvaluatorOutput{}
	}
	turnResultID, ok := turnID2TurnResultID[turnID]
	if !ok {
		return &entity.TurnEvaluatorOutput{}
	}

	evaluatorVersionID2Result, ok := e.turnResultID2EvaluatorVersionID2Result[turnResultID]
	if !ok {
		return &entity.TurnEvaluatorOutput{}
	}

	for _, evaluatorResult := range evaluatorVersionID2Result {
		if evaluatorResult == nil {
			continue
		}

		if evaluatorResult.EvaluatorOutputData != nil && evaluatorResult.EvaluatorOutputData.EvaluatorRunError != nil {
			evaluatorResult.EvaluatorOutputData.EvaluatorRunError.Message = errno.ServiceInternalErrMsg
		}
	}

	return &entity.TurnEvaluatorOutput{
		EvaluatorRecords: evaluatorVersionID2Result,
	}
}

func (e *ExptResultBuilder) buildEvalSet(ctx context.Context) error {
	if e.exptDO == nil {
		return fmt.Errorf("exptPO is nil")
	}
	evalSetID := e.exptDO.EvalSetID
	evalSetVersionID := e.exptDO.EvalSetVersionID

	param := &entity.BatchGetEvaluationSetItemsParam{
		SpaceID:         e.SpaceID,
		EvaluationSetID: evalSetID,
		ItemIDs:         e.ItemIDs,
	}
	if evalSetVersionID != evalSetID {
		param.VersionID = gptr.Of(evalSetVersionID)
	}

	items, err := e.evaluationSetItemService.BatchGetEvaluationSetItems(ctx, param)
	if err != nil {
		return err
	}

	itemIDTurnID2Turn := make(map[int64]map[int64]*entity.TurnEvalSet) // item_id -> turn_id -> turn
	for _, item := range items {
		for _, turn := range item.Turns {
			if itemIDTurnID2Turn[item.ItemID] == nil {
				itemIDTurnID2Turn[item.ItemID] = make(map[int64]*entity.TurnEvalSet)
			}
			turnEvalSet := &entity.TurnEvalSet{
				Turn: turn,
			}
			itemIDTurnID2Turn[item.ItemID][turn.ID] = turnEvalSet
		}
	}

	e.itemIDTurnID2Turn = itemIDTurnID2Turn

	return nil
}

func (e *ExptResultBuilder) getTurnEvalSet(ctx context.Context, itemID, turnID int64) *entity.TurnEvalSet {
	turnID2Turn, ok := e.itemIDTurnID2Turn[itemID]
	if !ok {
		return &entity.TurnEvalSet{}
	}
	turn, ok := turnID2Turn[turnID]
	if !ok {
		return &entity.TurnEvalSet{}
	}

	return turn
}

func (e *ExptResultBuilder) buildTargetOutput(ctx context.Context) error {
	if e.exptDO.ExptType == entity.ExptType_Online {
		return nil
	}
	targetResultIDs := make([]int64, 0)
	targetResultID2turnResultID := make(map[int64]int64)
	for _, turnResult := range e.turnResultDO {
		targetResultIDs = append(targetResultIDs, turnResult.TargetResultID)
		targetResultID2turnResultID[turnResult.TargetResultID] = turnResult.ID
	}
	targetRecords, err := e.evalTargetService.BatchGetRecordByIDs(ctx, e.SpaceID, targetResultIDs)
	if err != nil {
		return err
	}

	turnResultID2TargetOutput := make(map[int64]*entity.TurnTargetOutput) // turn_result_id -> version_id -> result
	for _, targetRecord := range targetRecords {
		turnResultID, ok := targetResultID2turnResultID[targetRecord.ID]
		if !ok {
			continue
		}

		turnResultID2TargetOutput[turnResultID] = &entity.TurnTargetOutput{
			EvalTargetRecord: targetRecord,
		}
	}

	e.turnResultID2TargetOutput = turnResultID2TargetOutput

	return nil
}

func (e *ExptResultBuilder) getTurnTargetOutput(ctx context.Context, itemID, turnID int64) *entity.TurnTargetOutput {
	if e.exptDO.ExptType == entity.ExptType_Online {
		return &entity.TurnTargetOutput{}
	}
	turnID2TurnResultID, ok := e.ItemIDTurnID2TurnResultID[itemID]
	if !ok {
		return &entity.TurnTargetOutput{}
	}
	turnResultID, ok := turnID2TurnResultID[turnID]
	if !ok {
		return &entity.TurnTargetOutput{}
	}

	turnTargetOutput, ok := e.turnResultID2TargetOutput[turnResultID]
	if !ok {
		return &entity.TurnTargetOutput{}
	}

	if turnTargetOutput.EvalTargetRecord != nil && turnTargetOutput.EvalTargetRecord.EvalTargetOutputData != nil && turnTargetOutput.EvalTargetRecord.EvalTargetOutputData.EvalTargetRunError != nil {
		turnTargetOutput.EvalTargetRecord.EvalTargetOutputData.EvalTargetRunError.Message = errno.ServiceInternalErrMsg
	}

	return turnTargetOutput
}

func (e *ExptResultBuilder) getTurnSystemInfo(ctx context.Context, itemID, turnID int64) *entity.TurnSystemInfo {
	turnResultID2TurnResult := gslice.ToMap(e.turnResultDO, func(t *entity.ExptTurnResult) (int64, *entity.ExptTurnResult) {
		return t.ID, t
	})

	turnID2TurnResultID, ok := e.ItemIDTurnID2TurnResultID[itemID]
	if !ok {
		return &entity.TurnSystemInfo{}
	}
	turnResultID, ok := turnID2TurnResultID[turnID]
	if !ok {
		return &entity.TurnSystemInfo{}
	}

	turnResult, ok := turnResultID2TurnResult[turnResultID]
	if !ok {
		return &entity.TurnSystemInfo{}
	}

	systemInfo := &entity.TurnSystemInfo{
		TurnRunState: entity.TurnRunState(turnResult.Status),
		LogID:        gptr.Of(turnResult.LogID),
	}

	if len(turnResult.ErrMsg) > 0 {
		// 仅吐出评估器和评估对象之外的error
		ok, errMsg := errno.ParseTurnOtherErr(errno.DeserializeErr([]byte(turnResult.ErrMsg)))
		if ok {
			systemInfo.Error = &entity.RunError{
				Detail: gptr.Of(errMsg),
			}
		}
	}

	return systemInfo
}

func (e ExptResultServiceImpl) MGetStats(ctx context.Context, exptIDs []int64, spaceID int64, session *entity.Session) ([]*entity.ExptStats, error) {
	models, err := e.ExptStatsRepo.MGet(ctx, exptIDs, spaceID)
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (e ExptResultServiceImpl) GetStats(ctx context.Context, exptID int64, spaceID int64, session *entity.Session) (*entity.ExptStats, error) {
	stats, err := e.MGetStats(ctx, []int64{exptID}, spaceID, session)
	if err != nil {
		return nil, err
	}
	return stats[0], nil
}

func (e ExptResultServiceImpl) CreateStats(ctx context.Context, exptStats *entity.ExptStats, session *entity.Session) error {
	return e.ExptStatsRepo.Create(ctx, exptStats)
}

func (e ExptResultServiceImpl) CalculateStats(ctx context.Context, exptID, spaceID int64, session *entity.Session) (*entity.ExptCalculateStats, error) {
	var (
		maxLoop = 10000
		limit   = 100
		offset  = 1
		total   = 0
		cnt     = 0

		pendingCnt      = 0
		failCnt         = 0
		successCnt      = 0
		processingCnt   = 0
		terminatedCnt   = 0
		incompleteTurns []*entity.ItemTurnID
	)

	for i := 0; i < maxLoop; i++ {
		logs.CtxInfo(ctx, "ExptStatsImpl.CalculateStats scan turn result, expt_id: %v, page: %v, limit: %v, cur_cnt: %v, total: %v",
			exptID, offset, limit, cnt, total)

		results, t, err := e.ExptTurnResultRepo.ListTurnResult(ctx, spaceID, exptID, nil, entity.NewPage(offset, limit), false)
		if err != nil {
			return nil, err
		}

		total = int(t)
		cnt += len(results)
		offset++

		for _, tr := range results {
			switch entity.TurnRunState(tr.Status) {
			case entity.TurnRunState_Success:
				successCnt++
			case entity.TurnRunState_Fail:
				failCnt++
			case entity.TurnRunState_Terminal:
				terminatedCnt++
			case entity.TurnRunState_Queueing:
				pendingCnt++
				incompleteTurns = append(incompleteTurns, &entity.ItemTurnID{
					TurnID: tr.TurnID,
					ItemID: tr.ItemID,
				})
			case entity.TurnRunState_Processing:
				processingCnt++
				incompleteTurns = append(incompleteTurns, &entity.ItemTurnID{
					TurnID: tr.TurnID,
					ItemID: tr.ItemID,
				})
			default:
			}
		}

		if cnt >= total || len(results) == 0 {
			break
		}

		time.Sleep(time.Millisecond * 30)
	}

	stats := &entity.ExptCalculateStats{
		PendingTurnCnt:    pendingCnt,
		FailTurnCnt:       failCnt,
		SuccessTurnCnt:    successCnt,
		ProcessingTurnCnt: processingCnt,
		TerminatedTurnCnt: terminatedCnt,
	}

	logs.CtxInfo(ctx, "ExptStatsImpl.CalculateStats scan turn result done, expt_id: %v, total_cnt: %v, incomplete_cnt: %v, total: %v, stats: %v", exptID, cnt, len(incompleteTurns), total, json.Jsonify(stats))

	return stats, nil
}
