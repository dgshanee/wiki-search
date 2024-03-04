package crawler

import (
	"fmt"
	"os"
	"testing"
)

var c *Crawler

func TestMain(m *testing.M) {
	c = NewCrawler()
	os.Exit(m.Run())
}

var blackhole []WordData

func Benchmark_crawl(b *testing.B) {
	for _, v := range []int{1, 100, 400, 900, 1000, 2000, 5000, 10000} {
		b.Run(fmt.Sprintf("Crawl-%d", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var url string = "https://www.wikipedia.org/wiki/Go_(programming_language)"
				result, err := c.Crawl(url, v, 0)
				if err != nil {
					b.Fatal(err)
				}
				blackhole = result
			}
		})
	}
}

func Test_valiateUrl(t *testing.T) {
	tests := []struct {
		name string
		url  string
		pass bool
	}{
		{"golang", "/wiki/Go_(programming_language)", true},
		{"shorkie", "/wiki/Yorkshire_Terrier", true},
		{"category", "/Category/all_dog_breeds:", false},
		{"no-slashes", "wiki/Yorkshire_Terrier", false},
		{"no-slashes-2", "wikiYorkshire_Terrier", false},
		{"file", "/wiki/File:poop.jpg", false},
		{"categiry", "/wiki/Category:poop", false},
		{"different language lol", "https://af.wikipedia.org/wiki/Go_(programmeertaal)", false},
		{"percent", "/wiki/go%ejdeadl", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ok := c.validUrl(tc.url)

			if ok && !tc.pass {
				t.Errorf("Unexpected pass at url %s", tc.url)
			}

			if !ok && tc.pass {
				t.Errorf("Expected pass at url %s, got fail", tc.url)
			}
		})
	}
}

func Test_fetch(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		errExpected bool
	}{
		{"google", "https://www.google.com", false},
		{"golang", "https://www.wikipedia.com/wiki/Go_(programming_language)", false},
		{"poop fail", "poop", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := c.fetch(tc.url)
			if err != nil {
				if !tc.errExpected {
					t.Fatalf("Unexpected error at url %s, %v", tc.url, err)
				}

				return
			}
			if tc.errExpected && err == nil {
				t.Errorf("Expected error at url %s", tc.url)
				return
			}
			if result.StatusCode != 200 {
				t.Errorf("Unstable connection at %s", tc.url)
				return
			}

		})
	}
}

func Test_crawl(t *testing.T) {
	url := "/wiki/Go_(programming_language)"
	_, err := c.Crawl(url, 2, 0)
	if err != nil {
		t.Error(err)
	}
}
