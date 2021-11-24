package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	a "github.com/logrusorgru/aurora"
	"github.com/yyewolf/termbox-go"
)

var bufferStdin []byte

func rePrintInput() {
	hostStr := a.Red(host)
	modeStr := a.Green("c")
	if connected {
		hostStr = a.Blue(host)
	}
	if currentMode == 1 {
		modeStr = a.Red("c")
	}

	fmt.Printf("%s:%s %s $ %s", a.Green(username), hostStr, modeStr, string(bufferStdin))
}

func clearInput() {
	fmt.Print("\x1b[2K\r")
}

func printf(format string, v ...interface{}) {
	clearInput()
	// NB : don't remove \n
	fmt.Printf(format+"\r\n", v...)
	rePrintInput()
}

func reflow() {
	clearInput()
	rePrintInput()
}

func askLoop() {
	for {
		evt := termbox.PollEvent()
		switch evt.Key {
		case termbox.KeyEnter:
			fmt.Print("\n")
			go commandHandler(bufferStdin)
			bufferStdin = make([]byte, 0)
		case termbox.KeyCtrlC:
			termbox.Close()
			os.Exit(0)
			return
		case termbox.KeyBackspace:
			if len(bufferStdin) > 0 {
				bufferStdin = bufferStdin[:len(bufferStdin)-1]
				reflow()
			}
		case termbox.KeyArrowLeft:
			currentMode -= 1
			if currentMode < 0 {
				currentMode = 1
			}
			reflow()
		case termbox.KeyArrowRight:
			currentMode += 1
			if currentMode < 1 {
				currentMode = 0
			}
			reflow()
		case termbox.KeySpace:
			bufferStdin = append(bufferStdin, ' ')
			fmt.Print(" ")
		default:
			if evt.Type == termbox.EventKey && byte(evt.Ch) != 0 {
				bufferStdin = append(bufferStdin, byte(evt.Ch))
				fmt.Print(string(evt.Ch))
			}
		}

	}
}

func commandHandler(data []byte) {
	input := string(data)
	splt := strings.Split(input, " ")
	if len(splt) == 0 {
		return
	}
	command := strings.ToLower(splt[0])

	switch currentMode {
	default:
		path, err := exec.LookPath(command)
		if err == nil {
			var args []string
			if len(splt) > 1 {
				args = append(args, splt[1:]...)
			}

			c := exec.Command(path, args...)
			var buff bytes.Buffer
			c.Stdout = &buff
			err = c.Run()
			if err != nil {
				printf("err : %s", a.Red(err))
				return
			}
			data := buff.String()
			data = strings.ReplaceAll(data, "\n", "\r\n")
			printf("%s", data)
			return
		} else {
			printf("")
		}
	case 0:
		call, found := commands[command]
		if !found {
			call, _ := commands[">"]
			args := []string{}
			if len(splt) > 0 {
				args = splt[0:]
			}
			ctx := &commandCtx{
				Args: args,
			}
			call(ctx)
			return
		}
		args := []string{}
		if len(splt) > 1 {
			args = splt[1:]
		}
		ctx := &commandCtx{
			Args: args,
		}
		call(ctx)
	}
}
