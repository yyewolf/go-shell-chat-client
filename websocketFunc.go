package main

import (
	a "github.com/logrusorgru/aurora"
)

func (i *ResponseIdentify) Handle() {
	switch i.Code {
	case codeError:
		printf("%s", a.Red("Error connecting to websocket."))
	case codeSuccess:
		connected = true
		printf("%s", a.Green("Connected to websocket."))
	}
}

func (m *ReceiveMessage) Handle() {
	hostStr := a.Red(host)
	if connected {
		hostStr = a.Blue(host)
	}
	switch m.Type {
	case messageClassic:
		printf("%s:%s $ %s", a.Green(m.User), hostStr, m.Message)
		break
	}
}
