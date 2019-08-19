package fetcher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func New(url string) Fetcher {
	return &fetcher{
		url,
	}
}

type Fetcher interface {
	Fetch() ([]byte, error)
}

type fetcher struct {
	url string
}

func (f *fetcher) Fetch() ([]byte, error) {
	resp, err := http.Get(f.url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Server returns %d instead of %d", resp.StatusCode, http.StatusOK))
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Failed to close 'response.Body'. Err: ", err)
		}
	}()

	return ioutil.ReadAll(resp.Body)
}
