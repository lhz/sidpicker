package hvsc

import (
	"github.com/lhz/considerate/cfg"

	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	SongLengthsFile = "DOCUMENTS/Songlengths.txt"
)


type SidTune struct {
	Path        string
	MD5         string
	NumSongs    int
	SongLengths []time.Duration
}

var sidTunes []SidTune

func BuildTuneInfo(config *cfg.Config) []SidTune {
	sidTunes := make([]SidTune, 0)

	file, err := os.Open(fmt.Sprintf("%s/%s", config.HvscPath, SongLengthsFile))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == ';' {
			tune := SidTune{Path: line[2:]}
			sidTunes = append(sidTunes, tune)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Read %d tunes from file %s.\n", len(sidTunes), file.Name())
	log.Printf("First tune: %q\n", sidTunes[0].Path)

	return sidTunes
}
