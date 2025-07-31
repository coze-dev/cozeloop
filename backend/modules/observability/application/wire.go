// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

//go:build wireinject
// +build wireinject

package application

import (
	"github.com/google/wire"

	"github.com/coze-dev/coze-loop/backend/infra/ck"
	"github.com/coze-dev/coze-loop/backend/infra/db"
	"github.com/coze-dev/coze-loop/backend/infra/external/benefit"
	"github.com/coze-dev/coze-loop/backend/infra/metrics"
	"github.com/coze-dev/coze-loop/backend/infra/mq"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/foundation/auth/authservice"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/foundation/file/fileservice"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/component/config"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/component/rpc"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/exporter"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/processor"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/collector/receiver"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/repo"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/service"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/service/collector/exporter/clickhouseexporter"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/service/collector/processor/queueprocessor"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/service/collector/receiver/rmqreceiver"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/service/trace/span_filter"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/service/trace/span_processor"
	"github.com/coze-dev/coze-loop/backend/modules/observability/infra/auth"
	obconfig "github.com/coze-dev/coze-loop/backend/modules/observability/infra/config"
	"github.com/coze-dev/coze-loop/backend/modules/observability/infra/file"
	obmetrics "github.com/coze-dev/coze-loop/backend/modules/observability/infra/metrics"
	mq2 "github.com/coze-dev/coze-loop/backend/modules/observability/infra/mq"
	obrepo "github.com/coze-dev/coze-loop/backend/modules/observability/infra/repo"
	ckdao "github.com/coze-dev/coze-loop/backend/modules/observability/infra/repo/ck"
	mysqldao "github.com/coze-dev/coze-loop/backend/modules/observability/infra/repo/mysql"
	"github.com/coze-dev/coze-loop/backend/pkg/conf"
)

var traceSet = wire.NewSet(
	NewTraceApplication,
	service.NewTraceServiceImpl,
	obrepo.NewTraceCKRepoImpl,
	ckdao.NewSpansCkDaoImpl,
	obrepo.NewViewRepoImpl,
	mysqldao.NewViewDaoImpl,
	auth.NewAuthProvider,
	obmetrics.NewTraceMetricsImpl,
	mq2.NewTraceProducerImpl,
	obconfig.NewTraceConfigCenter,
	NewTraceConfigLoader,
	NewTraceQueryProcessorBuilder,
	file.NewFileRPCProvider,
)

var traceIngestionSet = wire.NewSet(
	NewIngestionApplication,
	service.NewIngestionServiceImpl,
	obrepo.NewTraceCKRepoImpl,
	ckdao.NewSpansCkDaoImpl,
	obconfig.NewTraceConfigCenter,
	NewTraceConfigLoader,
	NewIngestionCollectorFactory,
)

func NewTraceQueryProcessorBuilder(
	traceConfig config.ITraceConfig,
	fileProvider rpc.IFileProvider,
	benefitSvc benefit.IBenefitService,
) service.TraceFilterProcessorBuilder {
	return service.NewTraceFilterProcessorBuilder(
		span_filter.NewPlatformFilterFactory(
			[]span_filter.Factory{
				span_filter.NewCozeLoopFilterFactory(),
				span_filter.NewPromptFilterFactory(traceConfig),
				span_filter.NewEvaluatorFilterFactory(),
				span_filter.NewEvalTargetFilterFactory(),
			}),
		// get trace processors
		[]span_processor.Factory{
			span_processor.NewPlatformProcessorFactory(traceConfig),
			span_processor.NewCheckProcessorFactory(),
			span_processor.NewAttrTosProcessorFactory(fileProvider),
			span_processor.NewExpireErrorProcessorFactory(benefitSvc),
		},
		// list spans processors
		[]span_processor.Factory{
			span_processor.NewAttrTosProcessorFactory(fileProvider),
			span_processor.NewExpireErrorProcessorFactory(benefitSvc),
		},
		// batch get advance info processors
		[]span_processor.Factory{
			span_processor.NewCheckProcessorFactory(),
		})
}

func NewIngestionCollectorFactory(mqFactory mq.IFactory, traceRepo repo.ITraceRepo) service.IngestionCollectorFactory {
	return service.NewIngestionCollectorFactory(
		[]receiver.Factory{
			rmqreceiver.NewFactory(mqFactory),
		},
		[]processor.Factory{
			queueprocessor.NewFactory(),
		},
		[]exporter.Factory{
			clickhouseexporter.NewFactory(traceRepo),
		},
	)
}

func NewTraceConfigLoader(confFactory conf.IConfigLoaderFactory) (conf.IConfigLoader, error) {
	return confFactory.NewConfigLoader("observability.yaml")
}

func InitTraceApplication(
	db db.Provider,
	mqFactory mq.IFactory,
	configFactory conf.IConfigLoaderFactory,
	fileClient fileservice.Client,
	ckDb ck.Provider,
	benefit benefit.IBenefitService,
	authClient authservice.Client,
	meter metrics.Meter,
) (ITraceApplication, error) {
	wire.Build(traceSet)
	return nil, nil
}

func InitTraceIngestionApplication(
	configFactory conf.IConfigLoaderFactory,
	ckDb ck.Provider,
	mqFactory mq.IFactory) (ITraceIngestionApplication, error) {
	wire.Build(traceIngestionSet)
	return nil, nil
}
