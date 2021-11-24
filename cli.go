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

func modeStr(mode int) string {
	modeSend = false
	modeFiles = false
	modeMultiline = false
	switch mode {
	case 1:
		modeSend = true
		return fmt.Sprintf("%s%s%s", a.Green("c"), a.Red("f"), a.Red("m"))
	case 2:
		modeSend = true
		modeFiles = true
		return fmt.Sprintf("%s%s%s", a.Green("c"), a.Green("f"), a.Red("m"))
	case 3:
		modeSend = true
		modeMultiline = true
		return fmt.Sprintf("%s%s%s", a.Green("c"), a.Red("f"), a.Green("m"))
	case 4:
		modeSend = true
		modeFiles = true
		modeMultiline = true
		return fmt.Sprintf("%s%s%s", a.Green("c"), a.Green("f"), a.Green("m"))
	}
	return fmt.Sprintf("%s%s%s", a.Red("c"), a.Red("f"), a.Red("m"))
}

func rePrintInput() {
	hostStr := a.Red(host)
	mode := modeStr(currentMode)
	if connected {
		hostStr = a.Blue(host)
	}
	if !modeMultiline || len(buffer) == 0 {
		fmt.Printf("%s:%s %s $ %s", a.Green(username), hostStr, mode, string(bufferStdin))
	} else {
		fmt.Printf("%s:%s %s $ %s\r\n", a.Green(username), hostStr, mode, buffer[0])
		for i := range buffer {
			if i != 0 {
				fmt.Printf("%s%s", buffer[i], "\r\n")
			}
		}
		fmt.Printf("%s", string(bufferStdin))
	}
}

func clearInput() {
	if modeMultiline && len(buffer) != 0 {
		fmt.Printf("\x1b[%dA", len(buffer))
	}
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
			if modeMultiline {
				if lastKey == termbox.KeyEnter {
					// We finish the multiline
					connection.WriteJSON(Message{
						Op: sendMessageOp,
						Data: SendMessage{
							Type:     messageMultiline,
							Messages: buffer,
						},
					})
					buffer = make([]string, 0)
				} else {
					buffer = append(buffer, string(bufferStdin))
					fmt.Print("\r\n")
					bufferStdin = make([]byte, 0)
				}
			} else {
				fmt.Print("\n")
				go commandHandler(bufferStdin)
				bufferStdin = make([]byte, 0)
			}
		case termbox.KeyCtrlC:
			termbox.Close()
			os.Exit(0)
			return
		case termbox.KeyBackspace:
			if len(bufferStdin) > 0 {
				bufferStdin = bufferStdin[:len(bufferStdin)-1]
				reflow()
			}
		case termbox.KeyBackspace2:
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
			if currentMode > 4 {
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
		lastKey = evt.Key
	}
}

func commandHandler(data []byte) {
	input := string(data)
	splt := strings.Split(input, " ")
	if len(splt) == 0 {
		return
	}
	command := strings.ToLower(splt[0])

	if !modeSend {
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
	} else {
		if connection == nil {
			printf("not connected yet")
			return
		}
		call, found := commands[command]
		if !found {
			call = commands[">"]
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
