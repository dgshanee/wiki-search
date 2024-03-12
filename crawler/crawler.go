package crawler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	client   *http.Client
	selector string
	regex    *regexp.Regexp
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

func (c *Crawler) Crawl(url string, maxCon int, depth int) ([]WordData, error) {

	if depth >= 5 {
		return nil, nil
	}
	res, err := c.fetch(fmt.Sprintf("https://en.wikipedia.org%s", url))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	//This is where we start to search
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxCon)

	//This highlights the main wikipedia page, not references
	mainBody := doc.Find("#bodyContent").Not(".reflist, .refbegin").First()
	mainBody.Find(c.selector).Each(func(i int, g *goquery.Selection) {
		//TODO: ADD EACH WORD TO DATABASE HERE
		g.Find("a").Not(".reference").Each(func(i int, ga *goquery.Selection) {
			if ga.ParentFiltered(".reference, .mw-editsection").Length() == 0 {
				href, ok := ga.Attr("href")
				//This gets all the links on the wikipedia article that aren't references
				if ok && c.validUrl(href) {
					wg.Add(1)
					semaphore <- struct{}{}

					go func() {
						time.Sleep(2 * time.Millisecond)
						defer wg.Done()

						//c.Crawl(href, 500, depth+1)
						<-semaphore
					}()
				}
			}
		})
	})
	wg.Wait()
	return nil, nil
}

func (c *Crawler) validUrl(url string) bool {

	ok := c.regex.Match([]byte(url))
	if !ok {
		return false
	}

	return !strings.ContainsAny(url, ":") && !strings.ContainsAny(url, "%")
}

func announceCall(url string, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()
	fmt.Println("Travelling to ", url)
	time.Sleep(2 * time.Millisecond)
	<-semaphore
}

func NewCrawler() *Crawler {
	//Initialize regex compilation to match URL
	pattern := `\/wiki\/[A-Za-z0-9_()]+`
	regex, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	selectorString := "p,h1,h2,h3,ul"
	newClient := &http.Client{
		Timeout: 2 * time.Second,
	}

	return &Crawler{
		client:   newClient,
		selector: selectorString,
		regex:    regex,
	}
}
