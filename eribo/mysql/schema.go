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
	if _, err := db.Exec(tableCmdLogs); err != nil {
		return err
	}
	if _, err := db.Exec(tableLothLogs); err != nil {
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
	if _, err := db.Exec(`DROP TABLE cmd_logs`); err != nil {
		return err
	}
	if _, err := db.Exec(`DROP TABLE loth_logs`); err != nil {
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
	tableCmdLogs = `
CREATE TABLE IF NOT EXISTS cmd_logs (
	id SERIAL,
	command VARCHAR(255) NOT NULL,
	args VARCHAR(2000) NOT NULL DEFAULT '',
	player VARCHAR(255) NOT NULL,
	channel VARCHAR(255) NOT NULL DEFAULT '',
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
)`
)

var tableLothLogs = `
CREATE TABLE IF NOT EXISTS loth_logs (
	id SERIAL,
	issuer VARCHAR(255) NOT NULL,
	channel VARCHAR(255) NOT NULL DEFAULT '',
	name VARCHAR(255) NOT NULL,
	role VARCHAR(20) NOT NULL,
	status VARCHAR(10) NOT NULL,
	is_new BOOL NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires TIMESTAMP NOT NULL ` + legacyDefault() + `,
	targets TEXT NOT NULL,
	PRIMARY KEY (id)
)`
