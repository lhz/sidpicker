package csdb

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
	ID     int      `json:"id"`
	Type   string   `json:"type"`
	Name   string   `json:"name"`
	Date   string   `json:"date"`
	Groups []string `json:"groups"`
	SIDs   []string `json:"sids"`
	elem   *etree.Element
}

func ReadReleases() {
	if _, err := os.Stat(hvsc.hvscPathTo(ReleasesUsedFile)); os.IsNotExist(err) {
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

	releases := make(map[string][]Release, 0)
	err = json.NewDecoder(r).Decode(&releases)
	if err != nil {
		log.Fatalf("Failed to read sid release used file: %s", err)
	}
	log.Printf("Read %d releases.", len(releases))

	for _, release := range release {
		for _, path := range release.SIDs {
			tuneIndex := TuneIndexByPath(path)
			if tuneIndex < 0 {
				log.Fatalf("Unknown path: %s", path)
			}
			Tunes[tuneIndex].Releases = append(Tunes[tuneIndex].Releases, release)
		}
	}
}

func ReadReleaseXML(path string) *Release {
	doc := etree.NewDocument()
	err := doc.ReadFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	elements := doc.FindElements("//CSDbData/Release")
	if len(elements) != 1 {
		log.Fatal("Expected 1 Release element, found %v", len(elements))
	}
	r := Release{elem: elements[0]}
	r.ID, _ = r.getInt("ID")
	r.Name, _ = r.getString("Name")
	r.Type, _ = r.getString("Type")
	r.Date = r.getDate()
	r.Groups = r.getGroups()
	r.SIDs = r.getSIDs()
	return &r
}

func (r *Release) getInt(name string) (int, error) {
	e := r.elem.SelectElement(name)
	if e == nil {
		return 0, fmt.Errorf("Element not set: %s", name)
	}
	v, err := strconv.Atoi(e.Text())
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (r *Release) getString(name string) (string, error) {
	e := r.elem.SelectElement(name)
	if e == nil {
		return "", fmt.Errorf("Element not set: %s", name)
	}
	return e.Text(), nil
}

func (r *Release) getDate() string {
	date := bytes.Buffer{}
	y, err := r.getInt("ReleaseYear")
	if err != nil {
		date.WriteString("????")
	} else {
		date.WriteString(fmt.Sprintf("%04d", y))
	}
	m, err := r.getInt("ReleaseMonth")
	if err != nil {
		date.WriteString("-??")
	} else {
		date.WriteString(fmt.Sprintf("-%02d", m))
	}
	d, err := r.getInt("ReleaseDay")
	if err != nil {
		date.WriteString("-??")
	} else {
		date.WriteString(fmt.Sprintf("-%02d", d))
	}
	return date.String()
}

func (r *Release) getGroups() []string {
	elems := r.elem.FindElements("ReleasedBy/Group/Name")
	//log.Printf("Found %d group names.", len(elems))
	groups := make([]string, 0)
	for _, e := range elems {
		groups = append(groups, e.Text())
	}
	return groups
}

func (r *Release) getSIDs() []string {
	elems := r.elem.FindElements("UsedSIDs/SID/HVSCPath")
	//log.Printf("Found %d used sids.", len(elems))
	sids := make([]string, 0)
	for _, e := range elems {
		sids = append(sids, e.Text())
	}
	return sids
}
