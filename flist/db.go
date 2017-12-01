package flist

import (
	"database/sql/driver"
	"fmt"
)

func (r Role) Value() (driver.Value, error) { return string(r), nil }
func (r *Role) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("attempt to scan nil role")
	}
	switch v := value.(type) {
	case string:
		*r = Role(v)
		return nil
	case []byte:
		*r = Role(string(v))
		return nil
	}
	return fmt.Errorf("cannot scan Role value")
}

func (s Status) Value() (driver.Value, error) { return string(s), nil }
func (s *Status) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("attempt to scan nil status")
	}
	switch v := value.(type) {
	case string:
		*s = Status(v)
		return nil
	case []byte:
		*s = Status(string(v))
		return nil
	}
	return fmt.Errorf("cannot scan Status value")
}
