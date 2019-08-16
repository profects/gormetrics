// Package gormi provides interface so forks of gorm can be used to work with gormetrics.
package gormi

import (
	"database/sql"
)

type DB interface {
	DB() *sql.DB
	Callback() Callback
	GetErrors() []error
}

type Callback interface {
	Create() CallbackProcessor
	Delete() CallbackProcessor
	Query() CallbackProcessor
	Update() CallbackProcessor
}

type CallbackProcessor interface {
	After(callbackName string) CallbackProcessor
	Register(callbackName string, callback func(scope Scope))
}

type Scope interface {
	DB() DB
	SQLQuery() string
}
