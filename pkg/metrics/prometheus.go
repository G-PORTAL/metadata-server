package metrics

import "github.com/prometheus/client_golang/prometheus"

var MetadataRequests = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "metadata_requests",
		Help: "Number of metadata API requests",
	}, []string{"url", "source"},
)

func init() {
	prometheus.MustRegister(MetadataRequests)
}
