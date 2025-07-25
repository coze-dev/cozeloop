// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { Navigate } from 'react-router-dom';

import { PERSONAL_ENTERPRISE_ID } from '@cozeloop/account';

interface Props {
  index?: boolean;
}

export function EnterpriseRoute({ index }: Props) {
  const enterpriseID = PERSONAL_ENTERPRISE_ID;
  const path = index ? `enterprise/${enterpriseID}` : enterpriseID;

  return <Navigate to={path} replace />;
}
