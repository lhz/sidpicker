package hvsc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const ReleasesUsedFile = "sid-releases-used.json"

type Release struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
	Year  string `json:"year"`
}

func readReleases() {
	if _, err := os.Stat(hvscPathTo(ReleasesUsedFile)); os.IsNotExist(err) {
		log.Println(err)
		return
	}

	log.Print("Reading sid release usage info.")
	content, err := ioutil.ReadFile(hvscPathTo(ReleasesUsedFile))
	if err != nil {
		log.Fatal(err)
	}

	uses := make(map[string][]Release, NumTunes)
	err = json.Unmarshal(content, &uses)
	if err != nil {
		log.Fatalf("Failed to read sid release used file: %s", err)
	}
	log.Printf("Release usage read for %d tunes (%d bytes)", len(uses), len(content))

	for path, releases := range uses {
		tuneIndex := TuneIndexByPath(path)
		if tuneIndex < 0 {
			log.Fatalf("Unknown path in file %s: %s", ReleasesUsedFile, path)
		}
		Tunes[tuneIndex].Releases = releases
	}
}
