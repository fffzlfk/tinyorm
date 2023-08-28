package tinyorm

import (
	"database/sql"
	"fmt"
	"strings"
	"tinyorm/dialect"
	"tinyorm/log"
	"tinyorm/session"
	"tinyorm/utils"
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

type TxFunc func(*session.Session) (any, error)

func (e *Engine) Transaction(f TxFunc) (result any, err error) {
	s := e.NewSession()
	if err = s.Begin(); err != nil {
		log.Error(err)
		return
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // 回滚后重新抛出 panic
		} else if err != nil {
			_ = s.Rollback() // 回滚
		} else {
			err = s.Commit() // 提交
		}
	}()

	return f(s)
}

func (e *Engine) Migrate(value any) error {
	_, err := e.Transaction(func(s *session.Session) (any, error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s not exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.RefTable()
		rows, err := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		if err != nil {
			return nil, err
		}
		cols, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		addCols := utils.Difference(table.FieldNames, cols, false)
		delCols := utils.Difference(cols, table.FieldNames, false)
		log.Infof("added cols: %v, deleted cols: %v", addCols, delCols)

		for _, col := range addCols {
			field := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table.Name, field.Name, field.Type)
			log.Info(sqlStr)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return nil, err
			}
		}

		if len(delCols) == 0 {
			return nil, nil
		}
		tmpTable := "__tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s;", tmpTable, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmpTable, table.Name))
		_, err = s.Exec()
		return nil, err
	})
	return err
}
