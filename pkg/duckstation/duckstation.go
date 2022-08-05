package duckstation

import (
	"os"
	"path"
	"regexp"

	"github.com/sumonado/duckover/pkg/helpers"
)

type DuckStation struct {
	Folder string
}

func NewDuckStation() *DuckStation {
	homeDir, err := os.UserHomeDir()
	helpers.HandleError(err)

	folder := path.Join(homeDir, "Documents", "DuckStation")
	return &DuckStation{Folder: folder}
}

func (d *DuckStation) GetCoversFolder() string {
	return path.Join(d.Folder, "covers")
}

func (d *DuckStation) GetCacheContent() string {
	cachePath := path.Join(d.Folder, "cache", "gamelist.cache")

	_, err := os.Stat(cachePath)
	helpers.HandleError(err)

	bytes, err := os.ReadFile(cachePath)
	helpers.HandleError(err)

	return string(bytes)
}

func (d *DuckStation) GetGames() (games []string) {
	regex, _ := regexp.Compile(`\x00{3}(\w{4}-\d{4,})`)
	matches := regex.FindAllStringSubmatch(d.GetCacheContent(), -1)

	for _, match := range matches {
		games = append(games, match[1])
	}

	return
}
