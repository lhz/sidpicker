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
var list *List
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

	list = NewList(h - 2)

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
		if player.Playing {
			player.Stop()
		} else {
			quit = true
		}
	case termbox.KeyPgup:
		list.PrevPage()
	case termbox.KeyPgdn:
		list.NextPage()
	case termbox.KeyArrowUp:
		list.PrevTune()
	case termbox.KeyArrowDown:
		list.NextTune()
	case termbox.KeyArrowLeft:
		player.PrevSong()
	case termbox.KeyArrowRight:
		player.NextSong()
	case termbox.KeyEnter:
		player.Play(list.CurrentItem().TuneIndex, 1)
	case termbox.KeySpace:
		list.NextTune()
		player.Play(list.CurrentItem().TuneIndex, 1)
	case termbox.KeyDelete:
		player.Stop()
	case 0:
		if n := strings.IndexRune("0123456789", ev.Ch); n > 0 {
			if player.Playing && n <= player.CurrentTune.Header.Songs {
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
		if hvsc.Filter(string(searchTerm)) {
			list = NewList(h - 2)
			mode = MODE_BROWSE
		}
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

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawHeader()
	drawList()
	drawTuneInfo()
	drawReleases()
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
		bg = termbox.ColorMagenta
		termbox.SetCursor(8+searchCursorPos, 0)
	}
	header = fmt.Sprintf(fmt.Sprintf("%%s%%%ds", w-len([]rune(header))), header, "")
	writeAt(0, 0, header, termbox.ColorWhite|termbox.AttrBold, bg)
}

func drawFooter() {
	footer := fmt.Sprintf("%d/%d  %s", list.CurrentItem().TuneIndex+1, hvsc.NumFilteredTunes, debugInfo)
	footer = fmt.Sprintf(fmt.Sprintf("%%s%%%ds", w-len([]rune(footer))), footer, "")
	writeAt(0, h-1, footer, termbox.ColorWhite, termbox.ColorBlack)
}

func drawTuneInfo() {
	if !player.Playing {
		return
	}
	tune := player.CurrentTune
	ox := 42
	oy := 2
	fg := termbox.ColorDefault
	bg := termbox.ColorDefault

	writeAt(ox, oy+0, fmt.Sprintf("Title:    %s", tune.Title()), fg, bg)
	writeAt(ox, oy+1, fmt.Sprintf("Author:   %s", tune.Header.Author), fg, bg)
	writeAt(ox, oy+2, fmt.Sprintf("Released: %s", tune.Header.Released), fg, bg)

	writeAt(ox, oy+4, fmt.Sprintf("Tune: %d/%d  Length: %s  Time: %s",
		player.CurrentSong, tune.Header.Songs, player.SongLength(), player.Elapsed()), fg, bg)

	if len(tune.Info) == 0 {
		return
	}

	writeAt(ox, oy+6, " STIL: ", termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)

	for i, line := range tune.Info {
		if oy+8+i > h-2 {
			break
		}
		writeAt(ox, oy+8+i, line, fg, bg)
	}
}

func drawReleases() {
	if !player.Playing {
		return
	}
	tune := player.CurrentTune
	if len(tune.Releases) == 0 {
		return
	}

	ox := 42
	oy := 11 + len(tune.Info)
	//fg := termbox.ColorDefault
	bg := termbox.ColorDefault

	writeAt(ox, oy, " RELEASES: ", termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)

	for i, r := range tune.Releases {
		if oy+2+i*4+2 > h-2 {
			break
		}
		credits := make([]string, 0)
		if r.Year != "" {
			credits = append(credits, r.Year)
		}
		if r.Group != "" {
			credits = append(credits, r.Group)
		}
		writeAt(ox, oy+2+i*4, r.Name, termbox.ColorWhite, bg)
		writeAt(ox, oy+2+i*4+1, strings.Join(credits, " by "), termbox.ColorMagenta, bg)
		writeAt(ox, oy+2+i*4+2, r.URL(), termbox.ColorBlue, bg)
	}
}

func drawList() {
	for i, item := range list.CurrentPage() {
		fg := termbox.ColorDefault
		bg := termbox.ColorDefault
		if i == list.PagePos {
			fg = termbox.ColorWhite
			bg = termbox.ColorBlue
		}
		if item.Type == ITEM_TUNE {
			tune := hvsc.FilteredTunes[item.TuneIndex]
			if player.CurrentTune != nil && player.CurrentTune.Index == tune.Index {
				fg |= termbox.AttrBold
			}
			writeAt(0, 1+i, fmt.Sprintf(" %-32s %4s", tune.ListName(), tune.Year()), fg, bg)
		} else {
			folder := []rune(item.Name[1:])
			if len(folder) > 40 {
				folder = folder[len(folder)-40 : len(folder)]
				folder[0] = '…'
			}
			writeAt(0, 1+i, string(folder), fg|termbox.AttrBold, bg)
		}
	}
}

func writeAt(x, y int, value string, fg, bg termbox.Attribute) {
	i := 0
	for _, c := range value {
		termbox.SetCell(x+i, y, c, fg, bg)
		i++
	}
}
