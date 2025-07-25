// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/gg/gptr"

	"github.com/coze-dev/cozeloop/backend/infra/backoff"
	"github.com/coze-dev/cozeloop/backend/infra/external/audit"
	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	"github.com/coze-dev/cozeloop/backend/infra/lock"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/component"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/component/idem"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/component/metrics"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/events"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/repo"
	"github.com/coze-dev/cozeloop/backend/pkg/ctxcache"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/json"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/goroutine"
	gslice "github.com/coze-dev/cozeloop/backend/pkg/lang/slices"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

type ExptSchedulerImpl struct {
	Manager                  IExptManager
	ExptRepo                 repo.IExperimentRepo
	Publisher                events.ExptEventPublisher
	ExptItemResultRepo       repo.IExptItemResultRepo
	ExptTurnResultRepo       repo.IExptTurnResultRepo
	ExptStatsRepo            repo.IExptStatsRepo
	ExptRunLogRepo           repo.IExptRunLogRepo
	Idem                     idem.IdempotentService
	Configer                 component.IConfiger
	QuotaRepo                repo.QuotaRepo
	Mutex                    lock.ILocker
	AuditClient              audit.IAuditService
	Metric                   metrics.ExptMetric
	Endpoints                SchedulerEndPoint
	ResultSvc                ExptResultService
	IDGen                    idgen.IIDGenerator
	evaluationSetItemService EvaluationSetItemService
	schedulerModeFactory     SchedulerModeFactory
}

func NewExptSchedulerSvc(
	manager IExptManager,
	exptRepo repo.IExperimentRepo,
	exptItemResultRepo repo.IExptItemResultRepo,
	exptTurnResultRepo repo.IExptTurnResultRepo,
	exptStatsRepo repo.IExptStatsRepo,
	exptRunLogRepo repo.IExptRunLogRepo,
	Idem idem.IdempotentService,
	configer component.IConfiger,
	quotaRepo repo.QuotaRepo,
	mutex lock.ILocker,
	publisher events.ExptEventPublisher,
	auditClient audit.IAuditService,
	metric metrics.ExptMetric,
	resultSvc ExptResultService,
	idGen idgen.IIDGenerator,
	evaluationSetItemService EvaluationSetItemService,
	schedulerModeFactory SchedulerModeFactory,
) ExptSchedulerEvent {
	i := &ExptSchedulerImpl{
		Manager:                  manager,
		ExptRepo:                 exptRepo,
		ExptItemResultRepo:       exptItemResultRepo,
		ExptTurnResultRepo:       exptTurnResultRepo,
		ExptStatsRepo:            exptStatsRepo,
		ExptRunLogRepo:           exptRunLogRepo,
		Idem:                     Idem,
		Configer:                 configer,
		QuotaRepo:                quotaRepo,
		Mutex:                    mutex,
		Publisher:                publisher,
		AuditClient:              auditClient,
		Metric:                   metric,
		ResultSvc:                resultSvc,
		IDGen:                    idGen,
		evaluationSetItemService: evaluationSetItemService,
		schedulerModeFactory:     schedulerModeFactory,
	}

	i.Endpoints = SchedulerChain(
		i.HandleEventErr,
		i.SysOps,
		i.HandleEventCheck,
		i.HandleEventLock,
		i.HandleEventEndpoint,
	)(func(_ context.Context, _ *entity.ExptScheduleEvent) error { return nil })

	return i
}

func (e *ExptSchedulerImpl) Schedule(ctx context.Context, event *entity.ExptScheduleEvent) error {
	ctx = ctxcache.Init(ctx)

	if err := e.Endpoints(ctx, event); err != nil {
		logs.CtxError(ctx, "[ExptScheduler] expt schedule fail, event: %v, err: %v", json.Jsonify(event), err)
		return err
	}

	return nil
}

type SchedulerEndPoint func(ctx context.Context, event *entity.ExptScheduleEvent) error

type SchedulerMiddleware func(next SchedulerEndPoint) SchedulerEndPoint

func SchedulerChain(mws ...SchedulerMiddleware) SchedulerMiddleware {
	return func(next SchedulerEndPoint) SchedulerEndPoint {
		for i := len(mws) - 1; i >= 0; i-- {
			next = mws[i](next)
		}
		return next
	}
}

func (e *ExptSchedulerImpl) SysOps(next SchedulerEndPoint) SchedulerEndPoint {
	return func(ctx context.Context, event *entity.ExptScheduleEvent) error {
		return next(ctx, event)
	}
}

