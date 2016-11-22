package ui

import (
	"fmt"
	"log"

	"github.com/lhz/considerate/hvsc"
	"github.com/lhz/considerate/player"
	"github.com/nsf/termbox-go"
)

const (
	MODE_BROWSE = iota
	MODE_SEARCH
)

var w, h int
var listOffset, listPos, lh, ly int
var mode int
var searchTerm, debugInfo string

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

	mode = MODE_BROWSE

	for quit := false; !quit; {
		debugInfo = ""
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
				switch mode {
				case MODE_BROWSE:
					selectTune()
				case MODE_SEARCH:
					hvsc.Filter(searchTerm)
					listOffset = 0
					listPos = 0
					mode = MODE_BROWSE
				}
			default:
				switch ev.Ch {
				case '/':
					switch mode {
					case MODE_BROWSE:
						mode = MODE_SEARCH
						searchTerm = ""
					}
				default:
					if mode == MODE_SEARCH {
						searchTerm = searchTerm + string(ev.Ch)
					}
				}
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
	if listOffset+listPos == hvsc.NumFilteredTunes-1 {
		return
	}
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
	if listOffset > hvsc.NumFilteredTunes-1 {
		listOffset -= lh
	}
}

func selectTune() {
	player.Play(hvsc.FilteredTunes[listOffset+listPos].FullPath(), 1)
}

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawHeader()
	drawList()
	drawFooter()
	termbox.Flush()
}

func drawHeader() {
	var header string
	bg := termbox.ColorBlack
	switch mode {
	case MODE_BROWSE:
		header = "Browse"
	case MODE_SEARCH:
		header = fmt.Sprintf("Search: %s_", searchTerm)
		bg = termbox.ColorGreen
	}
	header = fmt.Sprintf(fmt.Sprintf("%%s%%%ds", w-len(header)), header, "")
	writeAt(0, 0, header, termbox.ColorWhite, bg)
}

func drawFooter() {
	footer := fmt.Sprintf("%d/%d  %s", listOffset+listPos+1, hvsc.NumFilteredTunes, debugInfo)
	footer = fmt.Sprintf(fmt.Sprintf("%%s%%%ds", w-len(footer)), footer, "")
	writeAt(0, h-1, footer, termbox.ColorWhite, termbox.ColorBlack)
}

func drawList() {
	for y := 0; y < lh; y++ {
		if y+listOffset >= hvsc.NumFilteredTunes {
			break
		}
		tune := hvsc.FilteredTunes[y+listOffset]
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
