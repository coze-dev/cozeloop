// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { I18n } from '@cozeloop/i18n-adapter';
import { ResizeSidesheet, IDRender } from '@cozeloop/components';
import {
  type Experiment,
  type FieldSchema,
} from '@cozeloop/api-schema/evaluation';
import { Spin } from '@coze-arch/coze-design';

import { DetailItemStepSwitch } from '@/types';
import { type ExperimentDetailActiveItemStore } from '@/hooks/use-experiment-detail-active-item';

import { type ExperimentContrastItem } from '../../utils/tools';
import ContrastItemDetailTable from './contrast-item-detail-table';

export default function ExperimentContrastItemDetail({
  experiments = [],
  datasetFieldSchemas = [],
  spaceID,
  activeItemStore,
  onClose,
  onStepChange,
}: {
  experiments: Experiment[];
  datasetFieldSchemas: FieldSchema[];
  spaceID: Int64;
  activeItemStore: ExperimentDetailActiveItemStore<ExperimentContrastItem>;
  onClose?: () => void;
  onStepChange?: (step: DetailItemStepSwitch) => void;
}) {
  const experimentContrastItem = activeItemStore.activeItem;
  if (!experimentContrastItem) {
    return null;
  }
  return (
    <ResizeSidesheet
      onCancel={onClose}
      closable={false}
      title={
        <div className="flex items-center gap-2">
          {I18n.t('compare_experiment_detail')}{' '}
          <IDRender id={experimentContrastItem?.groupID ?? ''} />
        </div>
      }
      dragOptions={{
        defaultWidth: 880,
        minWidth: 448,
        maxWidth: 1382,
      }}
      visible={true}
      bodyStyle={{ padding: 0 }}
    >
      <div className="h-full overflow-hidden">
        <Spin
          spinning={activeItemStore.loading}
          wrapperClassName="!h-full overflow-hidden"
          childStyle={{ height: '100%' }}
        >
          <ContrastItemDetailTable
            experiments={experiments}
            datasetFieldSchemas={datasetFieldSchemas}
            datasetRow={experimentContrastItem?.datasetRow}
            experimentsDatasetRow={
              experimentContrastItem?.experimentsDatasetRow
            }
            experimentContrastItem={experimentContrastItem}
            spaceID={spaceID}
            onRefresh={() => onStepChange?.(DetailItemStepSwitch.Current)}
          />
        </Spin>
      </div>
    </ResizeSidesheet>
  );
}
