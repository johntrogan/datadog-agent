// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package server

import (
	workloadmeta "github.com/DataDog/datadog-agent/comp/core/workloadmeta/def"
	"github.com/DataDog/datadog-agent/comp/dogstatsd/packets"
	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/DataDog/datadog-agent/pkg/metrics"
	"github.com/DataDog/datadog-agent/pkg/util/option"
	utilstrings "github.com/DataDog/datadog-agent/pkg/util/strings"
)

var (
	// defaultSampleSize is the default allocation size used to store
	// samples when a message contains multiple values. This will
	// automatically be extended if needed when we append to it.
	defaultSampleSize = 1024
)

type worker struct {
	server *server
	// the batcher will be responsible of batching a few samples / events / service
	// checks and it will automatically forward them to the aggregator, meaning that
	// the flushing logic to the aggregator is actually in the batcher.
	batcher *batcher
	parser  *parser

	// we allocate it once per worker instead of once per packet. This will
	// be used to store the samples out a of packets. Allocating it every
	// time is very costly, especially on the GC.
	samples metrics.MetricSampleBatch

	packetsTelemetry *packets.TelemetryStore

	BlocklistUpdate chan utilstrings.Blocklist
	blocklist       utilstrings.Blocklist
}

func newWorker(s *server, workerNum int, wmeta option.Option[workloadmeta.Component], packetsTelemetry *packets.TelemetryStore, stringInternerTelemetry *stringInternerTelemetry) *worker {
	var batcher *batcher
	if s.ServerlessMode {
		batcher = newServerlessBatcher(s.demultiplexer, s.tlmChannel)
	} else {
		batcher = newBatcher(s.demultiplexer.(aggregator.DemultiplexerWithAggregator), s.tlmChannel)
	}

	return &worker{
		server:           s,
		batcher:          batcher,
		parser:           newParser(s.config, s.sharedFloat64List, workerNum, wmeta, stringInternerTelemetry),
		samples:          make(metrics.MetricSampleBatch, 0, defaultSampleSize),
		packetsTelemetry: packetsTelemetry,
		BlocklistUpdate:  make(chan utilstrings.Blocklist),
	}
}

func (w *worker) run() {
	for {
		select {
		case <-w.server.stopChan:
			return
		case <-w.server.health.C:
		case <-w.server.serverlessFlushChan:
			w.batcher.flush()
		case blocklist := <-w.BlocklistUpdate:
			w.blocklist = blocklist
		case ps := <-w.server.packetsIn:
			w.packetsTelemetry.TelemetryUntrackPackets(ps)
			w.samples = w.samples[0:0]
			// we return the samples in case the slice was extended
			// when parsing the packets
			w.samples = w.server.parsePackets(w.batcher, w.parser, ps, w.samples, &w.blocklist)
		}

	}
}
