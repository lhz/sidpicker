package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lhz/sidpicker/config"
	"github.com/lhz/sidpicker/csdb"
)

const lastReleaseId = 156656               // TODO: Store in file
const csdbPath = "/home/lars/src/lhz/csdb" // TODO: Add to config

// Check existing XML files and make list of releases that
// need to be updated due to non-existing paths

func main() {
	config.ReadConfig()

	updates := make([]int, 0)
	for i := 1; i <= lastReleaseId; i++ {
		if i%10000 == 0 {
			log.Print(i)
		}
		file := fmt.Sprintf("%s/xml/%06d/%06d.xml", csdbPath, i/100*100, i)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}
		release := csdb.ReadReleaseXML(file)
		for _, path := range release.SIDs {
			file = fmt.Sprintf("%s/%s", config.Config.HvscBase, path)
			if _, err := os.Stat(file); os.IsNotExist(err) {
				log.Printf("Release %d needs update due to non-existing path: %s\n", release.ID, path)
				updates = append(updates, release.ID)
				break
			}
		}
	}
	log.Printf("%d releases need to be updated.", len(updates))
	for _, id := range updates {
		fmt.Println(id)
	}
}
