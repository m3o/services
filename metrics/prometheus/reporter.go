package prometheus

import (
	"net/http"
	"strings"

	"github.com/m3o/services/metrics"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: These sensible defaults should really become config options
var (
	// The Prometheus metrics will be made available on this port:
	listenAddress = ":9000"
	// This is the endpoint where the Prometheus metrics will be made available ("/metrics" is the default with Prometheus):
	metricsEndpoint = "/metrics"
	// timingObjectives is the default spread of stats we maintain for timings / histograms:
	timingObjectives = map[float64]float64{0.0: 0, 0.5: 0.05, 0.75: 0.04, 0.90: 0.03, 0.95: 0.02, 0.98: 0.001, 1: 0}
)

// Reporter is an implementation of metrics.Reporter:
type Reporter struct {
	defaultLabels      prometheus.Labels
	prometheusRegistry *prometheus.Registry
	metrics            metricFamily
}

// New returns a configured prometheus reporter:
func New(service *service.Service) (*Reporter, error) {

	// Make a prometheus registry (this keeps track of any metrics we generate):
	prometheusRegistry := prometheus.NewRegistry()
	prometheusRegistry.Register(prometheus.NewGoCollector())
	prometheusRegistry.Register(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{Namespace: "goruntime"}))

	// Make a new Reporter:
	newReporter := &Reporter{
		prometheusRegistry: prometheusRegistry,
	}

	// Prepare some default tags which will be included with every metric (based on the service metadata):
	if service != nil {
		newReporter.defaultLabels = prometheus.Labels{
			"service_name":    service.Name(),
			"service_version": service.Version(),
		}
	}

	// Add metrics families for each type:
	newReporter.metrics = newReporter.newMetricFamily()

	// Handle the metrics endpoint with prometheus:
	log.Infof("Metrics/Prometheus [http] Listening on %s%s", listenAddress, metricsEndpoint)
	http.Handle(metricsEndpoint, promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}))
	go http.ListenAndServe(listenAddress, nil)

	return newReporter, nil
}

// convertTags turns Tags into prometheus labels:
func (r *Reporter) convertTags(tags metrics.Tags) prometheus.Labels {
	labels := prometheus.Labels{}
	for key, value := range tags {
		labels[key] = r.stripUnsupportedCharacters(value)
	}
	return labels
}

// listTagKeys returns a list of tag keys (we need to provide this to the Prometheus client):
func (r *Reporter) listTagKeys(tags metrics.Tags) (labelKeys []string) {
	for key := range tags {
		labelKeys = append(labelKeys, key)
	}
	return
}

// stripUnsupportedCharacters cleans up a metrics key or value:
func (r *Reporter) stripUnsupportedCharacters(metricName string) string {
	valueWithoutDots := strings.Replace(metricName, ".", "_", -1)
	valueWithoutCommas := strings.Replace(valueWithoutDots, ",", "_", -1)
	valueWIthoutSpaces := strings.Replace(valueWithoutCommas, " ", "", -1)
	return valueWIthoutSpaces
}
