// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/component/config"
	confmocks "github.com/coze-dev/cozeloop/backend/modules/observability/domain/component/config/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/component/metrics"
	metricmocks "github.com/coze-dev/cozeloop/backend/modules/observability/domain/component/metrics/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/component/mq"
	mqmocks "github.com/coze-dev/cozeloop/backend/modules/observability/domain/component/mq/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/loop_span"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/repo"
	repomocks "github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/repo/mocks"
	filtermocks "github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/service/trace/span_filter/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/service/trace/span_processor"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
)

func TestTraceServiceImpl_GetTracesAdvanceInfo(t *testing.T) {
	type fields struct {
		traceRepo     repo.ITraceRepo
		traceConfig   config.ITraceConfig
		traceProducer mq.ITraceProducer
		metrics       metrics.ITraceMetrics
		buildHelper   TraceFilterProcessorBuilder
	}
	type args struct {
		ctx context.Context
		req *GetTracesAdvanceInfoReq
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *GetTracesAdvanceInfoResp
		wantErr      bool
	}{
		{
			name: "get traces advance info successfully",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				repoMock := repomocks.NewMockITraceRepo(ctrl)
				repoMock.EXPECT().GetTrace(gomock.Any(), gomock.Any()).Return(loop_span.SpanList{{
					TraceID: "123",
					SpanID:  "234",
				}}, nil)
				filterFactoryMock := filtermocks.NewMockPlatformFilterFactory(ctrl)
				buildHelper := NewTraceFilterProcessorBuilder(filterFactoryMock, nil, nil, nil)
				metricsMock := metricmocks.NewMockITraceMetrics(ctrl)
				metricsMock.EXPECT().EmitGetTrace(gomock.Any(), gomock.Any(), gomock.Any()).Return()
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				return fields{
					traceRepo:   repoMock,
					traceConfig: confMock,
					buildHelper: buildHelper,
					metrics:     metricsMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &GetTracesAdvanceInfoReq{
					WorkspaceID:  1,
					PlatformType: loop_span.PlatformCozeLoop,
					Traces: []*TraceQueryParam{{
						TraceID:   "123",
						StartTime: 0,
						EndTime:   0,
					}},
				},
			},
			want: &GetTracesAdvanceInfoResp{
				Infos: []*loop_span.TraceAdvanceInfo{{
					TraceId:    "123",
					InputCost:  0,
					OutputCost: 0,
				}},
			},
		},
		{
			name: "get traces advance info successfully with processor",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				repoMock := repomocks.NewMockITraceRepo(ctrl)
				repoMock.EXPECT().GetTrace(gomock.Any(), gomock.Any()).Return(loop_span.SpanList{{
					TraceID:     "123",
					SpanID:      "234",
					WorkspaceID: "123",
				}}, nil)
				filterFactoryMock := filtermocks.NewMockPlatformFilterFactory(ctrl)
				buildHelper := NewTraceFilterProcessorBuilder(filterFactoryMock,
					nil,
					nil,
					[]span_processor.Factory{span_processor.NewCheckProcessorFactory()})
				metricsMock := metricmocks.NewMockITraceMetrics(ctrl)
				metricsMock.EXPECT().EmitGetTrace(gomock.Any(), gomock.Any(), gomock.Any()).Return()
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				return fields{
					traceRepo:   repoMock,
					traceConfig: confMock,
					buildHelper: buildHelper,
					metrics:     metricsMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &GetTracesAdvanceInfoReq{
					WorkspaceID:  123,
					PlatformType: loop_span.PlatformCozeLoop,
					Traces: []*TraceQueryParam{{
						TraceID:   "123",
						StartTime: 0,
						EndTime:   0,
					}},
				},
			},
			want: &GetTracesAdvanceInfoResp{
				Infos: []*loop_span.TraceAdvanceInfo{{
					TraceId:    "123",
					InputCost:  0,
					OutputCost: 0,
				}},
			},
		},
		{
			name: "get traces advance info failed due to repo error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				repoMock := repomocks.NewMockITraceRepo(ctrl)
				repoMock.EXPECT().GetTrace(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("repo error"))
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				filterFactoryMock := filtermocks.NewMockPlatformFilterFactory(ctrl)
				buildHelper := NewTraceFilterProcessorBuilder(filterFactoryMock, nil, nil, nil)
				metricsMock := metricmocks.NewMockITraceMetrics(ctrl)
				metricsMock.EXPECT().EmitGetTrace(gomock.Any(), gomock.Any(), gomock.Any()).Return()
				return fields{
					traceRepo:   repoMock,
					traceConfig: confMock,
					metrics:     metricsMock,
					buildHelper: buildHelper,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &GetTracesAdvanceInfoReq{
					WorkspaceID:  1,
					PlatformType: loop_span.PlatformCozeLoop,
					Traces: []*TraceQueryParam{{
						TraceID:   "123",
						StartTime: 0,
						EndTime:   0,
					}},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			r, _ := NewTraceServiceImpl(
				fields.traceRepo,
				fields.traceConfig,
				fields.traceProducer,
				fields.metrics,
				fields.buildHelper)
			got, err := r.GetTracesAdvanceInfo(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestTraceServiceImpl_IngestTraces(t *testing.T) {
	type fields struct {
		traceRepo     repo.ITraceRepo
		traceConfig   config.ITraceConfig
		traceProducer mq.ITraceProducer
		metrics       metrics.ITraceMetrics
		buildHelper   TraceFilterProcessorBuilder
	}
	type args struct {
		ctx context.Context
		req *IngestTracesReq
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantErr      bool
	}{
		{
			name: "ingest traces successfully",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				producerMock := mqmocks.NewMockITraceProducer(ctrl)
				producerMock.EXPECT().IngestSpans(gomock.Any(), gomock.Any()).Return(nil)
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetDefaultTraceTenant(gomock.Any()).Return("abc")
				return fields{
					traceProducer: producerMock,
					traceConfig:   confMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &IngestTracesReq{
					TTL: entity.TTL3d,
					Spans: loop_span.SpanList{{
						TraceID:     "123",
						SpanID:      "234",
						WorkspaceID: "1",
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "ingest traces failed due to producer error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				producerMock := mqmocks.NewMockITraceProducer(ctrl)
				producerMock.EXPECT().IngestSpans(gomock.Any(), gomock.Any()).Return(fmt.Errorf("producer error"))
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetDefaultTraceTenant(gomock.Any()).Return("abc")
				return fields{
					traceProducer: producerMock,
					traceConfig:   confMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &IngestTracesReq{
					TTL: entity.TTL3d,
					Spans: loop_span.SpanList{{
						TraceID:     "123",
						SpanID:      "234",
						WorkspaceID: "1",
					}},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			r := &TraceServiceImpl{
				traceRepo:     fields.traceRepo,
				traceConfig:   fields.traceConfig,
				traceProducer: fields.traceProducer,
				metrics:       fields.metrics,
				buildHelper:   fields.buildHelper,
			}
			err := r.IngestTraces(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestTraceServiceImpl_GetTracesMetaInfo(t *testing.T) {
	type fields struct {
		traceRepo     repo.ITraceRepo
		traceConfig   config.ITraceConfig
		traceProducer mq.ITraceProducer
		metrics       metrics.ITraceMetrics
		buildHelper   TraceFilterProcessorBuilder
	}
	type args struct {
		ctx context.Context
		req *GetTracesMetaInfoReq
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *GetTracesMetaInfoResp
		wantErr      bool
	}{
		{
			name: "get traces meta info successfully",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetTraceFieldMetaInfo(gomock.Any()).Return(&config.TraceFieldMetaInfoCfg{
					FieldMetas: map[loop_span.PlatformType]map[loop_span.SpanListType][]string{
						loop_span.PlatformCozeLoop: {
							loop_span.SpanListTypeAllSpan: {"field1", "field2"},
						},
					},
					AvailableFields: map[string]*config.FieldMeta{
						"field1": {FieldType: "string"},
						"field2": {FieldType: "int"},
					},
				}, nil)
				return fields{
					traceConfig: confMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &GetTracesMetaInfoReq{
					WorkspaceID:  1,
					PlatformType: loop_span.PlatformCozeLoop,
					SpanListType: loop_span.SpanListTypeAllSpan,
				},
			},
			want: &GetTracesMetaInfoResp{
				FilesMetas: map[string]*config.FieldMeta{
					"field1": {FieldType: "string"},
					"field2": {FieldType: "int"},
				},
			},
		},
		{
			name: "get traces meta info failed due to config error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetTraceFieldMetaInfo(gomock.Any()).Return(nil, fmt.Errorf("config error"))
				return fields{
					traceConfig: confMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &GetTracesMetaInfoReq{
					WorkspaceID:  1,
					PlatformType: loop_span.PlatformCozeLoop,
					SpanListType: loop_span.SpanListTypeAllSpan,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			r := &TraceServiceImpl{
				traceRepo:     fields.traceRepo,
				traceConfig:   fields.traceConfig,
				traceProducer: fields.traceProducer,
				metrics:       fields.metrics,
				buildHelper:   fields.buildHelper,
			}
			got, err := r.GetTracesMetaInfo(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestTraceServiceImpl_ListSpans(t *testing.T) {
	type fields struct {
		traceRepo     repo.ITraceRepo
		traceConfig   config.ITraceConfig
		traceProducer mq.ITraceProducer
		metrics       metrics.ITraceMetrics
		buildHelper   TraceFilterProcessorBuilder
	}
	type args struct {
		ctx context.Context
		req *ListSpansReq
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *ListSpansResp
		wantErr      bool
	}{
		{
			name: "list spans successfully",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				repoMock := repomocks.NewMockITraceRepo(ctrl)
				repoMock.EXPECT().ListSpans(gomock.Any(), gomock.Any()).Return(&repo.ListSpansResult{
					Spans: loop_span.SpanList{{
						TraceID: "123",
						SpanID:  "234",
					}},
					PageToken: "",
					HasMore:   false,
				}, nil)
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				filterMock := filtermocks.NewMockFilter(ctrl)
				filterMock.EXPECT().BuildBasicSpanFilter(gomock.Any(), gomock.Any()).Return([]*loop_span.FilterField{
					{
						FieldName: loop_span.SpanFieldSpaceId,
						FieldType: loop_span.FieldTypeString,
						Values:    []string{"123"},
						QueryType: ptr.Of(loop_span.QueryTypeEnumIn),
					},
				}, nil)
				filterMock.EXPECT().BuildALLSpanFilter(gomock.Any(), gomock.Any()).Return(nil, nil)
				filterFactoryMock := filtermocks.NewMockPlatformFilterFactory(ctrl)
				filterFactoryMock.EXPECT().GetFilter(gomock.Any(), gomock.Any()).Return(filterMock, nil)
				buildHelper := NewTraceFilterProcessorBuilder(filterFactoryMock, nil, nil, nil)
				metricsMock := metricmocks.NewMockITraceMetrics(ctrl)
				metricsMock.EXPECT().EmitListSpans(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return()
				return fields{
					traceRepo:   repoMock,
					traceConfig: confMock,
					buildHelper: buildHelper,
					metrics:     metricsMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &ListSpansReq{
					PlatformType: loop_span.PlatformCozeLoop,
					Limit:        10,
					SpanListType: loop_span.SpanListTypeAllSpan,
				},
			},
			want: &ListSpansResp{
				Spans: loop_span.SpanList{{
					TraceID: "123",
					SpanID:  "234",
				}},
			},
		},
		{
			name: "list spans successfully with specific filter",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				repoMock := repomocks.NewMockITraceRepo(ctrl)
				repoMock.EXPECT().ListSpans(gomock.Any(), gomock.Any()).Return(&repo.ListSpansResult{
					Spans: loop_span.SpanList{{
						TraceID: "123",
						SpanID:  "234",
					}},
					PageToken: "",
					HasMore:   false,
				}, nil)
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				filterMock := filtermocks.NewMockFilter(ctrl)
				filterMock.EXPECT().BuildBasicSpanFilter(gomock.Any(), gomock.Any()).Return([]*loop_span.FilterField{{}}, nil)
				filterMock.EXPECT().BuildALLSpanFilter(gomock.Any(), gomock.Any()).Return(nil, nil)
				filterFactoryMock := filtermocks.NewMockPlatformFilterFactory(ctrl)
				filterFactoryMock.EXPECT().GetFilter(gomock.Any(), gomock.Any()).Return(filterMock, nil)
				buildHelper := NewTraceFilterProcessorBuilder(filterFactoryMock, nil, nil, nil)
				metricsMock := metricmocks.NewMockITraceMetrics(ctrl)
				metricsMock.EXPECT().EmitListSpans(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return()
				return fields{
					traceRepo:   repoMock,
					traceConfig: confMock,
					buildHelper: buildHelper,
					metrics:     metricsMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &ListSpansReq{
					PlatformType: loop_span.PlatformCozeLoop,
					Limit:        10,
					SpanListType: loop_span.SpanListTypeAllSpan,
					Filters: &loop_span.FilterFields{
						QueryAndOr: nil,
						FilterFields: []*loop_span.FilterField{
							{
								FieldName: "status",
								FieldType: loop_span.FieldTypeString,
								Values:    []string{"success"},
								QueryType: ptr.Of(loop_span.QueryTypeEnumIn),
							},
							{
								FieldName: "status",
								FieldType: loop_span.FieldTypeString,
								Values:    []string{"success", "error"},
								QueryType: ptr.Of(loop_span.QueryTypeEnumIn),
							},
							{
								FieldName: "status",
								FieldType: loop_span.FieldTypeString,
								Values:    []string{"error"},
								QueryType: ptr.Of(loop_span.QueryTypeEnumIn),
							},
							{
								FieldName: loop_span.SpanFieldStartTimeFirstResp,
								FieldType: loop_span.FieldTypeLong,
								Values:    []string{"1234"},
								QueryType: ptr.Of(loop_span.QueryTypeEnumGte),
							},
						},
					},
				},
			},
			want: &ListSpansResp{
				Spans: loop_span.SpanList{{
					TraceID: "123",
					SpanID:  "234",
				}},
			},
		},
		{
			name: "list spans successfully with root span",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				repoMock := repomocks.NewMockITraceRepo(ctrl)
				repoMock.EXPECT().ListSpans(gomock.Any(), gomock.Any()).Return(&repo.ListSpansResult{
					Spans: loop_span.SpanList{{
						TraceID: "123",
						SpanID:  "234",
					}},
					PageToken: "",
					HasMore:   false,
				}, nil)
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				filterMock := filtermocks.NewMockFilter(ctrl)
				filterMock.EXPECT().BuildBasicSpanFilter(gomock.Any(), gomock.Any()).Return([]*loop_span.FilterField{{}}, nil)
				filterMock.EXPECT().BuildRootSpanFilter(gomock.Any(), gomock.Any()).Return(nil, nil)
				filterFactoryMock := filtermocks.NewMockPlatformFilterFactory(ctrl)
				filterFactoryMock.EXPECT().GetFilter(gomock.Any(), gomock.Any()).Return(filterMock, nil)
				buildHelper := NewTraceFilterProcessorBuilder(filterFactoryMock, nil, nil, nil)
				metricsMock := metricmocks.NewMockITraceMetrics(ctrl)
				metricsMock.EXPECT().EmitListSpans(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return()
				return fields{
					traceRepo:   repoMock,
					traceConfig: confMock,
					buildHelper: buildHelper,
					metrics:     metricsMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &ListSpansReq{
					PlatformType: loop_span.PlatformCozeLoop,
					Limit:        10,
					SpanListType: loop_span.SpanListTypeRootSpan,
				},
			},
			want: &ListSpansResp{
				Spans: loop_span.SpanList{{
					TraceID: "123",
					SpanID:  "234",
				}},
			},
		},
		{
			name: "list spans successfully with llm span",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				repoMock := repomocks.NewMockITraceRepo(ctrl)
				repoMock.EXPECT().ListSpans(gomock.Any(), gomock.Any()).Return(&repo.ListSpansResult{
					Spans: loop_span.SpanList{{
						TraceID: "123",
						SpanID:  "234",
					}},
					PageToken: "",
					HasMore:   false,
				}, nil)
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				filterMock := filtermocks.NewMockFilter(ctrl)
				filterMock.EXPECT().BuildBasicSpanFilter(gomock.Any(), gomock.Any()).Return([]*loop_span.FilterField{{}}, nil)
				filterMock.EXPECT().BuildLLMSpanFilter(gomock.Any(), gomock.Any()).Return(nil, nil)
				filterFactoryMock := filtermocks.NewMockPlatformFilterFactory(ctrl)
				filterFactoryMock.EXPECT().GetFilter(gomock.Any(), gomock.Any()).Return(filterMock, nil)
				buildHelper := NewTraceFilterProcessorBuilder(filterFactoryMock, nil, nil, nil)
				metricsMock := metricmocks.NewMockITraceMetrics(ctrl)
				metricsMock.EXPECT().EmitListSpans(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return()
				return fields{
					traceRepo:   repoMock,
					traceConfig: confMock,
					buildHelper: buildHelper,
					metrics:     metricsMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &ListSpansReq{
					PlatformType: loop_span.PlatformCozeLoop,
					Limit:        10,
					SpanListType: loop_span.SpanListTypeLLMSpan,
				},
			},
			want: &ListSpansResp{
				Spans: loop_span.SpanList{{
					TraceID: "123",
					SpanID:  "234",
				}},
			},
		},
		{
			name: "list spans successfully with processor",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				repoMock := repomocks.NewMockITraceRepo(ctrl)
				repoMock.EXPECT().ListSpans(gomock.Any(), gomock.Any()).Return(&repo.ListSpansResult{
					Spans: loop_span.SpanList{{
						TraceID:     "123",
						SpanID:      "234",
						WorkspaceID: "123",
					}},
					PageToken: "",
					HasMore:   false,
				}, nil)
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				filterMock := filtermocks.NewMockFilter(ctrl)
				filterMock.EXPECT().BuildBasicSpanFilter(gomock.Any(), gomock.Any()).Return([]*loop_span.FilterField{{}}, nil)
				filterMock.EXPECT().BuildALLSpanFilter(gomock.Any(), gomock.Any()).Return(nil, nil)
				filterFactoryMock := filtermocks.NewMockPlatformFilterFactory(ctrl)
				filterFactoryMock.EXPECT().GetFilter(gomock.Any(), gomock.Any()).Return(filterMock, nil)
				buildHelper := NewTraceFilterProcessorBuilder(filterFactoryMock,
					nil,
					[]span_processor.Factory{
						span_processor.NewCheckProcessorFactory(),
					},
					nil)
				metricsMock := metricmocks.NewMockITraceMetrics(ctrl)
				metricsMock.EXPECT().EmitListSpans(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return()
				return fields{
					traceRepo:   repoMock,
					traceConfig: confMock,
					buildHelper: buildHelper,
					metrics:     metricsMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &ListSpansReq{
					PlatformType: loop_span.PlatformCozeLoop,
					Limit:        10,
					SpanListType: loop_span.SpanListTypeAllSpan,
					WorkspaceID:  123,
				},
			},
			want: &ListSpansResp{
				Spans: loop_span.SpanList{{
					TraceID:     "123",
					SpanID:      "234",
					WorkspaceID: "123",
				}},
			},
		},
		{
			name: "list spans successfully with processor failed",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				repoMock := repomocks.NewMockITraceRepo(ctrl)
				repoMock.EXPECT().ListSpans(gomock.Any(), gomock.Any()).Return(&repo.ListSpansResult{
					Spans: loop_span.SpanList{{
						TraceID:     "123",
						SpanID:      "234",
						WorkspaceID: "1234",
					}},
					PageToken: "",
					HasMore:   false,
				}, nil)
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				filterMock := filtermocks.NewMockFilter(ctrl)
				filterMock.EXPECT().BuildBasicSpanFilter(gomock.Any(), gomock.Any()).Return([]*loop_span.FilterField{{}}, nil)
				filterMock.EXPECT().BuildALLSpanFilter(gomock.Any(), gomock.Any()).Return(nil, nil)
				filterFactoryMock := filtermocks.NewMockPlatformFilterFactory(ctrl)
				filterFactoryMock.EXPECT().GetFilter(gomock.Any(), gomock.Any()).Return(filterMock, nil)
				buildHelper := NewTraceFilterProcessorBuilder(filterFactoryMock,
					nil,
					[]span_processor.Factory{
						span_processor.NewCheckProcessorFactory(),
					},
					nil)
				metricsMock := metricmocks.NewMockITraceMetrics(ctrl)
				metricsMock.EXPECT().EmitListSpans(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return()
				return fields{
					traceRepo:   repoMock,
					traceConfig: confMock,
					buildHelper: buildHelper,
					metrics:     metricsMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &ListSpansReq{
					PlatformType: loop_span.PlatformCozeLoop,
					Limit:        10,
					SpanListType: loop_span.SpanListTypeAllSpan,
					WorkspaceID:  123,
				},
			},
			wantErr: true,
		},
		{
			name: "list spans failed due to invalid platform type",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				filterMock := filtermocks.NewMockFilter(ctrl)
				filterMock.EXPECT().BuildBasicSpanFilter(gomock.Any(), gomock.Any()).Return([]*loop_span.FilterField{{}}, nil)
				filterMock.EXPECT().BuildALLSpanFilter(gomock.Any(), gomock.Any()).Return(nil, nil)
				filterFactoryMock := filtermocks.NewMockPlatformFilterFactory(ctrl)
				filterFactoryMock.EXPECT().GetFilter(gomock.Any(), gomock.Any()).Return(filterMock, nil)
				buildHelper := NewTraceFilterProcessorBuilder(filterFactoryMock, nil, nil, nil)
				return fields{
					traceConfig: confMock,
					buildHelper: buildHelper,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &ListSpansReq{
					PlatformType: "abc",
					Limit:        10,
					SpanListType: loop_span.SpanListTypeAllSpan,
				},
			},
			wantErr: true,
		},
		{
			name: "list spans failed due to repo error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				repoMock := repomocks.NewMockITraceRepo(ctrl)
				repoMock.EXPECT().ListSpans(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("failed"))
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				filterMock := filtermocks.NewMockFilter(ctrl)
				filterMock.EXPECT().BuildBasicSpanFilter(gomock.Any(), gomock.Any()).Return([]*loop_span.FilterField{{}}, nil)
				filterMock.EXPECT().BuildALLSpanFilter(gomock.Any(), gomock.Any()).Return(nil, nil)
				filterFactoryMock := filtermocks.NewMockPlatformFilterFactory(ctrl)
				filterFactoryMock.EXPECT().GetFilter(gomock.Any(), gomock.Any()).Return(filterMock, nil)
				buildHelper := NewTraceFilterProcessorBuilder(filterFactoryMock, nil, nil, nil)
				metricsMock := metricmocks.NewMockITraceMetrics(ctrl)
				metricsMock.EXPECT().EmitListSpans(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return()
				return fields{
					traceRepo:   repoMock,
					traceConfig: confMock,
					metrics:     metricsMock,
					buildHelper: buildHelper,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &ListSpansReq{
					PlatformType: loop_span.PlatformCozeLoop,
					Limit:        10,
					SpanListType: loop_span.SpanListTypeAllSpan,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			r := &TraceServiceImpl{
				traceRepo:     fields.traceRepo,
				traceConfig:   fields.traceConfig,
				traceProducer: fields.traceProducer,
				metrics:       fields.metrics,
				buildHelper:   fields.buildHelper,
			}
			got, err := r.ListSpans(tt.args.ctx, tt.args.req)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestTraceServiceImpl_GetTrace(t *testing.T) {
	type fields struct {
		traceRepo     repo.ITraceRepo
		traceConfig   config.ITraceConfig
		traceProducer mq.ITraceProducer
		metrics       metrics.ITraceMetrics
		buildHelper   TraceFilterProcessorBuilder
	}
	type args struct {
		ctx context.Context
		req *GetTraceReq
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *GetTraceResp
		wantErr      bool
	}{
		{
			name: "get trace successfully",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				repoMock := repomocks.NewMockITraceRepo(ctrl)
				repoMock.EXPECT().GetTrace(gomock.Any(), gomock.Any()).Return(loop_span.SpanList{
					{
						TraceID: "123",
						SpanID:  "234",
					},
				}, nil)
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				filterFactoryMock := filtermocks.NewMockPlatformFilterFactory(ctrl)
				buildHelper := NewTraceFilterProcessorBuilder(filterFactoryMock, nil, nil, nil)
				metricsMock := metricmocks.NewMockITraceMetrics(ctrl)
				metricsMock.EXPECT().EmitGetTrace(gomock.Any(), gomock.Any(), gomock.Any()).Return()
				return fields{
					traceRepo:   repoMock,
					traceConfig: confMock,
					buildHelper: buildHelper,
					metrics:     metricsMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &GetTraceReq{
					PlatformType: loop_span.PlatformCozeLoop,
					TraceID:      "123",
				},
			},
			want: &GetTraceResp{
				TraceId: "123",
				Spans: loop_span.SpanList{
					{
						TraceID: "123",
						SpanID:  "234",
					},
				},
			},
		},
		{
			name: "get trace successfully with processor",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				repoMock := repomocks.NewMockITraceRepo(ctrl)
				repoMock.EXPECT().GetTrace(gomock.Any(), gomock.Any()).Return(loop_span.SpanList{
					{
						TraceID:     "123",
						SpanID:      "234",
						WorkspaceID: "123",
					},
				}, nil)
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				filterFactoryMock := filtermocks.NewMockPlatformFilterFactory(ctrl)
				buildHelper := NewTraceFilterProcessorBuilder(filterFactoryMock,
					[]span_processor.Factory{span_processor.NewCheckProcessorFactory()},
					nil,
					nil)
				metricsMock := metricmocks.NewMockITraceMetrics(ctrl)
				metricsMock.EXPECT().EmitGetTrace(gomock.Any(), gomock.Any(), gomock.Any()).Return()
				return fields{
					traceRepo:   repoMock,
					traceConfig: confMock,
					buildHelper: buildHelper,
					metrics:     metricsMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &GetTraceReq{
					PlatformType: loop_span.PlatformCozeLoop,
					TraceID:      "123",
					WorkspaceID:  123,
				},
			},
			want: &GetTraceResp{
				TraceId: "123",
				Spans: loop_span.SpanList{
					{
						TraceID:     "123",
						SpanID:      "234",
						WorkspaceID: "123",
					},
				},
			},
		},
		{
			name: "get failed due to invalid platform type",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				return fields{
					traceConfig: confMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &GetTraceReq{
					PlatformType: "abc",
					TraceID:      "123",
				},
			},
			wantErr: true,
		},
		{
			name: "get failed dur to repo error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				repoMock := repomocks.NewMockITraceRepo(ctrl)
				repoMock.EXPECT().GetTrace(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("failed"))
				confMock := confmocks.NewMockITraceConfig(ctrl)
				confMock.EXPECT().GetPlatformTenants(gomock.Any()).Return(&config.PlatformTenantsCfg{
					Config: map[string][]string{
						string(loop_span.PlatformCozeLoop): {"spans"},
					},
				}, nil)
				metricsMock := metricmocks.NewMockITraceMetrics(ctrl)
				metricsMock.EXPECT().EmitGetTrace(gomock.Any(), gomock.Any(), gomock.Any()).Return()
				return fields{
					traceRepo:   repoMock,
					traceConfig: confMock,
					metrics:     metricsMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &GetTraceReq{
					PlatformType: loop_span.PlatformCozeLoop,
					TraceID:      "123",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			r := &TraceServiceImpl{
				traceRepo:     fields.traceRepo,
				traceConfig:   fields.traceConfig,
				traceProducer: fields.traceProducer,
				metrics:       fields.metrics,
				buildHelper:   fields.buildHelper,
			}
			got, err := r.GetTrace(tt.args.ctx, tt.args.req)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}
