// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package entity

import (
	"context"
	"time"

	"github.com/coze-dev/coze-loop/backend/infra/middleware/session"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
	"github.com/coze-dev/coze-loop/backend/pkg/json"
)

type FieldType int64

const (
	FieldType_Unknown FieldType = 0
	// 评估器得分, FieldKey为evaluatorVersionID,value为score
	FieldType_EvaluatorScore     FieldType = 1
	FieldType_CreatorBy          FieldType = 2
	FieldType_ExptStatus         FieldType = 3
	FieldType_TurnRunState       FieldType = 4
	FieldType_TargetID           FieldType = 5
	FieldType_EvalSetID          FieldType = 6
	FieldType_EvaluatorID        FieldType = 7
	FieldType_TargetType         FieldType = 8
	FieldType_SourceTarget       FieldType = 9
	FieldType_EvaluatorVersionID FieldType = 20
	FieldType_TargetVersionID    FieldType = 21
	FieldType_EvalSetVersionID   FieldType = 22
)

// aggregate result
type UpdateExptAggrResultParam struct {
	SpaceID      int64
	ExperimentID int64
	FieldType    FieldType
	FieldKey     string
}

// AggregatorType 聚合器类型
type AggregatorType int

const (
	Average      AggregatorType = 1
	Sum          AggregatorType = 2
	Max          AggregatorType = 3
	Min          AggregatorType = 4
	Distribution AggregatorType = 5 // 得分的分布情况
)

type AggrResultDataType int

const (
	Double            AggrResultDataType = 0 // 默认，有小数的浮点数值类型
	ScoreDistribution AggrResultDataType = 1 // 得分分布
)

type ScoreDistributionData struct {
	ScoreDistributionItems []*ScoreDistributionItem
}

type ScoreDistributionItem struct {
	Score      string  // 得分,TOP5以外的聚合展示为“其他”
	Count      int64   // 此得分的数量
	Percentage float64 // 占总数的百分比
}

type AggregateData struct {
	DataType          AggrResultDataType
	Value             *float64
	ScoreDistribution *ScoreDistributionData
}

// AggregatorResult 一种聚合器类型的聚合结果
type AggregatorResult struct {
	AggregatorType AggregatorType
	Data           *AggregateData
}

// expt_aggr_result 表 aggr_result 字段blob结构
type AggregateResult struct {
	AggregatorResults []*AggregatorResult
}

func (a AggregatorResult) GetScore() float64 {
	if a.Data == nil {
		return 0
	}
	if a.Data.Value == nil {
		return 0
	}

	return *a.Data.Value
}

type ExptAggrResult struct {
	ID           int64
	SpaceID      int64
	ExperimentID int64
	FieldType    int32
	FieldKey     string
	Score        float64
	AggrResult   []byte
	Version      int64
	Status       int32
}

type ExptAggregateResult struct {
	ExperimentID     int64
	EvaluatorResults map[int64]*EvaluatorAggregateResult
	Status           int64
}

type EvaluatorAggregateResult struct {
	EvaluatorVersionID int64
	AggregatorResults  []*AggregatorResult
	Name               *string
	Version            *string
}

// item result
type ExptItemResult struct {
	ID        int64
	SpaceID   int64
	ExptID    int64
	ExptRunID int64
	ItemID    int64
	Status    ItemRunState
	ErrMsg    string
	ItemIdx   int32
	LogID     string
}

type ExptItemResultRunLog struct {
	ID          int64
	SpaceID     int64
	ExptID      int64
	ExptRunID   int64
	ItemID      int64
	Status      int32
	ErrMsg      []byte
	LogID       string
	ResultState int32
	UpdatedAt   *time.Time
}

type ExptItemEvalResult struct {
	ItemResultRunLog  *ExptItemResultRunLog
	TurnResultRunLogs map[int64]*ExptTurnResultRunLog
}

type ExptEvalItem struct {
	ExptID           int64
	EvalSetVersionID int64
	ItemID           int64
	State            ItemRunState
	UpdatedAt        *time.Time
}

type ExptEvalTurn struct {
	ExptID    int64
	ExptRunID int64
	ItemID    int64
	TurnID    int64
}

type ExptStats struct {
	ID                int64
	SpaceID           int64
	ExptID            int64
	PendingTurnCnt    int32
	SuccessTurnCnt    int32
	FailTurnCnt       int32
	ProcessingTurnCnt int32
	TerminatedTurnCnt int32
	CreditCost        float64
	InputTokenCost    int64
	OutputTokenCost   int64
}

type ExptTurnResult struct {
	ID               int64
	SpaceID          int64
	ExptID           int64
	ExptRunID        int64
	ItemID           int64
	TurnID           int64
	Status           int32
	TraceID          int64
	LogID            string
	TargetResultID   int64
	EvaluatorResults *EvaluatorResults
	ErrMsg           string
	TurnIdx          int32
}

func (tr *ExptTurnResult) ToRunLogDO() *ExptTurnResultRunLog {
	if tr == nil {
		return nil
	}
	return &ExptTurnResultRunLog{
		ID:                 tr.ID,
		SpaceID:            tr.SpaceID,
		ExptID:             tr.ExptID,
		ExptRunID:          tr.ExptRunID,
		ItemID:             tr.ItemID,
		TurnID:             tr.TurnID,
		Status:             TurnRunState(tr.Status),
		TraceID:            tr.TraceID,
		LogID:              tr.LogID,
		TargetResultID:     tr.TargetResultID,
		EvaluatorResultIds: tr.EvaluatorResults,
		ErrMsg:             tr.ErrMsg,
	}
}

