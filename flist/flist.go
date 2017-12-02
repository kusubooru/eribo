package flist

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	// ErrMsgTooLong is returned if there is an attempt to send a message
	// through MSG or PRI that exceeds the server's variables (chat_max and
	// priv_max respectively). The message is never send to the server. If the
	// message was sent, the server would reply with an ERR.
	ErrMsgTooLong = errors.New("message too long")
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

// LIS is a server command.
//
// Sends an array of all the online characters and their gender, status, and
// status message.
//
// Syntax
//
//     >> LIS { characters: [object] }
//
// Raw sample
//
//     LIS {"characters": [["Alexandrea", "Female", "online", ""], ["Fa Mulan",
//     "Female", "busy", "Away, check out my new alt Aya Kinjou!"], ["Adorkable
//     Lexi", "Female", "online", ""], ["Melfice Cyrum", "Male", "online", ""],
//     ["Jenasys Stryphe", "Female", "online", ""], ["Cassie Hazel", "Herm",
//     "looking", ""], ["Jun Watarase", "Male", "looking", "cute femmy boi
//     looking for a dominate partner"],["Motley Ferret", "Male", "online",
//     ""], ["Tashi", "Male", "online", ""], ["Viol", "Cunt-boy", "looking",
//     ""], ["Dorjan Kazyanenko", "Male", "looking", ""], ["Asaki", "Female",
//     "online", ""]]}
//
// Because of the large amount of data, this command is often sent out in
// batches of several LIS commands. Since you got a CON before LIS, you'll know
// when it has sent them all.
//
// The characters object has a syntax of ["Name", "Gender", "Status", "Status
// Message"].
type LIS struct {
	Characters [][]string `json:"characters"`
}

func (c *LIS) CmdEncode() ([]byte, error)  { return cmdEncode("LIS", c) }
func (c *LIS) CmdDecode(data []byte) error { return cmdDecode(data, c) }

// FLN is a server command.
//
// Sent by the server to inform the client a given character went offline.
//
// Syntax
//
//     >> FLN { "character": string }
//
// Raw sample
//
//     FLN {"character":"Hexxy"}
//
// Notes/Warnings
//
// Should be treated as a global LCH for this character.
type FLN struct {
	Character string `json:"character"`
}

func (c *FLN) CmdEncode() ([]byte, error)  { return cmdEncode("FLN", c) }
func (c *FLN) CmdDecode(data []byte) error { return cmdDecode(data, c) }

// NLN is a server command.
//
// A user connected.
//
// Syntax
//
//     >> NLN { "identity": string, "gender": enum, "status": enum }
//
// Raw sample
//
//     NLN {"status": "online", "gender": "Male", "identity": "Hexxy"}
//
// Parameters
//
// Identity: character name of the user connecting.
//
// Gender: a valid gender string.
//
// Status: a valid status, though since it is when signing on, the only
// possibility is online.
type NLN struct {
	Identity string `json:"identity"`
	Gender   string `json:"gender"`
	Status   Status `json:"status"`
}

func (c *NLN) CmdEncode() ([]byte, error)  { return cmdEncode("NLN", c) }
func (c *NLN) CmdDecode(data []byte) error { return cmdDecode(data, c) }

// ERR is a server command.
//
// Indicates that the given error has occurred.
//
// Syntax
//
//     >> ERR { "number": int, "message": string }
//
// Raw sample
//
//     ERR {"message": "You have already joined this channel.", "number": 28}
type ERR struct {
	Number  int    `json:"number"`
	Message string `json:"message"`
}

func (c *ERR) CmdEncode() ([]byte, error)  { return cmdEncode("ERR", c) }
func (c *ERR) CmdDecode(data []byte) error { return cmdDecode(data, c) }

type ICH struct {
	Channel string `json:"channel"`
	Mode    string `json:"mode"` // enum, can be "ads", "chat", or "both".
	Users   []struct {
		Identity string `json:"identity"`
	} `json:"users"`
}

