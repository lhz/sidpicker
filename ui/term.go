package ui

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/lhz/sidpicker/csdb"
	"github.com/lhz/sidpicker/hvsc"
	"github.com/lhz/sidpicker/player"
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

func Setup(initialSearch string) {
	err := termbox.Init()
	if err != nil {
		log.Panicln(err)
	}

	w, h = termbox.Size()

	if len(initialSearch) > 0 {
		searchTerm = []rune(initialSearch)
		hvsc.Filter(string(searchTerm))
	}

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
			case termbox.EventMouse:
				mouseEvent(ev)
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
		player.Play(list.CurrentItem().TuneIndex, -1)
		sort.Sort(csdb.ByDate(player.CurrentTune.Releases))
	case termbox.KeySpace:
		list.NextTune()
		player.Play(list.CurrentItem().TuneIndex, -1)
		sort.Sort(csdb.ByDate(player.CurrentTune.Releases))
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

func mouseEvent(ev termbox.Event) {
	x := ev.MouseX
	y := ev.MouseY
	if x < 38 && y > 0 && y < h - 1 {
		if list.TuneAtPos(y - 1) {
			player.Play(list.CurrentItem().TuneIndex, -1)
			sort.Sort(csdb.ByDate(player.CurrentTune.Releases))
		}
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
		header = fmt.Sprintf("Search: %s", string(searchTerm))
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

	writeAt(ox, oy+6, "STIL Information:", termbox.ColorYellow|termbox.AttrBold, termbox.ColorBlack)

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
	oy := 8
	mx := 0
	if len(tune.Info) > 0 {
		oy = oy + 3 + len(tune.Info)
	}

	//fg := termbox.ColorDefault
	bg := termbox.ColorDefault

	line := fmt.Sprintf("Known Releases: %d", len(tune.Releases))
	writeAt(ox, oy, line, termbox.ColorYellow|termbox.AttrBold, termbox.ColorBlack)

	oy += 2
	y := oy

	for _, r := range tune.Releases {

		credits := make([]string, 0)
		if r.Date != "" {
			credits = append(credits, r.Date)
		}
		if len(r.Groups) > 0 {
			credits = append(credits, strings.Join(r.Groups, ", "))
		}

		bylines := lineSplit(strings.Join(credits, " by "), 36)
		if y+len(bylines)+1 > h-2 {
			ox += mx + 3
			y = oy
			mx = 0
		}

		writeAt(ox, y, r.Name, termbox.ColorWhite, bg)
		if len(r.Name) > mx {
			mx = len(r.Name)
		}

		for j, byline := range bylines {
			writeAt(ox, y+j+1, byline, termbox.ColorMagenta, bg)
			if len(byline) > mx {
				mx = len(byline)
			}
		}

		url := r.URL()
		writeAt(ox, y+len(bylines)+1, url, termbox.ColorBlue, bg)
		if len(url) > mx {
			mx = len(url)
		}

		y += 3 + len(bylines)
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

func lineSplit(s string, w int) []string {
	words := strings.Fields(s)
	lines := make([]string, 0)
	line := bytes.Buffer{}
	for _, word := range words {
		if line.Len() > 0 && line.Len()+len(word) > w {
			lines = append(lines, line.String())
			line = bytes.Buffer{}
		}
		if line.Len() > 0 {
			line.WriteRune(' ')
		}
		line.WriteString(word)
	}
	lines = append(lines, line.String())
	return lines
}
