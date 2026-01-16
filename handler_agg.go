package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
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
		_, err = s.db.GetPostByUrl(context.Background(), item.Link)
		if err == nil {
			continue
		}
		pubDate := time.Now()
		pubDate, err = time.Parse(time.ANSIC, item.PubDate)
		if err != nil {
			pubDate, err = time.Parse(time.UnixDate, item.PubDate)
			if err != nil {
				pubDate, err = time.Parse(time.RubyDate, item.PubDate)
				if err != nil {
					pubDate, err = time.Parse(time.RFC822, item.PubDate)
					if err != nil {
						pubDate, err = time.Parse(time.RFC822Z, item.PubDate)
						if err != nil {
							pubDate, err = time.Parse(time.RFC850, item.PubDate)
							if err != nil {
								pubDate, err = time.Parse(time.RFC1123, item.PubDate)
								if err != nil {
									pubDate, err = time.Parse(time.RFC1123Z, item.PubDate)
									if err != nil {
										pubDate, err = time.Parse(time.RFC3339, item.PubDate)
										if err != nil {
											pubDate, err = time.Parse(time.RFC3339Nano, item.PubDate)
											if err != nil {
												return err
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubDate,
			FeedID:      DBfeed.ID,
		})
		if err != nil {
			return fmt.Errorf("couldn't create post: %w", err)
		}
	}

	return nil
}
