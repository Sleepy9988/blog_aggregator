package main

import (
	"blog_aggregator/internal/database"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
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
		log.Printf("Could not collect  feed %s", err)
		return
	}

	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}

		if parseTime, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  parseTime,
				Valid: true,
			}
		}

		_, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: sql.NullTime{
				Time:  time.Now().UTC(),
				Valid: true,
			},
			FeedID: feed.ID,
			Title:  item.Title,
			Url:    item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "dupicate key value violates unique constraint") {
				continue
			}
			log.Printf("Could not create post: %v", err)
			continue
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}
