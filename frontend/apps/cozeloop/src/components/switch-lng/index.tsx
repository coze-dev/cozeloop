import { useState } from 'react';

import { I18n } from '@cozeloop/i18n-adapter';
import { Switch } from '@coze-arch/coze-design';

import { ReactComponent as IconZhCN } from '@/assets/images/zh-cn.svg';
import { ReactComponent as IconEnUS } from '@/assets/images/en-us.svg';

// TODO
export function SwitchLang() {
  // 'zh-CN' | 'en-US'
  const [lang, setLang] = useState(() => I18n.lang);

  const toggleLang = async () => {
    const switchMap: Record<string, string> = {
      'zh-CN': 'en-US',
      'en-US': 'zh-CN',
    };
    const targetLang = switchMap[lang] ?? 'zh-CN';
    await I18n.setLang(switchMap[I18n.lang]);
    setLang(targetLang);
  };

  return (
    <div className="inline-block cursor-pointer" onClick={toggleLang}>
      <Switch
        checkedText={<IconZhCN className="h-5 w-5" />}
        uncheckedText={<IconEnUS className="h-5 w-5" />}
      />
      {lang === 'en-US' ? (
        <IconEnUS className="h-5 w-5" />
      ) : (
        <IconZhCN className="h-5 w-5" />
      )}
    </div>
  );
}
