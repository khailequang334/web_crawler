package crawler

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/gocolly/colly"
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

func DiscoverUrl(url string, rule CrawlerRule) []string {
	collectedUrls := []string{}
	c := colly.NewCollector(
		colly.UserAgent(RandomUserAgent()),
		colly.AllowedDomains(rule.Domain),
	)

	c.SetRequestTimeout(100 * time.Second)

	// Detect new urls
	for _, collector := range rule.LinkCollectors {
		c.OnHTML(collector.Selector, func(h *colly.HTMLElement) {
			collectedUrls = append(collectedUrls, collector.Callback(h))
		})
	}

	// Scrape content
	for _, parser := range rule.ContentParsers {
		c.OnHTML(parser.Selector, parser.Callback)
	}

	c.OnRequest(func(r *colly.Request) {})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got error:", e)
	})

	c.Visit(url)

	time.Sleep(500 * time.Millisecond)

	return collectedUrls
}

func StartCrawler(ruleType string) {
	rule, exists := CrawlerRules[ruleType]
	if !exists {
		fmt.Printf("Crawler rule type %s not found!\n", ruleType)
		return
	}

	workData := make(chan []string, 100) // Buffered channel
	var wg sync.WaitGroup

	visitedUrls := make(map[string]bool)

	// Add base urls to channel
	workData <- rule.BaseUrls

	for urls := range workData {
		for _, url := range urls {
			if !visitedUrls[url] {
				visitedUrls[url] = true
				wg.Add(1)
				go func(url string) {
					defer wg.Done()
					workData <- DiscoverUrl(url, rule)
				}(url)
			}
		}
	}

	// Wait for all remaining workers to finish
	go func() {
		wg.Wait()
		close(workData)
	}()
}
