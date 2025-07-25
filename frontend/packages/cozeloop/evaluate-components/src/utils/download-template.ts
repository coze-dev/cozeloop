// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
import Papa, { type UnparseObject } from 'papaparse';

export const downloadCSVTemplate = () => {
  try {
    const fields = ['input', 'reference_output'];
    const data = [
      ['世界上最大的动物是什么', '蓝鲸'],
      ['告诉我一些这个动物的生活习性', '吃鱼'],
    ];
    const templateJson: UnparseObject<string[]> = {
      fields,
      data,
    };
    const csv = Papa.unparse(templateJson);
    downloadCsv(csv, 'dataset template');
  } catch (error) {
    console.error(error);
  }
};
export function downloadCsv(csv: string, fileName: string) {
  try {
    const BOM = '\uFEFF';
    const file = new File([BOM, csv], fileName, {
      type: 'text/csv;charset=utf-8',
    });
    const anchor = document.createElement('a');
    anchor.download = fileName;
    anchor.href = URL.createObjectURL(file);
    anchor.click();
  } catch (err) {
    console.error(err);
  }
}
