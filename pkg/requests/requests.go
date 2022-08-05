package requests

import (
	"io"
	"net/http"
	"os"

	"github.com/sumonado/duckover/pkg/helpers"
)

func Get(url string) (response string) {
	resp, err := http.Get(url)
	helpers.HandleError(err)
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	helpers.HandleError(err)

	response = string(content)

	return
}

func Download(url string, dst string) (size int64, err error) {
	resp, err := http.Get(url)
	helpers.HandleError(err)
	defer resp.Body.Close()

	file, err := os.Create(dst)
	defer file.Close()
	helpers.HandleError(err)

	size, err = io.Copy(file, resp.Body)

	return
}
