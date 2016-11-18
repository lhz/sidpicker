package hvsc

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	Sid2Address byte
	Sid3Address byte
}

type SidTune struct {
	Path        string
	MD5         string
	NumSongs    int
	SongLengths []time.Duration
	Header      SidHeader
}

var header = make([]byte, 124)

// Read tunes data from cache file
func ReadTunesInfoCached() []SidTune {
	if _, err := os.Stat(hvscPathTo(TunesCacheFile)); os.IsNotExist(err) {
		return ReadTunesInfo()
	}

	content, err := ioutil.ReadFile(hvscPathTo(TunesCacheFile))
	if err != nil {
		log.Fatal(err)
	}

	sidTunes := make([]SidTune, 0)
	json.Unmarshal(content, &sidTunes)

	return sidTunes
}

// Build tunes data from .sid-files and various documents
func ReadTunesInfo() []SidTune {
	sidTunes := make([]SidTune, 0)

	file, err := os.Open(hvscPathTo(SongLengthsFile))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == ';' {
			tune := SidTune{Path: line[2:]}
			tune.Header = ReadSidHeader(hvscPathTo(tune.Path))
			sidTunes = append(sidTunes, tune)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Read %d tunes from file %s.\n", len(sidTunes), file.Name())

	b, err := json.MarshalIndent(sidTunes, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	jsonFile, err := os.Create(hvscPathTo(TunesCacheFile))
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(b)

	return sidTunes
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

	h := SidHeader{}

	enc := binary.BigEndian

	h.MagicID = string(header[0:4])
	h.Version = int(enc.Uint16(header[4:]))
	h.DataOffset = enc.Uint16(header[6:])
	h.LoadAddress = enc.Uint16(header[8:])
	h.InitAddress = enc.Uint16(header[10:])
	h.PlayAddress = enc.Uint16(header[12:])
	h.Songs = int(enc.Uint16(header[14:]))
	h.StartSong = int(enc.Uint16(header[16:]))
	h.Speed = enc.Uint32(header[18:])
	h.Name = stringExtract(header[22:54])
	h.Author = stringExtract(header[54:86])
	h.Released = stringExtract(header[86:118])

	return h
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
