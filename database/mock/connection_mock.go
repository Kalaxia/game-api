package mock

import(
	"reflect"
	"github.com/go-pg/pg"
	"kalaxia-game-api/database"
)

type TestConnection struct {
	pg.DB

	NextId uint8
}

func (tc *TestConnection) Insert(model ...interface{}) error {
	if id := reflect.ValueOf(model[0]).Elem().FieldByName("Id"); id.IsValid() {
		id.SetUint(uint64(tc.NextId))
		tc.NextId++
	}
	return nil
}

func (TestConnection) Update(model interface{}) error {
	return nil
}

func init() {
	database.Connection = &TestConnection{}
}