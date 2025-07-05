import {
  type MessageFormatElement,
  parse,
  TYPE,
} from '@formatjs/icu-messageformat-parser';

interface TypeInfo {
  key: string;
  type: 'string' | 'number' | 'Date' | string;
}

/**
 * Parse {@link MessageFormatElement}
 *
 * See {@link https://formatjs.github.io/docs/core-concepts/icu-syntax#basic-principles}
 */
function parseElement(el: MessageFormatElement): TypeInfo | undefined {
  switch (el.type) {
    case TYPE.literal:
      return undefined;
    case TYPE.argument:
      return { key: el.value, type: 'string' };
    case TYPE.number:
    case TYPE.date:
    case TYPE.time:
    case TYPE.plural:
      return { key: el.value, type: 'number' };
    case TYPE.select:
      return {
        key: el.value,
        type: Object.keys(el.options)
          .map(it => (it === 'other' ? 'undefined' : `'${it}'`))
          .join(' | '),
      };
    case TYPE.pound:
    case TYPE.tag:
    default:
      return undefined;
  }
}

export function icu2Type(message: string) {
  const elements = parse(message);

  return elements
    .map(it => parseElement(it))
    .filter((it): it is TypeInfo => Boolean(it));
}
