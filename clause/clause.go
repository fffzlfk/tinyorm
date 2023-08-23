package clause

import "strings"

type Type int

const (
	Insert Type = iota
	Values
	Select
	Limit
	Where
	OrderBy
)

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]any
}

func (c *Clause) Set(name Type, args ...any) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]any)
	}
	sql, args := generators[name](args)
	c.sql[name] = sql
	c.sqlVars[name] = args
}

func (c *Clause) Build(orders ...Type) (string, []any) {
	sqls := make([]string, 0, len(orders))
	vars := make([]any, 0, len(orders))
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
