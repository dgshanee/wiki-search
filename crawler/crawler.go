package crawler

import (
	"fmt"
	"net/http"
	"sync"
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

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 100)

	startTime := time.Now()
	mainBody := doc.Find("#bodyContent").Not(".reflist, .refbegin").First()
	mainBody.Find(c.selector).Each(func(i int, g *goquery.Selection) {
		//Index by word here
		g.Find("a").Not(".reference").Each(func(i int, ga *goquery.Selection) {
			if ga.ParentFiltered(".reference, .mw-editsection").Length() == 0 {
				redirect, ok := ga.Attr("href")
				//This gets all the links on the wikipedia article
				//Do something with this
				if ok {
					wg.Add(1)
					semaphore <- struct{}{}

					go announceCall(redirect, &wg, semaphore)
				}
			}
		})
	})
	wg.Wait()
	endTime := time.Now()
	fmt.Println("All routines finished in ", endTime.Sub(startTime), " seconds")
	return nil, nil
}

func announceCall(url string, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()
	fmt.Println("Travelling to ", url)
	time.Sleep(2 * time.Millisecond)
	<-semaphore
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
