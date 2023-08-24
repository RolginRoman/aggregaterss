package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

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
