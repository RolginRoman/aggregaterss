package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/rolginroman/aggregaterss/internal/database"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequests time.Duration) {
	log.Printf("Scraping on %v goroutines every %v duration", concurrency, timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Error fetching feeds for scrape %v", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()

	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error mark feed as fetched. Error: %v", err)
		return
	}

	rssfeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error getting rss feed by url %v. Error: %v", feed.Url, err)
		return
	}
	for _, item := range rssfeed.Channel.Item {
		log.Println("Found post: ", item.Title, "on feed", feed.Name)
	}
	log.Printf("Feed %v collected, %v posts found", feed.Name, len(rssfeed.Channel.Item))

}
