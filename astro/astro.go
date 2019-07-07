package astro

import (
	"encoding/json"
	"fmt"
	"net/http"
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

var signs = []Sign{
	Aries,
	Taurus,
	Gemini,
	Cancer,
	Leo,
	Virgo,
	Libra,
	Scorpio,
	Sagittarius,
	Capricorn,
	Aquarius,
	Pisces,
}

func validSign(s Sign) bool {
	for _, sign := range signs {
		if s == sign {
			return true
		}
	}
	return false
}

// For returns horoscope for a certain period and sign.
func For(period string, sign Sign) (string, error) {
	if period == "" {
		period = "today"
	}
	if period != "today" && period != "week" && period != "month" && period != "year" {
		period = "today"
	}
	if !validSign(sign) {
		return fmt.Sprintf("Valid signs are: %v", signs), nil
	}
	u := fmt.Sprintf("http://horoscope-api.herokuapp.com/horoscope/%s/%s", period, sign)
	resp, err := http.Get(u)
	if err != nil {
		return "", fmt.Errorf("GET %s:%v", u, err)
	}

	type response struct {
		Date      string `json:"date"`
		Week      string `json:"week"`
		Month     string `json:"month"`
		Year      string `json:"year"`
		Horoscope string `json:"horoscope"`
		Sunsign   string `json:"sunsign"`
	}
	v := new(response)
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return "", fmt.Errorf("decoding horoscope: %v", err)
	}
	return v.Horoscope, nil
}
