package main

import (
	"fmt"
	"sync"

	"github.com/sumonado/duckover/pkg/duckstation"
	"github.com/sumonado/duckover/pkg/repos"
)

var (
	downloader repos.Repository
)

func main() {

	var wg sync.WaitGroup

	duckstation := duckstation.NewDuckStation()
	downloader = repos.NewPSXDataCenter(duckstation)

	games := duckstation.GetGames()

	wg.Add(len(games))

	for _, serial := range games {
		fmt.Printf("Downloading %s\n", serial)

		go func(serial string, wg *sync.WaitGroup) {
			defer wg.Done()
			downloader.Download(serial)
		}(serial, &wg)
	}

	wg.Wait()
}
