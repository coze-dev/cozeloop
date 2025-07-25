// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { Panel, PanelGroup, PanelResizeHandle } from 'react-resizable-panels';
import { useState } from 'react';

import classNames from 'classnames';
import { Spin } from '@coze-arch/coze-design';

import { SpanDetail } from '@/trace-detail/components/span-detail';
import { VerticalTraceHeader } from '@/trace-detail/components/header';
import { TraceGraphs } from '@/trace-detail/components/graphs';
import { NodeDetailEmpty } from '@/trace-detail/components/common/empty-status';

import { type TraceDetailLayoutProps } from '../interface';

export const VerticalTraceDetail = ({
  loading,
  onCollapseChange,
  onSelect,
  rootNodes,
  advanceInfo,
  selectedSpan,
  selectedSpanId,
  spans,
  moduleName,
  headerConfig,
  spanDetailConfig,
  className,
  style,
  urlParams,
  hideTraceDetailHeader,
}: TraceDetailLayoutProps) => {
  const [dragging, setDragging] = useState(false);
  return (
    <div
      className={classNames('flex-1 flex flex-col  overflow-hidden', className)}
      style={style}
    >
      {!hideTraceDetailHeader && (
        <VerticalTraceHeader
          rootSpan={rootNodes?.[0]}
          advanceInfo={advanceInfo}
          showClose={headerConfig?.showClose}
          onClose={headerConfig?.onClose}
          minColWidth={headerConfig?.minColWidth}
          maxColNum={2}
          urlParams={urlParams}
        />
      )}
      <div className="flex-1 flex flex-col overflow-hidden">
        <PanelGroup direction="vertical">
          <Panel
            className="border-solid border border-[var(--coz-stroke-primary)] rounded"
            minSize={20}
            defaultSize={40}
            maxSize={60}
          >
            <TraceGraphs
              rootNodes={rootNodes}
              loading={loading}
              spans={spans}
              selectedSpanId={selectedSpanId}
              onSelect={onSelect}
              onCollapseChange={onCollapseChange}
            />
          </Panel>
          <PanelResizeHandle
            className="h-2 group hover:cursor-row-resize"
            onDragging={isDragging => {
              setDragging(isDragging);
            }}
          >
            <div
              className="h-[2px] box-border my-[3px] mx-2.5 transition group-hover:bg-[#336df4]"
              style={{ background: dragging ? '#336df4' : undefined }}
            />
          </PanelResizeHandle>
          <Panel className="border-solid border border-[var(--coz-stroke-primary)] rounded">
            <Spin
              spinning={loading}
              wrapperClassName="!w-full !h-full flex items-center justify-center max-h-full overflow-auto"
              childStyle={{ height: '100%', width: '100%' }}
            >
              {selectedSpan ? (
                <div className="flex">
                  <SpanDetail
                    showTags={spanDetailConfig?.showTags}
                    baseInfoPosition={spanDetailConfig?.baseInfoPosition}
                    maxColNum={spanDetailConfig?.maxColNum}
                    minColWidth={spanDetailConfig?.minColWidth}
                    span={selectedSpan}
                    moduleName={moduleName}
                    className="h-full overflow-auto max-h-full w-full"
                  />
                </div>
              ) : (
                <NodeDetailEmpty />
              )}
            </Spin>
          </Panel>
        </PanelGroup>
      </div>
    </div>
  );
};