type EvaluatorResults struct {
	EvalVerIDToResID map[int64]int64
}

func (e *EvaluatorResults) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil, errorx.Wrapf(err, "ExptTurnEvaluatorResultIDs json marshal fail")
	}
	return bytes, nil
}

type MGetExperimentResultParam struct {
	SpaceID    int64
	ExptIDs    []int64
	BaseExptID *int64
	Filters    map[int64]*ExptTurnResultFilter
	Page       Page
}

type ExptTurnResultRunLog struct {
	ID                 int64
	SpaceID            int64
	ExptID             int64
	ExptRunID          int64
	ItemID             int64
	TurnID             int64
	Status             TurnRunState
	TraceID            int64
	LogID              string
	TargetResultID     int64
	EvaluatorResultIds *EvaluatorResults
	ErrMsg             string
}

type ExptTurnEvaluatorResultRef struct {
	ID                 int64
	SpaceID            int64
	ExptTurnResultID   int64
	EvaluatorVersionID int64
	EvaluatorResultID  int64
	ExptID             int64
}

type ExptEvaluatorRef struct {
	ID                 int64
	SpaceID            int64
	ExptID             int64
	EvaluatorID        int64
	EvaluatorVersionID int64
}

// filter
type ExptListFilter struct {
	FuzzyName string
	Includes  *ExptFilterFields
	Excludes  *ExptFilterFields
}

type ExptFilterFields struct {
	CreatedBy    []string
	Status       []int64
	EvalSetIDs   []int64
	TargetIDs    []int64
	EvaluatorIDs []int64
	TargetType   []int64
	ExptType     []int64
	SourceType   []int64
	SourceID     []string
}

func (e *ExptFilterFields) IsValid() bool {
	if e == nil {
		return true
	}
	for _, slice := range [][]int64{e.Status, e.EvalSetIDs, e.TargetIDs, e.EvaluatorIDs, e.TargetType} {
		for _, item := range slice {
			if item < 0 {
				return false
			}
		}
	}
	for _, item := range e.CreatedBy {
		if len(item) <= 0 {
			return false
		}
	}
	return true
}

type ExptItemRunLogFilter struct {
	Status      []ItemRunState
	ResultState *ExptItemResultState
}

func (e *ExptItemRunLogFilter) GetResultState() ExptItemResultState {
	return *e.ResultState
}

func (e *ExptItemRunLogFilter) GetStatus() []int32 {
	res := make([]int32, 0, len(e.Status))
	for _, status := range e.Status {
		res = append(res, int32(status))
	}
	return res
}

const (
	defaultPage     = 1 // 页数从 1 开始
	defaultLimit    = 20
	defaultMaxLimit = 200 // 分页最大限制
)

type Page struct {
	offset int
	limit  int
}

func NewPage(offset, limit int) Page {
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > defaultMaxLimit {
		limit = defaultMaxLimit
	}
	if offset <= 0 {
		offset = defaultPage
	}

	return Page{
		offset: offset,
		limit:  limit,
	}
}

func (p Page) Offset() int {
	return (p.offset - 1) * p.limit
}

func (p Page) Limit() int {
	return p.limit
}

type Session struct {
	UserID string
	AppID  int32
}

func NewSession(ctx context.Context) *Session {
	userIDInContext := session.UserIDInCtxOrEmpty(ctx)
	return &Session{
		UserID: userIDInContext,
	}
}

type ExptTurnResultFilter struct {
	TrunRunStateFilters []*TurnRunStateFilter
	ScoreFilters        []*ScoreFilter
}

type ScoreFilter struct {
	Score              float64
	Operator           string
	EvaluatorVersionID int64
}

type TurnRunStateFilter struct {
	Status   []TurnRunState
	Operator string
}

type TurnTargetOutput struct {
	EvalTargetRecord *EvalTargetRecord
}

type TurnEvaluatorOutput struct {
	EvaluatorRecords map[int64]*EvaluatorRecord
}

type TurnEvalSet struct {
	Turn *Turn
}

type TurnSystemInfo struct {
	TurnRunState TurnRunState
	LogID        *string
	Error        *RunError
}

type ItemSystemInfo struct {
	RunState ItemRunState
	LogID    *string
	Error    *RunError
}

type RunError struct {
	Code    int64
	Message *string
	Detail  *string
}

type ItemResult struct {
	ItemID int64
	// row粒度实验结果详情
	TurnResults []*TurnResult
	SystemInfo  *ItemSystemInfo
	ItemIndex   *int64
}

type ExperimentTurnPayload struct {
	TurnID int64
	// 评测数据集数据
	EvalSet *TurnEvalSet
	// 评测对象结果
	TargetOutput *TurnTargetOutput
	// 评测规则执行结果
	EvaluatorOutput *TurnEvaluatorOutput
	// 评测系统相关数据日志、error
	SystemInfo *TurnSystemInfo
}

type ExperimentResult struct {
	ExperimentID int64
	Payload      *ExperimentTurnPayload
}

type TurnResult struct {
	TurnID int64
	// 参与对比的实验序列，对于单报告序列长度为1
	ExperimentResults []*ExperimentResult
	TurnIndex         *int64
}

type ColumnEvalSetField struct {
	Key         *string
	Name        *string
	Description *string
	ContentType ContentType
}

type ColumnEvaluator struct {
	EvaluatorVersionID int64
	EvaluatorID        int64
	EvaluatorType      EvaluatorType
	Name               *string
	Version            *string
	Description        *string
}
