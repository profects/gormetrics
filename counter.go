package gormetrics

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

var globalQueryCounters = map[string]*queryCounters{}
var globalDatabaseCounters = map[string]*databaseCounters{}
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
func newQueryCounters(namespace string) (*queryCounters, error) {
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

	qc := &queryCounters{
		all:     cc.new(metricAllTotal, helpAllTotal),
		creates: cc.new(metricCreatesTotal, helpCreatesTotal),
		deletes: cc.new(metricDeletesTotal, helpDeletesTotal),
		queries: cc.new(metricQueriesTotal, helpQueriesTotal),
		updates: cc.new(metricUpdatesTotal, helpUpdatesTotal),
	}

	if err := registerCollectors(
		qc.all,
		qc.creates,
		qc.deletes,
		qc.queries,
		qc.updates,
	); err != nil {
		return nil, errors.Wrap(err, "could not register collectors")
	}

	globalQueryCounters[namespace] = qc

	return qc, nil
}

type databaseCounters struct {
	idle  *prometheus.GaugeVec
	inUse *prometheus.GaugeVec
	open  *prometheus.GaugeVec
}

func newDatabaseCounters(namespace string) (*databaseCounters, error) {
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

	dc := &databaseCounters{
		idle:  vecCreator.new(metricIdleConnections, helpIdleConnections),
		inUse: vecCreator.new(metricInUseConnections, helpInUseConnections),
		open:  vecCreator.new(metricOpenConnections, helpOpenConnections),
	}

	if err := registerCollectors(
		dc.idle,
		dc.inUse,
		dc.open,
	); err != nil {
		return nil, err
	}

	globalDatabaseCounters[namespace] = dc

	return dc, nil
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
