package eribo

import "time"

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
	Command string
	Player  string
	Channel string
	Created time.Time
}

type Store interface {
	AddMessageWithURLs(m *Message, urls []string) error
	GetImages() ([]*Image, error)

	AddFeedback(f *Feedback) error
	GetFeedback() ([]*Feedback, error)

	Log(e *Event) error
}
