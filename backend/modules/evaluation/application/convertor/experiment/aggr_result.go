// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package experiment

import (
	"math"

	"github.com/bytedance/gg/gptr"

	domain_expt "github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/evaluation/domain/expt"
	"github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/entity"
)

func ExptAggregateResultDOToDTO(data *entity.ExptAggregateResult) *domain_expt.ExptAggregateResult_ {
	if data == nil {
		return nil
	}

	evaluatorResults := make(map[int64]*domain_expt.EvaluatorAggregateResult_, len(data.EvaluatorResults))
	for evaluatorVersionID, evaluatorResult := range data.EvaluatorResults {
		evaluatorResults[evaluatorVersionID] = EvaluatorResultsDOToDTO(evaluatorResult)
	}

	return &domain_expt.ExptAggregateResult_{
		ExperimentID:     data.ExperimentID,
		EvaluatorResults: evaluatorResults,
		Status:           domain_expt.ExptAggregateCalculateStatusPtr(domain_expt.ExptAggregateCalculateStatus(data.Status)),
	}
}

func EvaluatorResultsDOToDTO(result *entity.EvaluatorAggregateResult) *domain_expt.EvaluatorAggregateResult_ {
	if result == nil {
		return nil
	}

	return &domain_expt.EvaluatorAggregateResult_{
		EvaluatorVersionID: result.EvaluatorVersionID,
		AggregatorResults:  AggregatorResultDOsToDTOs(result.AggregatorResults),
		Name:               result.Name,
		Version:            result.Version,
	}
}

func AggregatorResultDOsToDTOs(result []*entity.AggregatorResult) []*domain_expt.AggregatorResult_ {
	if len(result) == 0 {
		return nil
	}
	results := make([]*domain_expt.AggregatorResult_, 0, len(result))
	for _, r := range result {
		results = append(results, AggregatorResultDOToDTO(r))
	}
	return results
}

func AggregatorResultDOToDTO(result *entity.AggregatorResult) *domain_expt.AggregatorResult_ {
	if result == nil {
		return nil
	}

	return &domain_expt.AggregatorResult_{
		AggregatorType: domain_expt.AggregatorType(result.AggregatorType),
		Data:           AggregateDataDOToDTO(result.Data),
	}
}

func AggregateDataDOToDTO(data *entity.AggregateData) *domain_expt.AggregateData {
	if data == nil {
		return nil
	}

	aggregateData := &domain_expt.AggregateData{
		DataType: domain_expt.DataType(data.DataType),
	}

	if data.Value != nil {
		aggregateData.Value = gptr.Of(math.Round(*data.Value*100) / 100)
	}

	if data.ScoreDistribution != nil {
		aggregateData.ScoreDistribution = &domain_expt.ScoreDistribution{
			ScoreDistributionItems: ScoreDistributionItemsDOToDTO(data.ScoreDistribution.ScoreDistributionItems),
		}
	}

	return aggregateData
}

func ScoreDistributionItemsDOToDTO(data []*entity.ScoreDistributionItem) []*domain_expt.ScoreDistributionItem {
	if len(data) == 0 {
		return nil
	}

	items := make([]*domain_expt.ScoreDistributionItem, 0, len(data))
	for _, item := range data {
		if item == nil {
			continue
		}
		items = append(items, &domain_expt.ScoreDistributionItem{
			Score:      item.Score,
			Count:      item.Count,
			Percentage: item.Percentage,
		})
	}

	return items
}