func (c *ICH) CmdEncode() ([]byte, error)  { return cmdEncode("ICH", c) }
func (c *ICH) CmdDecode(data []byte) error { return cmdDecode(data, c) }

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

type PRO struct {
	Character string `json:"character"`
}

func (c *PRO) CmdEncode() ([]byte, error)  { return cmdEncode("PRO", c) }
func (c *PRO) CmdDecode(data []byte) error { return cmdDecode(data, c) }

type PRDType string

const (
	PRDStart  PRDType = "start"
	PRDEnd            = "end"
	PRDInfo           = "info"
	PRDSelect         = "select"
)

type PRD struct {
	Type    PRDType `json:"type"`
	Message string  `json:"message,omitempty"`
	Key     string  `json:"key"`
	Value   string  `json:"value"`
}

func (c *PRD) CmdEncode() ([]byte, error)  { return cmdEncode("PRD", c) }
func (c *PRD) CmdDecode(data []byte) error { return cmdDecode(data, c) }

type Role string

const (
	RoleFullDom = "Always dominant"
	RoleSomeDom = "Usually dominant"
	RoleSwitch  = "Switch"
	RoleSomeSub = "Usually submissive"
	RoleFullSub = "Always submissive"
)

type Status string

func (s Status) IsActive() bool {
	switch s {
	case StatusOnline:
		return true
	case StatusLooking:
		return true
	case StatusBusy:
		return false
	case StatusDND:
		return false
	case StatusIdle:
		return false
	case StatusAway:
		return false
	default:
		return false
	}
}

const (
	StatusOnline  = "online"
	StatusLooking = "looking"
	StatusBusy    = "busy"
	StatusDND     = "dnd"
	StatusIdle    = "idle"
	StatusAway    = "away"
)

type STA struct {
	Status    Status `json:"status"`
	StatusMsg string `json:"statusmsg"`
	Character string `json:"character,omitempty"`
}

func (c *STA) CmdEncode() ([]byte, error)  { return cmdEncode("STA", c) }
func (c *STA) CmdDecode(data []byte) error { return cmdDecode(data, c) }

// JCH is a server and client command.
//
// Server
//
// Indicates the given user has joined the given channel. This may also be the
// client's character.
//
// Syntax
//
//     >> JCH { "channel": string, "character": object, "title": string }
//
// Raw sample
//
//     JCH {"character": {"identity": "Hexxy"}, "channel": "Frontpage",
//     "title": "Frontpage"}
//
// Notes/Warnings
//
// As with all commands that refer to a specific channel, official/public
// channels use the name, but unofficial/private/open private rooms use the
// channel ID, which can be gotten from ORS.
//
// Client
//
// Send a channel join request.
//
// Syntax
//
//     << JCH { "channel": string }
//
// Raw sample
//
//     JCH {"channel": "Frontpage"}
//
// Notes/Warnings
//
// As with all commands that refer to a specific channel, official/public
// channels use the name, but unofficial/private/open private rooms use the
// channel ID, which can be gotten from ORS.
type JCH struct {
	Channel   string `json:"channel"`
	Title     string `json:"title,omitempty"`
	Character struct {
		Identity string `json:"identity"`
	} `json:"character,omitempty"`
}

func (c *JCH) CmdEncode() ([]byte, error)  { return cmdEncode("JCH", c) }
func (c *JCH) CmdDecode(data []byte) error { return cmdDecode(data, c) }

// LCH is a server command.
//
// An indicator that the given character has left the channel. This may also be
// the client's character.
//
// Syntax
//
//     >> LCH { "channel": string, "character": character }
type LCH struct {
	Channel   string `json:"channel"`
	Character string `json:"character"`
}

func (c LCH) CmdEncode() ([]byte, error)   { return cmdEncode("LCH", c) }
func (c *LCH) CmdDecode(data []byte) error { return cmdDecode(data, c) }

type VAR struct {
	Variable string          `json:"variable"`
	Value    json.RawMessage `json:"value"`
	ChatMax  int
	PrivMax  int
}

