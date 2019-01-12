package mysql

import (
	"log"
	"time"

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
	// To resolve the invalid connection issue. See:
	//
	// https://github.com/go-sql-driver/mysql/issues/674#issuecomment-345661869
	db.SetConnMaxLifetime(10 * time.Second)
	if err := pingDatabase(db); err != nil {
		log.Fatalln("database ping attempts failed:", err)
	}
	store := &EriboStore{DB: db}
	if err := store.createSchema(); err != nil {
		return nil, err
	}
	return store, nil
}

func pingDatabase(db *sqlx.DB) (err error) {
	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			return
		}
		time.Sleep(time.Second)
	}
	return
}
