package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"mvdan.cc/xurls"

	_ "github.com/go-sql-driver/mysql"

	"github.com/kusubooru/eribo/eribo"
	"github.com/kusubooru/eribo/eribo/mysql"
	"github.com/kusubooru/eribo/flist"
	"github.com/kusubooru/eribo/rp"
)

func defaultAddr(addr string, testServer, insecure bool) string {
	switch {
	default:
		//log.Printf("Using encrypted production server address: %q", addr)
	case !testServer && insecure:
		addr = "ws://chat.f-list.net:9722"
		log.Printf("Using unencrypted production server address: %q", addr)
	case testServer && !insecure:
		addr = "wss://chat.f-list.net:8799"
		log.Printf("Using encrypted test server address: %q", addr)
	case testServer && insecure:
		addr = "ws://chat.f-list.net:8722"
		log.Printf("Using unencrypted test server address: %q", addr)
	}
	return addr
}

func splitRoomTitles(s string) ([]string, error) {
	var rooms []string
	err := json.Unmarshal([]byte(s), &rooms)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

var theVersion = "devel"

func main() {
	var (
		insecure    = flag.Bool("insecure", false, "use insecure ws:// websocket instead of wss://")
		testServer  = flag.Bool("testserver", false, "connect to test server instead of production")
		addr        = flag.String("addr", "wss://chat.f-list.net:9799", "websocket address to connect")
		account     = flag.String("account", "", "websocket address to connect")
		password    = flag.String("password", "", "websocket address to connect")
		character   = flag.String("character", "", "websocket address to connect")
		dataSource  = flag.String("datasource", "", "MySQL datasource")
		joinRooms   = flag.String("join", "", "open private `rooms` to join in JSON format e.g. "+`-join '["Room 1", "Room 2"]'`)
		showVersion = flag.Bool("v", false, "print program version")
		versionArg  bool
	)
	flag.Parse()

	versionArg = len(os.Args) > 1 && os.Args[1] == "version"
	if *showVersion || versionArg {
		fmt.Printf("%s %s (runtime: %s)\n", os.Args[0], theVersion, runtime.Version())
		return
	}

	roomTitles, err := splitRoomTitles(*joinRooms)
	if err != nil {
		log.Println(`-join [rooms] requires rooms to be in JSON format. Example: -join '["Room 1", "Room 2"]'`)
		log.Fatal("error decoding rooms to join:", err)
	}

	if *dataSource == "" {
		log.Println("Database datasource not provided, exiting...")
		log.Fatal("Use -datasource='username:password@(host:port)/database?parseTime=true'")
	}
	if *account == "" || *password == "" || *character == "" {
		log.Println("Account, password and character name needed for identification.")
		log.Fatal("Use -account=<username> -password=<password> -character=<char name>")
	}
	*addr = defaultAddr(*addr, *testServer, *insecure)

	store, err := mysql.NewEriboStore(*dataSource)
	if err != nil {
		log.Fatal("store error:", err)
	}

	// Connect to F-list.
	c, err := flist.Connect(*addr)
	if err != nil {
		log.Println("connect error:", err)
		return
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Println("close err:", err)
		}
	}()

	// Prepare channels to separate message types.
	idnch := make(chan *flist.IDN)
	msgch := make(chan *flist.MSG, 10)
	prich := make(chan *flist.PRI, 10)
	orsch := make(chan *flist.ORS, 10)
	pinch := make(chan *flist.PIN)
	quit := make(chan struct{})

	// The reader is responsible for closing the channels.
	go readMessages(c, idnch, msgch, prich, orsch, pinch, quit)

	// Login to F-list.
	if err := c.Identify(*account, *password, *character); err != nil {
		log.Println(err)
		return
	}
	// Wait for identification because: "If you send any commands before
	// identifying, you will be disconnected."
	//
	// https://wiki.f-list.net/F-Chat_Server_Commands#IDN
	<-idnch

	// Request open private rooms.
	if err := c.SendORS(); err != nil {
		log.Println(err)
		return
	}

	handleMessages(c, store, roomTitles, idnch, msgch, prich, orsch, pinch, quit)
}

// readMessages or "the reader" reads messages in a loop, sepearates them into
// different F-list server command types and sends them to the appropriate
// channels to be handled.
func readMessages(
	c *flist.Client,
	idnch chan<- *flist.IDN,
	msgch chan<- *flist.MSG,
	prich chan<- *flist.PRI,
	orsch chan<- *flist.ORS,
	pinch chan<- *flist.PIN,
	quit chan struct{},
) {
	defer close(idnch)
	defer close(msgch)
	defer close(prich)
	defer close(orsch)
	defer close(pinch)
	defer close(quit)
	for {
		message, err := c.ReadMessage()
		if err != nil {
			log.Println("read message error:", err)
			return
		}
		cmd, err := flist.DecodeCommand(message)
		if err == flist.ErrUnknownCmd && len(message) != 0 {
			//fmt.Println("got:", string(message))
		}
		if err != nil && err != flist.ErrUnknownCmd {
			log.Println("cmd decode error:", err)
		}
		switch t := cmd.(type) {
		case *flist.IDN:
			idnch <- t
		case *flist.MSG:
			msgch <- t
		case *flist.PRI:
			prich <- t
		case *flist.ORS:
			orsch <- t
		case *flist.PIN:
			pinch <- t
		}
	}
}

