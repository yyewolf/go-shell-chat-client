package main

import "strings"

func commandLoader() {
	commands = make(map[string]func(*commandCtx))
	commands[">"] = sendMsg
}

func sendMsg(c *commandCtx) {
	var message string
	for _, arg := range c.Args {
		message += arg + " "
	}
	message = strings.TrimSpace(message)
	connection.WriteJSON(Message{
		Op: sendMessageOp,
		Data: SendMessage{
			Type:    0,
			Message: message,
		},
	})
}
