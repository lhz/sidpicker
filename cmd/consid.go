package main

import (
	"github.com/lhz/considerate/config"
	"github.com/lhz/considerate/hvsc"
	"github.com/lhz/considerate/player"
	"github.com/lhz/considerate/ui"

	"log"
	"sync"
)

var workerGroup sync.WaitGroup

func main() {
	config.ReadConfig()

	hvsc.ReadTunesInfoCached()
	log.Printf("Read %d tunes.", hvsc.NumTunes)

	workerGroup.Add(1)
	go func() {
		defer workerGroup.Done()
		player.Run()
	}()

	ui.Setup()
	ui.Run()

	player.Quit()

	workerGroup.Wait()
}
