package main

import (
	"github.com/lhz/considerate/cfg"
	"github.com/lhz/considerate/hvsc"

	"fmt"
	"log"

	"github.com/nsf/termbox-go"
)

var config *cfg.Config
var w, h, listOffset int

var sidTunes []hvsc.SidTune


func main() {
	config = cfg.ReadConfig()

	sidTunes = hvsc.BuildTuneInfo(config)

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
		tune := sidTunes[y + listOffset]
		writeAt(  0, y, fmt.Sprintf("%05d", y + listOffset + 1), termbox.ColorRed, termbox.ColorDefault)
		writeAt( 12, y, tune.Header.Name, termbox.ColorBlue, termbox.ColorDefault)
		writeAt( 45, y, tune.Header.Author, termbox.ColorGreen, termbox.ColorDefault)
		writeAt( 78, y, tune.Header.Released, termbox.ColorYellow, termbox.ColorDefault)
		writeAt(112, y, tune.Path, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()
}

func writeAt(x, y int, value string, fg, bg termbox.Attribute) {
	for i, c := range value {
		termbox.SetCell(x+i, y, c, fg, bg)
	}
}
