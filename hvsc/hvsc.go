package hvsc

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lhz/sidpicker/config"
)

const (
	SongLengthsFile = "DOCUMENTS/Songlengths.txt"
	TunesCacheFile  = "cache-tunes.json"
	DefaultTitle    = "<?>"
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
	Releases    []Release       `json:"releases,omitempty"`
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
		if len(r.Group) > 0 {
			text.WriteString(", ")
			text.WriteString(r.Group)
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

var Tunes = make([]SidTune, 0)
var NumTunes = 0

var header = make([]byte, 124)

func TuneIndexByPath(path string) int {
	for i, tune := range Tunes {
		if tune.Path == path {
			return i
		}
	}
	return -1
}

// Read tunes data from cache file
func ReadTunesInfoCached() {
	if _, err := os.Stat(hvscPathTo(TunesCacheFile)); os.IsNotExist(err) {
		ReadTunesInfo()
		return
	}

	log.Print("Reading cached tunes info.")
	content, err := ioutil.ReadFile(hvscPathTo(TunesCacheFile))
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(content, &Tunes)
	NumTunes = len(Tunes)

	addDefaults()

	FilterAll()
}

// Build tunes data from .sid-files and various documents
func ReadTunesInfo() {
	file, err := os.Open(hvscPathTo(SongLengthsFile))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lr := regexp.MustCompile("[0-9]{1,2}:[0-9]{2}")

	log.Print("Building tunes info cache.")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == ';' {
			tune := SidTune{Index: len(Tunes), Path: line[2:]}
			tune.Header = ReadSidHeader(hvscPathTo(tune.Path))
			tune.SongLengths = make([]time.Duration, tune.Header.Songs)
			tune.YearMin = tune.CalcYearMin()
			tune.YearMax = tune.CalcYearMax()
			Tunes = append(Tunes, tune)
		} else {
			lengths := lr.FindAllString(line, -1)
			for i, l := range lengths {
				Tunes[len(Tunes)-1].SongLengths[i] = parseSongLength(l)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	NumTunes = len(Tunes)
	FilterAll()

	readSTIL()
	readReleases()

	removeDefaults()

	b, err := json.MarshalIndent(Tunes, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	cacheFile, err := os.Create(hvscPathTo(TunesCacheFile))
	if err != nil {
		log.Fatal(err)
	}
	defer cacheFile.Close()

	cacheFile.Write(b)
}

func ReadSidHeader(fileName string) SidHeader {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Read(header)
	if err != nil {
		log.Fatal(err)
	}

	enc := binary.BigEndian

	h := SidHeader{
		MagicID:     string(header[0:4]),
		Version:     int(enc.Uint16(header[4:])),
		DataOffset:  enc.Uint16(header[6:]),
		LoadAddress: enc.Uint16(header[8:]),
		InitAddress: enc.Uint16(header[10:]),
		PlayAddress: enc.Uint16(header[12:]),
		Songs:       int(enc.Uint16(header[14:])),
		StartSong:   int(enc.Uint16(header[16:])),
		Speed:       enc.Uint32(header[18:]),
		Name:        stringExtract(header[22:54]),
		Author:      stringExtract(header[54:86]),
		Released:    stringExtract(header[86:118]),
		Flags:       enc.Uint16(header[118:]),
		StartPage:   header[120],
		PageLength:  header[121],
	}
	if header[122] > 0 {
		h.Sid2Address = uint16(header[122])*16 + 0xD000
	}
	if header[123] > 0 {
		h.Sid3Address = uint16(header[123])*16 + 0xD000
	}
	return h
}

func parseYear(value string, defVal int) int {
	year, err := strconv.Atoi(value)
	if err != nil {
		return defVal
	}
	if year < 100 {
		if year < 70 {
			year += 2000
		} else {
			year += 1900
		}
	}
	return year
}

func stringExtract(slice []byte) string {
	codePoints := make([]rune, len(slice))
	pos := 0
	for ; pos < len(slice) && slice[pos] != 0; pos++ {
		codePoints[pos] = rune(slice[pos])
	}
	return string(codePoints[:pos])
}

func hvscPathTo(filePath string) string {
	return fmt.Sprintf("%s/%s", config.Config.HvscBase, filePath)
}

func parseSongLength(value string) time.Duration {
	parts := strings.Split(value, ":")
	dur, err := time.ParseDuration(fmt.Sprintf("%sm%ss", parts[0], parts[1]))
	if err != nil {
		return 0
	}
	return dur
}

// Set default tune/header fields to empty values to reduce marshalling size
func removeDefaults() {
	for i, tune := range Tunes {
		if tune.YearMax == tune.YearMin {
			tune.YearMax = 0
		}
		if tune.Header.MagicID == "PSID" {
			tune.Header.MagicID = ""
		}
		if tune.Header.Version == 2 {
			tune.Header.Version = 0
		}
		if tune.Header.DataOffset == 124 {
			tune.Header.DataOffset = 0
		}
		if tune.Header.Songs == 1 {
			tune.Header.Songs = 0
		}
		if tune.Header.StartSong == 1 {
			tune.Header.StartSong = 0
		}
		if tune.Header.Name == "<?>" {
			tune.Header.Name = ""
		}
		if tune.Header.Author == "<?>" {
			tune.Header.Author = ""
		}
		if tune.Header.Released == "<?>" {
			tune.Header.Released = ""
		}
		Tunes[i] = tune
	}
}

// Set empty tune/header fields to default values after unmarshalling
func addDefaults() {
	for i, tune := range Tunes {
		tune.Index = i
		if tune.YearMax == 0 {
			tune.YearMax = tune.YearMin
		}
		if tune.Header.MagicID == "" {
			tune.Header.MagicID = "PSID"
		}
		if tune.Header.Version == 0 {
			tune.Header.Version = 2
		}
		if tune.Header.DataOffset == 0 {
			tune.Header.DataOffset = 124
		}
		if tune.Header.Songs == 0 {
			tune.Header.Songs = 1
		}
		if tune.Header.StartSong == 0 {
			tune.Header.StartSong = 1
		}
		Tunes[i] = tune
	}
}
