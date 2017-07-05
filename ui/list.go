package ui

import (
	"math/rand"
	"path/filepath"
	"time"

	"github.com/lhz/sidpicker/hvsc"
)

type List struct {
	Items    []ListItem
	PageNum  int
	PagePos  int
	PageSize int
}

type ListItem struct {
	Type      int
	TuneIndex int
	Name      string
}

const (
	ITEM_FOLDER = iota
	ITEM_TUNE
)

func NewList(pageSize int) *List {
	items := make([]ListItem, 0)
	lastPath := ""
	for i, tune := range hvsc.FilteredTunes {
		path := filepath.Dir(tune.Path)
		if path != lastPath {
			lastPath = path
			items = append(items, ListItem{Type: ITEM_FOLDER, TuneIndex: -1, Name: path})
		}
		items = append(items, ListItem{Type: ITEM_TUNE, TuneIndex: i, Name: tune.Header.Name})
	}
	return &List{Items: items, PageNum: 0, PagePos: 1, PageSize: pageSize}
}

func (l *List) CurrentItem() ListItem {
	return l.Items[l.pos()]
}

func (l *List) CurrentPage() []ListItem {
	if l.pageOffset()+l.PageSize > l.maxPos() {
		return l.Items[l.pageOffset():len(l.Items)]
	} else {
		return l.Items[l.pageOffset() : l.pageOffset()+l.PageSize]
	}
}

func (l *List) PrevPage() {
	if l.PageNum == 0 {
		return
	}
	l.PageNum--
	if l.CurrentItem().Type != ITEM_TUNE {
		if l.pos() == 0 {
			l.NextTune()
		} else {
			l.PrevTune()
		}
	}
}

func (l *List) NextPage() {
	if l.pageOffset()+l.PageSize > l.maxPos() {
		return
	}
	l.PageNum++
	if l.pos() > l.maxPos() {
		l.PagePos = l.maxPos() - l.pageOffset()
	}
	if l.CurrentItem().Type != ITEM_TUNE {
		l.NextTune()
	}
}

func (l *List) PrevTune() {
	if l.pos() <= 1 {
		return
	}
	l.PagePos--
	if l.PagePos < 0 {
		if l.PageNum > 0 {
			l.PagePos = l.PageSize - 1
			l.PrevPage()
		} else {
			l.PagePos = 0
		}
	}
	if l.CurrentItem().Type != ITEM_TUNE {
		l.PrevTune()
	}
}

func (l *List) NextTune() {
	if l.pos() == l.maxPos() {
		return
	}
	l.PagePos++
	if l.PagePos >= l.PageSize {
		l.PagePos = 0
		l.NextPage()
	}
	if l.CurrentItem().Type != ITEM_TUNE {
		l.NextTune()
	}
}

func (l *List) TuneAtPos(pagePos int) bool {
	pos := l.pageOffset() + pagePos
	if pos > l.maxPos() {
		return false
	}
	if l.Items[pos].Type != ITEM_TUNE {
		return false
	}
	l.PagePos = pagePos
	return true
}

func (l *List) RandomTune() {
	rand.Seed(time.Now().Unix())
	for {
		n := rand.Intn(l.maxPos() + 1)
		if n != l.pos() && l.Items[n].Type == ITEM_TUNE {
			l.PageNum = n / l.PageSize
			l.PagePos = n - l.pageOffset()
			break
		}
	}
}

func (l *List) pos() int {
	return l.pageOffset() + l.PagePos
}

func (l *List) maxPos() int {
	return len(l.Items) - 1
}

func (l *List) pageOffset() int {
	return l.PageNum * l.PageSize
}
