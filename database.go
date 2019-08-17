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

	db *sql.DB
	sync.Mutex
}

// newDatabase creates a new database wrapper containing the name of the database,
// it's driver and the (sql) database itself.
func newDatabase(info extraInfo, db *sql.DB) *database {
	return &database{
		name:       info.dbName,
		driverName: info.driverName,
		db:         db,
	}
}

// collectConnectionStats collects database connections for Prometheus to scrape.
func (d *database) collectConnectionStats(counters *databaseGauges) {
	d.Lock()
	defer d.Unlock()

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

// databaseMetrics is a convenience struct for exporting database metrics to Prometheus.
type databaseMetrics struct {
	gauges *databaseGauges
	db     *database
}

// newDatabaseMetrics creates a new databaseMetrics instance with a database backing it
// for statistics. Use maintain to continuously collect statistics.
func newDatabaseMetrics(db *database, opts *pluginOpts) (*databaseMetrics, error) {
	gauges, err := newDatabaseGauges(opts.prometheusNamespace)
	if err != nil {
		return nil, errors.Wrap(err, "could not create database gauges")
	}

	return &databaseMetrics{
		gauges: gauges,
		db:     db,
	}, nil
}

// maintain collects connection statistics every 3 seconds and.
func (d *databaseMetrics) maintain() {
	ticker := time.NewTicker(time.Second * 3)

	for range ticker.C {
		d.db.collectConnectionStats(d.gauges)
	}
}
