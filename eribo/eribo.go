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
	exp := time.Now().Add(duration).UTC().Truncate(1 * time.Second)
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
	Done      bool
	Kuid      int
	Created   time.Time
	MessageID int64    `db:"message_id"`
	Message   *Message `db:"message"`
}

func (i Image) String() string {
	done := "0"
	if i.Done {
		done = "1"
	}
	player := "unknown player nil message"
	if i.Message != nil {
		player = i.Message.Player
	}
	kuid := ""
	if i.Kuid != 0 {
		kuid = fmt.Sprintf(" [url=https://kusubooru.com/post/view/%d]done[/url]", i.Kuid)
	}
	return fmt.Sprintf("%6d: %v> %s by %s: [url=%s]link[/url]%s",
		i.ID, i.Created.Format(time.Stamp), done, player, i.URL, kuid)
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
	Args    string
	Player  string
	Channel string
	Created time.Time
}

func (l CmdLog) String() string {
	return fmt.Sprintf("%4d: %v by %s - %s %s - %q", l.ID, l.Created.Format(time.Stamp), l.Player, l.Command, l.Args, l.Channel)
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

func (l LothLog) String() string {
	isNew := "0"
	if l.IsNew {
		isNew = "1"
	}
	return fmt.Sprintf("%4d: %v - %v> %s by %s [%s %v %v]",
		l.ID,
		l.Created.Format(time.Stamp),
		l.Expires.Format(time.Stamp),
		isNew,
		l.Issuer,
		l.Name, l.Role, l.Status)
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
	GetImages(limit, offset int, reverse bool) ([]*Image, error)
	ToggleImageDone(id int64) error
	SetImageKuid(id int64, kuid int) error

	AddFeedback(f *Feedback) error
	GetAllFeedback(limit, offset int) ([]*Feedback, error)
	GetRecentFeedback(limit, offset int) ([]*Feedback, error)

	AddCmdLog(e *CmdLog) error
	GetRecentCmdLogs(limit, offset int) ([]*CmdLog, error)

	AddLothLog(*LothLog) error
	GetRecentLothLogs(limit, offset int) ([]*LothLog, error)
}
