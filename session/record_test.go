package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	person1 = &Person{
		Name: "Tom",
		Age:  12,
	}
	person2 = &Person{
		Name: "Jack",
		Age:  22,
	}
)

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := newSession().Model(&Person{})
	err1 := s.DropTable()
	assert.NoError(t, err1)
	err2 := s.CreateTable()
	assert.NoError(t, err2)
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	affected, err := s.Insert(person1, person2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), affected)
}

func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)
	_, err := s.Insert(person1, person2)
	assert.NoError(t, err)
	var persons []Person
	err = s.Find(&persons)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(persons))
}
