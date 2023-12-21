package network

import (
	"io"
	"log"
	"net/http"
)

func GetSiteHtml(url string) ([]byte, error) {
	log.SetPrefix("network GetSiteHtml: ")
	log.SetFlags(0)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to request site over http. Error: %v", err)
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Response body read fail. Error: %v", err)
		return nil, err
	}
	if resp.StatusCode > 299 {
		log.Printf("Response returned with failed status. code %d body %s", resp.StatusCode, body)
		return nil, err
	}
	defer resp.Body.Close()
	return body, nil
}
