// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type TableProps, Table } from '@coze-arch/coze-design';

import styles from './index.module.less';

export const LoopTable: React.FC<TableProps> = ({ className, ...props }) => (
  <Table {...props} id={styles['loop-table']} />
);
