package advice

import (
	"encoding/json"
	"net/http"
)

type SlipResp struct {
	Slip struct {
		Advice string `json:"advice"`
		SlipID string `json:"slip_id"`
	} `json:"slip"`
}

func Random() (string, error) {
	resp, err := http.Get("https://api.adviceslip.com/advice")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	s := SlipResp{}
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return "", err
	}
	return s.Slip.Advice, nil
}
