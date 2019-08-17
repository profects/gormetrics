package gormetrics

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"sync"
)

type sqlDriverNames struct {
	byType map[reflect.Type]string

	sync.Mutex
}

var driverNames = sqlDriverNames{
	byType: make(map[reflect.Type]string),
}

// The database/sql API doesn't provide a way to get the registry name for
// a driver from the driver type.
// See https://github.com/golang/go/issues/12600#issuecomment-378363201
func sqlDriverToDriverName(driver driver.Driver) string {
	driverNames.Lock()
	defer driverNames.Unlock()

	if len(driverNames.byType) == 0 {
		for _, driverName := range sql.Drivers() {
			// We ignore this error because it will occur. We only need the
			// driver connected to the name.
			db, _ := sql.Open(driverName, "")

			if db != nil {
				driverType := reflect.TypeOf(db.Driver())
				driverNames.byType[driverType] = driverName
			}
		}
	}

	driverType := reflect.TypeOf(driver)
	if driverName, found := driverNames.byType[driverType]; found {
		return driverName
	}

	return ""
}
