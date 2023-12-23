package crawler

import (
	"io"
	"log"
	"net/http"
)

type Article struct {
	Url     string
	Header  string
	Content []string
}

func getSiteHTML(url string) ([]byte, error) {
	log.SetPrefix("network GetSiteHtml: ")
	log.SetFlags(0)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	return body, nil
}
