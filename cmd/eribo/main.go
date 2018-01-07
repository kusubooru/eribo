package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"mvdan.cc/xurls"

	_ "github.com/go-sql-driver/mysql"

	"github.com/kusubooru/eribo/dadjoke"
	"github.com/kusubooru/eribo/eribo"
	"github.com/kusubooru/eribo/eribo/mysql"
	"github.com/kusubooru/eribo/flist"
	"github.com/kusubooru/eribo/loot"
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
	if err := json.Unmarshal([]byte(s), &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

type lowLothNames []string

func (names lowLothNames) String() string { return fmt.Sprintf("%q", []string(names)) }
func (names *lowLothNames) Set(value string) error {
	*names = append(*names, value)
	return nil
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
		owner       = flag.String("owner", "Ryuunosuke Akashaka", "character name of the bot's owner")
		dataSource  = flag.String("datasource", "", "MySQL datasource")
		joinRooms   = flag.String("join", "", "open private `rooms` to join in JSON format e.g. "+`-join '["Room 1", "Room 2"]'`)
		statusMsg   = flag.String("status", "", "status message to be displayed")
		showVersion = flag.Bool("v", false, "print program version")
		lowNames    lowLothNames
		versionArg  bool
	)
	flag.Var(&lowNames, "lowname", "`name` of player for lower loth chance e.g. -lowname 'Name 1' -lowname 'Name 2'")
	flag.Parse()

	botVersion := fmt.Sprintf("%s %s (runtime: %s)", filepath.Base(os.Args[0]), theVersion, runtime.Version())
	versionArg = len(os.Args) > 1 && os.Args[1] == "version"
	if *showVersion || versionArg {
		fmt.Println(botVersion)
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
	ciuch := make(chan *flist.CIU, 10)
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
		ciuch,
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

	// Change bot status.
	sta := flist.STA{Status: flist.StatusBusy, StatusMsg: *statusMsg}
	if err := c.SendCmd(sta); err != nil {
		log.Println(err)
		return
	}

	handleMessages(
		c,
		*account,
		*password,
		*character,
		botVersion,
		*owner,
		lowNames,
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
		ciuch,
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
	ciuch chan<- *flist.CIU,
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
	defer close(ciuch)
	defer close(quit)
	for {
		message, err := c.ReadMessage()
		if err != nil {
			log.Println("read message error:", err)
			quit <- struct{}{}
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
		case *flist.CIU:
			ciuch <- t
		case *flist.VAR:
			switch t.Variable {
			case "chat_max":
				c.SetChatMax(t.ChatMax)
			case "priv_max":
				c.SetPrivMax(t.PrivMax)
			default:
			}
		case *flist.ERR:
			log.Println(fmt.Errorf("flist ERR %d: %s", t.Number, t.Message))
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
	botName string,
	botVersion string,
	owner string,
	lowNames []string,
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
	ciuch <-chan *flist.CIU,
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
			select {
			case <-quit:
			case <-time.After(5 * time.Second):
				log.Println("reader took too long")
			}
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
			respond(c, store, msg, channelMap, botName, owner, lowNames)
		case pri := <-prich:
			if err := gatherFeedback(c, store, pri); err != nil {
				log.Println("gather feedback err:", err)
			}
			respondPrivOwner(c, store, pri, channelMap, botName, botVersion, owner)
		case ors := <-orsch:
			flist.SortChannelsByTitle(ors.Channels)
			for _, title := range roomTitles {
				ch := flist.FindChannel(ors.Channels, title)
				if ch != nil {
					jch := flist.JCH{Channel: ch.Name}
					if err := c.SendCmd(jch); err != nil {
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
			go func(ticket string) {
				for _, u := range ich.Users {
					if p, ok := playerMap.GetPlayer(u.Identity); ok {
						channelMap.SetPlayer(ich.Channel, p)
					}
					actives := channelMap.GetActivePlayers()
					actives.ForEach(func(name string, p *eribo.Player) {
						// The JSON endpoint returns HTML with 405 error when lots
						// of requests are done concurrently.

						// go func(name, account, ticket string, p *eribo.Player, playerMap *eribo.PlayerMap, mappingList *flist.MappingList) {
						charData, err := flist.GetCharacterData(name, account, ticket)
						if err != nil {
							log.Printf("init channel could not get character data for %q: %v", name, err)
							if carrier, ok := err.(BodyCarrier); ok {
								log.Printf("-- %s character data begin --\n", name)
								log.Printf("%s\n", string(carrier.Body()))
								log.Printf("-- %s character data end --\n", name)
							}
							return
						}
						m := charData.HumanInfotags(mappingList)
						if role, ok := m["Dom/Sub Role"]; ok {
							playerMap.SetPlayerRole(p.Name, flist.Role(role))
						}
						playerMap.SetPlayerFave(name, charFavesTickling(charData, mappingList))
						// }(name, account, ticket, p, playerMap, mappingList)
					})
				}
			}(ticket)
		case sta := <-stach:
			name := sta.Character
			newStatus := flist.Status(sta.Status)
			player, _ := channelMap.GetPlayer(name)
			if player != nil && player.Role == "" && !player.Status.IsActive() && newStatus.IsActive() {
				//fmt.Printf("STA changed to active for char %q\n", name)
				if err := getCharDataAndSetRole(name, account, password, playerMap, mappingList); err != nil {
					log.Println("STA: %v", err)
				}
			}
			playerMap.SetPlayerStatus(name, newStatus)
		case jch := <-jchch:
			name := jch.Character.Identity
			player, _ := playerMap.GetPlayer(name)
			if player == nil {
				log.Printf("JCH: player %q not found in playerMap", name)
				return
			}
			if player.Role == "" && player.Status.IsActive() {
				//fmt.Printf("player %q joined, getting char data\n", name)
				if err := getCharDataAndSetRole(name, account, password, playerMap, mappingList); err != nil {
					log.Println("JCH: %v", err)
				}
			}
			channelMap.SetPlayer(jch.Channel, player)
		case lch := <-lchch:
			channelMap.DelPlayer(lch.Channel, lch.Character)
		case ciu := <-ciuch:
			jch := flist.JCH{Channel: ciu.Name}
			if err := c.SendCmd(jch); err != nil {
				log.Println("CIU error joining private room %q: %v", ciu.Title, err)
			}
		case prd := <-prdch:
			fmt.Println("got prd:", prd)
		case <-pinch:
			if err := c.SendCmd(flist.PIN{}); err != nil {
				log.Println("send PIN failed:", err)
			}
		case idn := <-idnch:
			// Expecting IDN only once during identification.
			log.Println("received IDN but shouldn't:", idn)
		}
	}
}

type BodyCarrier interface {
	Body() []byte
}

func getCharDataAndSetRole(name, account, password string, playerMap *eribo.PlayerMap, mappingList *flist.MappingList) error {
	ticket, err := flist.GetTicket(account, password)
	if err != nil {
		return fmt.Errorf("get ticket error: %v", err)
	}

	charData, err := flist.GetCharacterData(name, account, ticket)
	if err != nil {
		return fmt.Errorf("get character data for %q error: %v", name, err)
	}
	m := charData.HumanInfotags(mappingList)
	if role, ok := m["Dom/Sub Role"]; ok {
		playerMap.SetPlayerRole(name, flist.Role(role))
	}
	playerMap.SetPlayerFave(name, charFavesTickling(charData, mappingList))
	return nil
}

func charFavesTickling(char *flist.CharacterData, ml *flist.MappingList) bool {
	if char.HasFaveKink(ml.KinksMap(), "Tickling") || char.HasFaveCustomKink("tickling", "tickle") {
		return true
	}
	return false
}

type logAdder interface {
	AddCmdLog(*eribo.CmdLog) error
	AddLothLog(*eribo.LothLog) error
}

func respond(
	c *flist.Client,
	logAdder logAdder,
	m *flist.MSG,
	channelMap *eribo.ChannelMap,
	botName string,
	owner string,
	lowNames []string,
) {
	var msg string
	var rperr error
	cmd, args := eribo.ParseCommand(m.Message)
	switch cmd {
	case eribo.CmdTieup:
		msg = rp.RandTieUp(m.Character)
	case eribo.CmdTomato:
		msg = rp.Tomato(m.Character, owner)
	case eribo.CmdTktool:
		msg = rp.RandTktool(m.Character)
	case eribo.CmdVonprove:
		msg = rp.RandVonprove(m.Character)
	case eribo.CmdJojo:
		msg = rp.RandJojo(m.Character)
	case eribo.CmdTietool:
		toolType := ""
		if len(args) != 0 {
			toolType = args[0]
		}
		msg, rperr = rp.RandTietool(m.Character, toolType)
		if rperr != nil {
			log.Printf("RandTietool error: %v", rperr)
			return
		}
	case eribo.CmdDadJoke:
		j, err := dadjoke.Random()
		if err != nil {
			log.Printf("error getting dadjoke: %v", err)
			return
		}
		msg = j.Joke
	case eribo.CmdLoth:
		if len(args) > 0 && args[0] == "time" {
			loth := channelMap.Loth(m.Channel)
			msg = rp.LothTime(loth)
			break
		}
		loth, isNew, targets := channelMap.ChooseLoth(m.Character, m.Channel, botName, 1*time.Hour, lowNames)
		lothLog := &eribo.LothLog{Issuer: m.Character, Channel: m.Channel, Loth: loth, IsNew: isNew, Targets: targets}
		if err := logAdder.AddLothLog(lothLog); err != nil {
			log.Printf("error logging Loth: %v, isNew: %v, Targets: %v: %v", loth, isNew, targets, err)
		}
		msg = rp.Loth(m.Character, loth, isNew, targets)
	}

	if msg != "" {
		e := &eribo.CmdLog{Command: cmd, Args: strings.Join(args, " "), Player: m.Character, Channel: m.Channel}
		go func(e *eribo.CmdLog) {
			if err := logAdder.AddCmdLog(e); err != nil {
				log.Printf("error logging %v: %v", cmd, err)
			}
		}(e)

		resp := &flist.MSG{Channel: m.Channel, Message: msg}
		if err := c.SendMSG(resp); err != nil {
			log.Printf("error sending %v response: %v", cmd, err)
		}
	}
}

func atoiLimitOffset(args []string) (int, int) {
	limit, offset := 10, 0
	if len(args) > 0 {
		if lim, err := strconv.Atoi(args[0]); err == nil {
			limit = lim
		}
	}
	if len(args) > 1 {
		if off, err := strconv.Atoi(args[1]); err == nil {
			offset = off
		}
	}
	return limit, offset
}

func atoiFirstArg(args []string, def int) int {
	n := def
	if len(args) > 0 {
		if i, err := strconv.Atoi(args[0]); err == nil {
			n = i
		}
	}
	return n
}

func respondPrivOwner(c *flist.Client, store eribo.Store, pri *flist.PRI, channelMap *eribo.ChannelMap, botName, botVersion, owner string) {
	if pri.Character != owner {
		return
	}

	var msg string
	cmd, args := eribo.ParseCustomCommand(pri.Message)
	switch cmd {
	case "!version":
		msg = botVersion
	case "!status":
		sta := flist.STA{Status: flist.StatusBusy, StatusMsg: strings.Join(args, " ")}
		if err := c.SendCmd(sta); err != nil {
			log.Println("owner changing status:", err)
		}
	case "!simtktools":
		rolls := atoiFirstArg(args, 100)
		table := &loot.Table{}
		for _, t := range rp.Tktools() {
			table.Add(t, t.Weight)
		}
		drops, pr := table.Sim(rolls)
		var buf bytes.Buffer
		buf.WriteString("\n")
		for i, t := range rp.Tktools() {
			buf.WriteString(fmt.Sprintf("%s = %d, %.1f%%\n", t.Name, drops[i], pr[i]*100.0))
		}
		msg = buf.String()
	case "!simtietools":
		rolls := atoiFirstArg(args, 100)
		table := &loot.Table{}
		toolType := ""
		if len(args) > 1 {
			toolType = args[1]
		}
		for _, t := range rp.Tietools(toolType) {
			table.Add(t, t.Quality.Weight())
		}
		drops, pr := table.Sim(rolls)
		var buf bytes.Buffer
		buf.WriteString("\n")
		for i, t := range rp.Tietools(toolType) {
			buf.WriteString(fmt.Sprintf("%s = %d, %.1f%%\n", t.Name, drops[i], pr[i]*100.0))
		}
		msg = buf.String()
	case "!channelmap":
		var buf bytes.Buffer
		buf.WriteString("\n")
		channelMap.ForEach(func(channel string, pm *eribo.PlayerMap) {
			buf.WriteString(fmt.Sprintf("Channel: %q\n", channel))
			pm.ForEach(func(name string, p *eribo.Player) {
				buf.WriteString(fmt.Sprintf("|- %v\n", p))
			})
		})
		msg = buf.String()
	case "!feed":
		limit, offset := atoiLimitOffset(args)
		feedback, err := store.GetRecentFeedback(limit, offset)
		if err != nil {
			log.Printf("%v error getting feedback: %v", cmd, err)
		}
		var buf bytes.Buffer
		buf.WriteString("\n")
		for _, fb := range feedback {
			buf.WriteString(fmt.Sprintf("%v\n", fb))
		}
		msg = buf.String()
	case "!cmdlogs":
		limit, offset := atoiLimitOffset(args)
		logs, err := store.GetRecentCmdLogs(limit, offset)
		if err != nil {
			log.Printf("%v error getting cmd logs: %v", cmd, err)
		}
		var buf bytes.Buffer
		buf.WriteString("\n")
		for _, lg := range logs {
			buf.WriteString(fmt.Sprintf("%v\n", lg))
		}
		msg = buf.String()
	case "!lothlogs":
		limit, offset := atoiLimitOffset(args)
		logs, err := store.GetRecentLothLogs(limit, offset)
		if err != nil {
			log.Printf("%v error getting loth logs: %v", cmd, err)
		}
		var buf bytes.Buffer
		buf.WriteString("\n")
		for _, lg := range logs {
			buf.WriteString(fmt.Sprintf("%v\n", lg))
		}
		msg = buf.String()
	}

	if msg != "" {
		resp := &flist.PRI{
			Recipient: pri.Character,
			Message:   msg,
		}
		err := c.SendPRI(resp)
		switch err {
		case flist.ErrMsgTooLong:
			resp.Message = fmt.Sprintf("%v", flist.ErrMsgTooLong)
			if err2 := c.SendPRI(resp); err2 != nil {
				log.Printf("error sending PRI response: %v", err2)
			}
		case nil:
		default:
			log.Printf("error sending %v response: %v", cmd, err)
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
	e := &eribo.CmdLog{Command: eribo.CmdFeedback, Player: pri.Character}
	if err := store.AddCmdLog(e); err != nil {
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
