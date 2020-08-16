package prometheus

import (
	"testing"
	"time"

	"github.com/m3o/services/metrics"
	"github.com/micro/micro/v3/service"

	"github.com/stretchr/testify/assert"
)

func TestPrometheusReporter(t *testing.T) {

	// New Service:
	srv := service.New(
		service.Name("test-metrics"),
		service.Version("1.2.3"),
	)

	// Make a Reporter:
	reporter, err := New(srv)
	assert.NoError(t, err)
	assert.NotNil(t, reporter)

	// Test tag conversion:
	tags := metrics.Tags{
		"tag1": "false",
		"tag2": "true",
	}
	convertedTags := reporter.convertTags(tags)
	assert.Equal(t, "false", convertedTags["tag1"])
	assert.Equal(t, "true", convertedTags["tag2"])

	// Test tag enumeration:
	listedTags := reporter.listTagKeys(tags)
	assert.Contains(t, listedTags, "tag1")
	assert.Contains(t, listedTags, "tag2")

	// Test string cleaning:
	preparedMetricName := reporter.stripUnsupportedCharacters("some.kind,of tag")
	assert.Equal(t, "some_kind_oftag", preparedMetricName)

	// Test MetricFamilies:
	metricFamily := reporter.newMetricFamily()

	// Counters:
	assert.NotNil(t, metricFamily.getCounter("testCounter", []string{"test", "counter"}))
	assert.Len(t, metricFamily.counters, 1)

	// Gauges:
	assert.NotNil(t, metricFamily.getGauge("testGauge", []string{"test", "gauge"}))
	assert.Len(t, metricFamily.gauges, 1)

	// Timings:
	assert.NotNil(t, metricFamily.getTiming("testTiming", []string{"test", "timing"}))
	assert.Len(t, metricFamily.timings, 1)

	// Test submitting metrics through the interface methods:
	assert.NoError(t, reporter.Count("test.counter.1", 6, tags))
	assert.NoError(t, reporter.Count("test.counter.2", 19, tags))
	assert.NoError(t, reporter.Count("test.counter.1", 5, tags))
	assert.NoError(t, reporter.Gauge("test.gauge.1", 99, tags))
	assert.NoError(t, reporter.Gauge("test.gauge.2", 55, tags))
	assert.NoError(t, reporter.Gauge("test.gauge.1", 98, tags))
	assert.NoError(t, reporter.Timing("test.timing.1", time.Second, tags))
	assert.NoError(t, reporter.Timing("test.timing.2", time.Minute, tags))
	assert.Len(t, reporter.metrics.counters, 2)
	assert.Len(t, reporter.metrics.gauges, 2)
	assert.Len(t, reporter.metrics.timings, 2)

	// Test reading back the metrics:
	// This could be done by hitting the /metrics endpoint
}
