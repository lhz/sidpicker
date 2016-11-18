package main

import (
	"github.com/lhz/considerate/config"
	"github.com/lhz/considerate/hvsc"
	"github.com/lhz/considerate/player"
	"github.com/lhz/considerate/ui"

	"log"
)

func main() {
	config.ReadConfig()

	hvsc.ReadTunesInfoCached()
	log.Printf("Read %d tunes.", hvsc.NumTunes)

	player.Setup()
	go player.Run()

	ui.Setup()
	ui.Run()
}
