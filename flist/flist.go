package flist

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	clientName    = "eri"
	clientVersion = "0.2.0"
)

type CmdEncoder interface {
	CmdEncode() ([]byte, error)
}

type CmdDecoder interface {
	CmdDecode([]byte) error
}

type Command interface {
	CmdDecoder
	CmdEncoder
}

func cmdEncode(name string, body interface{}) ([]byte, error) {
	if body == nil {
		return []byte(name), nil
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	cmdType := []byte(name + " ")
	return append(cmdType, payload...), nil
}

func cmdDecode(data []byte, v interface{}) error {
	return json.Unmarshal(data[3:], v)
}

type PIN struct{}

func (c *PIN) CmdEncode() ([]byte, error)  { return cmdEncode("PIN", nil) }
func (c *PIN) CmdDecode(data []byte) error { return nil }

type IDN struct {
	Method        string `json:"method,omitempty"`
	Account       string `json:"account,omitempty"`
	Ticket        string `json:"ticket,omitempty"`
	Character     string `json:"character"`
	ClientName    string `json:"cname,omitempty"`
	ClientVersion string `json:"cversion,omitempty"`
}

func (c *Client) NewIDN(account, ticket, character string) *IDN {
	return &IDN{
		Method:        "ticket",
		Account:       account,
		Ticket:        ticket,
		Character:     character,
		ClientName:    c.Name,
		ClientVersion: c.Version,
	}
}

func (c *IDN) CmdEncode() ([]byte, error)  { return cmdEncode("IDN", c) }
func (c *IDN) CmdDecode(data []byte) error { return cmdDecode(data, c) }

type Channel struct {
	Name       string `json:"name"`
	Title      string `json:"title"`
	Characters int    `json:"characters"`
}

type byTitle []Channel

func (channels byTitle) Len() int           { return len(channels) }
func (channels byTitle) Swap(i, j int)      { channels[i], channels[j] = channels[j], channels[i] }
func (channels byTitle) Less(i, j int) bool { return channels[i].Title < channels[j].Title }

func SortChannelsByTitle(channels []Channel) {
	sort.Sort(byTitle(channels))
}

func FindChannel(channels []Channel, title string) *Channel {
	i := sort.Search(len(channels), func(i int) bool { return channels[i].Title >= title })
	if i < len(channels) && channels[i].Title == title {
		return &channels[i]
	}
	return nil
}

type ORS struct {
	Channels []Channel `json:"channels,omitempty"`
}

func (c *ORS) CmdEncode() ([]byte, error)  { return cmdEncode("ORS", c) }
func (c *ORS) CmdDecode(data []byte) error { return cmdDecode(data, c) }

type JCH struct {
	Channel string `json:"channel"`
}

func (c *JCH) CmdEncode() ([]byte, error)  { return cmdEncode("JCH", c) }
func (c *JCH) CmdDecode(data []byte) error { return cmdDecode(data, c) }

type MSG struct {
	Character string `json:"character,omitempty"`
	Message   string `json:"message"`
	Channel   string `json:"channel"`
}

func (m *MSG) CmdEncode() ([]byte, error)  { return cmdEncode("MSG", m) }
func (m *MSG) CmdDecode(data []byte) error { return json.Unmarshal(data[3:], m) }

type PRI struct {
	Character string `json:"character,omitempty"`
	Message   string `json:"message"`
	Recipient string `json:"recipient,omitempty"`
}

func (m *PRI) CmdEncode() ([]byte, error)  { return cmdEncode("PRI", m) }
func (m *PRI) CmdDecode(data []byte) error { return json.Unmarshal(data[3:], m) }

type Client struct {
	ws      *websocket.Conn
	Name    string
	Version string
}

func (c *Client) Close() error {
	return c.ws.Close()
}

func (c *Client) Disconnect() error {
	return c.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (c *Client) SendPIN() error {
	return c.ws.WriteMessage(websocket.TextMessage, []byte("PIN"))
}

func (c *Client) ReadMessage() ([]byte, error) {
	_, message, err := c.ws.ReadMessage()
	return message, err
}

type RawMessage struct {
	Data []byte
	Err  error
}

func Connect(url string) (*Client, error) {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}
	return &Client{ws: ws, Name: clientName, Version: clientVersion}, nil
}

func isCmd(data []byte, cmdType string) bool {
	return bytes.HasPrefix(data, []byte(cmdType))
}

var ErrUnknownCmd = errors.New("unknown command")

func DecodeCommand(data []byte) (Command, error) {
	switch {
	case isCmd(data, "IDN"):
		idn := new(IDN)
		if err := idn.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("IDN decode: %v", err)
		}
		return idn, nil
	case isCmd(data, "MSG"):
		msg := new(MSG)
		if err := msg.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("MSG decode: %v", err)
		}
		return msg, nil
	case isCmd(data, "PRI"):
		pri := new(PRI)
		if err := pri.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("PRI decode: %v", err)
		}
		return pri, nil
	case isCmd(data, "ORS"):
		ors := new(ORS)
		if err := ors.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("ORS decode: %v", err)
		}
		return ors, nil
	case isCmd(data, "PIN"):
		pin := new(PIN)
		return pin, nil
	default:
		return nil, ErrUnknownCmd
	}
}

