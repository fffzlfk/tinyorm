package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name string `tinyorm:"PRIMARY KEY"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	s := newSession().Model(&Person{})
	var err error
	err = s.DropTable()
	assert.NoError(t, err)
	err = s.CreateTable()
	assert.NoError(t, err)
	assert.True(t, s.HasTable())
}
