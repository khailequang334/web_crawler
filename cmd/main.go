package main

import (
	"github.com/khailequang334/web_crawler/crawler"
)

func main() {
	// Connect to MongoDB
	// db, err := database.ConnectMongoDB()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Disconnect()

	// Start the web crawler
	crawler.StartCrawler(nil)
}
