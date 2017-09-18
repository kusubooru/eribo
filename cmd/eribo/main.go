package main

import (
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
		log.Printf("Using encrypted production server address: %q", addr)
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
		joinRoom    = flag.String("join", "", "exact title of an open private room to join")
		showVersion = flag.Bool("v", false, "print program version")
		versionArg  bool
	)
	flag.Parse()

	versionArg = len(os.Args) > 1 && os.Args[1] == "version"
	if *showVersion || versionArg {
		fmt.Printf("%s %s (runtime: %s)\n", os.Args[0], theVersion, runtime.Version())
		return
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

	doneRead := make(chan struct{})
	defer close(doneRead)
	doneHandle := make(chan struct{})
	defer close(doneHandle)

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

	// Prepare channel for messages.
	msgch := make(chan *flist.MSG, 10)
	defer close(msgch)
	prich := make(chan *flist.PRI, 10)
	defer close(prich)
	orsch := make(chan *flist.ORS, 10)
	defer close(orsch)

	go readMessages(c, msgch, prich, orsch, doneRead)
	go handleMessages(c, store, *joinRoom, msgch, prich, orsch, doneHandle)

	// Login to F-list.
	if err := c.Identify(*account, *password, *character); err != nil {
		log.Println(err)
		return
	}

	// Request open private rooms.
	if err := c.SendORS(); err != nil {
		log.Println(err)
		return
	}

	waitForInterrupt(c, doneRead, doneHandle)
}

// waitForInterrupt blocks and waits either for interrupt signal or for the
// client to quit.
func waitForInterrupt(c *flist.Client, doneRead, doneHandle chan struct{}) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:
			log.Println("interrupt signal received, exiting...")
			doneRead <- struct{}{}
			doneHandle <- struct{}{}
			if err := c.Disconnect(); err != nil {
				log.Println("disconnect error:", err)
			}
			return
		case <-c.Quit:
			log.Println("disconnected")
			doneRead <- struct{}{}
			doneHandle <- struct{}{}
			return
		}
	}
}

func readMessages(
	c *flist.Client,
	msgch chan<- *flist.MSG,
	prich chan<- *flist.PRI,
	orsch chan<- *flist.ORS,
	done chan struct{},
) {
	for {
		select {
		case <-done:
			log.Println("done reading")
			return
		case message := <-c.Messenger:
			cmd, err := flist.DecodeCommand(message)
			if err == flist.ErrUnknownCmd && len(message) != 0 {
				//fmt.Println("got:", string(message))
			}
			if err != nil && err != flist.ErrUnknownCmd {
				log.Println("cmd decode error:", err)
			}
			switch t := cmd.(type) {
			case *flist.MSG:
				msgch <- t
			case *flist.PRI:
				prich <- t
			case *flist.ORS:
				orsch <- t
			case *flist.PIN:
				if err := c.SendPIN(); err != nil {
					log.Println("send PIN failed:", err)
				}
			}
		}
	}
}

func handleMessages(
	c *flist.Client,
	store eribo.Store,
	roomTitle string,
	msgch <-chan *flist.MSG,
	prich <-chan *flist.PRI,
	orsch <-chan *flist.ORS,
	done chan struct{},
) {
	for {
		select {
		case <-done:
			log.Println("done handling")
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
			for _, ch := range ors.Channels {
				if ch.Title == roomTitle {
					jch := &flist.JCH{Channel: ch.Name}
					if err := c.SendJCH(jch); err != nil {
						log.Println("error joining private room %q: %v", roomTitle, err)
					}
				}
			}
		}
	}
}

func respond(c *flist.Client, store eribo.Store, m *flist.MSG) {
	switch {
	case strings.Contains(m.Message, "!tieup"):
		resp := &flist.MSG{
			Channel: m.Channel,
			Message: rp.RandTieUp(m.Character),
		}
		if err := c.SendMSG(resp); err != nil {
			log.Println("error sending tieup response:", err)
		}
		e := &eribo.Event{Command: "!tieup", Player: m.Character, Channel: m.Channel}
		if err := store.Log(e); err != nil {
			log.Println("error storing to log:", err)
		}
	case strings.Contains(m.Message, "!tomato"):
		resp := &flist.MSG{
			Channel: m.Channel,
			Message: rp.Tomato(m.Character),
		}
		if err := c.SendMSG(resp); err != nil {
			log.Println("tomato error:", err)
		}
		e := &eribo.Event{Command: "!tomato", Player: m.Character, Channel: m.Channel}
		if err := store.Log(e); err != nil {
			log.Println("error storing to log:", err)
		}
	}
}

func gatherFeedback(c *flist.Client, store eribo.Store, pri *flist.PRI) error {
	if !strings.Contains(pri.Message, "!feedback") {
		return nil
	}
	f := &eribo.Feedback{
		Player:  pri.Character,
		Message: pri.Message,
	}
	if err := store.AddFeedback(f); err != nil {
		return fmt.Errorf("error storing feedback: %v", err)
	}
	resp := &flist.PRI{
		Recipient: pri.Character,
		Message:   rp.RandFeedback(pri.Character),
	}
	if err := c.SendPRI(resp); err != nil {
		return fmt.Errorf("error sending feedback response: %v", err)
	}
	e := &eribo.Event{Command: "!feedback", Player: pri.Character}
	if err := store.Log(e); err != nil {
		return fmt.Errorf("error storing to log: %v", err)
	}
	return nil
}