func (c VAR) CmdEncode() ([]byte, error) { return cmdEncode("VAR", c) }
func (c *VAR) CmdDecode(data []byte) error {
	if err := cmdDecode(data, c); err != nil {
		return err
	}
	switch c.Variable {
	case "chat_max":
		var chatMax int
		if err := json.Unmarshal(c.Value, &chatMax); err != nil {
			return err
		}
		c.ChatMax = chatMax
	case "priv_max":
		var privMax int
		if err := json.Unmarshal(c.Value, &privMax); err != nil {
			return err
		}
		c.PrivMax = privMax
	}
	return nil
}

type Client struct {
	ws      *websocket.Conn
	Name    string
	Version string
	mu      sync.Mutex
	chatMax int
	privMax int
}

func (c *Client) SetChatMax(max int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.chatMax = max
}

func (c *Client) SetPrivMax(max int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.privMax = max
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
	dialer := websocket.DefaultDialer
	// Sending more than 4096 bytes (which is the default) causes a silent
	// disconnect. Increasing the WriteBuffer fixes the issue.
	//
	// See: https://github.com/gorilla/websocket/issues/245
	dialer.WriteBufferSize = 52000
	ws, _, err := dialer.Dial(url, nil)
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
	case isCmd(data, "LIS"):
		lis := new(LIS)
		if err := lis.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("LIS decode: %v", err)
		}
		return lis, nil
	case isCmd(data, "NLN"):
		nln := new(NLN)
		if err := nln.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("NLN decode: %v", err)
		}
		return nln, nil
	case isCmd(data, "FLN"):
		fln := new(FLN)
		if err := fln.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("FLN decode: %v", err)
		}
		return fln, nil
	case isCmd(data, "ICH"):
		ich := new(ICH)
		if err := ich.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("ICH decode: %v", err)
		}
		return ich, nil
	case isCmd(data, "PRD"):
		prd := new(PRD)
		if err := prd.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("PRD decode: %v", err)
		}
		return prd, nil
	case isCmd(data, "STA"):
		sta := new(STA)
		if err := sta.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("STA decode: %v", err)
		}
		return sta, nil
	case isCmd(data, "JCH"):
		jch := new(JCH)
		if err := jch.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("JCH decode: %v", err)
		}
		return jch, nil
	case isCmd(data, "LCH"):
		lch := new(LCH)
		if err := lch.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("LCH decode: %v", err)
		}
		return lch, nil
	case isCmd(data, "VAR"):
		v := new(VAR)
		if err := v.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("VAR decode: %v", err)
		}
		return v, nil
	case isCmd(data, "ERR"):
		err := new(ERR)
		if err := err.CmdDecode(data); err != nil {
			return nil, fmt.Errorf("ERR decode: %v", err)
		}
		return err, nil
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
	if c.chatMax != 0 && len(data) > c.chatMax {
		return ErrMsgTooLong
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
	if c.privMax != 0 && len(data) > c.privMax {
		return ErrMsgTooLong
	}

	if err := c.writeMessage(data); err != nil {
		return fmt.Errorf("SendPRI error: %v", err)
	}
	return nil
}

