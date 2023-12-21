package crawler

import (
	"io"
	"log"
	"net/http"
)

func getSiteHtml(url string) ([]byte, error) {
	log.SetPrefix("network GetSiteHtml: ")
	log.SetFlags(0)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		return nil, err
	}
	defer resp.Body.Close()
	return body, nil
}

type Article struct {
	url     string
	header  string
	content []string
}
