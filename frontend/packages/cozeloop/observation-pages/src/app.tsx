// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { Routes, Route, Navigate } from 'react-router-dom';

import TraceList from '@cozeloop/trace-list';

const App = () => (
  <div className="text-sm h-full overflow-hidden">
    <Routes>
      <Route path="" element={<Navigate to="traces" replace />} />
      <Route path="traces" element={<TraceList />} />
    </Routes>
  </div>
);

export default App;
