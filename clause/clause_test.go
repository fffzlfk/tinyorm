package clause

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testSelect(t *testing.T) {
	var clause Clause
	clause.Set(Limit, 3)
	clause.Set(Select, "User", []string{"*"})
	clause.Set(Where, "Name = ?", "Tom")
	clause.Set(OrderBy, "Age ASC")
	sql, vars := clause.Build(Select, Where, OrderBy, Limit)
	assert.Equal(t, "SELECT * FROM User WHERE Name = $1 ORDER BY Age ASC LIMIT $2", sql)
	assert.Equal(t, []any{"Tom", 3}, vars)
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect(t)
	})
}
