package astro

import (
	"database/sql/driver"
	"fmt"
)

// Sign is a zodiac sign.
type Sign string

// All zodiac signs.
const (
	Aries       Sign = "aries"
	Taurus      Sign = "taurus"
	Gemini      Sign = "gemini"
	Cancer      Sign = "cancer"
	Leo         Sign = "leo"
	Virgo       Sign = "virgo"
	Libra       Sign = "libra"
	Scorpio     Sign = "scorpio"
	Sagittarius Sign = "sagittarius"
	Capricorn   Sign = "capricorn"
	Aquarius    Sign = "aquarius"
	Pisces      Sign = "pisces"
)

// Value tells the database driver how to store a Sign.
func (s Sign) Value() (driver.Value, error) { return string(s), nil }

// Scan tells the database driver how to scan the stored value into a Sign.
func (s *Sign) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("astro: cannot scan nil as Sign")
	}

	sv := fmt.Sprintf("%v", value)
	*s = Sign(sv)

	return nil
}
