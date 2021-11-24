package main

import "github.com/gorilla/websocket"

// Defines OP names
const (
	identifyOp = iota
	sendMessageOp
	receiveMessageOp
)

// Defines messages type
const (
	messageClassic = iota
	messageConnection
	messageDisconnection
	messageMultiline
)

// Defines codes
const (
	codeSuccess = 200
	codeError   = 400
)

// Defines the messages standard
type Message struct {
	Op   int         `json:"op"`
	Data interface{} `json:"data,omitempty"`
}

type SendIdentify struct {
	Username string `json:"username,omitempty"`
}

type SendMessage struct {
	Type     int      `json:"type"`
	Message  string   `json:"message,omitempty"`
	Messages []string `json:"messages,omitempty"`
}

type ResponseIdentify struct {
	Code int `json:"code,omitempty"`
}

type MessageAck struct {
	Code int `json:"code,omitempty"`
}

type ReceiveMessage struct {
	Type     int      `json:"type"`
	User     string   `json:"user,omitempty"`
	Message  string   `json:"message,omitempty"`
	Messages []string `json:"messages,omitempty"`
}

// Variables relevant to client
var host string
var port string
var username string
var connected bool
var currentMode int
var connection *websocket.Conn

// Client side commands
var commands map[string]func(*commandCtx)

type commandCtx struct {
	Args []string
}
