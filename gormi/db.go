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

// Package gormi provides interface so forks of gorm can be used to work with gormetrics.
// See github.com/profects/gormetrics/gormi/adapter/unforked for an example of
// how to use these interfaces.
package gormi

import "database/sql"

// DB is an interface which can be easily satisfied by wrapping gorm.DB.
type DB interface {
	DB() *sql.DB
	Callback() Callback
	Error() error
}

// Callback is an interface which can be satisfied by wrapping gorm.Callback.
type Callback interface {
	Create() CallbackProcessor
	Delete() CallbackProcessor
	Query() CallbackProcessor
	Update() CallbackProcessor
}

// CallbackProcessor is an interface which can be satisfied by wrapping
// gorm.CallbackProcessor.
type CallbackProcessor interface {
	After(callbackName string) CallbackProcessor
	Register(callbackName string, callback func(scope Scope))
}

// Scope is an interface which can be satisfied by wrapping gorm.Scope.
type Scope interface {
	DB() DB
}
