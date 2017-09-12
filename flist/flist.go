package flist

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	clientName    = "eri"
	clientVersion = "0.1.0"
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
	Method        string `json:"method"`
	Account       string `json:"account"`
	Ticket        string `json:"ticket"`
	Character     string `json:"character"`
	ClientName    string `json:"cname"`
	ClientVersion string `json:"cversion"`
}

func NewIDN(account, ticket, character, clientName, clientVersion string) *IDN {
	return &IDN{
		Method:        "ticket",
		Account:       account,
		Ticket:        ticket,
		Character:     character,
		ClientName:    clientName,
		ClientVersion: clientVersion,
	}
}

func (c *IDN) CmdEncode() ([]byte, error) {
	return cmdEncode("IDN", c)
}

func (c *IDN) CmdDecode(data []byte) error {
	return cmdDecode(data, c)
}

type channels struct {
	Channels []channel `json:"channels"`
}

type channel struct {
	Name       string `json:"name"`
	Title      string `json:"title"`
	Characters int    `json:"characters"`
}

type MSG struct {
	Character string `json:"character"`
	Message   string `json:"message"`
	Channel   string `json:"channel"`
}

func (m *MSG) CmdEncode() ([]byte, error) {
	return cmdEncode("MSG", m)
}

func (m *MSG) CmdDecode(data []byte) error {
	return json.Unmarshal(data[3:], m)
}

type Client struct {
	mu        sync.Mutex
	ws        *websocket.Conn
	Messenger <-chan []byte
	Quit      <-chan struct{}
	Name      string
	Version   string
}

func (c *Client) Close() error {
	return c.ws.Close()
}

func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (c *Client) SendPIN() error {
	c.mu.Lock()
	defer c.mu.Unlock()
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

	m := make(chan []byte, 100)
	q := make(chan struct{})
	go readMessages(ws, m, q)
	return &Client{ws: ws, Messenger: m, Quit: q, Name: clientName, Version: clientVersion}, nil
}

func readMessages(ws *websocket.Conn, messenger chan<- []byte, quit chan<- struct{}) {
	defer close(messenger)
	defer func() {
		quit <- struct{}{}
	}()
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			return
		}
		messenger <- msg
	}
}

func isCmd(data []byte, cmdType string) bool {
	return bytes.HasPrefix(data, []byte(cmdType))
}

var ErrUnknownCmd = errors.New("unknown command")

func DecodeCommand(data []byte) (Command, error) {
	switch {
	case isCmd(data, "MSG"):
		msg := new(MSG)
		if err := msg.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("MSG decode: %v", err)
		}
		return msg, nil
	case isCmd(data, "PIN"):
		pin := new(PIN)
		return pin, nil
	default:
		return nil, ErrUnknownCmd
	}
}

func (c *Client) Identify(account, password, character string) error {
	ticket, err := GetTicket(account, password)
	if err != nil {
		return fmt.Errorf("could not get ticket: %v", err)
	}

	idn := NewIDN(account, ticket, character, c.Name, c.Version)
	data, err := idn.CmdEncode()
	if err != nil {
		return fmt.Errorf("IDN encode failed: %v", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
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
