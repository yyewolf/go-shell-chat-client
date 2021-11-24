package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// We read passed params
func init() {
	flag.StringVar(&host, "h", "127.0.0.1", "hostname to connect to")
	flag.StringVar(&port, "p", "30", "port to connect to")
	flag.StringVar(&username, "u", "", "username to display")
	flag.Parse()
	if username == "" {
		panic("Set username using -u")
	}

	_, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
		return
	}
	terminal.NewTerminal(os.Stdin, ">")

	commandLoader()
	go askLoop()
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
	rePrintInput()

	connectGateway()
}

func main() {
	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
