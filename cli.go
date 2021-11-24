package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	a "github.com/logrusorgru/aurora"
)

var bufferStdin []byte

func rePrintInput() {
	hostStr := a.Red(host)
	if connected {
		hostStr = a.Blue(host)
	}
	fmt.Printf("%s:%s $ %s", a.Green(username), hostStr, string(bufferStdin))
}

func clearInput() {
	fmt.Print("\x1b[2K\r")
}

func printf(format string, v ...interface{}) {
	clearInput()
	// NB : don't remove \n
	fmt.Printf(format+"\n", v...)
	rePrintInput()
}

func askLoop() {
	in := bufio.NewReader(os.Stdin)
	for {
		c, _ := in.ReadByte()
		if c == '\015' {
			fmt.Print("\n")
			go commandHandler(bufferStdin)
			bufferStdin = make([]byte, 0)
			continue
		}
		if c == '\003' {
			os.Exit(0)
			return
		}
		bufferStdin = append(bufferStdin, c)
		fmt.Print(string(c))
	}
}

func commandHandler(data []byte) {
	input := string(data)
	splt := strings.Split(input, " ")
	if len(splt) == 0 {
		return
	}
	command := strings.ToLower(splt[0])

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
		printf("%s", string(data))
		return
	}

	call, found := commands[command]
	if !found {
		printf("Command %s does not exist.", a.Red(command))
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
