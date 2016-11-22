package ui

import (
	"fmt"
	"log"

	"github.com/lhz/considerate/hvsc"
	"github.com/lhz/considerate/player"
	"github.com/nsf/termbox-go"
)

var w, h int
var listOffset, listPos, lh, ly int

func Setup() {
	err := termbox.Init()
	if err != nil {
		log.Panicln(err)
	}

	w, h = termbox.Size()
	listOffset = 0
	listPos = 0
	ly = 1
	lh = h - 2

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	draw()
}

func Run() {
	defer termbox.Close()

	for quit := false; !quit; {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlC, termbox.KeyEsc:
				quit = true
			case termbox.KeyPgup:
				pageUp()
			case termbox.KeyPgdn:
				pageDown()
			case termbox.KeyArrowUp:
				moveUp()
			case termbox.KeyArrowDown:
				moveDown()
			case termbox.KeyEnter:
				selectTune()
			}
			draw()
		case termbox.EventResize:
			w, h = ev.Width, ev.Height
			draw()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func moveUp() {
	listPos--
	if listPos < 0 {
		if listOffset > 0 {
			listPos = lh - 1
			pageUp()
		} else {
			listPos = 0
		}
	}
}

func moveDown() {
	listPos++
	if listPos >= lh {
		listPos = 0
		pageDown()
	}
}

func pageUp() {
	listOffset -= lh
	if listOffset < 0 {
		listOffset = 0
	}
}

func pageDown() {
	listOffset += lh
	if listOffset > hvsc.NumTunes-1 {
		listOffset -= lh
	}
}

func selectTune() {
	player.Play(hvsc.Tunes[listOffset+listPos].FullPath(), 1)
}

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawHeader()
	drawList()
	drawFooter()
	termbox.Flush()
}

func drawHeader() {
	header := fmt.Sprintf(fmt.Sprintf("%%%ds", w), "")
	writeAt(0, 0, header, termbox.ColorWhite, termbox.ColorBlack)
}

func drawFooter() {
	footer := fmt.Sprintf("%d/%d", listOffset+listPos+1, hvsc.NumTunes)
	footer = fmt.Sprintf(fmt.Sprintf("%%s%%%ds", w-len(footer)), footer, "")
	writeAt(0, h-1, footer, termbox.ColorWhite, termbox.ColorBlack)
}

func drawList() {
	for y := 0; y < lh; y++ {
		tune := hvsc.Tunes[y+listOffset]
		fg := termbox.ColorDefault
		bg := termbox.ColorDefault
		if y == listPos {
			fg = termbox.ColorDefault | termbox.AttrBold
			bg = termbox.ColorBlue
		}
		writeAt(0, ly+y, tune.Path, fg, bg)
	}
}

func writeAt(x, y int, value string, fg, bg termbox.Attribute) {
	i := 0
	for _, c := range value {
		termbox.SetCell(x+i, y, c, fg, bg)
		i++
	}
}
