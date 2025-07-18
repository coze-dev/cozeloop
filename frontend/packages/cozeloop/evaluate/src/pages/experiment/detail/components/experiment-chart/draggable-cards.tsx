// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useEffect, useState } from 'react';

import { get } from 'lodash-es';
import { type Datum } from '@visactor/vchart/esm/typings';
import { type ISpec } from '@visactor/vchart';
import {
  EvaluatorPreview,
  Chart,
  ChartCardItemRender,
  DraggableGrid,
  type CustomTooltipProps,
  type ChartCardItem,
} from '@cozeloop/evaluate-components';
import {
  AggregatorType,
  type ScoreDistributionItem,
  type Evaluator,
  type EvaluatorAggregateResult,
} from '@cozeloop/api-schema/evaluation';
import { IconCozIllusAdd } from '@coze-arch/coze-design/illustrations';
import { EmptyState } from '@coze-arch/coze-design';
import { I18n } from '@cozeloop/i18n-adapter';

const spec: ISpec = {
  type: 'pie',
  crosshair: {
    xField: { visible: true },
  },
  color: ['#6D62EB', '#00B2B2', '#3377FF', '#FFB829', '#CA61FF', '#7DD600'],
  valueField: 'score',
  categoryField: 'name',
  outerRadius: 0.8,
  innerRadius: 0,
  height: 200,
  legends: {
    visible: true,
    orient: 'left',
    item: {
      shape: {
        style: {
          symbolType: 'square',
        },
      },
    },
  },
  // percent: true,
  tooltip: {
    mark: {
      content: [
        {
          // @ts-expect-error type
          key: (datum: Datum) => datum.name,
          // @ts-expect-error type
          value: (datum: Datum) => datum.score,
        },
      ],
    },
  },
};

function getScorePercentage(score: number | undefined) {
  if (typeof score !== 'number') {
    return '-';
  }
  const percent = score * 100;
  if (percent % 1 === 0) {
    return `${percent}%`;
  }
  return `${percent?.toFixed(1)}%`;
}

function ComplexTooltipContent(props: CustomTooltipProps) {
  const { params, actualTooltip } = props;
  // 获取hover目标柱状图数据
  const datum: Datum | undefined = params?.datum?.item
    ? params?.datum
    : get(actualTooltip, 'data[0].data[0].datum[0]');
  const item: ScoreDistributionItem | undefined = datum?.item;
  const prefixBgColor = actualTooltip?.title?.shapeFill;
  if (!item) {
    return null;
  }
  return (
    <div className="w-[220px] flex flex-col gap-2">
      <div className="text-sm font-medium">{I18n.t('score_details')}</div>
      <div className="flex items-center gap-2 text-xs">
        <div className="w-2 h-2" style={{ backgroundColor: prefixBgColor }} />
        <span>
          {I18n.t('score')} {item.score ?? '-'}
        </span>
        <span className="font-semibold ml-auto">
          <span className="font-medium text-[var(--coz-fg-primary)]">
            {item.count ?? '-'}
          </span>
          <span className="text-[var(--coz-fg-secondary)]">
            {I18n.t('item_unit')} ({getScorePercentage(item.percentage)})
          </span>
        </span>
      </div>
    </div>
  );
}

/**
 * 提取实验评估器统计分数的分布（如1分的8个，0.6分的4个等）
 * 第一层map_key为评估器版本id
 * 第二层map key为得分，值为得分的数量
 */
function getEvaluatorScoreMap(results: EvaluatorAggregateResult[]) {
  const map: Record<
    Int64,
    Record<number | string, ScoreDistributionItem> | undefined
  > = {};
  results.forEach(result => {
    const versionId = result?.evaluator_version_id ?? '';
    result?.aggregator_results?.forEach(item => {
      if (item.aggregator_type !== AggregatorType.Distribution) {
        return;
      }
      item.data?.score_distribution?.score_distribution_items?.forEach(
        scoreItem => {
          if (!map[versionId]) {
            map[versionId] = {};
          }
          map[versionId][scoreItem.score] = scoreItem;
        },
      );
    });
  });
  return map;
}

export function EvaluatorsDraggableCard({
  spaceID,
  evaluators = [],
  evaluatorAggregateResult = [],
  ready,
}: {
  spaceID: Int64;
  evaluators: Evaluator[];
  evaluatorAggregateResult: EvaluatorAggregateResult[];
  ready?: boolean;
}) {
  const [items, setItems] = useState<ChartCardItem[]>([]);

  useEffect(() => {
    const evaluatorScoreMap = getEvaluatorScoreMap(evaluatorAggregateResult);
    const newItems = evaluators.map(evaluator => {
      const versionId = evaluator?.current_version?.id ?? '';
      const scoreCountMap = evaluatorScoreMap[versionId] ?? {};
      const values = Object.entries(scoreCountMap).map(([score, item]) => ({
        name: `${I18n.t('score')}${score} - ${getScorePercentage(item?.percentage)}`,
        score: item?.count,
        item,
      }));
      const item: ChartCardItem = {
        id: versionId.toString(),
        title: (
          <EvaluatorPreview
            evaluator={evaluator}
            enableDescTooltip={false}
            tagProps={{ className: 'font-normal' }}
          />
        ),
        tooltip: evaluator?.description,
        content:
          ready && values.length === 0 ? (
            <div className="pt-10 pb-6">
              <EmptyState
                size="full_screen"
                icon={<IconCozIllusAdd />}
                title={I18n.t('no_data')}
                description={I18n.t('refresh_after_experiment')}
              />
            </div>
          ) : (
            <Chart
              className="h-[260px]"
              spec={spec}
              values={values}
              customTooltip={ComplexTooltipContent}
            />
          ),
      };
      return item;
    });
    setItems(newItems);
  }, [evaluators, evaluatorAggregateResult, spaceID]);
  return (
    <DraggableGrid<ChartCardItem>
      items={items}
      itemRender={ChartCardItemRender}
      onItemsChange={setItems}
    />
  );
}
