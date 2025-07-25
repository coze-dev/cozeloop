// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import {
  type ForwardedRef,
  forwardRef,
  useRef,
  useImperativeHandle,
} from 'react';

import {
  type Data,
  type InfiniteScrollOptions,
  type Service,
} from 'ahooks/lib/useInfiniteScroll/types';
import { useSize } from 'ahooks';
import { type TableProps } from '@coze-arch/coze-design';

import { LoopTable } from '@/table';
import { useInfiniteScroll } from '@/hooks/use-infinite-scroll';

interface ExpandData extends Data {
  hasMore?: boolean;
}
interface InfiniteScrollTableProps<TData extends ExpandData> {
  service: Service<TData>;
  options?: InfiniteScrollOptions<TData>;
}

export interface InfiniteScrollTableRef {
  hookRes: ReturnType<typeof useInfiniteScroll>;
}

// 定义组件的 Props 类型
type Props<TData extends ExpandData> = TableProps &
  InfiniteScrollTableProps<TData>;

// 显式指定 forwardRef 的类型
export const InfiniteScrollTable: <TData extends ExpandData>(
  props: Props<TData> & { ref?: ForwardedRef<InfiniteScrollTableRef> },
) => React.ReactElement | null = forwardRef(
  <TData extends ExpandData>(
    { service, options, ...restTableProps }: Props<TData>,
    ref: ForwardedRef<InfiniteScrollTableRef>,
  ): JSX.Element => {
    const containerRef = useRef<HTMLDivElement>(null);

    const hookRes = useInfiniteScroll(service, {
      target: containerRef.current,
      isNoMore: d => !d?.hasMore,
      ...options,
    });
    const { data, loading, loadingMore } = hookRes;

    const scrollSize = useSize(containerRef);
    const height = scrollSize?.height || 0;

    useImperativeHandle(ref, () => ({
      // @ts-expect-error type
      hookRes,
    }));

    return (
      <div className="w-full h-full overflow-hidden" ref={containerRef}>
        <LoopTable
          tableProps={{
            ...restTableProps.tableProps,
            loading: loading || loadingMore,
            pagination: false,
            dataSource: data?.list || [],
            scroll: {
              y: height - 48,
              ...restTableProps.tableProps?.scroll,
            },
          }}
          empty={restTableProps.empty}
        />
      </div>
    );
  },
) as <TData extends ExpandData>(
  props: Props<TData> & { ref?: ForwardedRef<InfiniteScrollTableRef> },
) => React.ReactElement | null;
