// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useRequest } from 'ahooks';
import { I18n } from '@cozeloop/i18n-adapter';
import { VersionItem } from '@cozeloop/components';
import {
  type EvaluatorVersion,
  type Evaluator,
} from '@cozeloop/api-schema/evaluation';
import { StoneEvaluationApi } from '@cozeloop/api-schema';
import { IconCozCross } from '@coze-arch/coze-design/icons';
import { IconButton, Spin } from '@coze-arch/coze-design';

export function VersionListPane({
  evaluator,
  selectedVersion,
  onSelectVersion,
  onClose,
  refreshFlag,
}: {
  evaluator: Evaluator;
  selectedVersion: EvaluatorVersion | undefined;
  onSelectVersion: (version: EvaluatorVersion | undefined) => void;
  onClose: () => void;
  refreshFlag: never[];
}) {
  const service = useRequest(
    async () =>
      StoneEvaluationApi.ListEvaluatorVersions({
        workspace_id: evaluator.workspace_id || '',
        evaluator_id: evaluator.evaluator_id,
        page_size: 10000,
      }),
    {
      refreshDeps: [refreshFlag],
    },
  );

  return (
    <div className="flex-shrink-0 w-[340px] h-full overflow-hidden flex flex-col border-0 border-l border-solid coz-stroke-primary">
      <div className="flex-shrink-0 h-12 px-6 flex flex-row items-center justify-between coz-mg-secondary border-0 border-b border-solid coz-stroke-primary">
        <div className="text-sm font-medium coz-fg-plus">
          {I18n.t('version_record')}
        </div>
        <IconButton
          className="flex-shrink-0"
          color="secondary"
          size="small"
          icon={<IconCozCross className="w-4 h-4 coz-fg-primary" />}
          onClick={onClose}
        />
      </div>
      <div className="flex-1 overflow-y-auto p-6 gap-3 styled-scrollbar pr-[18px]">
        {service.loading ? (
          <Spin spinning={true} wrapperClassName="!w-full" />
        ) : (
          <>
            <VersionItem
              key={'isDraft'}
              className="pb-3"
              version={{
                id: 'isDraft',
                isDraft: true,
                submitTime: evaluator?.base_info?.updated_at,
              }}
              active={!selectedVersion}
              onClick={() => onSelectVersion(undefined)}
            />
            {service.data?.evaluator_versions?.map(version => (
              <VersionItem
                key={version.id}
                version={{
                  id: version.id || '',
                  version: version.version,
                  description: version.description,
                  submitTime: version.base_info?.created_at,
                  submitter: version.base_info?.created_by,
                  isDraft: false,
                }}
                active={selectedVersion && selectedVersion?.id === version.id}
                onClick={() => onSelectVersion(version)}
              />
            ))}
          </>
        )}
      </div>
    </div>
  );
}
