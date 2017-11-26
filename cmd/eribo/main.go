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
	lisch := make(chan *flist.LIS, 10)
	flnch := make(chan *flist.FLN, 10)
	nlnch := make(chan *flist.NLN, 10)
	ichch := make(chan *flist.ICH)
	pinch := make(chan *flist.PIN)
	prdch := make(chan *flist.PRD, 100)
	stach := make(chan *flist.STA)
	jchch := make(chan *flist.JCH)
	lchch := make(chan *flist.LCH)
	quit := make(chan struct{})

	// The reader is responsible for closing the channels.
	go readMessages(
		c,
		idnch,
		msgch,
		prich,
		orsch,
		lisch,
		flnch,
		nlnch,
		ichch,
		pinch,
		prdch,
		stach,
		jchch,
		lchch,
		quit,
	)

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

	mappingList, err := flist.GetMappingList()
	if err != nil {
		log.Printf("could not get mapping list: %v", err)
		return
	}

	playerMap := eribo.NewPlayerMap()
	channelMap := eribo.NewChannelMap()

	handleMessages(
		c,
		*account,
		*password,
		mappingList,
		store,
		playerMap,
		channelMap,
		roomTitles,
		idnch,
		msgch,
		prich,
		orsch,
		lisch,
		flnch,
		nlnch,
		ichch,
		pinch,
		prdch,
		stach,
		jchch,
		lchch,
		quit,
	)
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
	lisch chan<- *flist.LIS,
	flnch chan<- *flist.FLN,
	nlnch chan<- *flist.NLN,
	ichch chan<- *flist.ICH,
	pinch chan<- *flist.PIN,
	prdch chan<- *flist.PRD,
	stach chan<- *flist.STA,
	jchch chan<- *flist.JCH,
	lchch chan<- *flist.LCH,
	quit chan struct{},
) {
	defer close(idnch)
	defer close(msgch)
	defer close(prich)
	defer close(orsch)
	defer close(lisch)
	defer close(flnch)
	defer close(nlnch)
	defer close(ichch)
	defer close(pinch)
	defer close(prdch)
	defer close(stach)
	defer close(jchch)
	defer close(lchch)
	defer close(quit)
	for {
		message, err := c.ReadMessage()
		if err != nil {
			log.Println("read message error:", err)
			return
		}
		cmd, err := flist.DecodeCommand(message)
		if err == flist.ErrUnknownCmd && len(message) != 0 {
			fmt.Println("got:", string(message))
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
		case *flist.LIS:
			lisch <- t
		case *flist.FLN:
			flnch <- t
		case *flist.NLN:
			nlnch <- t
		case *flist.ICH:
			ichch <- t
		case *flist.PIN:
			pinch <- t
		case *flist.PRD:
			prdch <- t
		case *flist.STA:
			stach <- t
		case *flist.JCH:
			jchch <- t
		case *flist.LCH:
			lchch <- t
		case *flist.ERR:
			log.Println(fmt.Errorf("Error %d: %s", t.Number, t.Message))
		}
	}
}

