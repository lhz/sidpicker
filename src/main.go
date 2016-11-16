package main

import (
	"log"

	"github.com/nsf/termbox-go"
)

var w, h, listOffset int

func main() {
	readSongLengths()

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
			switch ev.Key {
			case termbox.KeyCtrlC, termbox.KeyEsc:
				quit = true
			case termbox.KeyPgup:
				listOffset -= h
				if listOffset < 0 {
					listOffset = 0
				}
				draw()
			case termbox.KeyPgdn:
				listOffset += h
				if listOffset > len(sidTunes)-1 {
					listOffset = len(sidTunes)-1
				}
				draw()
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
	for y := 0; y < h; y++ {
		writeAt(2, y, sidTunes[y + listOffset].Path, termbox.ColorWhite, termbox.ColorBlue)
	}
	termbox.Flush()
}

func writeAt(x, y int, value string, fg, bg termbox.Attribute) {
	for i, c := range value {
		termbox.SetCell(x+i, y, c, fg, bg)
	}
}
