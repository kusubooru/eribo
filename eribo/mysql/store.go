package mysql

import (
	"fmt"

	"github.com/kusubooru/eribo/eribo"
)

func (db *EriboStore) AddMessageWithURLs(m *eribo.Message, urls []string) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rerr := tx.Rollback(); rerr != nil {
				err = fmt.Errorf("rollback failed: %v: %v", rerr, err)
			}
			return
		}
		err = tx.Commit()
	}()

	r, err := tx.Exec("insert into messages(message, player, channel) values (?, ?, ?)", m.Message, m.Player, m.Channel)
	if err != nil {
		return err
	}

	messageID, err := r.LastInsertId()
	if err != nil {
		return err
	}

	for _, u := range urls {
		_, err := tx.Exec("insert into images(url, message_id) values (?, ?)", u, messageID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *EriboStore) GetImages() ([]*eribo.Image, error) {
	images := []*eribo.Image{}
	const query = `
	SELECT
	  img.*,
	  m.id as "message.id",
	  m.player as "message.player",
	  m.channel as "message.channel",
	  m.message as "message.message",
	  m.created as "message.created"
	FROM images img
	  JOIN messages m ON img.message_id=m.id`
	if err := db.Select(&images, query); err != nil {
		return nil, err
	}
	return images, nil
}
