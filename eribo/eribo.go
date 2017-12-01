package eribo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Loth struct {
	*Player
	Expires time.Time
}

func NewLoth(p *Player, duration time.Duration) *Loth {
	exp := time.Now().Add(duration).UTC().Truncate(time.Microsecond)
	return &Loth{Player: p, Expires: exp}
}

func (l Loth) TimeLeft() string {
	d := time.Until(time.Time(l.Expires))
	rounded := time.Duration(d.Nanoseconds()/time.Second.Nanoseconds()) * time.Second
	return rounded.String()
}

func (l Loth) Expired() bool {
	return time.Now().After(time.Time(l.Expires))
}

// Store

type Image struct {
	ID        int64
	URL       string
	Created   time.Time
	MessageID int64    `db:"message_id"`
	Message   *Message `db:"message"`
}

type Message struct {
	ID      int64
	Message string
	Player  string
	Channel string
	Created time.Time
}

type Feedback struct {
	ID      int64
	Message string
	Player  string
	Created time.Time
}

func (f Feedback) String() string {
	return fmt.Sprintf("%4d: %v by %s - %q", f.ID, f.Created.Format(time.Stamp), f.Player, f.Message)
}

type CmdLog struct {
	ID      int64
	Command Command
	Player  string
	Channel string
	Created time.Time
}

func (l CmdLog) String() string {
	return fmt.Sprintf("%4d: %v by %s - %s - %q", l.ID, l.Created.Format(time.Stamp), l.Player, l.Command, l.Channel)
}

type LothLog struct {
	ID      int64
	Issuer  string
	Channel string
	Created time.Time
	*Loth
	IsNew   bool `db:"is_new"`
	Targets Targets
}

type Targets []*Player

func (t Targets) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (t *Targets) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("attempt to scan nil targets")
	}
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), t)
	case []byte:
		return json.Unmarshal(v, t)
	}
	return fmt.Errorf("cannot scan targets value")
}

type Store interface {
	AddMessageWithURLs(m *Message, urls []string) error
	GetImages() ([]*Image, error)

	AddFeedback(f *Feedback) error
	GetAllFeedback(limit, offset int) ([]*Feedback, error)
	GetRecentFeedback(limit, offset int) ([]*Feedback, error)

	AddCmdLog(e *CmdLog) error
	GetRecentCmdLogs(limit, offset int) ([]*CmdLog, error)

	AddLothLog(*LothLog) error
	GetRecentLothLogs(limit, offset int) ([]*LothLog, error)
}
