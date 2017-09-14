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

func (db *EriboStore) createSchema() error {
	if _, err := db.Exec(tableMessages); err != nil {
		return err
	}
	if _, err := db.Exec(tableImages); err != nil {
		return err
	}
	if _, err := db.Exec(tableFeedback); err != nil {
		return err
	}
	return nil
}

func (db *EriboStore) dropSchema() error {
	if _, err := db.Exec(`DROP TABLE images`); err != nil {
		return err
	}
	if _, err := db.Exec(`DROP TABLE messages`); err != nil {
		return err
	}
	if _, err := db.Exec(`DROP TABLE feedback`); err != nil {
		return err
	}
	return nil
}

const (
	tableMessages = `
CREATE TABLE IF NOT EXISTS messages (
	id SERIAL,
	message TEXT NOT NULL,
	player VARCHAR(255) NOT NULL,
	channel VARCHAR(255) NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
)`

	tableImages = `
CREATE TABLE IF NOT EXISTS images (
	id SERIAL,
	url VARCHAR(2000) NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	message_id BIGINT UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	FOREIGN KEY (message_id) REFERENCES messages(id)
)`

	tableFeedback = `
CREATE TABLE IF NOT EXISTS feedback (
	id SERIAL,
	message TEXT NOT NULL,
	player VARCHAR(255) NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
)`
)
