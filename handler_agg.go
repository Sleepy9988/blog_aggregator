package main

import (
	"blog_aggregator/internal/database"
	"context"
	"fmt"
	"log"
	"time"
)

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.Args) != 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %s <time_between_reqs", cmd.Name)
	}

	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Invalid duration: %w", err)
	}

	//feedURL := "https://www.wagslane.dev/index.xml"

	log.Printf("Collecting feeds every %s...", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	next_feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Could not fetch the next feed %w", err)
		return
	}

	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, next_feed)

}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Could not collect  feed %w", err)
		return
	}

	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}
