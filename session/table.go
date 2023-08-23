package session

import (
	"fmt"
	"reflect"
	"strings"
	"tinyorm/log"
	"tinyorm/schema"
)

func (s *Session) Model(value any) *Session {
	// update refTable
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.RefTable()
	columns := make([]string, 0, len(table.Fields))
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("create table %s (%s)", table.Name, desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("drop table if exists %s", s.RefTable().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	err := row.Scan(&tmp)
	if err != nil {
		log.Error(err)
		return false
	}
	return strings.EqualFold(tmp, s.refTable.Name)
}