func (e *ExptSchedulerImpl) HandleEventCheck(next SchedulerEndPoint) SchedulerEndPoint {
	return func(ctx context.Context, event *entity.ExptScheduleEvent) error {
		runLog, err := e.Manager.GetRunLog(ctx, event.ExptID, event.ExptRunID, event.SpaceID, event.Session)
		if err != nil {
			return err
		}

		if entity.IsExptFinished(entity.ExptStatus(runLog.Status)) {
			logs.CtxInfo(ctx, "ExptSchedulerConsumer consume finished expt run event, expt_id: %v, expt_run_id: %v", event.ExptID, event.ExptRunID)
			return nil
		}

		interval := int64(e.Configer.GetExptExecConf(ctx, event.SpaceID).GetZombieIntervalSecond())
		if time.Now().Unix()-event.CreatedAt >= interval {
			return fmt.Errorf("expt exec found timeout event, expt_id: %v, expt_run_id: %v", event.ExptID, event.ExptRunID)
		}

		return next(ctx, event)
	}
}

func (e *ExptSchedulerImpl) makeExptRunExecLockKey(exptID, exptRunID int64) string {
	return fmt.Sprintf("expt_run_exec_lock:%d:%d", exptID, exptRunID)
}

func (e *ExptSchedulerImpl) HandleEventLock(next SchedulerEndPoint) SchedulerEndPoint {
	return func(ctx context.Context, event *entity.ExptScheduleEvent) error {
		locked, ctx, unlock, err := e.Mutex.LockWithRenew(ctx, e.makeExptRunExecLockKey(event.ExptID, event.ExptRunID), time.Second*20, time.Second*60*3)
		if err != nil {
			return err
		}
		logs.CtxInfo(ctx, "ExptSchedulerConsumer.HandleEventLock locked expt eval event: %v", json.Jsonify(event))
		if !locked {
			logs.CtxWarn(ctx, "ExptSchedulerConsumer.HandleEventLock found locked expt eval event: %v. Abort event, err: %v", json.Jsonify(event), err)
			return nil
		}

		defer unlock()

		return next(ctx, event)
	}
}

func (e *ExptSchedulerImpl) HandleEventEndpoint(next SchedulerEndPoint) SchedulerEndPoint {
	return func(ctx context.Context, event *entity.ExptScheduleEvent) error {
		err := e.schedule(ctx, event)
		if err != nil {
			return err
		}

		return next(ctx, event)
	}
}

func (e *ExptSchedulerImpl) HandleEventErr(next SchedulerEndPoint) SchedulerEndPoint {
	return func(ctx context.Context, event *entity.ExptScheduleEvent) error {
		nextErr := func(ctx context.Context, event *entity.ExptScheduleEvent) (err error) {
			defer goroutine.Recover(ctx, &err)
			return next(ctx, event)
		}(ctx, event)

		if nextErr == nil {
			logs.CtxInfo(ctx, "[ExptEval] handle event success, event: %v", json.Jsonify(event))
			return nil
		}

		logs.CtxError(ctx, "[ExptEval] HandleEventErr found error: %v, event: %v", nextErr, json.Jsonify(event))

		completeCID := fmt.Sprintf("exptexec:onerr:%d", event.ExptRunID)

		if err := e.Manager.CompleteRun(ctx, event.ExptID, event.ExptRunID, event.ExptRunMode, event.SpaceID, event.Session, entity.WithCID(completeCID)); err != nil {
			return errorx.Wrapf(err, "terminate expt run fail, expt_id: %v, expt_run_id: %v", event.ExptID, event.ExptRunID)
		}

		if err := e.Manager.CompleteExpt(ctx, event.ExptID, event.SpaceID, event.Session, entity.WithStatus(entity.ExptStatus_Failed),
			entity.WithStatusMessage(nextErr.Error()), entity.WithCID(completeCID)); err != nil {
			return errorx.Wrapf(err, "complete expt fail, expt_id: %v, expt_run_id: %v", event.ExptID, event.ExptRunID)
		}

		return nil
	}
}

func (e *ExptSchedulerImpl) schedule(ctx context.Context, event *entity.ExptScheduleEvent) error {
	exptDetail, err := e.Manager.GetDetail(ctx, event.ExptID, event.SpaceID, event.Session)
	if err != nil {
		return err
	}

	mode, err := e.schedulerModeFactory.NewSchedulerMode(event.ExptRunMode)
	if err != nil {
		return err
	}

	err = mode.ExptStart(ctx, event, exptDetail)
	if err != nil {
		return err
	}

	err = mode.ScheduleStart(ctx, event, exptDetail)
	if err != nil {
		return err
	}

	toSubmit, incomplete, complete, err := mode.ScanEvalItems(ctx, event, exptDetail)
	if err != nil {
		return err
	}

	e.handleZombies(ctx, event, incomplete)

	if err = e.recordEvalItemRunLogs(ctx, event, complete); err != nil {
		return err
	}

	err = mode.ScheduleEnd(ctx, event, exptDetail, len(toSubmit), len(incomplete))
	if err != nil {
		return err
	}

	nextTick, err := mode.ExptEnd(ctx, event, exptDetail, len(toSubmit), len(incomplete))
	if err != nil {
		return err
	}

	if err = e.handleToSubmits(ctx, event, toSubmit); err != nil {
		return err
	}

	logs.CtxInfo(ctx, "[ExptEval] expt daemon with next tick, expt_id: %v, event: %v", event.ExptID, event)

	return mode.NextTick(ctx, event, nextTick)
}

