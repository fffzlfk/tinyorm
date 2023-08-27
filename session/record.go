package session

import (
	"errors"
	"reflect"
	"tinyorm/clause"
)

func (s *Session) Insert(values ...any) (int64, error) {
	recordValues := make([]any, 0, len(values))
	for _, value := range values {
		s.CallMethod(BeforeInsert, value)
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
		s.CallMethod(AfterQuery, dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

func (s *Session) First(value any) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("record not found")
	}
	dest.Set(destSlice.Index(0))
	return nil
}

func (s *Session) Update(kv ...any) (int64, error) {
	mp, ok := kv[0].(map[string]any)
	if !ok {
		mp = make(map[string]any)
		for i := 0; i < len(kv); i += 2 {
			mp[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.Update, s.refTable.Name, mp)
	sql, vars := s.clause.Build(clause.Update, clause.Where)
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.Delete, s.refTable.Name)
	sql, vars := s.clause.Build(clause.Delete, clause.Where)
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.Count, s.refTable.Name)
	sql, vars := s.clause.Build(clause.Count, clause.Where)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.Limit, num)
	return s
}

func (s *Session) Where(desc string, args ...any) *Session {
	s.clause.Set(clause.Where, append([]any{desc}, args...)...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.OrderBy, desc)
	return s
}
