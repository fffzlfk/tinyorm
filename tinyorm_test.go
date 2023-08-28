package tinyorm

import (
	"errors"
	"testing"
	"tinyorm/session"
	"tinyorm/utils/tests"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func openDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("postgres", "user=postgres dbname=test password=123456 sslmode=disable")
	assert.NoError(t, err)
	return engine
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func transactionRollback(t *testing.T) {
	engine := openDB(t)
	defer engine.Close()

	s := engine.NewSession()
	err1 := s.Model(&tests.Person{}).DropTable()
	assert.NoError(t, err1)

	_, err2 := engine.Transaction(func(s *session.Session) (any, error) {
		err3 := s.Model(&tests.Person{}).CreateTable()
		assert.NoError(t, err3)
		_, err4 := s.Insert(&tests.Person{Name: "Jack", Age: 19})
		assert.NoError(t, err4)
		return nil, errors.New("something wrong")
	})
	assert.Error(t, err2)
	assert.False(t, s.HasTable())
}

func transactionCommit(t *testing.T) {
	engine := openDB(t)
	defer engine.Close()

	s := engine.NewSession()
	err1 := s.Model(&tests.Person{}).DropTable()
	assert.NoError(t, err1)

	_, err2 := engine.Transaction(func(s *session.Session) (any, error) {
		err3 := s.Model(&tests.Person{}).CreateTable()
		assert.NoError(t, err3)
		_, err4 := s.Insert(&tests.Person{Name: "Jack", Age: 19})
		assert.NoError(t, err4)
		return nil, nil
	})
	assert.NoError(t, err2)

	assert.True(t, s.HasTable())
	var p tests.Person
	err5 := s.First(&p)
	assert.NoError(t, err5)
	assert.Equal(t, 19, p.Age)
	assert.Equal(t, "Jack", p.Name)
}
