package metrics

import "github.com/prometheus/client_golang/prometheus"

var MetadataRequests = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "metadata_requests",
	}, []string{"url", "source"},
)

func init() {
	prometheus.MustRegister(MetadataRequests)
}
