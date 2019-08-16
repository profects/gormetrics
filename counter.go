package gormetrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var globalQueryCounters = map[string]queryCounters{}
var globalDatabaseCounters = map[string]databaseCounters{}
var gqcMtx sync.Mutex
var gdcMtx sync.Mutex

// queryCounters contains all counters that are exported.
type queryCounters struct {
	all     *prometheus.CounterVec
	creates *prometheus.CounterVec
	deletes *prometheus.CounterVec
	queries *prometheus.CounterVec
	updates *prometheus.CounterVec
}

// newQueryCounters creates a new queryCounters instance with all counters valid.
func newQueryCounters(namespace string) (queryCounters, error) {
	gqcMtx.Lock()
	defer gqcMtx.Unlock()

	if gc, exists := globalQueryCounters[namespace]; exists {
		return gc, nil
	}

	cc := counterVecCreator{
		namespace: namespace,
		labels: []string{
			labelDatabase,
			labelDriver,
			labelStatus,
		},
	}

	qc := queryCounters{
		all:     cc.new(metricAllTotal, helpAllTotal),
		creates: cc.new(metricCreatesTotal, helpCreatesTotal),
		deletes: cc.new(metricDeletesTotal, helpDeletesTotal),
		queries: cc.new(metricQueriesTotal, helpQueriesTotal),
		updates: cc.new(metricUpdatesTotal, helpUpdatesTotal),
	}
	globalQueryCounters[namespace] = qc

	err := registerCollectors(
		qc.all,
		qc.creates,
		qc.deletes,
		qc.queries,
		qc.updates,
	)

	return qc, err
}

type databaseCounters struct {
	idle  *prometheus.GaugeVec
	inUse *prometheus.GaugeVec
	open  *prometheus.GaugeVec
}

func newDatabaseCounters(namespace string) (databaseCounters, error) {
	gdcMtx.Lock()
	defer gdcMtx.Unlock()

	if gc, exists := globalDatabaseCounters[namespace]; exists {
		return gc, nil
	}

	vecCreator := gaugeVecCreator{
		namespace: namespace,
		labels: []string{
			labelDatabase,
			labelDriver,
		},
	}

	dc := databaseCounters{
		idle:  vecCreator.new(metricIdleConnections, helpIdleConnections),
		inUse: vecCreator.new(metricInUseConnections, helpInUseConnections),
		open:  vecCreator.new(metricOpenConnections, helpOpenConnections),
	}
	globalDatabaseCounters[namespace] = dc

	err := registerCollectors(
		dc.idle,
		dc.inUse,
		dc.open,
	)

	return dc, err
}

// registerCollectors registers multiple instances of prometheus.Collector.
func registerCollectors(collectors ...prometheus.Collector) error {
	for _, c := range collectors {
		err := prometheus.Register(c)
		if err != nil {
			return err
		}
	}

	return nil
}
