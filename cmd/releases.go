package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/lhz/sidpicker/csdb"
)

func main() {
	releases := make([]csdb.Release, 0)
	for i := 1; i < 153000; i++ {
		if i%100 == 0 {
			log.Print(i)
		}
		file := fmt.Sprintf("/home/lars/c64/csdb/releases/xml/%06d/%06d.xml", i/100*100, i)
		if _, err := os.Stat(file); os.IsNotExist(err) { 
			continue
		}
		release := csdb.ReadReleaseXML(file)
		if len(release.SIDs) > 0 {
			releases = append(releases, *release)
		}
	}
	str, _ := json.MarshalIndent(releases, "", "  ")
	fmt.Printf("%s\n", str)
}
