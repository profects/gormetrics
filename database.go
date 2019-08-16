package gormetrics

import (
	"database/sql"
	"sync"
	"time"

	"github.com/pkg/errors"
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

func (d *database) collectConnectionStats(counters *databaseCounters) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	defaultLabels := prometheus.Labels{
		labelDatabase: d.name,
		labelDriver:   d.driverName,
	}

	stats := d.db.Stats()

	counters.idle.
		With(defaultLabels).
		Set(float64(stats.Idle))

	counters.inUse.
		With(defaultLabels).
		Set(float64(stats.InUse))

	counters.open.
		With(defaultLabels).
		Set(float64(stats.OpenConnections))
}

type databaseMetrics struct {
	counters *databaseCounters
	db       *database
}

func newDatabaseMetrics(db *database, opts *pluginOpts) (*databaseMetrics, error) {
	counters, err := newDatabaseCounters(opts.prometheusNamespace)
	if err != nil {
		return nil, errors.Wrap(err, "could not create database counters")
	}

	return &databaseMetrics{
		counters: counters,
		db:       db,
	}, nil
}

func (d *databaseMetrics) maintain() {
	ticker := time.NewTicker(time.Second * 3)

	// Collect connection statistics every 3 seconds.
	for range ticker.C {
		d.db.collectConnectionStats(d.counters)
	}
}
