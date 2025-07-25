// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0
/* eslint-disable @coze-arch/max-line-per-function */
import { useRef, useState } from 'react';

import { debounce } from 'lodash-es';
import cs from 'classnames';
import { I18n } from '@cozeloop/i18n-adapter';
import { GuardPoint, useGuard } from '@cozeloop/guard';
import { InfoTooltip } from '@cozeloop/components';
import { useSpace, useDataImportApi } from '@cozeloop/biz-hooks-adapter';
import { uploadFile } from '@cozeloop/biz-components-adapter';
import { type EvaluationSet } from '@cozeloop/api-schema/evaluation';
import { StorageProvider, FileFormat } from '@cozeloop/api-schema/data';
import { IconCozFileCsv } from '@coze-arch/coze-design/illustrations';
import { IconCozDownload, IconCozUpload } from '@coze-arch/coze-design/icons';
import {
  Button,
  // Button,
  Form,
  type FormApi,
  Modal,
  Typography,
  type UploadProps,
  withField,
} from '@coze-arch/coze-design';

import { getCSVHeaders } from '../../utils/upload';
import { getDefaultColumnMap } from '../../utils/import-file';
import { downloadCSVTemplate } from '../../utils/download-template';
import { useDatasetImportProgress } from './use-import-progress';
import { OverWriteField } from './overwrite-field';
import { ColumnMapField } from './column-map-field';

import styles from './index.module.less';
const FormColumnMapField = withField(ColumnMapField);
const FormOverWriteField = withField(OverWriteField);
export const DatasetImportItemsModal = ({
  onCancel,
  onOk,
  datasetDetail,
}: {
  onCancel: () => void;
  onOk: () => void;
  datasetDetail?: EvaluationSet;
}) => {
  const formRef = useRef<FormApi>();
  const { spaceID } = useSpace();
  const { importDataApi } = useDataImportApi();
  const [csvHeaders, setCsvHeaders] = useState<string[]>([]);
  const { startProgressTask, node } = useDatasetImportProgress(onOk);
  const [visible, setVisible] = useState(true);
  const [loading, setLoading] = useState(false);
  const guard = useGuard({ point: GuardPoint['eval.dataset.import'] });

  const handleUploadFile: UploadProps['customRequest'] = async ({
    fileInstance,
    file,
    onProgress,
    onSuccess,
    onError,
  }) => {
    await uploadFile({
      file: fileInstance,
      fileType: fileInstance.type?.includes('image') ? 'image' : 'object',
      onProgress,
      onSuccess,
      onError,
      spaceID,
    });
    getCSVHeaders(fileInstance, headers => {
      setCsvHeaders(headers);
      formRef?.current?.setValue(
        'fieldMappings',
        getDefaultColumnMap(datasetDetail, headers),
      );
    });
  };
  const onSubmit = async values => {
    setLoading(true);
    try {
      const res = await importDataApi({
        workspace_id: spaceID,
        dataset_id: datasetDetail?.id as string,
        file: {
          provider: StorageProvider.S3,
          path: values.file?.[0]?.response?.Uri,
          format: FileFormat.CSV,
        },
        field_mappings: values.fieldMappings?.filter(item => !!item?.source),
        option: {
          overwrite_dataset: values.overwrite,
        },
      });
      if (res.job_id) {
        startProgressTask(res.job_id);
        setVisible(false);
      }
    } finally {
      setLoading(false);
    }
  };
  const downloadCSV = debounce(downloadCSVTemplate, 400);

  return (
    <>
      <Modal
        title={I18n.t('import_data')}
        width={640}
        visible={visible}
        keepDOM={true}
        onCancel={onCancel}
        className={styles.modal}
        hasScroll={false}
        footer={null}
      >
        <Form
          initValues={{
            fieldMappings: getDefaultColumnMap(datasetDetail, csvHeaders),
            overwrite: false,
          }}
          getFormApi={formApi => {
            formRef.current = formApi;
          }}
          onValueChange={values => {
            console.log('values', values);
          }}
          onSubmit={onSubmit}
        >
          {({ formState, formApi }) => {
            const file = formState.values?.file;
            return (
              <>
                <div className={cs(styles.form, 'styled-scrollbar')}>
                  <Form.Upload
                    field="file"
                    label={I18n.t('upload_data')}
                    limit={1}
                    onChange={({ fileList }) => {
                      if (fileList.length === 0) {
                        setCsvHeaders([]);
                        formRef?.current?.setValue(
                          'fieldMappings',
                          getDefaultColumnMap(datasetDetail, []),
                        );
                      }
                    }}
                    draggable={true}
                    previewFile={() => (
                      <IconCozFileCsv className="w-[32px] h-[32px]" />
                    )}
                    className={styles.upload}
                    dragIcon={<IconCozUpload className="w-[32px] h-[32px]" />}
                    dragMainText={I18n.t('click_or_drag_file_to_upload')}
                    dragSubText={
                      <div>
                        <Typography.Text
                          className="!coz-fg-secondary"
                          size="small"
                        >
                          {I18n.t('recommend_template_upload_tip')}
                        </Typography.Text>
                        <Typography.Text
                          link
                          icon={<IconCozDownload />}
                          className="ml-[12px]"
                          size="small"
                          onClick={e => {
                            e.stopPropagation();
                            downloadCSV();
                          }}
                        >
                          {I18n.t('download_template')}
                        </Typography.Text>
                      </div>
                    }
                    action=""
                    accept=".csv"
                    customRequest={handleUploadFile}
                    rules={[
                      {
                        required: true,
                        message: I18n.t('upload_file'),
                      },
                    ]}
                  ></Form.Upload>
                  {file?.[0]?.response?.Uri ? (
                    <FormColumnMapField
                      extraTextPosition="middle"
                      field="fieldMappings"
                      extraText={
                        <Typography.Text
                          type="secondary"
                          size="small"
                          className="!coz-fg-secondary"
                        >
                          {I18n.t('no_mapping_no_import')}
                        </Typography.Text>
                      }
                      label={
                        <div className="inline-flex items-center gap-1 !coz-fg-primary">
                          <div>{I18n.t('column_mapping')}</div>
                          <InfoTooltip
                            className="h-[15px]"
                            content={I18n.t('source_column_mapping')}
                          />
                        </div>
                      }
                      sourceColumns={csvHeaders}
                      rules={[
                        {
                          required: true,
                          message: I18n.t('configure_column_mapping'),
                        },
                        {
                          validator: (_, data) => {
                            if (data?.every(item => item?.source === '')) {
                              return false;
                            }
                            return true;
                          },
                          message: I18n.t(
                            'configure_at_least_one_import_column',
                          ),
                        },
                      ]}
                    />
                  ) : null}
                  <FormOverWriteField
                    field="overwrite"
                    rules={[
                      {
                        required: true,
                        message: I18n.t('select_import_method'),
                      },
                    ]}
                    label={I18n.t('import_method')}
                  />
                </div>
                <div className="flex justify-end p-[24px] pb-0">
                  <Button
                    className="mr-2"
                    color="primary"
                    onClick={() => {
                      onCancel();
                    }}
                  >
                    {I18n.t('Cancel')}
                  </Button>
                  <Button
                    color="brand"
                    onClick={() => {
                      formRef.current?.submitForm();
                    }}
                    loading={loading}
                    disabled={guard.data.readonly || !file?.[0]?.response?.Uri}
                  >
                    {I18n.t('import')}
                  </Button>
                </div>
              </>
            );
          }}
        </Form>
      </Modal>
      {node}
    </>
  );
};
