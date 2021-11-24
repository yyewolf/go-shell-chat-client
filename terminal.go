package main

import (
	"bytes"
	"io/ioutil"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var CursorX, CursorY int
var ViewX, ViewY int

func tbPrint(fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(CursorX, CursorY, c, fg, bg)
		CursorX += runewidth.RuneWidth(c)
	}
}

func tbPrintln(fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(CursorX, CursorY, c, fg, bg)
		CursorX += runewidth.RuneWidth(c)
	}
	CursorX = 0
	CursorY += 1
}

func tbDraw(r *bytes.Buffer) {
	out, err := ioutil.ReadAll(r)
	if err != nil {

	}
	tbPrintln(termbox.ColorRed, termbox.ColorDefault, "got : "+string(out))
}
