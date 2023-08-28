package session

import (
	"database/sql"
	"strconv"
	"strings"

	"tinyorm/clause"
	"tinyorm/dialect"
	"tinyorm/log"
	"tinyorm/schema"
)

type CommonDB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	tx       *sql.Tx
	refTable *schema.Schema
	clause   clause.Clause
	sql      strings.Builder
	sqlVars  []any
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *Session) Raw(sql string, values ...any) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

func parseBraces(sqlStr string) string {
	for i := 1; strings.Contains(sqlStr, "?"); i++ {
		sqlStr = strings.Replace(sqlStr, "?", "$"+strconv.Itoa(i), 1)
	}
	return sqlStr
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	sqlStr := parseBraces(s.sql.String())
	log.Info(sqlStr, s.sqlVars)
	if result, err = s.DB().Exec(sqlStr, s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	sqlStr := parseBraces(s.sql.String())
	log.Info(sqlStr, s.sqlVars)
	return s.DB().QueryRow(sqlStr, s.sqlVars...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	sqlStr := parseBraces(s.sql.String())
	log.Info(sqlStr, s.sqlVars)
	if rows, err = s.DB().Query(sqlStr, s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
