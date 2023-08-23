package tinyorm

import (
	"database/sql"
	"tinyorm/dialect"
	"tinyorm/log"
	"tinyorm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s not found", driver)
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Infof("connect to %s success", source)
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error(err)
	}
	log.Info("close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}
