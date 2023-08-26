package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...any) (string, []any)

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[Insert] = _insert
	generators[Values] = _values
	generators[Select] = _select
	generators[Where] = _where
	generators[Limit] = _limit
	generators[OrderBy] = _orderBy
	generators[Update] = _update
	generators[Delete] = _delete
	generators[Count] = _count
}

func genBindVars(num int) string {
	vars := make([]string, 0, num)
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _insert(values ...any) (string, []any) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ", ")
	return fmt.Sprintf("INSERT INTO %s (%s)", tableName, fields), values[2:]
}

func _values(values ...any) (string, []any) {
	var sql strings.Builder
	var vars []any
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]any)
		bindStr := genBindVars(len(v))
		sql.WriteString(fmt.Sprintf("(%s)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}
func _select(values ...any) (string, []any) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ", ")
	return fmt.Sprintf("SELECT %s FROM %s", fields, tableName), []any{}
}

func _where(values ...any) (string, []any) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _limit(values ...any) (string, []any) {
	return "LIMIT ?", values
}

func _orderBy(values ...any) (string, []any) {
	return "ORDER BY " + values[0].(string), []any{}
}

func _update(values ...any) (string, []any) {
	tableName := values[0]
	m := values[1].(map[string]any)
	keys := make([]string, 0, len(m))
	vars := make([]any, 0, len(m))
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

func _delete(values ...any) (string, []any) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []any{}
}

func _count(values ...any) (string, []any) {
	return _select(values[0], []string{"COUNT(*)"})
}
