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
			case 0:
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
			case 0x7F:
				if mode == MODE_SEARCH {
					if len(searchTerm) > 0 {
						searchTerm = searchTerm[0 : len(searchTerm)-1]
					}
				}
			default:
				debugInfo = string(ev.Key)
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
	player.Play(listOffset+listPos, 1)
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
		termbox.HideCursor()
	case MODE_SEARCH:
		header = fmt.Sprintf("Search: %s", searchTerm)
		//bg = termbox.ColorGreen
		termbox.SetCursor(8+len(searchTerm), 0)
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
		if tune.Index == player.CurrentIndex {
			fg |= termbox.AttrBold
		}
		if y == listPos {
			bg = termbox.ColorBlue
		}
		writeAt(0, ly+y, listName(tune), fg, bg)
	}
}

func listName(tune hvsc.SidTune) string {
	return fmt.Sprintf("%4s %s", tune.Year(), tune.Path[1:len(tune.Path)-4])
}

func writeAt(x, y int, value string, fg, bg termbox.Attribute) {
	i := 0
	for _, c := range value {
		termbox.SetCell(x+i, y, c, fg, bg)
		i++
	}
}
