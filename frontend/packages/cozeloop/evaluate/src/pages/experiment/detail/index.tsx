// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useParams } from 'react-router-dom';
import { useCallback, useState } from 'react';

import { useRequest } from 'ahooks';
import { I18n } from '@cozeloop/i18n-adapter';
import { LoopTabs } from '@cozeloop/components';
import { useSpace } from '@cozeloop/biz-hooks-adapter';
import { useBreadcrumb } from '@cozeloop/base-hooks';
import { Layout, Spin } from '@coze-arch/coze-design';

import { batchGetExperiment } from '@/request/experiment';
import { ExperimentContextProvider } from '@/hooks/use-experiment';

import ExperimentHeader from './components/experiment-header';
import ExperimentTable from './components/experiment-detail-table';
import ExperimentDescription from './components/experiment-description';
import ExperimentChart from './components/experiment-chart';

export default function () {
  const { experimentID = '' } = useParams<{ experimentID: string }>();
  const { spaceID = '' } = useSpace();
  const [activeKey, setActiveKey] = useState('detail');
  const [refreshKey, setRefreshKey] = useState('');

  const {
    data: experiment,
    loading,
    refresh,
  } = useRequest(
    async () => {
      if (!experimentID) {
        return;
      }
      const res = await batchGetExperiment({
        workspace_id: spaceID,
        expt_ids: [experimentID],
      });
      return res.experiments?.[0];
    },
    {
      refreshDeps: [experimentID, refreshKey],
    },
  );

  useBreadcrumb({
    text: experiment?.name || '',
  });

  const onRefresh = useCallback(() => {
    setRefreshKey(Date.now().toString());
  }, [setRefreshKey]);

  return (
    <Layout className="h-full overflow-hidden flex flex-col">
      <ExperimentContextProvider experiment={experiment}>
        <ExperimentHeader
          experiment={experiment}
          spaceID={spaceID}
          onRefreshExperiment={refresh}
          onRefresh={onRefresh}
        />
        <Spin spinning={loading}>
          <div className="px-6 pt-3 pb-6 flex items-center text-sm">
            <ExperimentDescription experiment={experiment} spaceID={spaceID} />
          </div>
        </Spin>
        <LoopTabs
          type="card"
          activeKey={activeKey}
          onChange={setActiveKey}
          tabPaneMotion={false}
          keepDOM={false}
          tabList={[
            { tab: I18n.t('data_detail'), itemKey: 'detail' },
            { tab: I18n.t('measure_stat'), itemKey: 'chart' },
          ]}
        />
        <div className="grow overflow-hidden">
          {activeKey === 'detail' && (
            <div className="h-full overflow-hidden px-6 pt-4 pb-4">
              <ExperimentTable
                spaceID={spaceID}
                experimentID={experimentID}
                refreshKey={refreshKey}
                experiment={experiment}
                onRefreshPage={onRefresh}
              />
            </div>
          )}
          {activeKey === 'chart' && (
            <div className="h-full overflow-auto styled-scrollbar pl-6 pr-[18px] py-4">
              <ExperimentChart
                spaceID={spaceID}
                experiment={experiment}
                experimentID={experimentID}
                loading={loading}
              />
            </div>
          )}
        </div>
      </ExperimentContextProvider>
    </Layout>
  );
}
