package gormetrics

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"sync"
)

var sqlDriverNamesByType map[reflect.Type]string //nolint:gochecknoglobals
var sqlDriverNamesByTypeLock sync.Mutex          //nolint:gochecknoglobals

// The database/sql API doesn't provide a way to get the registry name for
// a driver from the driver type.
// See https://github.com/golang/go/issues/12600#issuecomment-378363201
func sqlDriverToDriverName(driver driver.Driver) string {
	sqlDriverNamesByTypeLock.Lock()
	defer sqlDriverNamesByTypeLock.Unlock()

	if sqlDriverNamesByType == nil {
		sqlDriverNamesByType = map[reflect.Type]string{}

		for _, driverName := range sql.Drivers() {
			// We ignore this error because it will occur. We only need the
			// driver connected to the name.
			db, _ := sql.Open(driverName, "")

			if db != nil {
				driverType := reflect.TypeOf(db.Driver())
				sqlDriverNamesByType[driverType] = driverName
			}
		}
	}

	driverType := reflect.TypeOf(driver)
	if driverName, found := sqlDriverNamesByType[driverType]; found {
		return driverName
	}

	return ""
}
