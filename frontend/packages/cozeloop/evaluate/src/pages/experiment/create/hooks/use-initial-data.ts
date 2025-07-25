// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useState } from 'react';

import { useRequest } from 'ahooks';

import type { CreateExperimentValues } from '@/types/experiment/experiment-create';
import { batchGetExperiment } from '@/request/experiment';
import { getEvaluationSetVersion } from '@/request/evaluation-set';

import {
  experimentToCreateExperimentValues,
  evaluationSetToCreateExperimentValues,
} from '../tools';

export interface UseInitialDataOptions {
  spaceID: string;
  copyExperimentID?: string;
  evaluationSetID?: string;
  evaluationSetVersionID?: string;
  setValue: (value: CreateExperimentValues) => void;
}

export const useInitialData = ({
  spaceID,
  copyExperimentID,
  evaluationSetID,
  evaluationSetVersionID,
  setValue,
}: UseInitialDataOptions) => {
  const [initValue, setInitValue] = useState<CreateExperimentValues>({
    workspace_id: spaceID,
  } satisfies CreateExperimentValues);

  // 加载数据
  const { loading } = useRequest(
    async () => {
      if (copyExperimentID) {
        // 复制实验
        const res = await batchGetExperiment({
          workspace_id: spaceID,
          expt_ids: [copyExperimentID],
        });
        const experiment = res.experiments?.[0];

        if (experiment) {
          const data = experimentToCreateExperimentValues({
            experiment,
            spaceID,
          });
          const payload = {
            ...data,
            name: `${experiment.name}_copy`,
          };
          setValue(payload);
          setInitValue(payload);
        }
      } else if (evaluationSetID && evaluationSetVersionID) {
        // 从评测集创建
        const { evaluation_set, version } = await getEvaluationSetVersion({
          workspace_id: spaceID,
          evaluation_set_id: evaluationSetID,
          version_id: evaluationSetVersionID,
        });

        if (evaluation_set && version) {
          const data = evaluationSetToCreateExperimentValues(
            evaluation_set,
            version,
            spaceID,
          );
          setValue(data);
          setInitValue(data);
        }
      }
    },
    {
      refreshDeps: [copyExperimentID, evaluationSetID, evaluationSetVersionID],
    },
  );

  return {
    loading,
    initValue,
  };
};
