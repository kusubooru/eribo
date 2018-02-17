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
