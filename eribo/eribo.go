package eribo

import (
	"time"
)

type Loth struct {
	*Player
	expires time.Time
}

func NewLoth(p *Player, duration time.Duration) *Loth {
	exp := time.Now().Add(duration)
	return &Loth{Player: p, expires: exp}
}

func (l Loth) TimeLeft() string {
	d := time.Until(time.Time(l.expires))
	rounded := time.Duration(d.Nanoseconds()/time.Second.Nanoseconds()) * time.Second
	return rounded.String()
}

func (l Loth) Expired() bool {
	return time.Now().After(time.Time(l.expires))
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

type Event struct {
	ID      int64
	Command Command
	Player  string
	Channel string
	Created time.Time
}

type Store interface {
	AddMessageWithURLs(m *Message, urls []string) error
	GetImages() ([]*Image, error)

	AddFeedback(f *Feedback) error
	GetAllFeedback(limit, offset int) ([]*Feedback, error)
	GetRecentFeedback(limit, offset int) ([]*Feedback, error)

	Log(e *Event) error
	GetRecentLogs(limit, offset int) ([]*Event, error)
}
