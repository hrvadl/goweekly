package main

import (
	"fmt"
	"golangweekly/network"
)

func main() {
	body, _ := network.GetSiteHtml("https://golangweekly.com/issues/489")
	fmt.Println(body)
}
