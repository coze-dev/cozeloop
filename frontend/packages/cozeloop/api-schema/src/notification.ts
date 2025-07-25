// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import { EventEmitter } from 'eventemitter3';
import { logger } from '@coze-arch/logger';

import { HttpStatusCode } from './http-codes';

export const $notification = new EventEmitter<{
  apiError: string;
}>();

// eslint-disable-next-line @typescript-eslint/no-explicit-any -- skip
export function checkResponseData(uri: string, data: any) {
  logger.info({
    namespace: 'API',
    scope: uri,
    message: '-',
    meta: data,
  });

  if (typeof data.code === 'number' && data.code !== 0) {
    const msg = data.msg || data.message || 'Unknown error';
    throw new Error(msg);
  }
}

export function checkFetchResponse(response: Response) {
  if (
    response.status >= HttpStatusCode.OK &&
    response.status < HttpStatusCode.MultipleChoices
  ) {
    return;
  }

  switch (response.status) {
    case HttpStatusCode.BadRequest:
      throw new Error('BadRequest');
    case HttpStatusCode.Unauthorized:
      throw new Error('AuthenticationError');
    case HttpStatusCode.Forbidden:
      throw new Error('PermissionDeniedError');
    case HttpStatusCode.NotFound:
      throw new Error('NotFound');
    case HttpStatusCode.TooManyRequests:
      throw new Error('RateLimitError');
    case HttpStatusCode.RequestTimeout:
      throw new Error('TimeoutError');
    case HttpStatusCode.BadGateway:
      throw new Error('BadGateway');
    default:
      throw new Error(
        response.status >= HttpStatusCode.InternalServerError
          ? 'InternalServerError'
          : 'NetworkError',
      );
  }
}

export function onClientError(uri: string, e: unknown) {
  const error =
    e instanceof SyntaxError
      ? 'Invalid JSON error'
      : e instanceof Error
        ? e.message
        : 'Unknown error';

  $notification.emit('apiError', error);
}
