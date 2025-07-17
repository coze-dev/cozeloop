// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { ErrorBoundary } from 'react-error-boundary';
import { useState } from 'react';

import classNames from 'classnames';
import { Spin, TabPane, Tabs } from '@coze-arch/coze-design';
import { I18n } from '@cozeloop/i18n-adapter';

import type { Span } from '@/trace-detail/typings/params';

import { RunTreeEmpty } from '../common/empty-status';
import { GraphTabEnum } from '../../consts/tab';
import { type SpanNode } from './trace-tree/type';
import { TraceTree } from './trace-tree';

import styles from './index.module.less';

interface TraceGraphsProps {
  rootNodes?: SpanNode[];
  spans: Span[];
  selectedSpanId: string;
  onSelect: (id: string) => void;
  onCollapseChange: (id: string) => void;
  loading?: boolean;
  className?: string;
}

export const TraceGraphs = ({
  rootNodes,
  loading = false,
  selectedSpanId,
  onSelect,
  onCollapseChange,
  className,
}: TraceGraphsProps) => {
  const [activeTab, setActiveTab] = useState<string>(GraphTabEnum.RunTree);

  return (
    <div className={classNames(className, styles['trace-graph'])}>
      <Tabs
        activeKey={activeTab}
        renderTabBar={({ activeKey, list }) => (
          <div className={styles['tabs-bar']}>
            {list?.map(({ tab, itemKey }) => (
              <div
                className={classNames(styles['tab-bar'], {
                  [styles.active]: activeKey === itemKey,
                })}
                key={itemKey}
                onClick={() => {
                  setActiveTab(itemKey);
                }}
              >
                {tab}
              </div>
            ))}
          </div>
        )}
        className={styles['trace-tabs']}
      >
        <TabPane
          tab={I18n.t('observation_tab_run_tree')}
          itemKey={GraphTabEnum.RunTree}
        >
          <ErrorBoundary fallback={<RunTreeEmpty />}>
            <Spin
              spinning={loading}
              wrapperClassName="!h-full"
              childStyle={{ height: '100%' }}
            >
              <div className={classNames(styles['run-tree-area'])}>
                {rootNodes
                  ? rootNodes.map(root => (
                      <TraceTree
                        key={root.span_id}
                        dataSource={root}
                        className={styles['run-tree']}
                        selectedSpanId={selectedSpanId}
                        onCollapseChange={onCollapseChange}
                        onSelect={({ node }) => onSelect(node.key)}
                      />
                    ))
                  : !loading && <RunTreeEmpty />}
              </div>
            </Spin>
          </ErrorBoundary>
        </TabPane>
      </Tabs>
    </div>
  );
};
