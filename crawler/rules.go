package crawler

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type ContentParser struct {
	Selector string
	Callback func(e *colly.HTMLElement)
}

type LinkCollector struct {
	Selector string
	Callback func(e *colly.HTMLElement) string
}

type CrawlerRule struct {
	Domain         string
	BaseUrls       []string
	ContentParsers []ContentParser
	LinkCollectors []LinkCollector
}

var CrawlerRules = map[string]CrawlerRule{
	"WikiContent": {
		Domain: "en.wikipedia.org",
		BaseUrls: []string{
			"https://en.wikipedia.org/wiki/FIFA_World_Cup",
		},
		ContentParsers: []ContentParser{
			{
				Selector: "div.mw-parser-output > p",
				Callback: WikiContentParser,
			},
		},
		LinkCollectors: []LinkCollector{
			{
				Selector: "a[href]",
				Callback: WikiLinkColector,
			},
		},
	},
	// Add more
}

func WikiContentParser(e *colly.HTMLElement) {
	content := e.Text
	// New timestamp-based filename
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("logfile-%s.txt", timestamp)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Write the content to the file
	if _, err := file.WriteString(content); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("Content appended successfully to: %s\n", filename)
}

func WikiLinkColector(e *colly.HTMLElement) string {
	link := e.Attr("href")
	absoluteURL := e.Request.AbsoluteURL(link)
	if strings.HasPrefix(absoluteURL, "https://en.wikipedia.org/wiki") {
		return absoluteURL
	}
	return ""
}
