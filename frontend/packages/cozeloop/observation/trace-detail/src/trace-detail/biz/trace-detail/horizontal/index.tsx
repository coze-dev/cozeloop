// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { Panel, PanelGroup, PanelResizeHandle } from 'react-resizable-panels';
import { useState } from 'react';

import classNames from 'classnames';
import { Spin } from '@coze-arch/coze-design';

import { SpanDetail } from '@/trace-detail/components/span-detail';
import { HorizontalTraceHeader } from '@/trace-detail/components/header';
import { TraceGraphs } from '@/trace-detail/components/graphs';
import { NodeDetailEmpty } from '@/trace-detail/components/common/empty-status';

import { type TraceDetailLayoutProps } from '../interface';

export const HorizontalTraceDetail = ({
  loading,
  onCollapseChange,
  onSelect,
  rootNodes,
  advanceInfo,
  selectedSpan,
  selectedSpanId,
  moduleName,
  spans,
  headerConfig,
  spanDetailConfig,
  switchConfig,
  className,
  style,
  urlParams,
  hideTraceDetailHeader,
}: TraceDetailLayoutProps) => {
  const [dragging, setDragging] = useState(false);

  return (
    <div
      className={classNames('flex-1 flex flex-col overflow-hidden', className)}
      style={style}
    >
      {!hideTraceDetailHeader && (
        <HorizontalTraceHeader
          rootSpan={rootNodes?.[0]}
          advanceInfo={advanceInfo}
          disableEnvTag={headerConfig?.disableEnvTag}
          showClose={headerConfig?.showClose}
          onClose={headerConfig?.onClose}
          minColWidth={headerConfig?.minColWidth}
          urlParams={urlParams}
          switchConfig={switchConfig}
        />
      )}

      <div className="flex-1 flex flex-col overflow-hidden">
        <PanelGroup
          direction="horizontal"
          className="border-solid border-0 border-t border-[var(--coz-stroke-primary)]"
        >
          <Panel minSize={24} defaultSize={30} maxSize={35}>
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
            className="w-[2px] group hover:cursor-col-resize"
            onDragging={isDragging => {
              setDragging(isDragging);
            }}
          >
            <div
              className="w-[1px] h-full box-border transition group-hover:bg-[rgb(var(--coze-up-brand-7))]"
              style={{
                background: dragging
                  ? 'rgb(var(--coze-up-brand-7))'
                  : 'var(--coz-stroke-primary)',
                width: dragging ? '2px' : '1px',
              }}
            />
          </PanelResizeHandle>
          <Panel>
            <Spin
              spinning={loading}
              wrapperClassName="!w-full !h-full flex items-center justify-center max-h-full overflow-auto"
              childStyle={{ height: '100%', width: '100%' }}
            >
              {selectedSpan ? (
                <SpanDetail
                  showTags={spanDetailConfig?.showTags}
                  baseInfoPosition={spanDetailConfig?.baseInfoPosition}
                  maxColNum={spanDetailConfig?.maxColNum}
                  minColWidth={spanDetailConfig?.minColWidth}
                  span={selectedSpan}
                  moduleName={moduleName}
                  className="h-full overflow-auto max-h-full w-full"
                />
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
