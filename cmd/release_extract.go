package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/lhz/sidpicker/config"
	"github.com/lhz/sidpicker/csdb"
)

const lastReleaseId = 156656               // TODO: Store in file
const csdbPath = "/home/lars/src/lhz/csdb" // TODO: Add to config

func main() {
	config.ReadConfig()

	releases := make([]csdb.Release, 0)
	for i := 1; i <= lastReleaseId; i++ {
		if i%10000 == 0 {
			log.Print(i)
		}
		file := fmt.Sprintf("%s/xml/%06d/%06d.xml", csdbPath, i/100*100, i)
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
