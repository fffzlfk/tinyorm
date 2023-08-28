package session

import (
	"testing"
	"tinyorm/utils/tests"

	"github.com/stretchr/testify/assert"
)

func TestSession_CreateTable(t *testing.T) {
	s := newSession().Model(&tests.Person{})
	defer s.DropTable()
	var err error
	err = s.DropTable()
	assert.NoError(t, err)
	err = s.CreateTable()
	assert.NoError(t, err)
	assert.True(t, s.HasTable())
}
