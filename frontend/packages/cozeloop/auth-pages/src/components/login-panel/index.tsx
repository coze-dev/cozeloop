// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0
import { useState } from 'react';

import { Input, Button, Typography } from '@coze-arch/coze-design';
import { ReactComponent as IconGithub } from '@/assets/github.svg';

import loopBanner from '@/assets/loop-banner.png';

import s from './index.module.less';

interface Props {
  loading?: boolean;
  onLogin?: (email: string, password: string) => void;
  onRegister?: (email: string, password: string) => void;
}

const { Text } = Typography;

export function LoginPanel({ loading, onLogin, onRegister }: Props) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  // const [consent, setConsent] = useState(false);
  const canSubmit = Boolean(email && password);

  const onClickRegister = () => {
    onRegister?.(email, password);
  };

  const onClickLogin = () => {
    onLogin?.(email, password);
  };

  return (
    <div className={s.container}>
      <div className="flex flex-col items-center">
        <img src={loopBanner} className={s.banner} />
        <div className="text-[18px] font-medium leading-[36px] my-[20px]">
          {'欢迎使用扣子罗盘-开源版'}
        </div>
      </div>
      <div className="w-[320px] flex flex-col items-stretch">
        <Input
          type="email"
          value={email}
          onChange={setEmail}
          placeholder={'请输入邮箱'}
        />
        <Input
          className="mt-[20px]"
          type="password"
          value={password}
          onChange={setPassword}
          placeholder={'请输入密码'}
        />
        <div className="mt-[20px] flex justify-between items-center">
          <Button
            className="w-[49%]"
            disabled={!canSubmit}
            onClick={onClickRegister}
            loading={loading}
            color="primary"
          >
            {'注册'}
          </Button>
          <Button
            className="w-[49%]"
            disabled={!canSubmit}
            onClick={onClickLogin}
            loading={loading}
          >
            {'登录'}
          </Button>
        </div>
        {/* <div className="mt-[20px] flex">
          <Checkbox
            checked={consent}
            onChange={e => setConsent(Boolean(e.target.checked))}
            disabled={loading}
          >
            {'请先同意'}
            <a
              href="" // 协议链接
              target="_blank"
              className="no-underline coz-fg-hglt"
              onClick={e => {
                e.stopPropagation();
              }}
            >
              用户协议
            </a>
          </Checkbox>
        </div> */}
      </div>
      <div className={s.copyright}>
        <Text component="div" type="secondary">
          ©2025 Coze Loop
        </Text>
        <Text type="secondary">
          基于开源代码部署
          <span> · </span>
          <Text
            link={{
              href: 'https://github.com/coze-dev/cozeloop?tab=Apache-2.0-1-ov-file',
              target: '_blank',
            }}
          >
            Apache 2.0 License
          </Text>
          <span> | </span>
          <Text
            link={{
              href: 'https://github.com/coze-dev/cozeloop',
              target: '_blank',
            }}
            icon={
              <IconGithub className="w-[14px] h-[14px] translate-y-[1px]" />
            }
          >
            coze-dev/cozeloop
          </Text>
        </Text>
      </div>
    </div>
  );
}
