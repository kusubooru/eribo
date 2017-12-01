package mysql

import (
	"time"

	"github.com/kusubooru/eribo/eribo"
	"github.com/kusubooru/eribo/flist"
)

func (db *EriboStore) Log(e *eribo.Event) error {
	_, err := db.Exec("INSERT INTO log(command, player, channel) VALUES (?, ?, ?)", e.Command, e.Player, e.Channel)
	if err != nil {
		return err
	}
	return nil
}

func (db *EriboStore) GetLog(id int64) (*eribo.Event, error) {
	e := &eribo.Event{}
	const query = `SELECT * FROM log where id = ?`
	if err := db.Get(e, query, id); err != nil {
		return nil, err
	}
	return e, nil
}

func (db *EriboStore) GetRecentLogs(limit, offset int) ([]*eribo.Event, error) {
	logs := []*eribo.Event{}
	const query = `SELECT * FROM log ORDER BY created DESC LIMIT ? OFFSET ?`
	if err := db.Select(&logs, query, limit, offset); err != nil {
		return nil, err
	}
	return logs, nil
}

func (db *EriboStore) LogLoth(l *eribo.LothLog) error {
	if l.Created == (time.Time{}) {
		l.Created = time.Now().UTC().Truncate(1 * time.Microsecond)
	}
	var (
		name    string
		role    flist.Role
		status  flist.Status
		expires time.Time = l.Created
	)
	if l.Loth != nil {
		name = l.Loth.Name
		role = l.Loth.Role
		status = l.Loth.Status
		expires = l.Loth.Expires
	}

	const query = `INSERT INTO
	loth_logs(issuer, channel, created, name, role, status, expires, is_new, targets)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, l.Issuer, l.Channel, l.Created, name, role, status, expires, l.IsNew, l.Targets)
	if err != nil {
		return err
	}
	return nil
}

func (db *EriboStore) GetRecentLothLogs(limit, offset int) ([]*eribo.LothLog, error) {
	logs := []*eribo.LothLog{}
	const query = `SELECT * FROM loth_logs ORDER BY created DESC LIMIT ? OFFSET ?`
	if err := db.Select(&logs, query, limit, offset); err != nil {
		return nil, err
	}
	return logs, nil
}
