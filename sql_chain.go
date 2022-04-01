package gqt

import (
	"github.com/pkg/errors"
	"github.com/suifengpiao14/gqt/v2/gqttpl"
)

type SQLChain struct {
	sqlRows       []*SQLRow
	sqlRepository func() *RepositorySQL
	err           error
}

// 批量获取sql记录
func NewSQLChain(sqlRepository func() *RepositorySQL) (s *SQLChain) {
	s = &SQLChain{
		sqlRows:       make([]*SQLRow, 0),
		sqlRepository: sqlRepository,
	}
	return
}

func (s *SQLChain) ParseSQL(t gqttpl.TplEntityInterface, result interface{}) *SQLChain {
	if s.sqlRepository == nil {
		s.err = errors.Errorf("want SQLChain.sqlRepository ,have %#v", s)
	}
	if s.err != nil {
		return s
	}
	sqlRow, err := s.sqlRepository().GetSQL(t)
	if err != nil {
		s.err = err
		return s
	}
	s.sqlRows = append(s.sqlRows, sqlRow)
	return s
}

func (s *SQLChain) ParseTpEntity(entity gqttpl.TplEntityInterface, result interface{}) *SQLChain {
	if s.sqlRepository == nil {
		s.err = errors.Errorf("want SQLChain.sqlRepository ,have %#v", s)
	}
	if s.err != nil {
		return s
	}
	sqlRow, err := s.sqlRepository().GetSQL(entity)
	if err != nil {
		s.err = err
		return s
	}
	s.sqlRows = append(s.sqlRows, sqlRow)
	return s
}

//GetAllSQL get all sql from SQLChain
func (s *SQLChain) SQLRows() (sqlRowList []*SQLRow, err error) {
	return s.sqlRows, s.err
}

//Exec exec sql
func (s *SQLChain) Exec(fn func(sqlRowList []*SQLRow) (e error)) (err error) {
	if s.err != nil {
		return s.err
	}
	s.err = fn(s.sqlRows)
	return s.err
}

//Exec exec sql ,get data
func (s *SQLChain) Scan(fn func(sqlRowList []*SQLRow) (e error)) (err error) {
	if s.err != nil {
		return
	}
	s.err = fn(s.sqlRows)
	return s.err
}

//AddSQL add one sql to SQLChain
func (s *SQLChain) AddSQL(namespace string, name string, sql string, result interface{}) {
	sqlRow := &SQLRow{
		Name:      name,
		Namespace: name,
		SQL:       sql,
		Result:    result,
	}
	s.sqlRows = append(s.sqlRows, sqlRow)
}

func (s *SQLChain) SetError(err error) {
	if s.err != nil {
		return
	}
	if err != nil {
		err = errors.WithStack(err)
		s.err = err
	}
}

func (s *SQLChain) Error() (err error) {
	return s.err
}
