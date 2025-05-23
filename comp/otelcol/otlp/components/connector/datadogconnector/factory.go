// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:generate mdatagen metadata.yaml

package datadogconnector // import "github.com/DataDog/datadog-agent/comp/otelcol/otlp/components/connector/datadogconnector"

import (
	"context"
	"github.com/DataDog/datadog-agent/comp/core/tagger/types"
	"github.com/DataDog/datadog-agent/pkg/util/option"
	"time"

	"github.com/DataDog/datadog-agent/comp/otelcol/otlp/components/metricsclient"
	datadogconfig "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/datadog/config"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
)

type factory struct {
	tagger   types.TaggerClient
	hostname option.Option[string]
}

// SourceProviderFunc is a function that returns the source of the host.
type SourceProviderFunc func(context.Context) (string, error)

// NewFactoryForAgent creates a factory for datadog connector for use in OTel agent
func NewFactoryForAgent(tagger types.TaggerClient, hostGetter SourceProviderFunc) connector.Factory {
	f := &factory{
		tagger: tagger,
	}

	if hostGetter != nil {
		if hostname, err := hostGetter(context.Background()); err == nil {
			f.hostname = option.New(hostname)
		}
	}

	//  OTel connector factory to make a factory for connectors
	return connector.NewFactory(
		Type,
		createDefaultConfig,
		connector.WithTracesToMetrics(f.createTracesToMetricsConnector, TracesToMetricsStability),
		connector.WithTracesToTraces(f.createTracesToTracesConnector, TracesToTracesStability))
}

// NewFactory creates a factory for datadog connector.
func NewFactory() connector.Factory {
	//  OTel connector factory to make a factory for connectors
	return NewFactoryForAgent(nil, nil)
}

func createDefaultConfig() component.Config {
	return &Config{
		Traces: datadogconfig.TracesConnectorConfig{
			TracesConfig: datadogconfig.TracesConfig{
				IgnoreResources:        []string{},
				PeerServiceAggregation: true,
				PeerTagsAggregation:    true,
				ComputeStatsBySpanKind: true,
			},

			TraceBuffer:    1000,
			BucketInterval: 10 * time.Second,
		},
	}
}

// defines the consumer type of the connector
// we want to consume traces and export metrics therefore define nextConsumer as metrics, consumer is the next component in the pipeline
func (f *factory) createTracesToMetricsConnector(_ context.Context, params connector.Settings, cfg component.Config, nextConsumer consumer.Metrics) (c connector.Traces, err error) {
	metricsClient, err := metricsclient.InitializeMetricClient(params.MeterProvider, metricsclient.ConnectorSourceTag)
	if err != nil {
		return nil, err
	}

	c, err = newTraceToMetricConnector(params.TelemetrySettings, cfg, nextConsumer, metricsClient, f.tagger, f.hostname)

	if err != nil {
		return nil, err
	}
	return c, nil
}

func (f *factory) createTracesToTracesConnector(_ context.Context, params connector.Settings, _ component.Config, nextConsumer consumer.Traces) (connector.Traces, error) {
	return newTraceToTraceConnector(params.Logger, nextConsumer), nil
}
