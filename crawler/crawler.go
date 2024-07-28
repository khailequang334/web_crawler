package crawler

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/khailequang334/web_crawler/database"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/91.0.864.64",
}

func RandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

var count = 1

func ScrapeContent(content string) {
	// fmt.Println("Content:", content)
	filename := "file_" + strconv.Itoa(count) + ".txt"
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
	count++
}

func DiscoverUrl(url string, domain string) []string {
	collectedUrls := []string{}
	c := colly.NewCollector(
		colly.UserAgent(RandomUserAgent()),
		colly.AllowedDomains(domain),
	)
	c.SetRequestTimeout(100 * time.Second)

	// Detect new urls
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)
		if strings.HasPrefix(absoluteURL, "https://"+domain) {
			fmt.Println("Content:", absoluteURL)
			collectedUrls = append(collectedUrls, absoluteURL)
		}
	})

	// Scrape content
	c.OnHTML("div.mw-parser-output > p", func(e *colly.HTMLElement) {
		ScrapeContent(e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		time.Sleep(2 * time.Second)
	})

	// c.OnResponse(func(r *colly.Response) {
	// 	fmt.Println("Got a response from", r.Request.URL)
	// })

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got error:", e)
	})

	c.Visit(url)

	return collectedUrls
}

func Worker(url string, domain string, workData chan<- []string, db *database.MongoDB) {
	collectedUrls := DiscoverUrl(url, domain)

	workData <- collectedUrls

}

func StartCrawler(db *database.MongoDB) {
	workData := make(chan []string, 100) // Buffered channel
	var wg sync.WaitGroup

	visitedUrls := make(map[string]bool)

	domain := "en.wikipedia.org"
	baseUrls := []string{
		"https://en.wikipedia.org/wiki/FIFA_World_Cup",
	}

	// Add base urls to channel
	go func() {
		workData <- baseUrls
	}()

	for urls := range workData {
		for _, url := range urls {
			if !visitedUrls[url] {
				visitedUrls[url] = true
				wg.Add(1)
				go func(url string) {
					defer wg.Done()
					Worker(url, domain, workData, db)
				}(url)
			}
		}
	}

	// Wait for all remaining workers to finish
	wg.Wait()
}
