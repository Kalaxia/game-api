package api

import(
	"reflect"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

type TestDatabase struct {
	pg.DB

	NextId uint8

	Data []interface{}
}

func (tc *TestDatabase) Insert(model ...interface{}) error {
	if id := reflect.ValueOf(model[0]).Elem().FieldByName("Id"); id.IsValid() {
		id.SetUint(uint64(tc.NextId))
		tc.NextId++
	}
	return nil
}

func (tc *TestDatabase) Update(model interface{}) error {
	return nil
}

func (tc *TestDatabase) Delete(model interface{}) error {
	return nil
}

func (tc *TestDatabase) Select(model interface{}) error {
	return nil
}

func (tc *TestDatabase) Model(model ...interface{}) *orm.Query {
	return orm.NewQuery(nil)
}

func InitDatabaseMock() {
	Database = &TestDatabase{}
}