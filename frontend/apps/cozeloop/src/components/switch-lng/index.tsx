// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useI18nStore } from '@/stores';
import { ReactComponent as IconZhCN } from '@/assets/images/zh-cn.svg';
import { ReactComponent as IconEnUS } from '@/assets/images/en-us.svg';

export function SwitchLang() {
  const toggleLang = useI18nStore(s => s.toggleLng);
  const lang = useI18nStore(s => s.lng);

  return (
    <div className="inline-block cursor-pointer" onClick={toggleLang}>
      {lang === 'en-US' ? (
        <IconEnUS className="h-5 w-5" />
      ) : (
        <IconZhCN className="h-5 w-5" />
      )}
    </div>
  );
}
