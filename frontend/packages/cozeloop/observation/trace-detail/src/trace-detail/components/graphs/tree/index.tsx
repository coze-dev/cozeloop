// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type FC, useState } from 'react';

import { isFunction } from 'lodash-es';
import cs from 'classnames';
import { SpanStatus } from '@cozeloop/api-schema/observation';

import { type SpanNode } from '../trace-tree/type';
import { flattenTreeData, checkIsNodeOrChildSelected } from './util';
import type {
  TreeProps,
  TreeNode,
  TreeNodeExtra,
  MouseEventParams,
  LineStyle,
  GlobalStyle,
} from './typing';
import { PathEnum } from './typing';

import styles from './index.module.less';

export type {
  TreeProps,
  TreeNode,
  TreeNodeExtra,
  MouseEventParams,
  LineStyle,
  GlobalStyle,
};
const Tree: FC<TreeProps> = ({
  treeData,
  selectedKey,
  disableDefaultHover,
  hoverKey: customHoverKey,
  indentDisabled = false,
  className,
  onMouseMove,
  onMouseEnter,
  onMouseLeave,
  onClick,
  onSelect,
  lineStyle,
}) => {
  const [hoverKey, setHoverKey] = useState<string>('');

  const controlledHoverKey = disableDefaultHover ? customHoverKey : hoverKey;

  const { nodes } = flattenTreeData(treeData, {
    indentDisabled,
  });
  const normalLineColor = lineStyle?.normal?.stroke;
  const selectLineColor = lineStyle?.select?.stroke;

  return (
    <div className={`${styles.tree} ${className ?? ''}`}>
      <div className={styles['tree-container']}>
        <div className={styles['tree-node-list']}>
          {nodes.map(node => {
            const {
              key,
              title,
              selectEnabled = true,
              linePath,
              isLastChild,
            } = node;
            const selected = selectedKey === key;
            const isCurrentNodeOrChildSelected = checkIsNodeOrChildSelected(
              node,
              selectedKey,
            );
            const nodeExtra: TreeNodeExtra = {
              ...node,
              selected,
              lineStyle,
              isCurrentNodeOrChildSelected,
              hover: controlledHoverKey === key,
            };
            const spanNode = (node?.extra as { spanNode: SpanNode })?.spanNode;
            const isError = spanNode?.status !== SpanStatus.Success;
            return (
              <div
                className={cs(styles['tree-node'])}
                style={
                  selected
                    ? {
                        backgroundColor: isError ? '#FFF8F7' : '#EFF1FF',
                      }
                    : {}
                }
                key={node.key}
                onClick={event => {
                  if (selectEnabled) {
                    onSelect?.({ node: nodeExtra });
                  }
                  onClick?.({ event, node: nodeExtra });
                }}
              >
                {selected ? (
                  <div
                    className="absolute top-0 bottom-0 left-0 w-[2px] bg-[rgb(87 105 227)]"
                    style={
                      selected
                        ? {
                            backgroundColor: isError ? '#D0292F' : '#5A4DED',
                          }
                        : {}
                    }
                  ></div>
                ) : null}
                {linePath?.map((line, index) => {
                  const isLast = index === linePath.length - 1;
                  const isActive = line === PathEnum.Active;
                  return (
                    <div className="w-[24px] relative " key={index}>
                      {isLast ? (
                        <div
                          className="absolute left-[12px] top-0 -ml-[0.5px] w-[14px] rounded-bl-[4px]  border-b  border-solid border-l border-t-[0px] border-r-[0px] border-current z-[1] coz-fg-dim"
                          style={{
                            height: 16,
                            borderColor: isCurrentNodeOrChildSelected
                              ? selectLineColor
                              : normalLineColor,
                            zIndex: isCurrentNodeOrChildSelected ? 3 : 1,
                          }}
                        ></div>
                      ) : null}
                      {!(isLastChild && isLast) && line !== PathEnum.Hidden ? (
                        <div
                          className="absolute inset-y-0 left-1/2 -ml-[0.5px] w-[1px] coz-fg-dim"
                          style={{
                            zIndex: 2,
                            backgroundColor:
                              isActive && !isCurrentNodeOrChildSelected
                                ? selectLineColor
                                : normalLineColor,
                          }}
                        ></div>
                      ) : null}
                    </div>
                  );
                })}
                <div
                  className={styles['tree-node-box']}
                  onMouseMove={event => {
                    onMouseMove?.({ event, node: nodeExtra });
                  }}
                  onMouseEnter={event => {
                    if (selectEnabled) {
                      setHoverKey(key);
                    }
                    onMouseEnter?.({
                      event,
                      node: { ...nodeExtra, hover: true },
                    });
                  }}
                  onMouseLeave={event => {
                    if (selectEnabled) {
                      setHoverKey('');
                    }
                    onMouseLeave?.({
                      event,
                      node: { ...nodeExtra, hover: false },
                    });
                  }}
                >
                  {isFunction(title) ? title(nodeExtra) : title}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};

export default Tree;
