package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/kusubooru/eribo/flist"
)

const (
	targetChannel = "ADH-c63dd350865f6eb33043"
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

func main() {
	var (
		insecure   = flag.Bool("insecure", false, "use insecure ws:// websocket instead of wss://")
		testServer = flag.Bool("testserver", false, "connect to test server instead of production")
		addr       = flag.String("addr", "wss://chat.f-list.net:9799", "websocket address to connect")
		account    = flag.String("account", "", "websocket address to connect")
		password   = flag.String("password", "", "websocket address to connect")
		character  = flag.String("character", "", "websocket address to connect")
	)
	flag.Parse()
	if *account == "" || *password == "" || *character == "" {
		log.Println("Account, password and character name needed for identification.")
		log.Fatal("Use -account=<username> -password=<password> -character=<char name>")
	}
	*addr = defaultAddr(*addr, *testServer, *insecure)

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

	go readMessages(c.Messenger, msgch, doneRead)
	go handleMessages(msgch, doneHandle)

	// Login to F-list.
	if err := c.Identify(*account, *password, *character); err != nil {
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

func readMessages(receiver <-chan []byte, sender chan<- *flist.MSG, done chan struct{}) {
	for {
		select {
		case <-done:
			log.Println("done reading")
			return
		case message := <-receiver:
			cmd, err := flist.DecodeCommand(message)
			if err == flist.ErrUnknownCmd && len(message) != 0 {
				fmt.Println("got:", string(message))
			}
			if err != nil && err != flist.ErrUnknownCmd {
				log.Println("cmd decode error:", err)
				return
			}
			switch t := cmd.(type) {
			case *flist.MSG:
				sender <- t
			}
		}
	}
}

func handleMessages(messages <-chan *flist.MSG, done chan struct{}) {
	for {
		select {
		case <-done:
			log.Println("done handling")
			return
		case msg := <-messages:
			fmt.Println("--->", msg)
		}
	}
}
