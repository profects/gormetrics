package unforked

import (
	"database/sql"

	"github.com/profects/gormetrics/gormi"

	"github.com/jinzhu/gorm"
)

// New creates a new gormi.DB interface from the unforked version of GORM for
// gormetrics to use.
func New(db *gorm.DB) gormi.DB {
	return wrappedDB{db}
}

type wrappedDB struct {
	db *gorm.DB
}

type wrappedCallback struct {
	callback *gorm.Callback
}

type wrappedCallbackProcessor struct {
	processor *gorm.CallbackProcessor
}

type wrappedScope struct {
	scope *gorm.Scope
}

func (w wrappedDB) Callback() gormi.Callback {
	return wrappedCallback{
		callback: w.db.Callback(),
	}
}

func (w wrappedDB) DB() *sql.DB {
	return w.db.DB()
}

func (w wrappedDB) GetErrors() []error {
	return w.db.GetErrors()
}

func (w wrappedCallback) Create() gormi.CallbackProcessor {
	return wrappedCallbackProcessor{
		processor: w.callback.Create(),
	}
}

func (w wrappedCallback) Delete() gormi.CallbackProcessor {
	return wrappedCallbackProcessor{
		processor: w.callback.Delete(),
	}
}

func (w wrappedCallback) Query() gormi.CallbackProcessor {
	return wrappedCallbackProcessor{
		processor: w.callback.Query(),
	}
}

func (w wrappedCallback) Update() gormi.CallbackProcessor {
	return wrappedCallbackProcessor{
		processor: w.callback.Update(),
	}
}

func (w wrappedCallbackProcessor) After(callbackName string) gormi.CallbackProcessor {
	return wrappedCallbackProcessor{
		processor: w.processor.After(callbackName),
	}
}

func (w wrappedCallbackProcessor) Register(callbackName string, callback func(scope gormi.Scope)) {
	w.processor.Register(callbackName, func(scope *gorm.Scope) {
		callback(wrappedScope{
			scope: scope,
		})
	})
}

func (w wrappedScope) DB() gormi.DB {
	return wrappedDB{w.scope.DB()}
}
