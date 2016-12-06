package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"encoding/json"

	"github.com/PuerkitoBio/goquery"
)

const basePath = "/home/lars/c64/sidinfo"

type Release struct {
	CSDbId int
	Name   string
	Groups string
	Year   string
}

var sidUses map[string][]Release

func main() {
	if len(os.Args) < 1 {
		log.Fatal("Usage: usedex <maxId>")
	}

	maxId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Usage: usedex <maxId>")
	}

	sidUses = make(map[string][]Release, 0)

	for i := 1; i <= maxId; i++ {
		path := fmt.Sprintf("%s/html/%05d/%05d.html", basePath, (i / 100) * 100, i)
		if _, err = os.Stat(path); os.IsNotExist(err) {
			//log.Fatalf("No such file: %s", path)
			continue
		}
		html, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal("File %s could not be read: %v", path, err)
		}
		if len(html) < 256 {
			continue
		}
		parseContent(bytes.NewReader(html))
	}

	json, _ := json.MarshalIndent(sidUses, "", "  ")
	fmt.Printf(string(json))
}

func parseContent(r io.Reader) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	s := doc.Find("a[href*='http://hvsc.perff.dk/']").First()
	if s.Length() != 1 {
		log.Fatal("No hvsc path found.")
	}
	path := s.Text()
	releases := make([]Release, 0)
	doc.Find("tr td a[href*='/release/?id=']").Each(func(i int, s *goquery.Selection) {
		release := Release{Name: s.Text()}
		if href, ok := s.Attr("href"); ok {
			release.CSDbId, _ = strconv.Atoi(href[13:])
		}
		extras := s.Parent().Parent().Find("td font").Map(func(i int, s *goquery.Selection) string {
			return s.Text()
		})
		release.Groups = strings.Join(extras[:len(extras)-1], ", ")
		release.Year   = extras[len(extras)-1]
		releases = append(releases, release)
	})
	if len(releases) > 0 {
		sidUses[path] = releases
	}
}
