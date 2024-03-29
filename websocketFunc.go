package main

import (
	"io/ioutil"
	"os"

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

func (i *MessageAck) Handle() {
	switch i.Code {
	case codeError:
		printf("%s", a.Red("Error sending message."))
	case codeSuccess:
		printf("\033[A")
	}
}

func (m *ReceiveMessage) Handle() {
	hostStr := a.Red(host)
	if connected {
		hostStr = a.Blue(host)
	}
	switch m.Type {
	case messageClassic:
		if m.User != username {
			printf("%s:%s $ %s", a.Green(m.User), hostStr, m.Message)
		}
	case messageDM:
		printf("%s%s:%s $ %s", a.Blue("p."), a.Green(m.User), hostStr, m.Message)
	case messageConnection:
		printf("%s:%s has connected.", a.Green(m.User), hostStr)
	case messageDisconnection:
		printf("%s:%s has disconnected.", a.Green(m.User), hostStr)
	case messageMultiline:
		m.HandleMulti()
	}
}

func (m *ReceiveMessage) HandleMulti() {
	hostStr := a.Red(host)
	if connected {
		hostStr = a.Blue(host)
	}
	var str string
	for _, msg := range m.Messages {
		str = str + msg + "\r\n"
	}
	str += "\033[A"
	if m.User != username {
		printf("%s:%s $ %s", a.Green(m.User), hostStr, str)
	}
}

func (m *File) Handle() {
	if modeFiles {
		err := os.Mkdir("dl", os.ModePerm)
		if err != nil {
			printf("%v", err)
		}
		ioutil.WriteFile("dl/"+m.Name, m.Data, 0777)
		printf("Received dl/%s from %s", m.Name, a.Blue(m.User))
	} else {
		printf("Ignored %s from %s", m.Name, a.Blue(m.User))
	}
}
