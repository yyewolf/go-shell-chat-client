package main

import (
	"os"
	"strings"
)

func commandLoader() {
	commands = make(map[string]func(*commandCtx))
	commands[">"] = sendMsg
	commands["send"] = sendFile
}

func sendMsg(c *commandCtx) {
	var message string
	for _, arg := range c.Args {
		message += arg + " "
	}
	message = strings.TrimSpace(message)
	if message != "" {
		connection.WriteJSON(Message{
			Op: sendMessageOp,
			Data: SendMessage{
				Type:    0,
				Message: message,
			},
		})
	} else {
		printf("")
	}
}

func sendFile(c *commandCtx) {
	var file string
	var target string
	if len(c.Args) < 2 {
		printf("Usage : send <filepath> <target>.")
		return
	}
	target = c.Args[1]
	file = c.Args[0]
	data, err := os.ReadFile(file)
	if err != nil {
		printf("Error : %v.", err)
		return
	}
	connection.WriteJSON(Message{
		Op: sendFileOp,
		Data: File{
			Data: data,
			User: target,
			Name: file,
		},
	})
}
