import { type IPlugin, type Program, on } from '@coze-arch/idl2ts-plugin';
import {
  type Identifier,
  type FunctionDefinition,
  type FieldDefinition,
  type IParseResultItem,
  isStructDefinition,
  isIdentifier,
  type FunctionType,
  type IMeta,
  type IHttpRpcMapping,
  findDefinition,
  type Position,
  type ServiceDefinition,
  getStatementById,
  parseIdFiledType,
  getFieldsAlias,
  parseId,
  isFullBody,
  getAstFromNamespace,
  getSchemaRootByPath,
  getAnnotation,
} from '@coze-arch/idl2ts-helper';

import { type Contexts, HOOK, type IProcessMetaItemCtx } from '../context';

interface IOptions {
  outputDir: string;
  idlRoot: string;
}

export class MetaPlugin implements IPlugin {
  options: IOptions;
  constructor(options: IOptions) {
    this.options = options;
  }
  apply(program: Program<Contexts>): void {
    program.register(
      on(HOOK.PARSE_FUN_META),
      ctx => {
        const node = ctx.node as ServiceDefinition;
        node.functions.forEach(fun => {
          // 过滤非泛化接口
          if (!fun.extensionConfig?.method) {
            return;
          }
          const { meta } = program.trigger<IProcessMetaItemCtx>(
            HOOK.PARSE_FUN_META_ITEM,
            {
              ast: ctx.ast,
              service: node,
              method: fun,
            } as IProcessMetaItemCtx,
          );
          ctx.meta.push(meta);
        });
        return ctx;
      },
      0,
    );
    program.register(on(HOOK.PARSE_FUN_META_ITEM), ctx => {
      const { ast, service, method } = ctx;
      const item = this.parseFunAnnotation(method, ast, service.name.value);
      ctx.meta = item;
      return ctx;
    });
  }
  parseFunAnnotation(
    params: FunctionDefinition,
    ast: IParseResultItem,
    service: string,
  ) {
    const { name, returnType, fields, extensionConfig } = params;
    const reqType = fields[0].fieldType as any;
    const reqMapping = this.processPayloadFields(
      reqType,
      extensionConfig?.method === 'GET' ? 'query' : 'body',
      ast,
    );
    const res = {
      url: extensionConfig?.uri,
      method: extensionConfig?.method ?? 'POST',
      name: name.value,
      reqType: parseId(reqType.value),
      reqMapping,
      resType: parseId(this.processReqResPramsType(returnType, ast)),
      schemaRoot: getSchemaRootByPath(ast.idlPath, this.options.idlRoot),
      service,
    } as IMeta;
    // 不是 json 时，需要加上 serializer 标识
    if (extensionConfig?.serializer && extensionConfig?.serializer !== 'json') {
      res.serializer = extensionConfig?.serializer;
    }

    return res;
  }

  private processReqResPramsType(id: FunctionType, ast: IParseResultItem) {
    if (isIdentifier(id)) {
      const statement = getStatementById(id, ast);
      if (isStructDefinition(statement)) {
        const wholeBody = statement.fields.find(isFullBody);
        if (wholeBody) {
          // 处理 api.body="." 以及 api.full_body=''
          return `${id.value}['${getFieldsAlias(wholeBody)}']`;
        } else {
          return id.value;
        }
      }
      throw new Error('params must be identifier');
    } else {
      return 'void';
    }
  }

  private processPayloadFields(
    id: Identifier,
    defaultPosition: 'query' | 'body',
    entry: IParseResultItem,
  ): IHttpRpcMapping {
    const { namespace, refName } = parseIdFiledType(id);

    if (namespace) {
      const ast = getAstFromNamespace(namespace, entry);
      const struct = findDefinition(ast, refName);
      if (!struct || !isStructDefinition(struct)) {
        throw new Error(`can not find Struct: ${refName} `);
      }
      return this.createMapping(struct.fields, defaultPosition);
    }

    const struct = findDefinition(entry.statements, id.value);
    if (!struct || !isStructDefinition(struct)) {
      throw new Error(`can not find Struct: ${id.value} `);
    }

    return this.createMapping(struct.fields, defaultPosition);
  }

  private createMapping(
    fields: FieldDefinition[],
    defaultPosition: 'query' | 'body',
  ): IHttpRpcMapping {
    const specificPositionFiled = new Set<string>();
    const mapping = {} as IHttpRpcMapping;
    fields.forEach(filed => {
      const jsonAnnotation = getAnnotation(filed.annotations, 'api.json');
      if (jsonAnnotation) {
        filed.extensionConfig = filed.extensionConfig || {};
        filed.extensionConfig.key = jsonAnnotation;
      }
      const { extensionConfig } = filed;
      const alias = getFieldsAlias(filed);
      if (extensionConfig) {
        if (isFullBody(filed)) {
          mapping.entire_body = [alias];
          return;
        }
        Object.keys(extensionConfig).forEach(key => {
          if (key === 'position' && extensionConfig.position) {
            const filedMapping = this.processMapping(
              mapping,
              extensionConfig.position,
              alias,
            );
            mapping[extensionConfig.position] = filedMapping;
            specificPositionFiled.add(alias);
          }
        });
      }
      // 如果没有指定，根据method默认指定为query 或者 body
      if (!specificPositionFiled.has(alias)) {
        const filedMapping = mapping[defaultPosition];
        mapping[defaultPosition] = filedMapping
          ? [...filedMapping, alias]
          : [alias];
      }
    });
    return mapping;
  }

  private processMapping(
    mapping: IHttpRpcMapping,
    position: Position,
    filedName: string,
  ): string[] {
    const mappingKeys = [
      'path',
      'query',
      'status_code',
      'header',
      'cookie',
      'entire_body',
      'body',
    ];
    if (mappingKeys.find(i => i === position)) {
      const data = mapping[position];
      return data ? [...data, filedName] : [filedName];
    } else {
      return mapping[position] || [];
    }
  }
}
