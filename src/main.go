package main

import (
	"fmt"
	"log"

	"github.com/nsf/termbox-go"
)

var w, h int

func main() {
	err := termbox.Init()
	if err != nil {
		log.Panicln(err)
	}
	defer termbox.Close()

	w, h = termbox.Size()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	draw()

	for quit := false; !quit; {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyEsc {
				quit = true
			}
		case termbox.EventResize:
			w, h = ev.Width, ev.Height
			draw()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	writeAt(5, 3, fmt.Sprintf("Terminal size: %dx%d", w, h), termbox.ColorWhite, termbox.ColorBlue)
	termbox.Flush()
}

func writeAt(x, y int, value string, fg, bg termbox.Attribute) {
	for i, c := range value {
		termbox.SetCell(x+i, y, c, fg, bg)
	}
}
