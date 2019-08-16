package gormetrics

const labelStatus = "status"
const labelDatabase = "database"
const labelDriver = "driver"

// Statuses for metrics (values of labelStatus).
const metricStatusFail = "fail"
const metricStatusSuccess = "success"

const metricOpenConnections = "connections_open"
const metricIdleConnections = "connections_idle"
const metricInUseConnections = "connections_in_use"

const helpOpenConnections = `Currently open connections to the database`
const helpIdleConnections = `Currently idle connections to the database`
const helpInUseConnections = `Currently in use connections`

const metricAllTotal = "all_total"
const metricCreatesTotal = "creates_total"
const metricDeletesTotal = "deletes_total"
const metricQueriesTotal = "queries_total"
const metricUpdatesTotal = "updates_total"

const helpAllTotal = `All queries requested`
const helpCreatesTotal = `All create queries requested`
const helpDeletesTotal = `All delete queries requested`
const helpQueriesTotal = `All select queries requested`
const helpUpdatesTotal = `All update queries requested`
