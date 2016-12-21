package hvsc

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/beevik/etree"
)

const ReleasesUsedFile = "releases.json.gz"

type Release struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
	Date  string `json:"date"`
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

func readReleaseXML(path) {
	doc := etree.NewDocument()
	err := doc.ReadFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	tunePaths := doc.FindElements("/CSDbData/Release/UsedSIDs/SID/HVSCPath")
	if len(tunePaths) < 1 {
		return
	}

	rel := doc.SelectElement("CSDbData/Release")

}

func parseRelease(e *etree.Element) *Release {
	var err error
	r := Release{}
	r.Id, err = strconv.Atoi(e.SelectElement("ID").Text())
	r.Name = e.SelectElement("Name").Text()
	r.Type = e.SelectElement("Type").Text()
	date := bytes.Buffer{}
	y := e.SelectElement("ReleaseYear")
	if y != nil {
		date.WriteString(y.Text())
	}
	m := e.SelectElement("ReleaseMonth")
	if m != nil {
		date.WriteString(fmt.Sprintf("-%02s", m.Text()))
	} else {
		date.WriteString("-xx")
	}
	d := e.SelectElement("ReleaseDay")
	if d != nil {
		date.WriteString(fmt.Sprintf("-%02s", d.Text()))
	} else {
		date.WriteString("-xx")
	}
	r.Date = date.String()
}
