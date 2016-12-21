package hvsc

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const ReleasesUsedFile = "releases.json.gz"

type Release struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
	Year  string `json:"year"`
}

func (r *Release) URL() string {
	return fmt.Sprintf("http://csdb.dk/release/?id=%d", r.Id)
}

func readReleases() {
	if _, err := os.Stat(hvscPathTo(ReleasesUsedFile)); os.IsNotExist(err) {
		log.Println(err)
		return
	}

	log.Print("Reading sid release usage info.")

	dataGzip, err := ioutil.ReadFile(hvscPathTo(ReleasesUsedFile))
	if err != nil {
		log.Fatal(err)
	}
	r, err := gzip.NewReader(bytes.NewBuffer(dataGzip))
	if err != nil {
		log.Fatal(err)
	}

	uses := make(map[string][]Release, NumTunes)
	err = json.NewDecoder(r).Decode(&uses)
	if err != nil {
		log.Fatalf("Failed to read sid release used file: %s", err)
	}
	log.Printf("Release usage read for %d tunes.", len(uses))

	for path, releases := range uses {
		tuneIndex := TuneIndexByPath(path)
		if tuneIndex < 0 {
			log.Fatalf("Unknown path in file %s: %s", ReleasesUsedFile, path)
		}
		Tunes[tuneIndex].Releases = releases
	}
}
