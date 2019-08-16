package gormetrics

import "github.com/prometheus/client_golang/prometheus"

// counterVecCreator allows for mass creation of counter vectors in the same Prometheus namespace
// and with equal constant labels.
type counterVecCreator struct {
	namespace string
	labels    []string
}

// new creates a new prometheus.CounterVec based on the specified name and values in the counterVecCreator.
func (c counterVecCreator) new(
	name string,
	help string,
) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: c.namespace,
			Name:      name,
			Help:      help,
		},
		c.labels,
	)
}

// gaugeVecCreator allows for mass creation of gauge vectors in the same Prometheus namespace
// and with equal constant labels.
type gaugeVecCreator struct {
	namespace string
	labels    []string
}

// new creates a new prometheus.Counter based on the specified name and values in the counterCreator.
func (c gaugeVecCreator) new(
	name string,
	help string,
) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: c.namespace,
			Name:      name,
			Help:      help,
		},
		c.labels,
	)
}
