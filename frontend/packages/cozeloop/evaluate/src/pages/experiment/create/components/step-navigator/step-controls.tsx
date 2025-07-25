// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import React from 'react';

import { I18n } from '@cozeloop/i18n-adapter';
import { type GuardPoint, Guard } from '@cozeloop/guard';
import { Button } from '@coze-arch/coze-design';

import { type StepConfig } from '../../constants/steps';

interface StepControlsProps {
  currentStep: number;
  steps: StepConfig[];
  onNext: () => void;
  onPrevious: () => void;
  isNextLoading?: boolean;
}

export const StepControls: React.FC<StepControlsProps> = ({
  currentStep,
  steps,
  onNext,
  onPrevious,
  isNextLoading = false,
}) => {
  const currentStepConfig = steps[currentStep];

  return (
    <div className="flex-shrink-0 p-6">
      <div className="w-[800px] mx-auto flex flex-row items-center justify-end gap-2">
        {currentStep > 0 && (
          <Button color="primary" onClick={onPrevious}>
            {I18n.t('prev_step')}
          </Button>
        )}

        <Guard
          point={currentStepConfig.guardPoint as GuardPoint}
          ignore={!currentStepConfig.isLast}
        >
          <Button onClick={onNext} loading={isNextLoading}>
            {currentStepConfig.nextStepText}
          </Button>
        </Guard>
      </div>
    </div>
  );
};
