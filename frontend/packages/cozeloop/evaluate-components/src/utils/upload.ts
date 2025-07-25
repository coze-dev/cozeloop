// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import Papa from 'papaparse';

export const getCSVHeaders = (file, callback) => {
  const reader = new FileReader();
  reader.onload = function (e) {
    const text = e.target?.result as string;
    const lines = text?.split('\n');
    if (lines?.length > 0) {
      Papa.parse(lines[0], {
        header: true,
        skipEmptyLines: true,
        transformHeader(header) {
          return header.trim(); // 去除列名前后的空白
        },
        beforeFirstChunk(chunk) {
          try {
            // 分割第一行（标题行）
            const chunkLines = chunk?.split('\n') || [];
            const headers = chunkLines?.[0]?.split(',');

            // 过滤掉空的和自动生成的列名
            const validHeaders = headers?.filter(
              header =>
                header?.trim() !== '' && !header?.trim()?.match(/^_\d+$/),
            );

            // 重建第一行
            chunkLines[0] = validHeaders?.join(',');
            return chunkLines.join('\n');
            // eslint-disable-next-line @coze-arch/use-error-in-catch
          } catch (error) {
            return chunk;
          }
        },
        preview: 1,
        complete(results) {
          callback(results.meta.fields?.filter(field => !!field) ?? []);
        },
      });
    }
  };
  reader.readAsText(file.slice(0, 10240)); // 读取文件的前10KB
};
