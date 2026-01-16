package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/seancrowe1/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	dur, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't parse duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", dur)

	ticker := time.NewTicker(dur)

	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("couldn't scrape feed: %w", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	DBfeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: DBfeed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't mark feed as fetched: %w", err)
	}

	feed, err := fetchFeed(context.Background(), DBfeed.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	for _, item := range feed.Channel.Item {
		fmt.Printf(" * Title: %s\n", item.Title)
	}

	return nil
}

func printFeed(f *RSSFeed) {
	fmt.Printf(" * Title: 		%v\n", f.Channel.Title)
	fmt.Printf(" * Link: 		%v\n", f.Channel.Link)
	fmt.Printf(" * Description: %v\n\n", f.Channel.Description)
	for _, item := range f.Channel.Item {
		fmt.Printf(" ** Title: 			%v\n", item.Title)
		fmt.Printf(" ** Link: 			%v\n", item.Link)
		fmt.Printf(" ** Description: 	\n%v\n", item.Description)
		fmt.Printf(" ** PubDate: 		%v\n\n", item.PubDate)
	}
}
