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

type Store interface {
	AddMessageWithURLs(m *Message, urls []string) error
	GetImages() ([]*Image, error)
}
