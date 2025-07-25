// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type EvalTarget } from '@cozeloop/api-schema/evaluation';

import { useEvalTargetDefinition } from '@/stores/eval-target-store';

/** 评测对象预览 */
export function EvalTargetPreview({
  evalTarget,
  spaceID,
  enableLinkJump,
  size,
  jumpBtnClassName,
  showIcon,
}: {
  evalTarget: EvalTarget | undefined;
  spaceID: Int64;
  enableLinkJump?: boolean;
  size?: 'small' | 'medium';
  jumpBtnClassName?: string;
  showIcon?: boolean;
}) {
  const { getEvalTargetDefinition } = useEvalTargetDefinition();
  const { eval_target_type } = evalTarget ?? {};
  const target = getEvalTargetDefinition(eval_target_type ?? '');
  const Preview = target?.preview;
  if (evalTarget && Preview) {
    return (
      <Preview
        evalTarget={evalTarget}
        spaceID={spaceID}
        enableLinkJump={enableLinkJump}
        size={size}
        jumpBtnClassName={jumpBtnClassName}
        showIcon={showIcon}
      />
    );
  }
  return <>-</>;
}
