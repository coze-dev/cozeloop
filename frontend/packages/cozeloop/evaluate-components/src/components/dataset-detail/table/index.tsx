// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @coze-arch/max-line-per-function */

import { useEffect } from 'react';

import cs from 'classnames';
import { I18n } from '@cozeloop/i18n-adapter';
import { formatTimestampToString } from '@cozeloop/toolkit';
import { GuardPoint, useGuard } from '@cozeloop/guard';
import { type Version } from '@cozeloop/components';
import { TableColActions, TableWithPagination } from '@cozeloop/components';
import { type EvaluationSet } from '@cozeloop/api-schema/evaluation';
import { StoneEvaluationApi } from '@cozeloop/api-schema';
import { IconCozIllusAdd } from '@coze-arch/coze-design/illustrations';
import {
  type ColumnProps,
  EmptyState,
  Modal,
  Typography,
} from '@coze-arch/coze-design';

import { useVersionManage } from '../version-manage/use-version-manage';
import { TextEllipsis } from '../../text-ellipsis';
import { DatasetItemPanel } from '../../dataset-item-panel';
import {
  type EvaluationSetItemTableData,
  useDatasetItemList,
} from './use-dataset-item-list';
import { useBatchSelect } from './use-batch-select';
import { TableHeader } from './header';

