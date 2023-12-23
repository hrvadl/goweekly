package main

import (
	"fmt"

	"github.com/hrvadl/go-weekly/crawler"
)

func main() {
	arts, err := crawler.CrawlSite("https://golangweekly.com/issues/489")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(arts)
}
