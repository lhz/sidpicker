package ui

import (
	"path/filepath"

	"github.com/lhz/considerate/hvsc"
)

type List struct {
	Items      []ListItem
	PageOffset int
	PagePos    int
	PageSize   int
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
	return &List{Items: items, PageOffset: 0, PagePos: 1, PageSize: pageSize}
}

func (l *List) CurrentItem() ListItem {
	return l.Items[l.PageOffset+l.PagePos]
}

func (l *List) CurrentPage() []ListItem {
	if l.PageOffset+l.PageSize > l.maxPos() {
		return l.Items[l.PageOffset:len(l.Items)]
	} else {
		return l.Items[l.PageOffset : l.PageOffset+l.PageSize]
	}
}

func (l *List) PrevPage() {
	l.PageOffset -= l.PageSize
	if l.PageOffset < 0 {
		l.PageOffset = 0
	}
	if l.CurrentItem().Type != ITEM_TUNE {
		if l.pos() == 0 {
			l.NextTune()
		} else {
			l.PrevTune()
		}
	}
}

func (l *List) NextPage() {
	if l.PageOffset+l.PageSize > l.maxPos() {
		return
	}
	l.PageOffset += l.PageSize
	if l.pos() > l.maxPos() {
		l.PagePos = l.maxPos() - l.PageOffset
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
		if l.PageOffset > 0 {
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

func (l *List) pos() int {
	return l.PageOffset + l.PagePos
}

func (l *List) maxPos() int {
	return len(l.Items) - 1
}
