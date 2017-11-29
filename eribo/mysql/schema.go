package mysql

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
	if _, err := db.Exec(tableLog); err != nil {
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
	if _, err := db.Exec(`DROP TABLE log`); err != nil {
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
	created TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	PRIMARY KEY (id)
)`

	tableImages = `
CREATE TABLE IF NOT EXISTS images (
	id SERIAL,
	url VARCHAR(2000) NOT NULL,
	created TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	message_id BIGINT UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	FOREIGN KEY (message_id) REFERENCES messages(id)
)`

	tableFeedback = `
CREATE TABLE IF NOT EXISTS feedback (
	id SERIAL,
	message TEXT NOT NULL,
	player VARCHAR(255) NOT NULL,
	created TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	PRIMARY KEY (id)
)`
	tableLog = `
CREATE TABLE IF NOT EXISTS log (
	id SERIAL,
	command VARCHAR(255) NOT NULL,
	player VARCHAR(255) NOT NULL,
	channel VARCHAR(255) NOT NULL DEFAULT '',
	created TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	PRIMARY KEY (id)
)`
)
