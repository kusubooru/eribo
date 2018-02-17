package mysql

import (
	"fmt"
	"time"

	"github.com/kusubooru/eribo/eribo"
)

const timeTruncate = 1 * time.Second

func (db *EriboStore) AddMessageWithURLs(m *eribo.Message, urls []string) (err error) {
	if (m.Created == time.Time{}) {
		m.Created = time.Now().UTC().Truncate(timeTruncate)
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
		_, err := tx.Exec("INSERT INTO images(url, done, message_id, created) VALUES (?, ?, ?, ?)", u, false, messageID, m.Created)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *EriboStore) GetImages(limit, offset int, reverse bool) ([]*eribo.Image, error) {
	desc := ""
	if reverse {
		desc = "DESC"
	}

	var query = `
	SELECT
	  img.*,
	  m.id as "message.id",
	  m.player as "message.player",
	  m.channel as "message.channel",
	  m.message as "message.message",
	  m.created as "message.created"
	FROM images img
	  JOIN messages m ON img.message_id=m.id
	  ORDER BY created ` + desc + ` LIMIT ?, ?`

	images := []*eribo.Image{}
	if err := db.Select(&images, query, offset, limit); err != nil {
		return nil, err
	}
	return images, nil
}

func (db *EriboStore) ToggleImageDone(id int64) (err error) {
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

	const query = `SELECT * FROM images WHERE id = ?`

	img := &eribo.Image{}
	if err := tx.Get(img, query, id); err != nil {
		return err
	}

	const update = `update images set done = ? where id = ?`
	if _, err := tx.Exec(update, !img.Done, id); err != nil {
		return err
	}

	return nil
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
	if (f.Created == time.Time{}) {
		f.Created = time.Now().UTC().Truncate(timeTruncate)
	}
	const query = `INSERT INTO feedback(message, player, created) VALUES (?, ?, ?)`
	_, err := db.Exec(query, f.Message, f.Player, f.Created)
	return err
}
