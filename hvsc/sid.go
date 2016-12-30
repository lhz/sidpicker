package hvsc

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lhz/sidpicker/config"
	"github.com/lhz/sidpicker/csdb"
)

type SidHeader struct {
	MagicID     string `json:"magic,omitempty"`
	Version     int    `json:"version,omitempty"`
	DataOffset  uint16 `json:"offset,omitempty"`
	LoadAddress uint16 `json:"load,omitempty"`
	InitAddress uint16 `json:"init,omitempty"`
	PlayAddress uint16 `json:"play,omitempty"`
	Songs       int    `json:"songs,omitempty"`
	StartSong   int    `json:"start,omitempty"`
	Speed       uint32 `json:"speed,omitempty"`
	Name        string `json:"name,omitempty"`
	Author      string `json:"author,omitempty"`
	Released    string `json:"released,omitempty"`
	Flags       uint16 `json:"flags,omitempty"`
	StartPage   byte   `json:"page,omitempty"`
	PageLength  byte   `json:"pages,omitempty"`
	Sid2Address uint16 `json:"s2addr,omitempty"`
	Sid3Address uint16 `json:"s3addr,omitempty"`
}

type SidTune struct {
	Index       int             `json:"-"`
	Path        string          `json:"path"`
	SongLengths []time.Duration `json:"lengths"`
	Info        []string        `json:"info,omitempty"`
	Releases    []csdb.Release  `json:"releases,omitempty"`
	YearMin     int             `json:"year"`
	YearMax     int             `json:"ymax,omitempty"`
	Header      SidHeader       `json:"header"`
}

func (tune *SidTune) FullPath() string {
	return fmt.Sprintf("%s/%s", config.Config.HvscBase, tune.Path)
}

func (tune *SidTune) Filename() string {
	return filepath.Base(tune.Path)
}

func (tune *SidTune) Title() string {
	if len(tune.Header.Name) > 0 {
		return tune.Header.Name
	} else {
		return DefaultTitle
	}
}

func (tune *SidTune) ListName() string {
	if len(tune.Header.Name) > 0 {
		return tune.Header.Name
	} else {
		name := tune.Filename()
		ext := filepath.Ext(name)
		return name[0 : len(name)-len(ext)]
	}
}

func (tune *SidTune) InfoFilterText() string {
	return strings.Join(tune.Info, " ")
}

func (tune *SidTune) ReleasesFilterText() string {
	text := bytes.Buffer{}
	for _, r := range tune.Releases {
		if text.Len() > 0 {
			text.WriteString(", ")
		}
		text.WriteString(r.Name)
		if len(r.Groups) > 0 {
			text.WriteString(strings.Join(r.Groups, ", "))
		}
	}
	return text.String()
}

func (tune *SidTune) Year() string {
	return tune.Header.Released[0:4]
}

func (tune *SidTune) CalcYearMin() int {
	value := strings.Split(tune.Header.Released, " ")[0]
	value = strings.Split(value, "-")[0]
	value = strings.Replace(value, "?", "0", -1)
	v, err := strconv.Atoi(value)
	if err != nil {
		v = 1900
	}
	return v
}

func (tune *SidTune) CalcYearMax() int {
	value := strings.Split(tune.Header.Released, " ")[0]
	values := strings.Split(value, "-")
	value = values[len(values)-1]
	value = strings.Replace(value, "?", "9", -1)
	if len(value) == 2 {
		if value[0] < '7' {
			value = "20" + value
		} else {
			value = "19" + value
		}
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		v = 9999
	}
	return v
}
