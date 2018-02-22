package mysql

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kusubooru/eribo/eribo"
)

const timeTruncate = 1 * time.Second

func (db *EriboStore) Tx(fn func(*sqlx.Tx) error) (err error) {
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
	return fn(tx)
}

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

	const insertMessage = `INSERT INTO messages(message, player, channel, created) VALUES (?, ?, ?, ?)`
	r, err := tx.Exec(insertMessage, m.Message, m.Player, m.Channel, m.Created)
	if err != nil {
		return err
	}

	messageID, err := r.LastInsertId()
	if err != nil {
		return err
	}

	const insertImage = "INSERT INTO images(url, done, kuid, message_id, created) VALUES (?, ?, ?, ?, ?)"
	for _, u := range urls {
		_, err := tx.Exec(insertImage, u, false, 0, messageID, m.Created)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *EriboStore) GetImages(limit, offset int, reverse, filterDone bool) ([]*eribo.Image, error) {
	desc := ""
	if reverse {
		desc = "DESC"
	}

	filter := ""
	if filterDone {
		filter = " AND img.done = 0 "
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
	  JOIN messages m ON img.message_id=m.id` + filter + `
	  ORDER BY created ` + desc + ` LIMIT ?, ?`

	images := []*eribo.Image{}
	if err := db.Select(&images, query, offset, limit); err != nil {
		return nil, err
	}
	return images, nil
}

func (db *EriboStore) GetImage(id int64) (*eribo.Image, error) {
	img := &eribo.Image{}
	err := db.Tx(func(tx *sqlx.Tx) error {
		const query = `
	    SELECT
	      img.*,
	      m.id as "message.id",
	      m.player as "message.player",
	      m.channel as "message.channel",
	      m.message as "message.message",
	      m.created as "message.created"
	    FROM images img
	      JOIN messages m ON img.message_id=m.id
	      WHERE img.id = ?`
		return tx.Get(img, query, id)
	})
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (db *EriboStore) ToggleImageDone(id int64) (err error) {
	return db.Tx(func(tx *sqlx.Tx) error {
		const query = `SELECT * FROM images WHERE id = ?`

		img := &eribo.Image{}
		err = tx.Get(img, query, id)
		if err != nil {
			return err
		}

		const update = `update images set done = ? where id = ?`
		_, err = tx.Exec(update, !img.Done, id)
		return err
	})
}

func (db *EriboStore) SetImageKuid(id int64, kuid int) error {
	done := false
	if kuid != 0 {
		done = true
	}
	return db.Tx(func(tx *sqlx.Tx) error {
		const query = `UPDATE images SET kuid = ?, done = ? WHERE id = ?`
		_, err := tx.Exec(query, kuid, done, id)
		return err
	})
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
