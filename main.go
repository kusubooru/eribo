package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"mvdan.cc/xurls"
)

var addr = flag.String("addr", "chat.f-list.net:9799", "http service address") // production server
//var addr = flag.String("addr", "chat.f-list.net:9722", "http service address") // production server unencrypted
//var addr = flag.String("addr", "chat.f-list.net:8799", "http service address") // test server
//var addr = flag.String("addr", "chat.f-list.net:8722", "http service address") // unencrypted test server

const (
	account       = ""
	password      = ""
	clientName    = "testbot2"
	clientVersion = "0.0.1"
	targetChannel = "ADH-c63dd350865f6eb33043"
)

type IDN struct {
	Method        string `json:"method"`
	Account       string `json:"account"`
	Ticket        string `json:"ticket"`
	Character     string `json:"character"`
	ClientName    string `json:"cname"`
	ClientVersion string `json:"cversion"`
}

func NewIDN(account, ticket, character string) *IDN {
	return &IDN{
		Method:        "ticket",
		Account:       account,
		Ticket:        ticket,
		Character:     character,
		ClientName:    clientName,
		ClientVersion: clientVersion,
	}
}

func (c IDN) Command() string {
	return fmt.Sprintf(
		"IDN { \"method\": \"ticket\", \"account\": %q, \"ticket\": %q, \"character\": %q, \"cname\": %q, \"cversion\": %q }",
		c.Account, c.Ticket, c.Character, c.ClientName, c.ClientVersion,
	)
}

type channels struct {
	Channels []channel `json:"channels"`
}

type channel struct {
	Name       string `json:"name"`
	Title      string `json:"title"`
	Characters int    `json:"characters"`
}

type message struct {
	Character string `json:"character"`
	Message   string `json:"message"`
	Channel   string `json:"channel"`
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	ticket, err := getTicket()
	if err != nil {
		log.Fatalf("could not get ticket: %v", err)
	}
	fmt.Println("ticket:", ticket)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	u := url.URL{Scheme: "wss", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws.Close()

	msgch := make(chan message, 10)
	go func(msgch chan message) {
		defer close(msgch)
		for {
			select {
			default:
			case msg := <-msgch:
				urls := xurls.Strict.FindAllString(msg.Message, -1)
				for _, u := range urls {
					fmt.Println("found url:", u)
				}
			}
		}
	}(msgch)

	done := make(chan struct{})
	go func() {
		defer ws.Close()
		defer close(done)
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			switch {
			default:
				log.Printf("recv: %s", message)
			case bytes.Equal(message, []byte("PIN")):
				if err := ws.WriteMessage(websocket.TextMessage, []byte("PIN")); err != nil {
					log.Println("write PIN failed:", err)
				}
			case bytes.HasPrefix(message, []byte("ORS")):
				chans := new(channels)
				if err := json.Unmarshal(message[3:], chans); err != nil {
					log.Println("Error decoding channels:", err)
				}
				for _, c := range chans.Channels {
					if strings.Contains(strings.ToLower(c.Title), "tickling") {
						fmt.Println("found:", c)
					}
				}
			case bytes.HasPrefix(message, []byte("MSG")):
				msg := new(message)
				if err := json.Unmarshal(message[3:], msg); err != nil {
					log.Println("Error decoding message:", err)
				}
				msgch <- msg
			}
		}
	}()

	idn := NewIDN(account, ticket, "testbot2")
	idncmd := idn.Command()
	if err := ws.WriteMessage(websocket.TextMessage, []byte(idncmd)); err != nil {
		log.Println("write message failed:", err)
		return
	}

	if err := ws.WriteMessage(websocket.TextMessage, []byte("ORS")); err != nil {
		log.Println("write message failed:", err)
		return
	}

	//if err := c.WriteMessage(websocket.TextMessage, []byte("CHA")); err != nil {
	//	log.Println("write message failed:", err)
	//	return
	//}

	//channelName := "Politics"
	//joinChannelCmd := fmt.Sprintf("JCH { \"channel\": %q }", channelName)
	//if err := ws.WriteMessage(websocket.TextMessage, []byte(joinChannelCmd)); err != nil {
	//	log.Println("write message failed:", err)
	//	return
	//}

	//time.Sleep(5 * time.Second)
	//msgCmd := fmt.Sprintf("MSG { \"channel\": %q, \"message\": %q }", channelName, "hello!")
	//fmt.Println("sending msg:", msgCmd)
	//if err := c.WriteMessage(websocket.TextMessage, []byte(msgCmd)); err != nil {
	//	log.Println("write message failed:", err)
	//	return
	//}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			time.Sleep(time.Second)
			//err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			//if err != nil {
			//	log.Println("write:", err)
			//	return
			//}
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			ws.Close()
			return
		}
	}
}

func getTicket() (string, error) {
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
