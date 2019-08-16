package gormetrics

import (
	"database/sql"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type database struct {
	name       string
	driverName string

	db  *sql.DB
	mtx sync.Mutex
}

func databaseFrom(info extraInfo, db *sql.DB) *database {
	return &database{
		name:       info.dbName,
		driverName: info.driverName,
		db:         db,
	}
}

func (d *database) collectConnectionStats(counters databaseCounters) {
	d.mtx.Lock()

	stats := d.db.Stats()

	defaultLabels := prometheus.Labels{
		labelDatabase: d.name,
		labelDriver:   d.driverName,
	}

	counters.idle.
		With(defaultLabels).
		Set(float64(stats.Idle))

	counters.inUse.
		With(defaultLabels).
		Set(float64(stats.InUse))

	counters.open.
		With(defaultLabels).
		Set(float64(stats.OpenConnections))

	d.mtx.Unlock()
}

type databaseMetrics struct {
	counters databaseCounters
	db       *database
}

func newDatabaseMetrics(db *database, opts *pluginOpts) (*databaseMetrics, error) {
	counters, err := newDatabaseCounters(opts.prometheusNamespace)
	if err != nil {
		return nil, err
	}

	return &databaseMetrics{
		counters: counters,
		db:       db,
	}, nil
}

func (d *databaseMetrics) maintain() {
	ticker := time.NewTicker(time.Second * 3)

	for range ticker.C {
		d.db.collectConnectionStats(d.counters)
	}
}
