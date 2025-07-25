// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import {
  type Span,
  type TraceAdvanceInfo,
  type UrlParams,
} from '@/trace-detail/typings/params';

import { type SwitchConfig } from '../../biz/trace-detail/interface';
export interface TraceHeaderProps {
  moduleName?: string;
  showClose?: boolean;
  onClose?: () => void;
  rootSpan?: Span;
  advanceInfo?: TraceAdvanceInfo;
  disableEnvTag?: boolean;
  className?: string;
  minColWidth?: number;
  maxColNum?: number;
  urlParams: UrlParams;
  switchConfig?: SwitchConfig;
}
