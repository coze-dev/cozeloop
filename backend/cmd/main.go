// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/coze-dev/cozeloop-go"
	goredis "github.com/redis/go-redis/v9"

	"github.com/coze-dev/cozeloop/backend/api"
	"github.com/coze-dev/cozeloop/backend/api/handler/coze/loop/apis"
	"github.com/coze-dev/cozeloop/backend/infra/ck"
	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/infra/external/audit"
	"github.com/coze-dev/cozeloop/backend/infra/external/benefit"
	"github.com/coze-dev/cozeloop/backend/infra/fileserver"
	"github.com/coze-dev/cozeloop/backend/infra/i18n"
	"github.com/coze-dev/cozeloop/backend/infra/i18n/goi18n"
	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	"github.com/coze-dev/cozeloop/backend/infra/idgen/redis_gen"
	"github.com/coze-dev/cozeloop/backend/infra/limiter"
	"github.com/coze-dev/cozeloop/backend/infra/limiter/dist"
	"github.com/coze-dev/cozeloop/backend/infra/looptracer"
	"github.com/coze-dev/cozeloop/backend/infra/looptracer/rpc"
	"github.com/coze-dev/cozeloop/backend/infra/metrics"
	"github.com/coze-dev/cozeloop/backend/infra/mq"
	"github.com/coze-dev/cozeloop/backend/infra/mq/registry"
	"github.com/coze-dev/cozeloop/backend/infra/mq/rocketmq"
	"github.com/coze-dev/cozeloop/backend/infra/redis"
	"github.com/coze-dev/cozeloop/backend/loop_gen/coze/loop/foundation/lofile"
	"github.com/coze-dev/cozeloop/backend/loop_gen/coze/loop/observability/lotrace"
	"github.com/coze-dev/cozeloop/backend/pkg/conf"
	"github.com/coze-dev/cozeloop/backend/pkg/conf/viper"
	"github.com/coze-dev/cozeloop/backend/pkg/file"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

func main() {
	ctx := context.Background()
	c, err := newComponent(ctx)
	if err != nil {
		panic(err)
	}

	handler, err := api.Init(ctx, c.idgen, c.db, c.redis, c.cfgFactory, c.mqFactory, c.objectStorage, c.batchObjectStorage, c.benefitSvc, c.auditClient, c.metric, c.limiterFactory, c.ckDb, c.translater)
	if err != nil {
		panic(err)
	}

	if err := initTracer(handler); err != nil {
		panic(err)
	}

	if err := registry.NewConsumerRegistry(c.mqFactory).Register(mustInitConsumerWorkers(c.cfgFactory, handler, handler)).StartAll(ctx); err != nil {
		panic(err)
	}

	api.Start(handler)
}