func (c *Client) SendPRO(pro *PRO) error {
	data, err := pro.CmdEncode()
	if err != nil {
		return fmt.Errorf("PRO encode failed: %v", err)
	}

	if err := c.writeMessage(data); err != nil {
		return fmt.Errorf("SendPRO error: %v", err)
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

func (c *Client) SendSTA(sta *STA) error {
	data, err := sta.CmdEncode()
	if err != nil {
		return fmt.Errorf("STA encode failed: %v", err)
	}

	if err := c.writeMessage(data); err != nil {
		return fmt.Errorf("SendSTA error: %v", err)
	}
	return nil
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

type CharacterData struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Views        int      `json:"views"`
	CustomsFirst bool     `json:"customs_first"`
	CustomTitle  string   `json:"custom_title"`
	CreatedAt    int64    `json:"created_at"`
	UpdatedAt    int64    `json:"updated_at"`
	Infotags     Infotags `json:"infotags"`
}

func (d CharacterData) HumanInfotags(ml *MappingList) map[string]string {
	infotags := ml.InfotagsMap()
	listitems := ml.ListitemsMap()
	m := make(map[string]string)
	for k, v := range d.Infotags {
		m[infotags[k]] = listitems[v]
	}
	return m
}

type Infotags map[string]string

func (it *Infotags) UnmarshalJSON(data []byte) error {
	if bytes.HasPrefix(data, []byte("{")) {
		m := new(map[string]string)
		if err := json.Unmarshal(data, m); err != nil {
			return err
		}
		*it = make(map[string]string)
		for k, v := range *m {
			(*it)[k] = v
		}
		return nil
	}
	return nil
}

func (it Infotags) MarshalJSON() ([]byte, error) {
	return json.Marshal(it)
}

func GetCharacterData(name, account, ticket string) (*CharacterData, error) {
	u := "https://www.f-list.net/json/api/character-data.php"

	v := url.Values{}
	v.Add("name", name)
	v.Add("account", account)
	v.Add("ticket", ticket)

	body := strings.NewReader(v.Encode())
	resp, err := http.Post(u, "application/x-www-form-urlencoded", body)
	if err != nil {
		return nil, fmt.Errorf("post character data failed: %v", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)
	}
	defer resp.Body.Close()
	r := bytes.NewReader(b)

	d := new(CharacterData)
	if err := json.NewDecoder(r).Decode(d); err != nil {
		return nil, ErrorResponse{"could not decode character data", err, b}
	}
	return d, nil
}

type ErrorResponse struct {
	Message string
	Cause   error
	body    []byte
}

func (e ErrorResponse) Body() []byte {
	return e.body
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

type MappingList struct {
	Kinks []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		GroupID     string `json:"group_id"`
	} `json:"kinks"`
	KinkGroups []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"kink_groups"`
	Infotags []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		List    string `json:"list"`
		GroupID string `json:"group_id"`
	} `json:"infotags"`
	InfotagGroups []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"infotag_groups"`
	Listitems []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"listitems"`
	Error string `json:"error"`
}

func (ml MappingList) InfotagsMap() map[string]string {
	m := make(map[string]string)
	for _, it := range ml.Infotags {
		m[it.ID] = it.Name
	}
	return m
}

func (ml MappingList) ListitemsMap() map[string]string {
	m := make(map[string]string)
	for _, li := range ml.Listitems {
		m[li.ID] = li.Value
	}
	return m
}

func GetMappingList() (*MappingList, error) {
	u := "https://www.f-list.net/json/api/mapping-list.php"

	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("get mapping list failed: %v", err)
	}

	d := new(MappingList)
	if err := json.NewDecoder(resp.Body).Decode(d); err != nil {
		return nil, fmt.Errorf("could not decode mapping list: %v", err)
	}
	return d, nil
}

func GetAccountCharacters(account, ticket string) ([]string, error) {
	u := "https://www.f-list.net/json/api/character-list.php"

	v := url.Values{}
	v.Add("account", account)
	v.Add("ticket", ticket)

	body := strings.NewReader(v.Encode())
	resp, err := http.Post(u, "application/x-www-form-urlencoded", body)
	if err != nil {
		return nil, fmt.Errorf("post character data failed: %v", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	type AccountCharacters struct {
		Characters []string
		Error      string
	}
	ac := new(AccountCharacters)
	if err := json.NewDecoder(resp.Body).Decode(ac); err != nil {
		return nil, fmt.Errorf("could not decode account characters: %v", err)
	}
	chars := ac.Characters
	if ac.Error != "" {
		return chars, fmt.Errorf(ac.Error)
	}
	return chars, nil
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("%v %v: %d %s",
			r.Request.Method, r.Request.URL,
			r.StatusCode, fmt.Sprintf("error reading response body: %v", err))
	}
	return fmt.Errorf("%v %v: %d %s",
		r.Request.Method, r.Request.URL,
		r.StatusCode, string(data))
}
