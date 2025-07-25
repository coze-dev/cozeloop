// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type FC, useEffect, useMemo, useState } from 'react';

import classNames from 'classnames';

import Tree, { type TreeNode } from '../tree';
import { spanNode2TreeNode, dealTreeNodeHighlight } from './utils';
import { type TraceTreeProps } from './type';
import { defaultProps } from './config';

import styles from './index.module.less';

const TraceTree: FC<TraceTreeProps> = props => {
  const [treeData, setTreeData] = useState<TreeNode>();
  const {
    dataSource: spanNode,
    selectedSpanId,
    indentDisabled,
    lineStyle: _lineStyle,
    globalStyle: _globalStyle,
    className,
    onCollapseChange,
    ...restProps
  } = props;

  const lineStyle = useMemo(
    () => ({
      normal: Object.assign(
        {},
        defaultProps.lineStyle?.normal,
        _lineStyle?.normal,
      ),
      select: Object.assign(
        {},
        defaultProps.lineStyle?.select,
        _lineStyle?.select,
      ),
      hover: Object.assign(
        {},
        defaultProps.lineStyle?.hover,
        _lineStyle?.hover,
      ),
    }),
    [_lineStyle],
  );

  const globalStyle = useMemo(
    () => Object.assign({}, defaultProps.globalStyle, _globalStyle),
    [_globalStyle],
  );

  useEffect(() => {
    if (spanNode) {
      const treeNode = spanNode2TreeNode({
        spanNode,
        onCollapseChange,
      });
      const treeNodeWithHighlight = dealTreeNodeHighlight(
        treeNode,
        selectedSpanId,
      );
      setTreeData(treeNodeWithHighlight);
    }
  }, [spanNode, selectedSpanId]);

  return treeData ? (
    <Tree
      className={classNames(styles['trace-tree'], className)}
      treeData={treeData}
      selectedKey={selectedSpanId}
      indentDisabled={indentDisabled}
      lineStyle={lineStyle}
      globalStyle={globalStyle}
      {...restProps}
    />
  ) : null;
};

export { TraceTree };
