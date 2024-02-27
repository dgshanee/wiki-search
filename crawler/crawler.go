package crawler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	client   *http.Client
	selector string
}

type WordData struct {
	Word string
	Url  string
}

func (c *Crawler) fetch(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "crawler-name")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Crawler) Crawl(url string) ([]WordData, error) {
	res, err := c.fetch(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	mainBody := doc.Find("#bodyContent").Not(".reflist, .refbegin").First()

	mainBody.Find(c.selector).Each(func(i int, g *goquery.Selection) {
		g.Find("a").Not(".reference").Each(func(i int, ga *goquery.Selection) {
			if ga.ParentFiltered(".reference, .mw-editsection").Length() == 0 {
				fmt.Print(strings.TrimSpace(ga.Text()))
			}
		})
		fmt.Println()
	})

	return nil, nil
}

func NewCrawler() *Crawler {
	selectorString := "p,h1,h2,h3,ul"
	newClient := &http.Client{
		Timeout: 2 * time.Second,
	}
	return &Crawler{
		client:   newClient,
		selector: selectorString,
	}
}
