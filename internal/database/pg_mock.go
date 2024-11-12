package database

import (
	"github.com/go-pg/pg/v10"
)

func NewDBClientMock(err error) DBClient {
	return &mockDB{
		err: err,
	}
}

type mockDB struct {
	err error
}

func (m *mockDB) Conn() *pg.DB {
	return &pg.DB{}
}
