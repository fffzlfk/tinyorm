package session

import (
	"reflect"
	"tinyorm/clause"
)

func (s *Session) Insert(values ...any) (int64, error) {
	recordValues := make([]any, 0, len(values))
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.Insert, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}
	s.clause.Set(clause.Values, recordValues...)
	sql, vars := s.clause.Build(clause.Insert, clause.Values)
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Session) Find(values any) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()
	s.clause.Set(clause.Select, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.Select, clause.Where, clause.OrderBy, clause.Limit)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		values := make([]any, 0, len(table.Fields))
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
