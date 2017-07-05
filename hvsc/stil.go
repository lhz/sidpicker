package hvsc

import (
	"bufio"
	"log"
	"strings"

	"github.com/lhz/sidpicker/util"
)

const SongInfoFile = "DOCUMENTS/STIL.txt"

func readSTIL() {
	content, err := util.ReadLatin1File(hvscPathTo(SongInfoFile))
	if err != nil {
		log.Fatalf("Failed to read STIL file: %s", err)
	}

	var tune *SidTune
	var info []string
	var tuneIndex = -1

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if line[0] == '/' {
			if tune != nil {
				tune.Info = info
				tune = nil
			}
			if line[len(line)-4:] == ".sid" {
				tuneIndex = tuneIndexByPath(line) //, tuneIndex+1)
				if tuneIndex < 0 {
					log.Fatalf("Unknown path in file %s: %s", SongInfoFile, line)
				}
				tune = &Tunes[tuneIndex]
				info = make([]string, 0)
			}
		} else if tune != nil {
			info = append(info, line)
		}
	}
}

func tuneIndexByPath(path string) int {
	for i, tune := range Tunes {
		if tune.Path == path {
			return i
		}
	}
	return -1
}

func tuneIndexByPathIndex(path string, index int) int {
	for i, tune := range Tunes[index:] {
		if tune.Path == path {
			//log.Printf("Found path at index %d from index %d: %s [%s]", index+i, index, path, Tunes[index].Path)
			return index + i
		}
	}
	log.Fatalf("Unknown path from index %d: %s [%s]", index, path, Tunes[index].Path)
	return -1
}
