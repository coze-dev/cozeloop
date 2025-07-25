// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useState } from 'react';

import { sendEvent, EVENT_NAMES } from '@cozeloop/tea-adapter';
import { I18n } from '@cozeloop/i18n-adapter';
import { Guard, GuardPoint } from '@cozeloop/guard';
import { useSpace } from '@cozeloop/biz-hooks-adapter';
import {
  type EvaluationSet,
  type EvaluationSetItem,
} from '@cozeloop/api-schema/evaluation';
import { StoneEvaluationApi } from '@cozeloop/api-schema';
import {
  Button,
  Checkbox,
  Modal,
  Typography,
  type ColumnProps,
} from '@coze-arch/coze-design';

export const useBatchSelect = ({
  itemList,
  onDelete,
  datasetDetail,
}: {
  itemList?: EvaluationSetItem[];
  onDelete: () => void;
  datasetDetail?: EvaluationSet | undefined;
}) => {
  const { spaceID } = useSpace();
  const [batchSelectItems, setBatchSelectedItems] = useState<Set<string>>(
    new Set(),
  );
  const [batchSelectVisible, setBatchSelectVisible] = useState(false);

  const handleBatchSelect = e => {
    if (e.target.checked) {
      setBatchSelectedItems(
        new Set([
          ...(itemList?.map(item => item.item_id as string) || []),
          ...batchSelectItems,
        ]),
      );
    } else {
      const newSet = Array.from(batchSelectItems).filter(
        item => !itemList?.some(set => set.item_id === item),
      );
      setBatchSelectedItems(new Set(newSet));
    }
  };

  const handleSingleSelect = (e, itemId: string) => {
    if (e.target.checked) {
      setBatchSelectedItems(new Set([...batchSelectItems, itemId]));
    } else {
      setBatchSelectedItems(
        new Set(Array.from(batchSelectItems).filter(item => item !== itemId)),
      );
    }
  };

  const selectColumn: ColumnProps = {
    title: (
      <Checkbox
        checked={itemList?.every(item =>
          batchSelectItems.has(item.item_id as string),
        )}
        onChange={handleBatchSelect}
      />
    ),
    key: 'check',
    width: 50,
    fixed: 'left',
    render: (_, record) => (
      <div onClick={e => e.stopPropagation()}>
        <Checkbox
          checked={batchSelectItems.has(record.item_id as string)}
          onChange={e => {
            handleSingleSelect(e, record.item_id as string);
          }}
        />
      </div>
    ),
  };
  const EnterBatchSelectButton = (
    <Button
      color="primary"
      onClick={() => {
        setBatchSelectVisible(true);
        setBatchSelectedItems(new Set());
        sendEvent(EVENT_NAMES.cozeloop_dataset_batch_action);
      }}
    >
      {I18n.t('batch_select')}
    </Button>
  );

  const handleDelete = () => {
    Modal.confirm({
      title: I18n.t('delete_data_item'),
      content: I18n.t('confirm_to_delete_selected_data_item', {
        num: batchSelectItems.size,
      }),
      okText: I18n.t('delete'),
      cancelText: I18n.t('Cancel'),
      okButtonProps: {
        color: 'red',
      },
      autoLoading: true,
      onOk: async () => {
        await StoneEvaluationApi.BatchDeleteEvaluationSetItems({
          workspace_id: spaceID,
          evaluation_set_id: datasetDetail?.id as string,
          item_ids: Array.from(batchSelectItems),
        });
        setBatchSelectVisible(false);
        setBatchSelectedItems(new Set());
        onDelete();
      },
    });
  };
  const BatchSelectHeader = (
    <div className="flex items-center justify-end gap-2">
      <Typography.Text size="small">
        {I18n.t('x_data_item_selected', {
          num: (
            <Typography.Text size="small" className="mx-[2px]  font-medium">
              {batchSelectItems.size}
            </Typography.Text>
          ),
        })}
      </Typography.Text>
      <Typography.Text
        link
        onClick={() => {
          setBatchSelectVisible(false);
          setBatchSelectedItems(new Set());
        }}
      >
        {I18n.t('unselect')}
      </Typography.Text>
      <Guard point={GuardPoint['eval.dataset.batch_delete']}>
        <Button
          color="red"
          disabled={batchSelectItems.size === 0}
          onClick={handleDelete}
        >
          {I18n.t('delete')}
        </Button>
      </Guard>
    </div>
  );

  return {
    selectColumn,
    setBatchSelectedItems,
    EnterBatchSelectButton,
    BatchSelectHeader,
    batchSelectVisible,
    setBatchSelectVisible,
  };
};