import styles from './index.module.less';
interface DatasetItemListProps {
  datasetDetail?: EvaluationSet;
  spaceID: string;
  refreshDatasetDetail: () => void;
  setCurrentVersion: (version: Version) => void;
}
export const DatasetItemList: React.FC<DatasetItemListProps> = ({
  datasetDetail,
  spaceID,
  refreshDatasetDetail,
  setCurrentVersion: setCurrentVersionProps,
}) => {
  const {
    service,
    columns,
    setColumns,
    defaultColumnsItems,
    fieldSchemas,
    selectedItem,
    setSelectedItem,
    currentVersion,
    setCurrentVersion,
    ExpandNode,
    switchConfig,
    setOrderBy,
  } = useDatasetItemList({
    datasetDetail,
    spaceID,
    refreshDatasetDetail,
  });
  const {
    EnterBatchSelectButton,
    selectColumn,
    BatchSelectHeader,
    batchSelectVisible,
  } = useBatchSelect({
    itemList: service?.data?.list,
    onDelete: refreshDatasetDetail,
    datasetDetail,
  });

  useEffect(() => {
    setCurrentVersionProps(currentVersion);
  }, [currentVersion]);
  const guard = useGuard({
    point: GuardPoint['eval.dataset.delete'],
  });

  const { VersionPanel, VersionChangeButton } = useVersionManage({
    datasetDetail,
    currentVersion,
    setCurrentVersion,
  });

  const isDraftVersion = currentVersion?.id === 'draft';

  const handleDeleteItem = (item: EvaluationSetItemTableData) => {
    Modal.error({
      width: 420,
      title: I18n.t('delete_data_item'),
      type: 'dialog',
      content: (
        <Typography.Text className="break-all">
          {I18n.t('confirm_to_delete_data_item', {
            name: (
              <Typography.Text className="!font-medium">
                #{(item.item_id as string)?.slice(-5)}
              </Typography.Text>
            ),
          })}
        </Typography.Text>
      ),
      autoLoading: true,
      onOk: async () => {
        await StoneEvaluationApi.BatchDeleteEvaluationSetItems({
          workspace_id: spaceID,
          evaluation_set_id: datasetDetail?.id as string,
          item_ids: [item.item_id as string],
        });
        refreshDatasetDetail();
      },
      showCancelButton: true,
      cancelText: I18n.t('Cancel'),
      okText: I18n.t('delete'),
    });
  };
  const columnsItems: ColumnProps[] = [
    ...(batchSelectVisible ? [selectColumn] : []),
    ...(columns?.filter(column => !!column.checked) || []),
    {
      title: I18n.t('update_time'),
      key: 'updated_at',
      displayName: I18n.t('update_time'),
      sorter: true,
      width: 180,
      dataIndex: 'base_info.updated_at',
      render: (record: string) =>
        record ? (
          <TextEllipsis>
            {formatTimestampToString(record, 'YYYY-MM-DD HH:mm:ss')}
          </TextEllipsis>
        ) : (
          '-'
        ),
    },
    {
      title: I18n.t('create_time'),
      key: 'created_at',
      displayName: I18n.t('create_time'),
      width: 180,
      dataIndex: 'base_info.created_at',
      sorter: true,
      render: (record: string) =>
        record ? (
          /** 查看版本时，创建时间作为最后一项会被默认右对齐，这里通过flex修改为左对齐 */
          <div className="flex">
            <TextEllipsis>
              {formatTimestampToString(record, 'YYYY-MM-DD HH:mm:ss')}
            </TextEllipsis>
          </div>
        ) : (
          '-'
        ),
    },
    ...(isDraftVersion
      ? ([
          {
            title: I18n.t('operation'),
            key: 'action',
            width: 140,
            fixed: 'right',
            disabled: true,
            render: (row: EvaluationSetItemTableData, _, index: number) => (
              <TableColActions
                actions={[
                  {
                    label: I18n.t('edit'),
                    onClick: () => {
                      setSelectedItem({
                        item: row,
                        isEdit: true,
                        index,
                      });
                    },
                  },

                  {
                    label: I18n.t('view'),
                    onClick: () => {
                      setSelectedItem({
                        item: row,
                        isEdit: false,
                        index,
                      });
                    },
                  },
                  {
                    label: I18n.t('delete'),
                    type: 'danger',
                    disabled: guard.data.readonly,
                    onClick: () => {
                      handleDeleteItem(row);
                    },
                  },
                ]}
                maxCount={2}
              />
            ),
          },
        ] as ColumnProps[])
      : []),
  ];

  return (
    <div className="h-full w-full flex overflow-hidden">
      <div
        className={cs(
          styles.table,
          'flex-1 h-full px-6 py-4 gap-4 w-full overflow-hidden',
        )}
      >
        <TableWithPagination
          service={service}
          heightFull={true}
          showTableWhenEmpty
          tableProps={{
            rowKey: 'id',
            columns: columnsItems as ColumnProps[],
            sticky: { top: 0 },
            onRow: (record: EvaluationSetItemTableData, index) => ({
              onClick: () => {
                setSelectedItem({
                  item: record,
                  isEdit: false,
                  index: index || 0,
                });
              },
            }),
            onChange: data => {
              if (data.extra?.changeType === 'sorter') {
                setOrderBy(
                  data.sorter?.sortOrder === false
                    ? undefined
                    : {
                        field: data.sorter?.key,
                        is_asc: data.sorter?.sortOrder === 'ascend',
                      },
                );
              }
            },
          }}
          empty={
            <EmptyState
              size="full_screen"
              icon={<IconCozIllusAdd />}
              title={I18n.t('no_data')}
              description={I18n.t('click_to_add_data')}
            />
          }
          header={
            batchSelectVisible ? (
              BatchSelectHeader
            ) : (
              <TableHeader
                isDraftVersion={isDraftVersion}
                currentVersion={currentVersion}
                defaultColumnsItems={defaultColumnsItems}
                datasetDetail={datasetDetail}
                columns={columns}
                refreshDatasetDetail={refreshDatasetDetail}
                batchSelectNode={EnterBatchSelectButton}
                versionChangeNode={VersionChangeButton}
                datasetItemExpandNode={ExpandNode}
                setColumns={setColumns}
                totalItemCount={service?.data?.total}
              />
            )
          }
        />
        {selectedItem.item ? (
          <DatasetItemPanel
            datasetItem={selectedItem.item}
            fieldSchemas={fieldSchemas}
            isEdit={selectedItem.isEdit}
            onCancel={() => {
              setSelectedItem({
                item: undefined,
                isEdit: false,
                index: 0,
              });
            }}
            onSave={() => {
              setSelectedItem({
                item: undefined,
                isEdit: false,
                index: 0,
              });
              refreshDatasetDetail();
            }}
            switchConfig={switchConfig}
          />
        ) : null}
      </div>
      {VersionPanel}
    </div>
  );
};
