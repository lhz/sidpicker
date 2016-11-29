package ui

import (
	"fmt"
	"log"
	"strings"
	"time"

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
var debugInfo string
var searchTerm []rune
var searchCursorPos int
var quit = false

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

	tickChan := time.NewTicker(1 * time.Second).C

	eventChan := make(chan termbox.Event)
	go func() {
		for {
			eventChan <- termbox.PollEvent()
		}
	}()

	for !quit {
		debugInfo = ""
		select {
		case ev := <-eventChan:
			switch ev.Type {
			case termbox.EventKey:
				keyEvent(ev)
				draw()
			case termbox.EventResize:
				resizeEvent(ev)
				draw()
			case termbox.EventError:
				panic(ev.Err)
			}
		case <-tickChan:
			if player.Playing {
				draw()
			}
		}
	}
}

func resizeEvent(ev termbox.Event) {
	w, h = ev.Width, ev.Height
	draw()
}

func keyEvent(ev termbox.Event) {
	if mode == MODE_SEARCH {
		keyEventSearch(ev)
		return
	}
	switch ev.Key {
	case termbox.KeyCtrlC, termbox.KeyEsc:
		player.Stop()
		quit = true
	case termbox.KeyPgup:
		pageUp()
	case termbox.KeyPgdn:
		pageDown()
	case termbox.KeyArrowUp:
		moveUp()
	case termbox.KeyArrowDown:
		moveDown()
	case termbox.KeyArrowLeft:
		player.PrevSong()
	case termbox.KeyArrowRight:
		player.NextSong()
	case termbox.KeyEnter:
		player.Play(listOffset+listPos, 1)
	case termbox.KeySpace:
		moveDown()
		player.Play(listOffset+listPos, 1)
	case termbox.KeyDelete:
		player.Stop()
	case 0:
		if n := strings.IndexRune("0123456789", ev.Ch); n > 0 {
			if n <= player.CurrentTune.Header.Songs {
				player.PlaySub(n)
			}
			return
		}
		if ev.Ch == '/' {
			mode = MODE_SEARCH
		}
	default:
		debugInfo = string(ev.Key)
	}
}

func keyEventSearch(ev termbox.Event) {
	switch ev.Key {
	case termbox.KeyCtrlC, termbox.KeyEsc:
		mode = MODE_BROWSE
	case termbox.KeyArrowLeft:
		if searchCursorPos > 0 {
			searchCursorPos--
		}
	case termbox.KeyArrowRight:
		if searchCursorPos < len(searchTerm) {
			searchCursorPos++
		}
	case termbox.KeyEnter:
		hvsc.Filter(string(searchTerm))
		listOffset = 0
		listPos = 0
		mode = MODE_BROWSE
	case termbox.KeySpace:
		searchInsert(rune(' '))
	case 0:
		searchInsert(ev.Ch)
	case 0x7F:
		if searchCursorPos > 0 {
			searchTerm = append(searchTerm[0:searchCursorPos-1], searchTerm[searchCursorPos:len(searchTerm)]...)
			searchCursorPos--
		}
	default:
		debugInfo = string(ev.Key)
	}
}

func searchInsert(ch rune) {
	searchTerm = append(searchTerm, rune(' '))
	copy(searchTerm[searchCursorPos+1:], searchTerm[searchCursorPos:])
	searchTerm[searchCursorPos] = ch
	searchCursorPos++
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

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawHeader()
	drawList()
	drawTuneInfo()
	drawFooter()
	termbox.Flush()
}

func drawHeader() {
	var header string
	bg := termbox.ColorBlack
	switch mode {
	case MODE_BROWSE:
		header = fmt.Sprintf("Browse: %s", string(searchTerm))
		termbox.HideCursor()
	case MODE_SEARCH:
		header = fmt.Sprintf("Search: %s", string(searchTerm))
		bg = termbox.ColorGreen
		termbox.SetCursor(8+searchCursorPos, 0)
	}
	header = fmt.Sprintf(fmt.Sprintf("%%s%%%ds", w-len([]rune(header))), header, "")
	writeAt(0, 0, header, termbox.ColorWhite, bg)
}

func drawFooter() {
	footer := fmt.Sprintf("%d/%d  %s", listOffset+listPos+1, hvsc.NumFilteredTunes, debugInfo)
	footer = fmt.Sprintf(fmt.Sprintf("%%s%%%ds", w-len([]rune(footer))), footer, "")
	writeAt(0, h-1, footer, termbox.ColorWhite, termbox.ColorBlack)
}

func drawTuneInfo() {
	if !player.Playing {
		return
	}
	tune := player.CurrentTune
	ox := w - 80
	oy := 2
	fg := termbox.ColorDefault
	bg := termbox.ColorDefault

	writeAt(ox, oy+0, fmt.Sprintf("Title:    %s", tune.Header.Name), fg, bg)
	writeAt(ox, oy+1, fmt.Sprintf("Author:   %s", tune.Header.Author), fg, bg)
	writeAt(ox, oy+2, fmt.Sprintf("Released: %s", tune.Header.Released), fg, bg)

	writeAt(ox, oy+4, fmt.Sprintf("Tune: %d/%d  Length: %s  Time: %s",
		player.CurrentSong, tune.Header.Songs, player.SongLength(), player.Elapsed()), fg, bg)
}

func drawList() {
	for y := 0; y < lh; y++ {
		if y+listOffset >= hvsc.NumFilteredTunes {
			break
		}
		tune := hvsc.FilteredTunes[y+listOffset]
		fg := termbox.ColorDefault
		bg := termbox.ColorDefault
		if player.CurrentTune != nil && player.CurrentTune.Index == tune.Index {
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