func (c *Client) writeMessage(data []byte) error {
	return c.ws.WriteMessage(websocket.TextMessage, data)
}

func (c *Client) SendMSG(msg *MSG) error {
	data, err := msg.CmdEncode()
	if err != nil {
		return fmt.Errorf("MSG encode failed: %v", err)
	}

	if err := c.writeMessage(data); err != nil {
		return fmt.Errorf("SendMSG error: %v", err)
	}
	return nil
}

func (c *Client) SendPRI(pri *PRI) error {
	data, err := pri.CmdEncode()
	if err != nil {
		return fmt.Errorf("PRI encode failed: %v", err)
	}

	if err := c.writeMessage(data); err != nil {
		return fmt.Errorf("SendPRI error: %v", err)
	}
	return nil
}

func (c *Client) SendJCH(jch *JCH) error {
	data, err := jch.CmdEncode()
	if err != nil {
		return fmt.Errorf("JCH encode failed: %v", err)
	}

	if err := c.writeMessage(data); err != nil {
		return fmt.Errorf("SendJCH error: %v", err)
	}
	return nil
}

func (c *Client) SendORS() error {
	return c.writeMessage([]byte("ORS"))
}

func (c *Client) Identify(account, password, character string) error {
	ticket, err := GetTicket(account, password)
	if err != nil {
		return fmt.Errorf("could not get ticket: %v", err)
	}

	idn := c.NewIDN(account, ticket, character)
	data, err := idn.CmdEncode()
	if err != nil {
		return fmt.Errorf("IDN encode failed: %v", err)
	}

	if err := c.ws.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("identify error: %v", err)
	}
	return nil
}

func GetTicket(account, password string) (string, error) {
	u := "https://www.f-list.net/json/getApiTicket.php"

	v := url.Values{}
	v.Add("account", account)
	v.Add("password", password)
	v.Add("no_characters", "true")
	v.Add("no_friends", "true")
	v.Add("no_bookmarks", "true")

	body := strings.NewReader(v.Encode())
	resp, err := http.Post(u, "application/x-www-form-urlencoded", body)
	if err != nil {
		return "", fmt.Errorf("post failed: %v", err)
	}

	type ticket struct {
		Ticket string `json:"ticket"`
		Error  string `json:"error"`
	}

	t := new(ticket)
	if err := json.NewDecoder(resp.Body).Decode(t); err != nil {
		return "", fmt.Errorf("could not decode ticket: %v", err)
	}
	if t.Error != "" {
		return "", fmt.Errorf("ticket contains error: %s", t.Error)
	}
	return t.Ticket, nil
}
