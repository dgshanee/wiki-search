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
				result, err := c.Crawl(url, v)
				if err != nil {
					b.Fatal(err)
				}
				blackhole = result
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
					t.Fatalf("Unexpected error at url %s", tc.url)
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
