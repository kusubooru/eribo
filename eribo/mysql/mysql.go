package mysql

import (
	"github.com/jmoiron/sqlx"
)

type EriboStore struct {
	*sqlx.DB
}

func NewEriboStore(dataSource string) (*EriboStore, error) {
	db, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	store := &EriboStore{DB: db}
	if err := store.createSchema(); err != nil {
		return nil, err
	}
	return store, nil
}
