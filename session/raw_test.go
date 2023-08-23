package session

import (
	"database/sql"
	"os"
	"testing"
	"tinyorm/dialect"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var (
	testDB      *sql.DB
	testDial, _ = dialect.GetDialect("postgres")
)

func newSession() *Session {
	return New(testDB, testDial)
}

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open("postgres", "user=postgres dbname=test password=123456 sslmode=disable")
	if err != nil {
		panic(err)
	}
	code := m.Run()
	testDB.Close()
	os.Exit(code)
}

func TestSession_Exec(t *testing.T) {
	s := newSession()
	_, err := s.Raw("DROP TABLE IF EXISTS users").Exec()
	assert.NoError(t, err)
	_, err = s.Raw("CREATE TABLE users (name TEXT)").Exec()
	assert.NoError(t, err)
	res, err := s.Raw("INSERT INTO users VALUES ($1), ($2)", "Tom", "Sam").Exec()
	assert.NoError(t, err)
	count, err := res.RowsAffected()
	assert.Equal(t, int64(2), count)
	assert.NoError(t, err)
}

func TestSession_QueryRows(t *testing.T) {
	s := newSession()
	_, err := s.Raw("drop table if exists users").Exec()
	assert.NoError(t, err)
	_, err = s.Raw("create table users (name text)").Exec()
	assert.NoError(t, err)
	row := s.Raw("select count(*) from users").QueryRow()
	var count int
	if err := row.Scan(&count); err != nil {
		assert.NoError(t, err, "failed to query db")
	}
	assert.Equal(t, 0, count)
}
