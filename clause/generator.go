package clause

import (
	"fmt"
	"strings"
)

type generator func(values []any) (string, []any)

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[Insert] = func(values []any) (string, []any) {
		tableName := values[0]
		fields := strings.Join(values[1].([]string), ", ")
		return fmt.Sprintf("INSERT INTO %s (%s)", tableName, fields), values[2:]
	}
	genBindVars := func(start int, num int) string {
		var vars []string
		for i := start; i < start+num; i++ {
			vars = append(vars, fmt.Sprintf("$%d", i))
		}
		return strings.Join(vars, ", ")
	}
	generators[Values] = func(values []any) (string, []any) {
		var sql strings.Builder
		var vars []any
		sql.WriteString("VALUES ")
		index := 1
		for i, value := range values {
			v := value.([]any)
			bindStr := genBindVars(index, len(v))
			index += len(v)
			sql.WriteString(fmt.Sprintf("(%s)", bindStr))
			if i+1 != len(values) {
				sql.WriteString(", ")
			}
			vars = append(vars, v...)
		}
		return sql.String(), vars
	}
	generators[Select] = func(values []any) (string, []any) {
		tableName := values[0]
		fields := strings.Join(values[1].([]string), ",")
		return fmt.Sprintf("SELECT %s FROM %s", fields, tableName), []any{}
	}
	generators[Limit] = func(values []any) (string, []any) {
		return "LIMIT ?", values
	}
	generators[Where] = func(values []any) (string, []any) {
		desc, vars := values[0], values[1:]
		return fmt.Sprintf("WHERE %s", desc), vars
	}
	generators[OrderBy] = func(values []any) (string, []any) {
		return "ORDER BY " + values[0].(string), []any{}
	}
}
