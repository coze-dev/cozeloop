// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { type BreadcrumbItemConfig } from '@cozeloop/stores';

import { type MenuConfig } from '../navbar/menu-config';

export type BreadcrumbConfig = Record<string, BreadcrumbItemConfig>;

export function getBreadcrumbMap(config: MenuConfig[]) {
  const breadcrumbMap: BreadcrumbConfig = {};

  const collect = (
    menus: MenuConfig[],
    cfg: BreadcrumbConfig = {},
    level = 0,
  ) => {
    menus.forEach(menu => {
      cfg[menu.itemKey] = {
        text: menu.text,
        path: menu.itemKey,
      };
      if (menu.items) {
        collect(menu.items, cfg, level + 1);
      }
    });
  };

  collect(config, breadcrumbMap);
  return breadcrumbMap;
}
