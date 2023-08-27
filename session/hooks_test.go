package session

import (
	"testing"
	"tinyorm/log"

	"github.com/stretchr/testify/assert"
)

type Account struct {
	ID       int `tinyorm:"PRIMARY KEY"`
	Password string
}

func (a *Account) BeforeInsert(s *Session) error {
	log.Info("before insert", a)
	a.ID += 1000
	return nil
}

func (a *Account) AfterQuery(s *Session) error {
	log.Info("after query", a)
	a.Password = "******"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	s := newSession().Model(&Account{})
	err1 := s.DropTable()
	assert.NoError(t, err1)
	err2 := s.CreateTable()
	assert.NoError(t, err2)
	affected, err3 := s.Insert(&Account{1, "123456"}, &Account{2, "mypassword"})
	assert.NoError(t, err3)
	assert.Equal(t, affected, int64(2))
	acc := Account{}
	err4 := s.First(&acc)
	assert.NoError(t, err4)
	assert.Equal(t, 1001, acc.ID)
	assert.Equal(t, "******", acc.Password)
}
