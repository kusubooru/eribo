package mysql

import (
	"fmt"
	"time"

	"github.com/kusubooru/eribo/eribo"
)

func (db *EriboStore) AddMessageWithURLs(m *eribo.Message, urls []string) (err error) {
	if (m.Created == time.Time{}) {
		m.Created = time.Now().UTC().Truncate(1 * time.Microsecond)
	}
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

	const query = `INSERT INTO messages(message, player, channel, created) VALUES (?, ?, ?, ?)`
	r, err := tx.Exec(query, m.Message, m.Player, m.Channel, m.Created)
	if err != nil {
		return err
	}

	messageID, err := r.LastInsertId()
	if err != nil {
		return err
	}

	for _, u := range urls {
		_, err := tx.Exec("INSERT INTO images(url, message_id, created) VALUES (?, ?, ?)", u, messageID, m.Created)
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

func (db *EriboStore) GetAllFeedback(limit, offset int) ([]*eribo.Feedback, error) {
	feedback := []*eribo.Feedback{}
	const query = `SELECT * FROM feedback LIMIT ? OFFSET ?`
	if err := db.Select(&feedback, query, limit, offset); err != nil {
		return nil, err
	}
	return feedback, nil
}

func (db *EriboStore) GetRecentFeedback(limit, offset int) ([]*eribo.Feedback, error) {
	feedback := []*eribo.Feedback{}
	const query = `SELECT * FROM feedback ORDER BY created DESC LIMIT ? OFFSET ?`
	if err := db.Select(&feedback, query, limit, offset); err != nil {
		return nil, err
	}
	return feedback, nil
}

func (db *EriboStore) AddFeedback(f *eribo.Feedback) error {
	_, err := db.Exec("INSERT INTO feedback(message, player) VALUES (?, ?)", f.Message, f.Player)
	if err != nil {
		return err
	}
	return nil
}
