package dialect

import "reflect"

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []any)
}

var dialectsMap = map[string]Dialect{}

func RegisterDialect(name string, d Dialect) {
	dialectsMap[name] = d
}

func GetDialect(name string) (d Dialect, ok bool) {
	d, ok = dialectsMap[name]
	return
}
