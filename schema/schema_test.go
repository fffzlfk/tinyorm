package schema

import (
	"testing"
	"tinyorm/dialect"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name string `tinyorm:"PRIMARY KEY"`
	Age  int
}

var testDial, _ = dialect.GetDialect("postgres")

func TestParse(t *testing.T) {
	schema := Parse(&Person{}, testDial)
	assert.Equal(t, "Person", schema.Name)
	assert.Equal(t, []string{"Name", "Age"}, schema.FieldNames)
	assert.Equal(t, "PRIMARY KEY", schema.GetField("Name").Tag)
}
