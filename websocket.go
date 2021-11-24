package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	a "github.com/logrusorgru/aurora"
)

func connectGateway() {
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%s", host, port), Path: "/"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		printf("Error : %s", a.Red(err))
		return
	}
	connection = c
	go receiveGateway()
	time.Sleep(1 * time.Second)
	connection.WriteJSON(Message{
		Op: identifyOp,
		Data: SendIdentify{
			Username: username,
		},
	})
}

func receiveGateway() {
	for {
		msg := &Message{}
		err := connection.ReadJSON(msg)
		if err != nil {
			break
		}
		handleMessages(msg)
	}
}

func handleMessages(m *Message) {
	d, err := json.Marshal(m.Data)
	if err != nil {
		return
	}
	switch m.Op {
	case identifyOp:
		packet := &ResponseIdentify{}
		err = json.Unmarshal(d, &packet)
		if err != nil {
			return
		}
		packet.Handle()
	case sendMessageOp:
		packet := &MessageAck{}
		err = json.Unmarshal(d, &packet)
		if err != nil {
			return
		}
		packet.Handle()
	case receiveMessageOp:
		packet := &ReceiveMessage{}
		err = json.Unmarshal(d, &packet)
		if err != nil {
			return
		}
		packet.Handle()
	}
}
