package session

import (
	"testing"
	"tinyorm/utils/tests"

	"github.com/stretchr/testify/assert"
)

var (
	person1 = &tests.Person{
		Name: "Tom",
		Age:  12,
	}
	person2 = &tests.Person{
		Name: "Jack",
		Age:  22,
	}
	person3 = &tests.Person{
		Name: "Rose",
		Age:  32,
	}
)

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := newSession().Model(&tests.Person{})
	err1 := s.DropTable()
	assert.NoError(t, err1)
	err2 := s.CreateTable()
	assert.NoError(t, err2)
	_, err3 := s.Insert(person1, person2)
	assert.NoError(t, err3)
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	affected, err := s.Insert(person3)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), affected)
}

func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)
	defer s.DropTable()
	var persons []tests.Person
	err := s.Find(&persons)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(persons))
}

func TestSession_Limit(t *testing.T) {
	s := testRecordInit(t)
	defer s.DropTable()
	var persons []tests.Person
	err := s.Limit(1).Find(&persons)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(persons))
}

func TestSession_Update(t *testing.T) {
	s := testRecordInit(t)
	defer s.DropTable()
	affected, err1 := s.Where("Name = ?", "Tom").Update("Age", 30)
	assert.NoError(t, err1)
	p := tests.Person{}
	err2 := s.OrderBy("Age DESC").First(&p)
	assert.NoError(t, err2)
	assert.Equal(t, int64(1), affected)
	assert.Equal(t, 30, p.Age)
}

func TestSession_DeleteAndCount(t *testing.T) {
	s := testRecordInit(t)
	defer s.DropTable()
	affected, err := s.Where("Name = ?", "Tom").Delete()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), affected)
	count, err := s.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}
