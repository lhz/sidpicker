package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	//"github.com/ghodss/yaml"
	"github.com/PuerkitoBio/goquery"
)

const basePath = "/home/lars/c64/sidinfo"

type Release struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
	Year  string `json:"year"`
}

var sidUses map[string][]Release

func main() {
	minId := 1
	maxId := 99999
	var err error

	if len(os.Args) > 1 {
		maxId, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal("Usage: usedex [id]")
		}
		minId = maxId
	}

	sidUses = make(map[string][]Release, 0)

	for i := minId; i <= maxId; i++ {
		path := fmt.Sprintf("%s/html/%05d/%05d.html", basePath, (i/100)*100, i)
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

	json, err := json.MarshalIndent(sidUses, "", "  ")
	//json, err := json.Marshal(sidUses)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %s", err)
	}
	fmt.Printf(string(json))
	//yaml, _ := yaml.Marshal(sidUses)
	//fmt.Printf(string(yaml))
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
	doc.Find("table[cellpadding='0'] tr td a[href*='/release/?id=']").Each(func(i int, s *goquery.Selection) {
		release := Release{Name: sanitize(s.Text())}
		if href, ok := s.Attr("href"); ok {
			release.Id, _ = strconv.Atoi(href[13:])
		}
		groups := s.Parent().Parent().Find("td font[color='#2575ff']").Map(func(i int, s *goquery.Selection) string {
			return sanitize(s.Text())
		})
		year := sanitize(s.Parent().Parent().Find("td font[size='1']").First().Text())
		if len(groups) > 0 || year != "" {
			release.Group = strings.Join(groups, ", ")
			release.Year = year
			releases = append(releases, release)
		}
	})
	if len(releases) > 0 {
		sidUses[path] = releases
	}
}

func sanitize(s string) string {
	s = strings.Replace(s, "%", "%%", -1)
	return s
}
