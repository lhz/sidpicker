package main

import (
	"os"
	"strings"

	"github.com/lhz/sidpicker/config"
	"github.com/lhz/sidpicker/hvsc"
	"github.com/lhz/sidpicker/player"
	"github.com/lhz/sidpicker/ui"

	"sync"
)

var workerGroup sync.WaitGroup

func main() {
	config.ReadConfig()

	hvsc.ReadTunesInfoCached()
	//log.Printf("Read %d tunes.", hvsc.NumTunes)

	workerGroup.Add(1)
	go func() {
		defer workerGroup.Done()
		player.Run()
	}()

	searchTerm := strings.Join(os.Args[1:], " ")
	ui.Setup(searchTerm)

	ui.Run()

	player.Quit()

	workerGroup.Wait()
}
