package hvsc

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lhz/considerate/config"
)

const (
	SongLengthsFile = "DOCUMENTS/Songlengths.txt"
	TunesCacheFile  = "cache-tunes.json"
)

type SidHeader struct {
	MagicID     string
	Version     int
	DataOffset  uint16
	LoadAddress uint16
	InitAddress uint16
	PlayAddress uint16
	Songs       int
	StartSong   int
	Speed       uint32
	Name        string
	Author      string
	Released    string
	Flags       uint16
	StartPage   byte
	PageLength  byte
	Sid2Address uint16
	Sid3Address uint16
}

type SidTune struct {
	Index       int
	Path        string
	MD5         string
	SongLengths []time.Duration
	YearMin     int
	YearMax     int
	Header      SidHeader
}

func (tune *SidTune) FullPath() string {
	return fmt.Sprintf("%s/%s", config.Config.HvscPath, tune.Path)
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

	b, err := json.MarshalIndent(Tunes, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	jsonFile, err := os.Create(hvscPathTo(TunesCacheFile))
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(b)
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
	return fmt.Sprintf("%s/%s", config.Config.HvscPath, filePath)
}

func parseSongLength(value string) time.Duration {
	parts := strings.Split(value, ":")
	dur, err := time.ParseDuration(fmt.Sprintf("%sm%ss", parts[0], parts[1]))
	if err != nil {
		return 0
	}
	return dur
}
