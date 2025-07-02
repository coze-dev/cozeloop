// Original Files: open-telemetry/opentelemetry-collector
// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// This file may have been modified by ByteDance Ltd.

package queueprocessor

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alitto/pond/v2"
	"golang.org/x/sync/errgroup"

	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/component"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/consumer"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity/collector/processor"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/goroutine"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

func createDefaultConfig() component.Config {
	return &Config{}
}

func createTracesProcessor(ctx context.Context, set processor.CreateSettings, cfg component.Config, nextConsumer consumer.Consumer) (processor.Processor, error) {
	config := cfg.(*Config)
	logs.CtxInfo(ctx, "queue processor config: %v", *config)
	q := &queueProcessor{
		nextConsumer: nextConsumer,
		config:       config,
		goPool:       pond.NewPool(int(config.MaxPoolSize)),
	}
	shards := make([]*shard, 0, config.ShardCount)
	for i := 0; i < config.ShardCount; i++ {
		shards = append(shards, &shard{
			index:      i,
			shutdownCh: make(chan struct{}),
			doneCh:     make(chan struct{}),
			queue:      make(chan *consumer.Traces, config.QueueSize),
			tracesData: make(map[string]*consumer.Traces),
			processor:  q,
		})
	}
	q.batch = &batch{
		shards: shards,
	}
	return q, nil
}

type queueProcessor struct {
	nextConsumer consumer.Consumer
	config       *Config
	goPool       pond.Pool
	batch        *batch
}

func (q *queueProcessor) Start(ctx context.Context) error {
	logs.Info("queue processor start")
	for _, tmp := range q.batch.shards {
		s := tmp
		goroutine.Go(ctx, s.startWorker)
	}
	return nil
}

func (q *queueProcessor) Shutdown(ctx context.Context) error {
	logs.Info("queue processor shutting down")
	err := q.batch.shutdown(ctx)
	if err != nil {
		logs.CtxError(ctx, "shutdown batch failed, %v", err)
		return err
	}
	q.goPool.StopAndWait()
	logs.Info("queue processor shutted down")
	return nil
}

func (q *queueProcessor) ConsumeTraces(ctx context.Context, td consumer.Traces) error {
	logs.Debug("queue processor handle trace data: %+v", td)
	sd := q.batch.getShard()
	if sd == nil {
		logs.CtxError(ctx, "queue processor receive req without shard to process")
		return nil
	}
	sd.receive(ctx, td)
	return nil
}

func (q *queueProcessor) process(ctx context.Context, td *consumer.Traces) func() {
	return func() {
		defer goroutine.Recovery(ctx)
		err := q.nextConsumer.ConsumeTraces(ctx, *td)
		if err != nil {
			logs.CtxError(ctx, "next consumer consume trace failed, %v", err)
		}
	}
}

func (b *batch) getShard() *shard {
	i := int(atomic.AddUint64(&b.current, uint64(1)) % uint64(len(b.shards)))
	if i < len(b.shards) {
		return b.shards[i]
	}
	return nil
}

type batch struct {
	current uint64
	shards  []*shard
}

type shard struct {
	index      int
	spansCount atomic.Uint64
	mutex      sync.Mutex
	tracesData map[string]*consumer.Traces
	shutdownCh chan struct{}
	queue      chan *consumer.Traces
	processor  *queueProcessor
	doneCh     chan struct{}
}

func (b *shard) receive(ctx context.Context, td consumer.Traces) {
	b.queue <- &td
}

// 预期只有一个Tenant, 不同租户最佳实践下隔离使用
// 统计Spans的数量，而不是Traces
func (b *shard) append(td *consumer.Traces) {
	if td == nil {
		return
	}
	traces, ok := b.tracesData[td.Tenant]
	if ok {
		traces.TraceData = append(traces.TraceData, td.TraceData...)
	} else {
		b.tracesData[td.Tenant] = td
	}
	b.spansCount.Add(uint64(td.SpansCount()))
}

func (b *shard) count() int {
	return int(b.spansCount.Load())
}

func (b *shard) startWorker() {
	ticker := time.NewTicker(time.Duration(b.processor.config.TickIntervalsMs) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-b.shutdownCh:
			logs.Info("shard %d receive shutdown signal", b.index)
			for td := range b.queue {
				b.append(td)
			}
			b.doneCh <- struct{}{}
			return
		case traceData := <-b.queue:
			if traceData == nil {
				continue
			}
			b.append(traceData)
			if b.count() >= b.processor.config.MaxBatchSize {
				logs.Debug("shard %d spans count %d exceeds limit %d, start flushing",
					b.index, b.count(), b.processor.config.MaxBatchSize)
				b.flush()
			}
		case <-ticker.C:
			logs.Debug("shard %d receive tick signal, count %d", b.index, b.count())
			if b.count() > 0 {
				b.flush()
			}
		}
	}
}

func (b *shard) flush() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for _, tmp := range b.tracesData {
		td := tmp
		ctx := context.Background()
		// no context needed yet
		if err := b.processor.goPool.Go(b.processor.process(ctx, td)); err != nil {
			// not supposed to be here
			logs.Error("fail to submit task, %v", err)
		}
	}
	b.tracesData = make(map[string]*consumer.Traces)
	b.spansCount.Store(0)
}

func (b *batch) shutdown(ctx context.Context) error {
	group := errgroup.Group{}
	for i, s := range b.shards {
		ti := i
		ts := s
		group.Go(func() error {
			defer goroutine.Recovery(ctx)
			close(ts.shutdownCh)
			close(ts.queue)
			<-ts.doneCh
			logs.CtxInfo(ctx, "shutdown batch %d", ti)
			ts.flush()
			return nil
		})
	}
	err := group.Wait()
	if err != nil {
		return fmt.Errorf("shutdown batch failed err=%+v", err)
	}
	return nil
}

//func init() {
//	logs.SetLogLevel(logs.DebugLevel)
//}
