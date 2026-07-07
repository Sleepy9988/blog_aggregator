package main

import (
	"blog_aggregator/internal/database"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func handlerFeed(s *state, cmd command) error {
	current_user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed_params := database.CreateFeedsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    current_user.ID,
	}

	feed, err := s.db.CreateFeeds(context.Background(), feed_params)
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	feed_follow_params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    current_user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feed_follow_params)
	if err != nil {
		return fmt.Errorf("Could not generate feed follow - %w", err)
	}
	fmt.Println("Feed created successfully:")
	printFeed(feed, current_user)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=============================================")

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:        %s\n", user.Name)
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Could not get feeds %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))

	for _, feed := range feeds {
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(feed.Name_2.String)
		fmt.Println("================================================")
	}
	return nil
}

func handlerFollowFeed(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	current_user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	url := cmd.Args[0]
	feed_url, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	feed_follow_params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    current_user.ID,
		FeedID:    feed_url.ID,
	}

	feed_follow, err := s.db.CreateFeedFollow(context.Background(), feed_follow_params)
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed follow created:")
	printFeedFollow(feed_follow.UserName, feed_follow.FeedName)

	return nil
}

func handlerFeedFollowing(s *state, cmd command) error {
	current_user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), current_user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("Feed follows for user %s\n", current_user.Name)
	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed.FeedName)
	}
	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:		%s\n", username)
	fmt.Printf("* Feed:		%s\n", feedname)
}
