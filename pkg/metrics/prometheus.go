package metrics

import "github.com/prometheus/client_golang/prometheus"

var MetadataRequests = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "metadata_requests",
		Help: "No of request handled by Ping handler",
	}, []string{"url", "source"},
)

func init() {
	prometheus.MustRegister(MetadataRequests)
}
