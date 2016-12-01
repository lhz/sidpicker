package hvsc

import (
	"bufio"
	"log"
	"os"
)

const SongInfoFile = "DOCUMENTS/STIL.txt"

func readSTIL() {
	file, err := os.Open(hvscPathTo(SongInfoFile))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var tune *SidTune
	var info []string
	var tuneIndex = -1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if line[0] == '/' && line[len(line)-4:] == ".sid" {
			tuneIndex = tuneIndexByPath(line) //, tuneIndex+1)
			if tuneIndex < 0 {
				log.Fatalf("Unknown path in file %s: %s", SongInfoFile, line)
			}
			if tune != nil {
				tune.Info = info
			}
			tune = &Tunes[tuneIndex]
			info = make([]string, 0)
		} else {
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
			log.Printf("Found path at index %d from index %d: %s [%s]", index+i, index, path, Tunes[index].Path)
			return index + i
		}
	}
	log.Fatalf("Unknown path from index %d: %s [%s]", index, path, Tunes[index].Path)
	return -1
}