// handleMessages receives different command types from the reader's channels
// and responds accordingly. It also listens for the interrupt signal or for
// the reader quitting due to error.
func handleMessages(
	c *flist.Client,
	account string,
	password string,
	mappingList *flist.MappingList,
	store eribo.Store,
	playerMap *eribo.PlayerMap,
	channelMap *eribo.ChannelMap,
	roomTitles []string,
	idnch <-chan *flist.IDN,
	msgch <-chan *flist.MSG,
	prich <-chan *flist.PRI,
	orsch <-chan *flist.ORS,
	lisch <-chan *flist.LIS,
	flnch <-chan *flist.FLN,
	nlnch <-chan *flist.NLN,
	ichch <-chan *flist.ICH,
	pinch <-chan *flist.PIN,
	prdch <-chan *flist.PRD,
	stach <-chan *flist.STA,
	jchch <-chan *flist.JCH,
	lchch <-chan *flist.LCH,
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
		case lis := <-lisch:
			for _, c := range lis.Characters {
				pl := &eribo.Player{Name: c[0], Status: flist.Status(c[2])}
				playerMap.SetPlayer(pl)
			}
		case fln := <-flnch:
			playerMap.DelPlayer(fln.Character)
			channelMap.DelPlayerAllChannels(fln.Character)
		case nln := <-nlnch:
			pl := &eribo.Player{Name: nln.Identity, Status: nln.Status}
			playerMap.SetPlayer(pl)
		case ich := <-ichch:
			ticket, err := flist.GetTicket(account, password)
			if err != nil {
				log.Printf("init channel get ticket error: %v", err)
				return
			}
			for _, u := range ich.Users {
				if p, ok := playerMap.GetPlayer(u.Identity); ok {
					channelMap.SetPlayer(ich.Channel, p)
				}
				actives := channelMap.GetActivePlayers()
				actives.ForEach(func(name string, p *eribo.Player) {
					charData, err := flist.GetCharacterData(name, account, ticket)
					if err != nil {
						log.Printf("init channel could not get character data for %q: %v", name, err)
						return
					}
					m := charData.HumanInfotags(mappingList)
					if role, ok := m["Dom/Sub Role"]; ok {
						p.Role = flist.Role(role)
					}
				})
				//pro := &flist.PRO{Character: u.Identity}
				//if err := c.SendPRO(pro); err != nil {
				//	log.Println("error sending PRO command for identity %q: %v", u.Identity, err)
				//}
				//fmt.Println("sent pro for:", u.Identity)
				//time.Sleep(11 * time.Second)
			}
		case sta := <-stach:
			name := sta.Character
			newStatus := flist.Status(sta.Status)
			player, _ := channelMap.GetPlayer(name)
			if player != nil && player.Role == "" && !player.Status.IsActive() && newStatus.IsActive() {
				fmt.Printf("STA changed to active for char %q\n", name)
				ticket, err := flist.GetTicket(account, password)
				if err != nil {
					log.Printf("STA change get ticket error: %v", err)
					return
				}

				charData, err := flist.GetCharacterData(name, account, ticket)
				if err != nil {
					log.Printf("STA change could not get character data for %q: %v", name, err)
					return
				}
				m := charData.HumanInfotags(mappingList)
				if role, ok := m["Dom/Sub Role"]; ok {
					player.Role = flist.Role(role)
				}
			}
			playerMap.SetPlayerStatus(name, newStatus)
		case jch := <-jchch:
			name := jch.Character.Identity
			player, _ := playerMap.GetPlayer(name)
			if player == nil {
				log.Printf("JCH player %q not found in playerMap", name)
				return
			}
			if player.Role == "" && player.Status.IsActive() {
				fmt.Println("player %q joined, getting char data", name)
				ticket, err := flist.GetTicket(account, password)
				if err != nil {
					log.Printf("JCH get ticket error: %v", err)
					return
				}

				charData, err := flist.GetCharacterData(name, account, ticket)
				if err != nil {
					log.Printf("JCH could not get character data for %q: %v", name, err)
					return
				}
				m := charData.HumanInfotags(mappingList)
				if role, ok := m["Dom/Sub Role"]; ok {
					player.Role = flist.Role(role)
				}
			}
			channelMap.SetPlayer(jch.Channel, player)
		case lch := <-lchch:
			channelMap.DelPlayer(lch.Channel, lch.Character)
		case prd := <-prdch:
			fmt.Println("got prd:", prd)
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
		e := &eribo.Event{Command: eribo.CmdTieup, Player: m.Character, Channel: m.Channel}
		if err := store.Log(e); err != nil {
			log.Printf("error logging %v: %v", eribo.CmdTieup, err)
		}
		if err := c.SendMSG(resp); err != nil {
			log.Printf("error sending %v response: %v", eribo.CmdTieup, err)
		}
	case strings.HasPrefix(m.Message, eribo.CmdTomato.String()):
		resp := &flist.MSG{
			Channel: m.Channel,
			Message: rp.Tomato(m.Character),
		}
		e := &eribo.Event{Command: eribo.CmdTomato, Player: m.Character, Channel: m.Channel}
		if err := store.Log(e); err != nil {
			log.Printf("error logging %v: %v", eribo.CmdTomato, err)
		}
		if err := c.SendMSG(resp); err != nil {
			log.Printf("error sending %v response: %v", eribo.CmdTomato, err)
		}
	case strings.HasPrefix(m.Message, eribo.CmdTktool.String()):
		resp := &flist.MSG{
			Channel: m.Channel,
			Message: rp.RandTktool(m.Character),
		}
		e := &eribo.Event{Command: eribo.CmdTktool, Player: m.Character, Channel: m.Channel}
		if err := store.Log(e); err != nil {
			log.Printf("error logging %v: %v", eribo.CmdTktool, err)
		}
		if err := c.SendMSG(resp); err != nil {
			log.Printf("error sending %v response: %v", eribo.CmdTktool, err)
		}
	case strings.HasPrefix(m.Message, eribo.CmdVonprove.String()):
		resp := &flist.MSG{
			Channel: m.Channel,
			Message: rp.RandVonprove(m.Character),
		}
		e := &eribo.Event{Command: eribo.CmdVonprove, Player: m.Character, Channel: m.Channel}
		if err := store.Log(e); err != nil {
			log.Printf("error logging %v: %v", eribo.CmdVonprove, err)
		}
		if err := c.SendMSG(resp); err != nil {
			log.Printf("error sending %v response: %v", eribo.CmdVonprove, err)
		}
	case strings.HasPrefix(m.Message, eribo.CmdJojo.String()):
		resp := &flist.MSG{
			Channel: m.Channel,
			Message: rp.RandJojo(m.Character),
		}
		e := &eribo.Event{Command: eribo.CmdJojo, Player: m.Character, Channel: m.Channel}
		if err := store.Log(e); err != nil {
			log.Printf("error logging %v: %v", eribo.CmdJojo, err)
		}
		if err := c.SendMSG(resp); err != nil {
			log.Printf("error sending %v response: %v", eribo.CmdJojo, err)
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
	e := &eribo.Event{Command: eribo.CmdFeedback, Player: pri.Character}
	if err := store.Log(e); err != nil {
		log.Printf("error logging %v: %v", eribo.CmdFeedback, err)
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
	return nil
}
