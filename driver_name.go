// Copyright 2019 Profects Group B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// driverNames stores all mapped drivers by name.
var driverNames = sqlDriverNames{
	byType: make(map[reflect.Type]string),
}

// The database/sql API doesn't provide a way to get the registry name for
// a driver from the driver type.
// Adapted from https://github.com/golang/go/issues/12600#issuecomment-378363201.
func sqlDriverToDriverName(driver driver.Driver) string {
	driverNames.Lock()
	defer driverNames.Unlock()

	driverType := reflect.TypeOf(driver)

	if len(driverNames.byType) > 0 {
		if driverName, found := driverNames.byType[driverType]; found {
			return driverName
		}
	}

	for _, driverName := range sql.Drivers() {
		// We ignore this error because it will occur. We only need the
		// driver connected to the name.
		db, _ := sql.Open(driverName, "")

		if db != nil {
			driverType := reflect.TypeOf(db.Driver())
			driverNames.byType[driverType] = driverName
		}
	}

	if len(driverNames.byType) > 0 {
		if driverName, found := driverNames.byType[driverType]; found {
			return driverName
		}
	}

	return ""
}
