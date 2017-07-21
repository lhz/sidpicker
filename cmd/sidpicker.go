package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/lhz/sidpicker/config"
	"github.com/lhz/sidpicker/hvsc"
	"github.com/lhz/sidpicker/player"
	"github.com/lhz/sidpicker/ui"
)

var workerGroup sync.WaitGroup

func main() {
	reindexFlag := flag.Bool("i", false, "Rebuild tune index.")
	versionFlag := flag.Bool("v", false, "Output version string then exit.")
	flag.Parse()

	config.ReadConfig()

	// Output version string and exit
	if *versionFlag {
		fmt.Printf("sidpicker v%s\n", config.Version)
		os.Exit(0)
	}

	// Rebuild tunes index and exit
	if *reindexFlag {
		hvsc.BuildTunesIndex()
		fmt.Printf("Indexed %d tunes.\n", hvsc.NumTunes)
		os.Exit(0)
	}

	hvsc.ReadTunesIndex()

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
