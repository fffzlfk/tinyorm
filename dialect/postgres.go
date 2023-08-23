package dialect

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type PostgreSQL struct{}

func init() {
	RegisterDialect("postgres", &PostgreSQL{})
}

const intSize = 32 << (^uint(0) >> 63)

// DataTypeOf implements Dialect.
func (*PostgreSQL) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Int16:
		return "smallint"
	case reflect.Int32:
		return "integer"
	case reflect.Int64:
		return "bigint"
	case reflect.Int:
		if intSize == 64 {
			return "bigint"
		}
		return "integer"
	case reflect.Float32:
		return "real"
	case reflect.Float64:
		return "double precision"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "bytea"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "timestamp with time zone"
		}
		return "json"
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

// TableExistSQL implements Dialect.
func (*PostgreSQL) TableExistSQL(tableName string) (string, []any) {
	args := []any{strings.ToLower(tableName)}
	return "SELECT table_name FROM information_schema.tables WHERE table_schema='public' and table_name=$1", args
}

var _ Dialect = &PostgreSQL{}
