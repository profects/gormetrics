package gormetrics

import (
	"fmt"

	"github.com/profects/gormetrics/gormi"

	"github.com/prometheus/client_golang/prometheus"
)

type callbackHandler struct {
	opts          pluginOpts
	counters      queryCounters
	defaultLabels map[string]string
}

func (h *callbackHandler) registerCallback(cb gormi.Callback) {
	cb.Create().After("gorm:create").Register(
		h.opts.callbackName("after_create"),
		h.afterCreate,
	)

	cb.Delete().After("gorm:delete").Register(
		h.opts.callbackName("after_delete"),
		h.afterDelete,
	)

	cb.Query().After("gorm:query").Register(
		h.opts.callbackName("after_query"),
		h.afterQuery,
	)

	cb.Update().After("gorm:update").Register(
		h.opts.callbackName("after_update"),
		h.afterUpdate,
	)
}

func (h *callbackHandler) afterCreate(scope gormi.Scope) {
	h.addToAfterScope(scope, h.counters.creates)
}

func (h *callbackHandler) afterDelete(scope gormi.Scope) {
	h.addToAfterScope(scope, h.counters.deletes)
}

func (h *callbackHandler) afterQuery(scope gormi.Scope) {
	h.addToAfterScope(scope, h.counters.queries)
}

func (h *callbackHandler) afterUpdate(scope gormi.Scope) {
	h.addToAfterScope(scope, h.counters.updates)
}

// addToAfterScope registers one or more of prometheus.CounterVec to increment
// after a scope (any type of query). If any errors are in
// scope.DB().GetErrors(), a status "fail" will be assigned to the increment.
// Otherwise, a status "success" will be assigned.
// Increments h.counters.all (gormetrics_all_total) by default.
func (h *callbackHandler) addToAfterScope(scope gormi.Scope, vectors ...*prometheus.CounterVec) {
	vectors = append(vectors, h.counters.all)

	hasErrors := len(scope.DB().GetErrors()) > 0
	status := metricStatusFail
	if !hasErrors {
		status = metricStatusSuccess
	}

	for _, counter := range vectors {
		if counter == nil {
			continue
		}
		counter.With(
			mergeLabels(prometheus.Labels{
				labelStatus: status,
			}, h.defaultLabels),
		).Inc()
	}
}

// extraInfo contains information for filtering the provided metrics.
type extraInfo struct {
	// The name of the database in use.
	dbName string

	// The name of the driver powering database/sql (underlying database for GORM).
	driverName string
}

// newCallbackHandler creates a new callback handler configured with info and opts.
// info does not contain any mandatory information for the functioning of the
// function, but sets label values which can be useful in the usage of
// the provided metrics (driver, database, connection).
// Automatically registers metrics.
func newCallbackHandler(info extraInfo, opts *pluginOpts) (*callbackHandler, error) {
	counters, err := newQueryCounters(
		opts.prometheusNamespace,
	)
	if err != nil {
		return nil, err
	}

	return &callbackHandler{
		opts:     *opts,
		counters: counters,
		defaultLabels: map[string]string{
			labelDriver:   info.driverName,
			labelDatabase: info.dbName,
		},
	}, nil
}

// callbackName creates a GORM callback name based on the configured plugin
// scope and callback name.
func (c *pluginOpts) callbackName(callback string) string {
	return fmt.Sprintf("%v:%v", c.gormPluginScope, callback)
}

// Merges maps a and b. a is returned with extra values from b. Existing items
// in a with a matching key in b will not get overwritten.
func mergeLabels(a, b prometheus.Labels) prometheus.Labels {
	for k, v := range b {
		if _, exists := a[k]; !exists {
			a[k] = v
		}
	}
	return a
}