func (e *ExptSchedulerImpl) recordEvalItemRunLogs(ctx context.Context, event *entity.ExptScheduleEvent, completeItems []*entity.ExptEvalItem) error {
	for _, item := range completeItems {
		if item.State != entity.ItemRunState_Fail && item.State != entity.ItemRunState_Success {
			return fmt.Errorf("recordEvalItemRunLogs found invalid item run state: %v", item.State)
		}
		if err := backoff.RetryFiveMin(ctx, func() error {
			return e.ResultSvc.RecordItemRunLogs(ctx, event.ExptID, event.ExptRunID, item.ItemID, event.SpaceID, event.Session)
		}); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 50)
	}
	return nil
}

func (e *ExptSchedulerImpl) handleToSubmits(ctx context.Context, event *entity.ExptScheduleEvent, toSubmits []*entity.ExptEvalItem) error {
	if len(toSubmits) == 0 {
		return nil
	}

	now := time.Now().Unix()
	itemIDs := make([]int64, 0, len(toSubmits))
	itemEvalEvents := make([]*entity.ExptItemEvalEvent, 0, len(toSubmits))
	for _, ts := range toSubmits {
		if entity.IsItemRunFinished(ts.State) {
			continue
		}
		itemIDs = append(itemIDs, ts.ItemID)
		itemEvalEvents = append(itemEvalEvents, &entity.ExptItemEvalEvent{
			SpaceID:       event.SpaceID,
			ExptID:        event.ExptID,
			ExptRunID:     event.ExptRunID,
			ExptRunMode:   event.ExptRunMode,
			EvalSetItemID: ts.ItemID,
			CreateAt:      now,
			Ext:           event.Ext,
			Session:       event.Session,
		})
	}

	logs.CtxInfo(ctx, "submit item eval events: %v", json.Jsonify(itemEvalEvents))

	interval := e.Configer.GetExptExecConf(ctx, event.SpaceID).GetExptItemEvalConf().GetInterval()
	if err := e.Publisher.BatchPublishExptRecordEvalEvent(ctx, itemEvalEvents, gptr.Of(interval)); err != nil {
		return err
	}

	defer e.Metric.EmitItemExecEval(event.SpaceID, int64(event.ExptRunMode), len(toSubmits))

	if err := e.ExptItemResultRepo.UpdateItemRunLog(ctx, event.ExptID, event.ExptRunID, itemIDs, map[string]any{"status": int32(entity.ItemRunState_Processing)},
		event.SpaceID); err != nil {
		return err
	}

	if err := e.ExptItemResultRepo.UpdateItemsResult(ctx, event.SpaceID, event.ExptID, itemIDs, map[string]any{"status": int32(entity.ItemRunState_Processing)}); err != nil {
		return err
	}

	if err := e.ExptTurnResultRepo.UpdateTurnResultsWithItemIDs(ctx, event.ExptID, itemIDs, event.SpaceID, map[string]any{"status": int32(entity.TurnRunState_Processing)}); err != nil {
		return err
	}

	turnResults, err := e.ExptTurnResultRepo.BatchGet(ctx, event.SpaceID, event.ExptID, itemIDs)
	if err != nil {
		return err
	}

	if err := e.ExptStatsRepo.ArithOperateCount(ctx, event.ExptID, event.SpaceID, &entity.StatsCntArithOp{
		OpStatusCnt: map[entity.TurnRunState]int{
			entity.TurnRunState_Processing: len(turnResults),
			entity.TurnRunState_Queueing:   0 - len(turnResults),
		},
	}); err != nil {
		return err
	}

	return nil
}

func (e *ExptSchedulerImpl) handleZombies(ctx context.Context, event *entity.ExptScheduleEvent, items []*entity.ExptEvalItem) {
	var (
		zombies      []*entity.ExptEvalItem
		zombieSecond = e.Configer.GetConsumerConf(ctx).GetExptExecConf(event.SpaceID).GetExptItemEvalConf().GetZombieSecond()
	)

	for _, item := range items {
		if item.State == entity.ItemRunState_Processing && item.UpdatedAt != nil && !gptr.Indirect(item.UpdatedAt).IsZero() {
			if time.Since(gptr.Indirect(item.UpdatedAt)).Seconds() > float64(zombieSecond) {
				zombies = append(zombies, item)
				continue
			}
		}
	}

	if len(zombies) > 0 {
		logs.CtxWarn(ctx, "[ExptEval] found zombie items: %v, expt_id: %v, expt_run_id: %v",
			gslice.Transform(zombies, func(e *entity.ExptEvalItem, _ int) int64 { return e.ItemID }), event.ExptID, event.ExptRunID)
	}
}
