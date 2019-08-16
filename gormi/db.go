// Package gormi provides interface so forks of gorm can be used to work with gormetrics.
// See github.com/profects/gormetrics/gormi/adapter/unforked for an example of
// how to use these interfaces.
package gormi

import "database/sql"

// DB is an interface which can be easily satisfied by wrapping gorm.DB.
type DB interface {
	DB() *sql.DB
	Callback() Callback
	GetErrors() []error
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
	SQLQuery() string
}
