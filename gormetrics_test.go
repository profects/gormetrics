package gormetrics

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

type attr struct {
	gorm.Model
	Name  string
	Obj   obj
	ObjID int
}

type obj struct {
	gorm.Model
	Name  string
	Attrs []attr
}

func TestGormetrics(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")
	assert.NoError(t, err)

	db.LogMode(true)

	defer func() {
		err = os.Remove("test.db")
		assert.NoError(t, err)
	}()

	err = Register(db, "testdb")
	assert.NoError(t, err)

	prepareTests(t, db)

	objs := make([]obj, 0)
	err = db.Preload("Attrs").Find(&objs, obj{}).Error
	assert.NoError(t, err)

	err = db.Unscoped().Delete(&attr{
		Model: gorm.Model{
			ID: 1,
		},
		Name: "obj1",
	}).Error
	assert.NoError(t, err)

	counter := collectors.query["gormetrics"].queries.With(prometheus.Labels{
		labelDriver:   "sqlite3",
		labelDatabase: "testdb",
		labelStatus:   metricStatusSuccess,
	})

	assert.Equal(t, *metricValue(counter).Counter.Value, 2.)
}

func metricValue(c prometheus.Collector) dto.Metric {
	ch := make(chan prometheus.Metric, 1)
	c.Collect(ch)

	m := dto.Metric{}
	_ = (<-ch).Write(&m)

	return m
}

func prepareTests(t *testing.T, db *gorm.DB) {
	err := db.AutoMigrate(obj{}, attr{}).Error
	assert.NoError(t, err)

	assert.NoError(t, db.Create(&obj{
		Name: "obj1",
		Attrs: []attr{
			{Name: "attr1"},
		},
	}).Error)
	assert.NoError(t, db.Create(&obj{
		Name: "obj2",
		Attrs: []attr{
			{Name: "attr2"},
			{Name: "attr3"},
			{Name: "attr4"},
			{Name: "attr5"},
			{Name: "attr6"},
		},
	}).Error)
}
