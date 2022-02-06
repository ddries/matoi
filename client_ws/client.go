package client_ws

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"net/url"

	"github.com/ddries/matoi/util"
)

var verbose bool
var running bool

var termChan *chan bool
var conn *websocket.Conn

func InitializeClient(url *url.URL, terminateChan *chan bool, v bool) {
	verbose = v
	termChan = terminateChan

	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)

	if err != nil {
		util.ThrowError("error connecting to given host")
	}

	conn = c
	defer conn.Close()

	running = true
	util.Verbose("ws", "connected to " + url.Host)

	clientMainLoop()
}

func clientMainLoop() {
	recvMsgs := make(chan string)
	sendMsgs := make(chan string)

	go recvLoop(&recvMsgs)
	go sendLoop(&sendMsgs)

	for running {
		select {
		case msg := <- recvMsgs:
			if verbose {
				util.Verbose("recv", "received message from " + conn.RemoteAddr().String())
			}

			fmt.Println(msg)
			break

		case msg := <- sendMsgs:
			if verbose {
				util.Verbose("ws", "sending: " + msg)
			}

			conn.WriteMessage(1, []byte(msg))
			break
		
		case <- *termChan:
			running = false
			break
		}
	}

	util.Verbose("ws", "terminated client main loop")
	TerminateClient()
}

func recvLoop(recvChan *chan string) {
	defer close(*recvChan)

	for running {
		_, message, err := conn.ReadMessage()

		if err != nil && verbose {
			util.Verbose("ws", "error reading recv message")
		} else {
			*recvChan <- string(message)
		}
	}
}

func sendLoop(sendChan *chan string) {
	defer close(*sendChan)

	reader := bufio.NewReader(os.Stdin)

	for running {
		t, err := reader.ReadString('\n')

		if err != nil {
			if verbose {
				util.Verbose("ws", "error reading from console")
			}
		} else {
			t = strings.Replace(t, "\n", "", -1)

			if len(t) > 0 {
				*sendChan <- t
			}
		}
	}
}

func TerminateClient() {
	if verbose {
		util.Verbose("ws", "terminating client connection...")
	}
	conn.Close()
}