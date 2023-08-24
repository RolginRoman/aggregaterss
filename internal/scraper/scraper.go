package scraper

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/rolginroman/aggregaterss/internal/database"
)

func StartScraping(db *database.Queries, concurrency int, timeBetweenRequests time.Duration) {
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

	// TODO check the most recent posts and don't recreate all posts if it is older than last_fetched_at
	for _, item := range rssfeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			pubDate, err = time.Parse(time.RFC1123, item.PubDate)
			if err != nil {
				log.Printf("Cannot parse PublishedAt date %v with error %v", pubDate, err)
				continue
			}
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubDate,
			Url:         item.Link,
			FeedID:      feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Failed to create post %v with error %v", item.Link, err)
			continue
		}
	}
	log.Printf("Feed %v collected, %v posts found", feed.Name, len(rssfeed.Channel.Item))

}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	httpResponse, err := httpClient.Get(url)

	if err != nil {
		return RSSFeed{}, errors.New(fmt.Sprintf("Cannot parse RSS Feed with url %v. Error %v", url, err))
	}

	defer httpResponse.Body.Close()

	data, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return RSSFeed{}, errors.New(fmt.Sprintf("Cannot process RSS Feed with url %v. Error %v", url, err))
	}
	rssFeed := RSSFeed{}

	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return RSSFeed{}, errors.New(fmt.Sprintf("Cannot handle RSS Format from url %v. Error %v", url, err))
	}

	return rssFeed, nil
}
