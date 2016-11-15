package main

import (
	"fmt"
	"log"

	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		log.Panicln(err)
	}
	defer termbox.Close()

	w, h := termbox.Size()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()

	writeAt(5, 3, fmt.Sprintf("Terminal size: %dx%d", w, h), termbox.ColorWhite, termbox.ColorBlue)
	termbox.Flush()

	for quit := false; !quit; {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyEsc {
				quit = true
			}
		case termbox.EventResize:
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func writeAt(x, y int, value string, fg, bg termbox.Attribute) {
	for i, c := range value {
		termbox.SetCell(x+i, y, c, fg, bg)
	}
}
