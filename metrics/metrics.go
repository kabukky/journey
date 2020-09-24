package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	JourneyHandler = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "journey_handler",
			Help: "Number of visits for the handler.",
		},
		[]string{"handler"},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(JourneyHandler)
}
