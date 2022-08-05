package repos

import (
	"fmt"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/sumonado/duckover/pkg/cover"
	"github.com/sumonado/duckover/pkg/duckstation"
	"github.com/sumonado/duckover/pkg/requests"
)

var (
	psxDataCenterURL = "https://psxdatacenter.com"
)

type PSXDataCenter struct {
	URL         string
	Catalog     map[string]cover.Cover
	DuckStation *duckstation.DuckStation
}

func NewPSXDataCenter(ds *duckstation.DuckStation) *PSXDataCenter {
	return &PSXDataCenter{
		DuckStation: ds,
		URL:         psxDataCenterURL,
		Catalog:     GenerateCatalog(),
	}
}

func GenerateCatalog() map[string]cover.Cover {
	var catalog = map[string]cover.Cover{}

	var regions = []string{"J", "P", "U"}
	var letters = strings.Split(`0-9,A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z`, ",")

	var wg sync.WaitGroup
	wg.Add(len(regions) * len(letters))

	covers := make(chan cover.Cover)

	go func(catalog *map[string]cover.Cover, ch chan cover.Cover) {
		for cover := range ch {
			(*catalog)[cover.Name] = cover
		}
	}(&catalog, covers)

	for _, region := range regions {
		for _, letter := range letters {
			url := fmt.Sprintf("%s/games/%s/%s", psxDataCenterURL, region, letter)

			go func(region, letter string, covers chan cover.Cover, wg *sync.WaitGroup) {
				defer wg.Done()
				response := requests.Get(url)

				regex, _ := regexp.Compile(`href="(\w{4}-\d{4,})`)
				matches := regex.FindAllStringSubmatch(response, -1)

				for _, match := range matches {
					serial := match[1]
					covers <- cover.Cover{
						Name: serial,
						Path: fmt.Sprintf(
							"%s/images/covers/%s/%s/%s.jpg",
							psxDataCenterURL, region, letter, serial,
						),
					}
				}
			}(region, letter, covers, &wg)

		}
	}

	wg.Wait()
	return catalog
}

func (p *PSXDataCenter) Download(serial string) error {
	dst := path.Join(
		p.DuckStation.GetCoversFolder(),
		fmt.Sprintf("%s.jpg", serial),
	)

	if cover, ok := p.Catalog[serial]; ok {
		requests.Download(cover.Path, dst)
	}

	return nil
}