type ComponentConfig struct {
	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
	} `mapstructure:"redis"`
	RDS struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		DB       string `mapstructure:"db"`
	} `mapstructure:"rds"`
	S3Config struct {
		Region          string `mapstructure:"region"`
		Endpoint        string `mapstructure:"endpoint"`
		Bucket          string `mapstructure:"bucket"`
		AccessKey       string `mapstructure:"access_key"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
	} `mapstructure:"s3_config"`
	CKConfig struct {
		Host        string `mapstructure:"host"`
		Database    string `mapstructure:"database"`
		UserName    string `mapstructure:"username"`
		Password    string `mapstructure:"password"`
		DialTimeout int    `mapstructure:"dial_timeout"`
		ReadTimeout int    `mapstructure:"read_timeout"`
	} `mapstructure:"ck_config"`
	IDGen struct {
		ServerIDs []int64 `mapstructure:"server_ids"`
	} `mapstructure:"idgen"`
	LogLevel string `mapstructure:"log_level"`
}

func getComponentConfig(configFactory conf.IConfigLoaderFactory) (*ComponentConfig, error) {
	ctx := context.Background()
	componentConfigLoader, err := configFactory.NewConfigLoader("infrastructure.yaml")
	if err != nil {
		return nil, err
	}
	componentConfig := &ComponentConfig{}
	err = componentConfigLoader.UnmarshalKey(ctx, "infra", componentConfig)
	if err != nil {
		return nil, err
	}
	return componentConfig, nil
}

type component struct {
	idgen              idgen.IIDGenerator
	db                 db.Provider
	redis              redis.Cmdable
	cfgFactory         conf.IConfigLoaderFactory
	mqFactory          mq.IFactory
	objectStorage      fileserver.ObjectStorage
	batchObjectStorage fileserver.BatchObjectStorage
	benefitSvc         benefit.IBenefitService
	auditClient        audit.IAuditService
	metric             metrics.Meter
	limiterFactory     limiter.IRateLimiterFactory
	ckDb               ck.Provider
	translater         i18n.ITranslater
}

func initTracer(handler *apis.APIHandler) error {
	rpc.SetLoopTracerHandler(
		lofile.NewLocalFileService(handler.FileService),
		lotrace.NewLocalTraceService(handler.ITraceApplication),
	)

	client, err := cozeloop.NewClient(
		cozeloop.WithWorkspaceID("0"),
		cozeloop.WithAPIToken("0"),
		cozeloop.WithExporter(&looptracer.MultiSpaceSpanExporter{}),
	)
	if err != nil {
		return err
	}
	looptracer.InitTracer(looptracer.NewTracer(client))

	return nil
}

func newComponent(ctx context.Context) (*component, error) {
	c := new(component)
	cfgFactory := viper.NewFileConfigLoaderFactory(viper.WithFactoryConfigPath("conf"))
	componentConfig, err := getComponentConfig(cfgFactory)
	if err != nil {
		return c, err
	}
	switch componentConfig.LogLevel {
	case "debug":
		logs.SetLogLevel(logs.DebugLevel)
	case "info":
		logs.SetLogLevel(logs.InfoLevel)
	case "warn":
		logs.SetLogLevel(logs.WarnLevel)
	case "error":
		logs.SetLogLevel(logs.ErrorLevel)
	case "fatal":
		logs.SetLogLevel(logs.FatalLevel)
	}
	cmdable, err := redis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%d", componentConfig.Redis.Host, componentConfig.Redis.Port),
		Password: componentConfig.Redis.Password,
	})
	if err != nil {
		return nil, err
	}

	redisCli, ok := redis.Unwrap(cmdable)
	if !ok {
		return c, errors.New("unwrap redis cli fail")
	}

	db, err := db.NewDBFromConfig(&db.Config{
		User:         componentConfig.RDS.User,
		Password:     componentConfig.RDS.Password,
		DBHostname:   componentConfig.RDS.Host,
		DBPort:       componentConfig.RDS.Port,
		DBName:       componentConfig.RDS.DB,
		Loc:          "Local",
		DBCharset:    "utf8mb4",
		Timeout:      time.Minute,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		DSNParams:    url.Values{"clientFoundRows": []string{"true"}},
	})
	if err != nil {
		return nil, err
	}

	s3Config := fileserver.NewS3Config(func(cfg *fileserver.S3Config) {
		cfg.Endpoint = componentConfig.S3Config.Endpoint
		cfg.Region = componentConfig.S3Config.Region
		cfg.Bucket = componentConfig.S3Config.Bucket
		cfg.AccessKeyID = componentConfig.S3Config.AccessKey
		cfg.SecretAccessKey = componentConfig.S3Config.SecretAccessKey
	})
	objectStorage, err := fileserver.NewS3Client(s3Config)
	if err != nil {
		return nil, err
	}

	ckDb, err := ck.NewCKFromConfig(&ck.Config{
		Host:              componentConfig.CKConfig.Host,
		Database:          componentConfig.CKConfig.Database,
		Username:          componentConfig.CKConfig.UserName,
		Password:          componentConfig.CKConfig.Password,
		CompressionMethod: ck.CompressionMethodZSTD,
		CompressionLevel:  3,
		Protocol:          ck.ProtocolNative,
		DialTimeout:       time.Duration(componentConfig.CKConfig.DialTimeout) * time.Second,
		ReadTimeout:       time.Duration(componentConfig.CKConfig.ReadTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}

	idgenerator, err := redis_gen.NewIDGenerator(redisCli, componentConfig.IDGen.ServerIDs)
	if err != nil {
		return nil, err
	}

	localeDir, err := file.FindSubDir(os.Getenv("PWD"), "runtime/locales")
	if err != nil {
		return nil, err
	}
	translater, err := goi18n.NewTranslater(localeDir)
	if err != nil {
		return nil, err
	}

	return &component{
		idgen:              idgenerator,
		db:                 db,
		redis:              cmdable,
		cfgFactory:         cfgFactory,
		mqFactory:          rocketmq.NewFactory(),
		objectStorage:      objectStorage,
		batchObjectStorage: objectStorage,
		benefitSvc:         benefit.NewNoopBenefitService(),
		auditClient:        audit.NewNoopAuditService(),
		metric:             metrics.GetMeter(),
		limiterFactory:     dist.NewRateLimiterFactory(cmdable),
		ckDb:               ckDb,
		translater:         translater,
	}, nil
}
