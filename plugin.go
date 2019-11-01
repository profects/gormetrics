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
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/profects/gormetrics/gormi"
	"github.com/profects/gormetrics/gormi/adapter/unforked"
)

// Register gormetrics. Options (opts) can be used to configure the Prometheus
// namespace and GORM plugin scope.
func Register(db *gorm.DB, dbName string, opts ...RegisterOpt) error {
	if db == nil {
		return ErrDbIsNil
	}
	return RegisterInterface(unforked.New(db), dbName, opts...)
}

// RegisterInterface registers gormetrics with a gormi.DB interface, which can
// be created using one of the adapters in gormi/adapter. This can be useful if
// you use a forked version of GORM.
// Options (opts) can be used to configure the Prometheus namespace and
// GORM plugin scope.
func RegisterInterface(db gormi.DB, dbName string, opts ...RegisterOpt) error {
	if v := reflect.ValueOf(db); v.Kind() == reflect.Ptr && v.IsNil() {
		return ErrDbIsNil
	}

	driverName := sqlDriverToDriverName(db.DB().Driver())
	handlerOpts := getOpts(opts)
	info := extraInfo{
		dbName:     dbName,
		driverName: driverName,
	}

	handler, err := newCallbackHandler(info, handlerOpts)
	if err != nil {
		return errors.Wrap(err, "could not create callback handler")
	}
	handler.registerCallback(db.Callback())

	dbInterface := newDatabase(info, db.DB())
	dbMetrics, err := newDatabaseMetrics(dbInterface, handlerOpts)
	if err != nil {
		return errors.Wrap(err, "could not create database metrics exporter")
	}
	go dbMetrics.maintain()

	return nil
}
