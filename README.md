# gormetrics

[![Go Report Card](https://goreportcard.com/badge/github.com/profects/gormetrics)](https://goreportcard.com/report/github.com/profects/gormetrics)
[![GoDoc](https://godoc.org/github.com/profects/gormetrics?status.svg)](http://godoc.org/github.com/profects/gormetrics)

A plugin for GORM providing metrics using Prometheus.

Warning: this plugin is still in an early stage of development. APIs may change.

## Usage

```go
import "github.com/profects/gormetrics"

err := gormetrics.Register(db, <database name>)
if err != nil {
	// handle the error
}
```

gormetrics does not expose the metrics endpoint using promhttp, you have to do this yourself.
You can use the following snippet for exposing metrics on port 2112 at `/metrics`:

```go
go func() {
    http.Handle("/metrics", promhttp.Handler())
    log.Fatal(http.ListenAndServe(":2112", nil))
}()
```

## Exported metrics

gormetrics exports the following metrics (counter vectors):
* `gormetrics_all_total`
* `gormetrics_creates_total`
* `gormetrics_deletes_total`
* `gormetrics_queries_total`
* `gormetrics_updates_total`

These all have the following labels:
* `database`: the name of the database
* `driver`: the driver for the database (e.g. pq)
* `status`: fail or success

It also export the following metrics (gauge vectors):
* `gormetrics_connections_idle`
* `gormetrics_connections_in_use`
* `gormetrics_connections_open`

These all have the following labels:
* `database`: the name of the database
* `driver`: the driver for the database (e.g. pq)
