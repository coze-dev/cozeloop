// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { useEffect, useMemo, useRef } from 'react';

import {
  type Params,
  type PaginationResult,
} from 'ahooks/lib/usePagination/types';
import { useSize } from 'ahooks';
import { CozPagination, type TableProps } from '@coze-arch/coze-design';

import LoopTableSortIcon from './sort-icon';
import { LoopTable } from './index';

/** 获取本地存储的表格分页数量 */
export function getStoragePageSize(pageSizeStorageKey: string | undefined) {
  if (!pageSizeStorageKey) {
    return undefined;
  }
  const pageSize = localStorage.getItem(pageSizeStorageKey);
  if (pageSize && !isNaN(Number(pageSize))) {
    return Number(pageSize);
  }
  return undefined;
}

export const PAGE_SIZE_OPTIONS = [10, 20, 50];
export const DEFAULT_PAGE_SIZE = 20;

// eslint-disable-next-line complexity
export function TableWithPagination<RecordItem>(
  props: TableProps & {
    heightFull?: boolean;
    service: PaginationResult<{ total: number; list: RecordItem[] }, Params>;
    pageSizeOpts?: number[];
    header?: React.ReactNode;
    pageSizeStorageKey?: string;
    showSizeChanger?: boolean;
  },
) {
  const {
    pageSizeOpts,
    service,
    header,
    heightFull = false,
    pageSizeStorageKey,
    showSizeChanger = true,
  } = props;
  const { columns } = props.tableProps ?? {};
  const tableContainerRef = useRef<HTMLDivElement>(null);
  const size = useSize(tableContainerRef.current);
  const tableHeaderSize = useSize(
    tableContainerRef.current?.querySelector('.semi-table-header'),
  );

  const tablePagination = useMemo(
    () => ({
      currentPage: service.pagination.current,
      pageSize:
        getStoragePageSize(pageSizeStorageKey) || service.pagination.pageSize,
      total: Number(service.pagination.total),
      onChange: (page: number, pageSize: number) => {
        service.pagination.onChange(page, pageSize);
      },
      onPageSizeChange(newPageSize: number) {
        if (pageSizeStorageKey) {
          localStorage.setItem(pageSizeStorageKey, String(newPageSize));
        }
      },
      showSizeChanger,
      pageSizeOpts: pageSizeOpts ?? PAGE_SIZE_OPTIONS,
    }),
    [service.pagination, pageSizeOpts, pageSizeStorageKey, showSizeChanger],
  );

  useEffect(() => {
    if (service.pagination.current > 1 && service?.data?.list?.length === 0) {
      service.pagination.changeCurrent(1);
    }
  }, [service.pagination.current, service?.data?.list]);

  const tableHeaderHeight = tableHeaderSize?.height ?? 56;

  return (
    <div
      className={`${heightFull ? 'h-full flex overflow-hidden' : ''} flex flex-col gap-3`}
    >
      {header ? header : null}
      <div
        ref={tableContainerRef}
        className={heightFull ? 'flex-1 overflow-hidden' : ''}
      >
        <LoopTable
          {...props}
          tableProps={{
            empty: <></>,
            ...(props.tableProps ?? {}),
            scroll: {
              // 表格容器的高度减去表格头的高度
              y:
                size?.height === undefined || !heightFull
                  ? undefined
                  : size.height - tableHeaderHeight,
              ...(props.tableProps?.scroll ?? {}),
            },
            loading: service?.loading,
            columns: columns
              ?.filter(
                column => column.hidden !== true && column.checked !== false,
              )
              ?.map(column => ({
                ...column,
                ...(column.sorter && !column.sortIcon
                  ? { sortIcon: LoopTableSortIcon }
                  : {}),
              })),
            dataSource: service?.data?.list ?? [],
          }}
        />
      </div>
      {service.pagination.current > 1 ||
      (service?.data?.list?.length && service?.data?.list?.length > 0) ? (
        <div className="shrink-0 flex flex-row-reverse">
          <CozPagination showTotal {...tablePagination}></CozPagination>
        </div>
      ) : null}
    </div>
  );
}
