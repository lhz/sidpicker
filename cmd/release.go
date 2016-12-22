package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lhz/sidpicker/csdb"
)

func main() {
	release := csdb.ReadRelease(os.Args[1])
	str, _ := json.MarshalIndent(release, "", "  ")
	fmt.Printf("%s\n", str)
}