// handleMessages receives different command types from the reader's channels
// and responds accordingly. It also listens for the interrupt singal or for
// the reader quitting due to error.
func handleMessages(
	c *flist.Client,
	store eribo.Store,
	roomTitles []string,
	idnch <-chan *flist.IDN,
	msgch <-chan *flist.MSG,
	prich <-chan *flist.PRI,
	orsch <-chan *flist.ORS,
	pinch <-chan *flist.PIN,
	quit <-chan struct{},
) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:
			log.Println("interrupt signal received...")
			if err := c.Disconnect(); err != nil {
				log.Println("disconnect error:", err)
			}
			log.Println("waiting for reader to quit...")
			<-quit
			log.Println("exiting...")
			return
		case <-quit:
			// If the reader quits with an error, there's no point for the
			// program to continue so it exists.
			log.Println("reader quit")
			return
		case msg := <-msgch:
			urls := xurls.Strict.FindAllString(msg.Message, -1)
			if len(urls) != 0 {
				m := &eribo.Message{Channel: msg.Channel, Player: msg.Character, Message: msg.Message}
				if err := store.AddMessageWithURLs(m, urls); err != nil {
					log.Println("error storing message:", err)
				}
			}
			respond(c, store, msg)
		case pri := <-prich:
			if err := gatherFeedback(c, store, pri); err != nil {
				log.Println("gather feedback err:", err)
			}
		case ors := <-orsch:
			flist.SortChannelsByTitle(ors.Channels)
			for _, title := range roomTitles {
				ch := flist.FindChannel(ors.Channels, title)
				if ch != nil {
					jch := &flist.JCH{Channel: ch.Name}
					if err := c.SendJCH(jch); err != nil {
						log.Println("error joining private room %q: %v", title, err)
					}
				}
			}
		case <-pinch:
			if err := c.SendPIN(); err != nil {
				log.Println("send PIN failed:", err)
			}
		case idn := <-idnch:
			// Expecting IDN only once during identification.
			log.Println("received IDN but shouldn't:", idn)
		}
	}
}

func respond(c *flist.Client, store eribo.Store, m *flist.MSG) {
	switch {
	case strings.HasPrefix(m.Message, eribo.CmdTieup.String()):
		resp := &flist.MSG{
			Channel: m.Channel,
			Message: rp.RandTieUp(m.Character),
		}
		if err := c.SendMSG(resp); err != nil {
			log.Printf("error sending %v response: %v", eribo.CmdTieup, err)
		}
		e := &eribo.Event{Command: eribo.CmdTieup, Player: m.Character, Channel: m.Channel}
		if err := store.Log(e); err != nil {
			log.Printf("error logging %v: %v", eribo.CmdTieup, err)
		}
	case strings.HasPrefix(m.Message, eribo.CmdTomato.String()):
		resp := &flist.MSG{
			Channel: m.Channel,
			Message: rp.Tomato(m.Character),
		}
		if err := c.SendMSG(resp); err != nil {
			log.Printf("error sending %v response: %v", eribo.CmdTomato, err)
		}
		e := &eribo.Event{Command: eribo.CmdTomato, Player: m.Character, Channel: m.Channel}
		if err := store.Log(e); err != nil {
			log.Printf("error logging %v: %v", eribo.CmdTomato, err)
		}
	case strings.HasPrefix(m.Message, eribo.CmdTktool.String()):
		resp := &flist.MSG{
			Channel: m.Channel,
			Message: rp.RandTktool(m.Character),
		}
		if err := c.SendMSG(resp); err != nil {
			log.Printf("error sending %v response: %v", eribo.CmdTktool, err)
		}
		e := &eribo.Event{Command: eribo.CmdTktool, Player: m.Character, Channel: m.Channel}
		if err := store.Log(e); err != nil {
			log.Printf("error logging %v: %v", eribo.CmdTktool, err)
		}
	case strings.HasPrefix(m.Message, eribo.CmdVonprove.String()):
		resp := &flist.MSG{
			Channel: m.Channel,
			Message: rp.RandVonprove(m.Character),
		}
		if err := c.SendMSG(resp); err != nil {
			log.Printf("error sending %v response: %v", eribo.CmdVonprove, err)
		}
		e := &eribo.Event{Command: eribo.CmdVonprove, Player: m.Character, Channel: m.Channel}
		if err := store.Log(e); err != nil {
			log.Printf("error logging %v: %v", eribo.CmdVonprove, err)
		}
	}
}

func gatherFeedback(c *flist.Client, store eribo.Store, pri *flist.PRI) error {
	if !strings.HasPrefix(pri.Message, eribo.CmdFeedback.String()+" ") {
		return nil
	}
	message := strings.TrimPrefix(pri.Message, eribo.CmdFeedback.String()+" ")
	message = strings.TrimSpace(message)
	if message == "" {
		return nil
	}
	f := &eribo.Feedback{
		Player:  pri.Character,
		Message: message,
	}
	if err := store.AddFeedback(f); err != nil {
		return fmt.Errorf("error storing feedback: %v", err)
	}
	resp := &flist.PRI{
		Recipient: pri.Character,
		Message:   rp.RandFeedback(pri.Character),
	}
	if err := c.SendPRI(resp); err != nil {
		return fmt.Errorf("error sending %v response: %v", eribo.CmdFeedback, err)
	}
	e := &eribo.Event{Command: eribo.CmdFeedback, Player: pri.Character}
	if err := store.Log(e); err != nil {
		log.Printf("error logging %v: %v", eribo.CmdFeedback, err)
	}
	return nil
}
