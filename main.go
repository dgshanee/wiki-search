package main

import (
	"github.com/dgshanee/search-engine-demo/crawler"
)

func main() {
	crl := crawler.NewCrawler()

	crl.Crawl("https://en.wikipedia.org/wiki/Go_(programming_language)")
}
